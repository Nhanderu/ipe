package main

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
	"github.com/Nhanderu/trena"
	"github.com/Nhanderu/tuyo/text"
	"github.com/fatih/color"
)

const (
	kilobyte = 1024
	megabyte = kilobyte * 1024
	gigabyte = megabyte * 1024
	terabyte = gigabyte * 1024
)

var (
	gridView, treeView, longView, longTreeView bool
	bgstMode, bgstSize, bgstUser, bgstInode    int
	srcs                                       []srcInfo
	outBuffer                                  bytes.Buffer
	direction                                  gridt.Direction

	osWindows = runtime.GOOS == "windows"
)

type srcInfo struct {
	file   ipe.File
	grid   *gridt.Grid
	buffer *bytes.Buffer
}

func (s srcInfo) writec(column func(ipe.File) string, sep string, size int, conds ...bool) {
	s.buffer.WriteString(fmtColumn(column(s.file), sep, size, conds...))
}

func (s srcInfo) write(str string) {
	s.buffer.WriteString(str)
}

func main() {
	args := parseArgs()

	if args.color != colorAuto {
		color.NoColor = args.color == colorNever
	}

	if args.across {
		direction = gridt.LeftToRight
	} else {
		direction = gridt.TopToBottom
	}

	gridView = !args.long && !args.tree
	treeView = !args.long && args.tree
	longView = args.long && !args.tree
	longTreeView = args.long && args.tree

	srcs = make([]srcInfo, 0)
	width, _, err := trena.Size()
	if err != nil {
		endWithErr(err.Error())
	}

	for _, src := range args.source {
		src = fixSrc(src)
		f, err := ipe.Read(src)
		if err != nil {
			endWithErr(err.Error())
		}
		printDir(src, f, []bool{}, args)
	}

	if args.tree {
		os.Stdout.WriteString(outBuffer.String())
	} else {
		writeNames := len(srcs) > 1
		for _, src := range srcs {
			if writeNames {
				os.Stdout.WriteString(src.file.FullName())
				os.Stdout.WriteString("\n")
			}
			if gridView {
				g, ok := src.grid.FitIntoWidth(width)
				if !ok || args.oneLine {
					for _, cell := range src.grid.Cells() {
						os.Stdout.WriteString(cell)
						os.Stdout.WriteString("\n")
					}
				} else {
					os.Stdout.WriteString(g.String())
				}
			} else {
				os.Stdout.WriteString(src.buffer.String())
			}
			if writeNames {
				os.Stdout.WriteString("\n")
			}
		}
	}
}

func printDir(src string, file ipe.File, corners []bool, args argsInfo) {
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}

	var buffer *bytes.Buffer
	if args.tree {
		buffer = &outBuffer
	} else {
		buffer = bytes.NewBuffer([]byte{})
	}
	grid := gridt.New(direction, args.separator)
	srcs = append(srcs, srcInfo{file, grid, buffer})

	if args.reverse {
		reverse(fs)
	}

	// First loop: preparation.
	for _, f := range fs {
		checkBiggestValues(f, args)
	}

	// Second loop: printing.
	for ii, f := range fs {
		printFile(srcInfo{f, grid, buffer}, append(corners, ii+1 == len(fs)), args)
	}
}

func printFile(src srcInfo, corners []bool, args argsInfo) {
	if !show(src.file, args) {
		return
	}

	var name string
	if args.classify {
		name = src.file.ClassifiedName()
	} else {
		name = src.file.Name()
	}

	if !gridView {
		if args.long {
			if args.inode {
				src.writec(fmtInode, args.separator, bgstInode, !osWindows)
			}
			src.writec(fmtMode, args.separator, bgstMode)
			src.writec(fmtSize, args.separator, bgstSize)
			src.writec(fmtTime, args.separator, 0)
			src.writec(fmtUser, args.separator, bgstUser, !osWindows)
		}
		if args.tree {
			src.write(makeTree(corners))
		}
		src.write(name)
		src.write("\n")
	} else {
		src.grid.Add(name)
	}

	if args.recursive && src.file.IsDir() && (args.depth == 0 || args.depth >= len(corners)) {
		printDir(src.file.Name(), src.file, corners, args)
	}
}

func checkBiggestValues(f ipe.File, args argsInfo) {
	if !show(f, args) {
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
	if args.recursive {
		for _, ff := range f.Children() {
			checkBiggestValues(ff, args)
		}
	}
}

func show(f ipe.File, args argsInfo) bool {
	return (args.all || !f.IsDotfile()) &&
		(args.ignore == nil || !args.ignore.MatchString(f.Name())) &&
		(args.filter == nil || args.filter.MatchString(f.Name()))
}

func fmtColumn(column, sep string, size int, conds ...bool) string {
	for _, c := range conds {
		if !c {
			return ""
		}
	}
	return addSep(text.PadLeft(column, " ", size), sep)
}

func fmtInode(f ipe.File) string {
	if osWindows {
		return ""
	}
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
	} else if s < megabyte {
		return fmt.Sprintf("%.1dKB", s/kilobyte)
	} else if s < gigabyte {
		return fmt.Sprintf("%.1dMB", s/megabyte)
	} else if s < terabyte {
		return fmt.Sprintf("%.1dGB", s/gigabyte)
	} else {
		return fmt.Sprintf("%.1dTB", s/terabyte)
	}
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

func addSep(s, sep string) string {
	return fmt.Sprint(s, sep)
}

func reverse(a []ipe.File) {
	for l, r := 0, len(a)-1; l < r; l, r = l+1, r-1 {
		a[l], a[r] = a[r], a[l]
	}
}

func endWithErr(err string) {
	os.Stderr.WriteString(os.Args[0])
	os.Stderr.WriteString(" error: ")
	os.Stderr.WriteString(err)
	os.Stderr.WriteString("\n")
	os.Exit(1)
}

func fixSrc(src string) string {
	if osWindows {
		return strings.Replace(src, "~", os.Getenv("USERPROFILE"), -1)
	}
	return src
}
