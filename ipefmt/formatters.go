package ipefmt

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
	"github.com/Nhanderu/tuyo/text"
	"github.com/fatih/color"
)

const (
	kilobyte = 1024
	megabyte = kilobyte * 1024
	gigabyte = megabyte * 1024
	terabyte = gigabyte * 1024

	osWindows = runtime.GOOS == "windows"
)

var (
	bgstMode, bgstSize, bgstUser, bgstInode int
)

type srcInfoBuffer struct {
	file   ipe.File
	err    error
	buffer *bytes.Buffer
}

type srcInfoGrid struct {
	file ipe.File
	err  error
	grid *gridt.Grid
}

func NewFormatter(args ArgsInfo) fmt.Stringer {
	if args.Color != ArgColorAuto {
		color.NoColor = args.Color == ArgColorNever
	}

	if args.Long && args.Tree {
		return newLongTreeFormatter(args)
	}
	if args.Long {
		return newLongFormatter(args)
	}
	if args.Tree {
		return newTreeFormatter(args)
	}
	return newGridFormatter(args)
}

// Grid!

type gridFormatter struct {
	srcs      []srcInfoGrid
	direction gridt.Direction
	width     int
	oneLine   bool
}

func newGridFormatter(args ArgsInfo) *gridFormatter {
	var direction gridt.Direction
	if args.Across {
		direction = gridt.LeftToRight
	} else {
		direction = gridt.TopToBottom
	}
	f := &gridFormatter{make([]srcInfoGrid, 0), direction, args.Width, args.OneLine}
	for _, src := range args.Sources {
		file, err := ipe.Read(fixInSrc(src))
		if err != nil {
			f.srcs = append(f.srcs, srcInfoGrid{file, err, nil})
		} else {
			f.getDir(file, args, 0)
		}
	}
	return f
}

func (f *gridFormatter) getDir(file ipe.File, args ArgsInfo, depth uint8) {
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}

	grid := gridt.New(f.direction, args.Separator)
	f.srcs = append(f.srcs, srcInfoGrid{file, nil, grid})

	if args.Reverse {
		reverse(fs)
	}

	// First loop: preparation.
	for _, file := range fs {
		checkBiggestValues(file, args)
	}

	// Second loop: printing.
	for _, file := range fs {
		f.getFile(file, grid, args, depth+1)
	}
}

func (f *gridFormatter) getFile(file ipe.File, grid *gridt.Grid, args ArgsInfo, depth uint8) {
	if !shouldShow(file, args) {
		return
	}

	if args.Classify {
		grid.Add(file.ClassifiedName())
	} else {
		grid.Add(file.Name())
	}

	if args.Recursive && file.IsDir() && (args.Depth == 0 || args.Depth >= depth) {
		f.getDir(file, args, depth)
	}
}

func (f *gridFormatter) String() string {
	var buffer bytes.Buffer
	writeNames := len(f.srcs) > 1
	for _, src := range f.srcs {
		if writeNames {
			buffer.WriteString(src.file.FullName())
			buffer.WriteString("\n")
		}
		if src.err != nil {
			buffer.WriteString("Error: ")
			buffer.WriteString(src.err.Error())
		} else {
			d, ok := src.grid.FitIntoWidth(f.width)
			if !ok || f.oneLine {
				for _, cell := range src.grid.Cells() {
					buffer.WriteString(cell)
					buffer.WriteString("\n")
				}
			} else {
				buffer.WriteString(d.String())
			}
		}
		if writeNames {
			buffer.WriteString("\n")
		}
	}
	return buffer.String()
}

// Tree!

type treeFormatter struct {
	srcs []srcInfoBuffer
}

func newTreeFormatter(args ArgsInfo) *treeFormatter {
	f := &treeFormatter{make([]srcInfoBuffer, 0)}
	for _, src := range args.Sources {
		file, err := ipe.Read(fixInSrc(src))
		if err != nil {
			f.srcs = append(f.srcs, srcInfoBuffer{file, err, nil})
		} else {
			f.getDir(file, bytes.NewBuffer([]byte{}), args, []bool{})
		}
	}
	return f
}

func (f *treeFormatter) getDir(file ipe.File, buffer *bytes.Buffer, args ArgsInfo, corners []bool) {
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}

	if len(corners) == 0 {
		f.srcs = append(f.srcs, srcInfoBuffer{file, nil, buffer})
	}

	if args.Reverse {
		reverse(fs)
	}

	// First loop: preparation.
	for _, file := range fs {
		checkBiggestValues(file, args)
	}

	// Second loop: printing.
	for ii, file := range fs {
		f.getFile(file, buffer, args, append(corners, ii+1 == len(fs)))
	}
}

func (f *treeFormatter) getFile(file ipe.File, buffer *bytes.Buffer, args ArgsInfo, corners []bool) {
	if !shouldShow(file, args) {
		return
	}

	buffer.WriteString(makeTree(corners))
	if args.Classify {
		buffer.WriteString(file.ClassifiedName())
	} else {
		buffer.WriteString(file.Name())
	}
	buffer.WriteRune('\n')

	if args.Recursive && file.IsDir() && (args.Depth == 0 || int(args.Depth) >= len(corners)) {
		f.getDir(file, buffer, args, corners)
	}
}

func (f *treeFormatter) String() string {
	var buffer bytes.Buffer
	writeNames := len(f.srcs) > 1
	for _, src := range f.srcs {
		if writeNames {
			buffer.WriteString(src.file.FullName())
			buffer.WriteString("\n")
		}
		if src.err != nil {
			buffer.WriteString("Error: ")
			buffer.WriteString(src.err.Error())
		} else {
			buffer.WriteString(src.buffer.String())
		}
		if writeNames {
			buffer.WriteString("\n")
		}
	}
	return buffer.String()
}

// Long!

type longFormatter struct {
	srcs []srcInfoBuffer
}

func newLongFormatter(args ArgsInfo) *longFormatter {
	f := &longFormatter{make([]srcInfoBuffer, 0)}
	for _, src := range args.Sources {
		file, err := ipe.Read(fixInSrc(src))
		if err != nil {
			f.srcs = append(f.srcs, srcInfoBuffer{file, err, nil})
		} else {
			f.getDir(file, args, 0)
		}
	}
	return f
}

func (f *longFormatter) getDir(file ipe.File, args ArgsInfo, depth uint8) {
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	f.srcs = append(f.srcs, srcInfoBuffer{file, nil, buffer})

	if args.Reverse {
		reverse(fs)
	}

	// First loop: preparation.
	for _, file := range fs {
		checkBiggestValues(file, args)
	}

	// Second loop: printing.
	for _, file := range fs {
		f.getFile(file, buffer, args, depth+1)
	}
}

func (f *longFormatter) getFile(file ipe.File, buffer *bytes.Buffer, args ArgsInfo, depth uint8) {
	if !shouldShow(file, args) {
		return
	}

	if args.Inode && !osWindows {
		buffer.WriteString(fmtColumn(fmtInode(file), args.Separator, bgstInode))
	}
	buffer.WriteString(fmtColumn(fmtMode(file), args.Separator, bgstMode))
	buffer.WriteString(fmtColumn(fmtSize(file), args.Separator, bgstSize))
	buffer.WriteString(fmtColumn(fmtTime(file), args.Separator, 0))
	if !osWindows {
		buffer.WriteString(fmtColumn(fmtUser(file), args.Separator, bgstUser))
	}
	if args.Classify {
		buffer.WriteString(file.ClassifiedName())
	} else {
		buffer.WriteString(file.Name())
	}
	buffer.WriteRune('\n')

	if args.Recursive && file.IsDir() && (args.Depth == 0 || args.Depth >= depth) {
		f.getDir(file, args, depth)
	}
}

func (f *longFormatter) String() string {
	var buffer bytes.Buffer
	writeNames := len(f.srcs) > 1
	for _, src := range f.srcs {
		if writeNames {
			buffer.WriteString(src.file.FullName())
			buffer.WriteString("\n")
		}
		if src.err != nil {
			buffer.WriteString("Error: ")
			buffer.WriteString(src.err.Error())
		} else {
			buffer.WriteString(src.buffer.String())
		}
		if writeNames {
			buffer.WriteString("\n")
		}
	}
	return buffer.String()
}

// Long Tree!

type longTreeFormatter struct {
	srcs []srcInfoBuffer
}

func newLongTreeFormatter(args ArgsInfo) *longTreeFormatter {
	f := &longTreeFormatter{make([]srcInfoBuffer, 0)}
	for _, src := range args.Sources {
		file, err := ipe.Read(fixInSrc(src))
		if err != nil {
			f.srcs = append(f.srcs, srcInfoBuffer{file, err, nil})
		} else {
			f.getDir(file, bytes.NewBuffer([]byte{}), args, []bool{})
		}
	}
	return f
}

func (f *longTreeFormatter) getDir(file ipe.File, buffer *bytes.Buffer, args ArgsInfo, corners []bool) {
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}

	if len(corners) == 0 {
		f.srcs = append(f.srcs, srcInfoBuffer{file, nil, buffer})
	}

	if args.Reverse {
		reverse(fs)
	}

	// First loop: preparation.
	for _, file := range fs {
		checkBiggestValues(file, args)
	}

	// Second loop: printing.
	for ii, file := range fs {
		f.getFile(file, buffer, args, append(corners, ii+1 == len(fs)))
	}
}

func (f *longTreeFormatter) getFile(file ipe.File, buffer *bytes.Buffer, args ArgsInfo, corners []bool) {
	if !shouldShow(file, args) {
		return
	}

	if args.Inode && !osWindows {
		buffer.WriteString(fmtColumn(fmtInode(file), args.Separator, bgstInode))
	}
	buffer.WriteString(fmtColumn(fmtMode(file), args.Separator, bgstMode))
	buffer.WriteString(fmtColumn(fmtSize(file), args.Separator, bgstSize))
	buffer.WriteString(fmtColumn(fmtTime(file), args.Separator, 0))
	if !osWindows {
		buffer.WriteString(fmtColumn(fmtUser(file), args.Separator, bgstUser))
	}
	buffer.WriteString(makeTree(corners))
	if args.Classify {
		buffer.WriteString(file.ClassifiedName())
	} else {
		buffer.WriteString(file.Name())
	}
	buffer.WriteRune('\n')

	if args.Recursive && file.IsDir() && (args.Depth == 0 || int(args.Depth) >= len(corners)) {
		f.getDir(file, buffer, args, corners)
	}
}

func (f *longTreeFormatter) String() string {
	var buffer bytes.Buffer
	writeNames := len(f.srcs) > 1
	for _, src := range f.srcs {
		if writeNames {
			buffer.WriteString(src.file.FullName())
			buffer.WriteString("\n")
		}
		if src.err != nil {
			buffer.WriteString("Error: ")
			buffer.WriteString(src.err.Error())
		} else {
			buffer.WriteString(src.buffer.String())
		}
		if writeNames {
			buffer.WriteString("\n")
		}
	}
	return buffer.String()
}

// Helpers.

func checkBiggestValues(f ipe.File, args ArgsInfo) {
	if !shouldShow(f, args) {
		return
	}
	if m := len(fmtMode(f)); m > bgstMode {
		bgstMode = m
	}
	if s := len(fmtSize(f)); s > bgstSize {
		bgstSize = s
	}
	if u := len(fmtUser(f)); u > bgstUser {
		bgstUser = u
	}
	if i := len(fmtInode(f)); i > bgstInode {
		bgstInode = i
	}
	if args.Recursive {
		for _, ff := range f.Children() {
			checkBiggestValues(ff, args)
		}
	}
}

func shouldShow(f ipe.File, args ArgsInfo) bool {
	return (args.All || !f.IsDotfile()) &&
		(args.Ignore == nil || !args.Ignore.MatchString(f.Name())) &&
		(args.Filter == nil || args.Filter.MatchString(f.Name()))
}

func fmtColumn(column, sep string, size int) string {
	var buf bytes.Buffer
	buf.WriteString(text.PadLeft(column, " ", size))
	buf.WriteString(sep)
	return buf.String()
}

func fmtInode(f ipe.File) string {
	return strconv.FormatUint(f.Inode(), 10)
}

func fmtMode(f ipe.File) string {
	return f.Mode().String()
}

func fmtSize(f ipe.File) string {
	if !f.IsRegular() {
		return "-"
	}
	s := f.Size()
	if s < kilobyte {
		return fmt.Sprintf("%dB", s)
	}
	if s < megabyte {
		return fmt.Sprintf("%.1dKB", s/kilobyte)
	}
	if s < gigabyte {
		return fmt.Sprintf("%.1dMB", s/megabyte)
	}
	if s < terabyte {
		return fmt.Sprintf("%.1dGB", s/gigabyte)
	}
	return fmt.Sprintf("%.1dTB", s/terabyte)
}

func fmtTime(f ipe.File) string {
	t := f.ModTime()
	year, month, day := t.Date()
	str := fmt.Sprintf("%2d %s ", day, month.String()[:3])
	if year == time.Now().Year() {
		return fmt.Sprintf("%s%2d:%02d", str, t.Hour(), t.Minute())
	}
	return fmt.Sprintf("%s%d ", str, year)
}

func fmtUser(f ipe.File) string {
	if osWindows {
		return ""
	}
	return f.User().Username
}

func fixInSrc(src string) string {
	if osWindows {
		return strings.Replace(src, "~", os.Getenv("USERPROFILE"), -1)
	}
	return src
}

func reverse(a []ipe.File) {
	for l, r := 0, len(a)-1; l < r; l, r = l+1, r-1 {
		a[l], a[r] = a[r], a[l]
	}
}

func makeTree(corners []bool) string {
	var s string
	arrowTree := map[bool]map[bool]string{
		true: {
			true:  "└──",
			false: "   ",
		}, false: {
			true:  "├──",
			false: "│  ",
		},
	}
	for i, c := range corners {
		s = fmt.Sprint(s, arrowTree[c][i+1 == len(corners)])
	}
	return s
}
