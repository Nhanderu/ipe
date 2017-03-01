package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"strings"

	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
	"github.com/Nhanderu/trena"
	"github.com/Nhanderu/tuyo/text"
	"github.com/fatih/color"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	kilobyte = 1024
	megabyte = kilobyte * 1024
	gigabyte = megabyte * 1024
	terabyte = gigabyte * 1024

	colorNever  = "never"
	colorAlways = "always"
	colorAuto   = "auto"
)

var (
	sourceArg     = kingpin.Arg("source", "the directory to list contents").Default(".").Strings()
	separatorFlag = kingpin.Flag("separator", "separator of the columns").Short('S').Default("  ").String()
	acrossFlag    = kingpin.Flag("across", "writes the entries by lines instead of by columns").Short('x').Bool()
	allFlag       = kingpin.Flag("all", "do not hide entries starting with .").Short('a').Bool()
	colorFlag     = kingpin.Flag("color", "control whether color is used to distinguish file types").Enum(colorNever, colorAlways, colorAuto)
	classifyFlag  = kingpin.Flag("classify", "append indicator to the entries").Short('F').Bool()
	depthFlag     = kingpin.Flag("depth", "maximum depth of recursion").Short('D').Int()
	filterFlag    = kingpin.Flag("filter", "only show entries that matches the pattern").Short('f').Regexp()
	ignoreFlag    = kingpin.Flag("ignore", "do not show entries that matches the pattern").Short('I').Regexp()
	inodeFlag     = kingpin.Flag("inode", "show entry inode").Short('i').Bool()
	longFlag      = kingpin.Flag("long", "show entries in the \"long view\"").Short('l').Bool()
	oneLine       = kingpin.Flag("one-line", "show one entry per line").Short('1').Bool()
	reverseFlag   = kingpin.Flag("reverse", "reverse order of entries").Short('r').Bool()
	recursiveFlag = kingpin.Flag("recursive", "list subdirectories recursively").Short('R').Bool()
	treeFlag      = kingpin.Flag("tree", "shows the entries in the tree view").Short('t').Bool()

	gridView                                bool
	bgstMode, bgstSize, bgstUser, bgstInode int
	srcs                                    []srcInfo
	direction                               gridt.Direction

	osWindows = runtime.GOOS == "windows"
)

type srcInfo struct {
	file   ipe.File
	grid   *gridt.Grid
	buffer *bytes.Buffer
}

func main() {
	kingpin.Parse()

	if *colorFlag != colorAuto {
		color.NoColor = *colorFlag == colorNever
	}

	if *acrossFlag {
		direction = gridt.LeftToRight
	} else {
		direction = gridt.TopToBottom
	}

	gridView = !*longFlag && !*treeFlag
	srcs = make([]srcInfo, 0)
	width, _, err := trena.Size()
	if err != nil {
		endWithErr(err.Error())
	}

	for _, src := range *sourceArg {
		src = fixSrc(src)
		f, err := ipe.Read(src)
		if err != nil {
			endWithErr(err.Error())
		}
		printDir(src, f, []bool{})
	}

	writeNames := len(srcs) > 1 && !*treeFlag
	for _, src := range srcs {
		if writeNames {
			os.Stdout.WriteString(src.file.FullName())
			os.Stdout.WriteString("\n")
		}
		if gridView {
			g, ok := src.grid.FitIntoWidth(width)
			if !ok || *oneLine {
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

func printDir(src string, file ipe.File, corners []bool) {
	if !file.IsDir() {
		return
	}
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}

	var buffer bytes.Buffer
	grid := gridt.New(direction, *separatorFlag)
	srcs = append(srcs, srcInfo{file, grid, &buffer})

	if *reverseFlag {
		reverse(fs)
	}

	// First loop: preparation.
	for _, f := range fs {
		checkBiggestValues(f, *allFlag, *ignoreFlag, *filterFlag)
	}

	// Second loop: printing.
	for ii, f := range fs {
		println(ii + 1)
		println(len(fs))
		printFile(srcInfo{f, grid, &buffer}, append(corners, ii+1 == len(fs)))
	}
}

func printFile(src srcInfo, corners []bool) {
	if !show(src.file, *allFlag, *ignoreFlag, *filterFlag) {
		return
	}

	var name string
	if *classifyFlag {
		name = src.file.ClassifiedName()
	} else {
		name = src.file.Name()
	}

	if !gridView {
		if *longFlag {
			if *inodeFlag {
				src.buffer.WriteString(fmtColumn(fmtInode(src.file), *separatorFlag, bgstInode, !osWindows))
			}
			src.buffer.WriteString(fmtColumn(fmtMode(src.file), *separatorFlag, bgstMode))
			src.buffer.WriteString(fmtColumn(fmtSize(src.file), *separatorFlag, bgstSize))
			src.buffer.WriteString(fmtColumn(fmtTime(src.file), *separatorFlag, 0))
			src.buffer.WriteString(fmtColumn(fmtUser(src.file), *separatorFlag, bgstUser, !osWindows))
		}
		if *treeFlag {
			src.buffer.WriteString(makeTree(corners))
		}
		src.buffer.WriteString(name)
		src.buffer.WriteString("\n")
	} else {
		src.grid.Add(name)
	}

	if *recursiveFlag && src.file.IsDir() && (*depthFlag == 0 || *depthFlag >= len(corners)) {
		printDir(src.file.Name(), src.file, corners)
	}
}

func checkBiggestValues(f ipe.File, all bool, ignore, filter *regexp.Regexp) {
	if !show(f, all, ignore, filter) {
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
	if *recursiveFlag {
		for _, ff := range f.Children() {
			checkBiggestValues(ff, all, ignore, filter)
		}
	}
}

func show(f ipe.File, all bool, ignore, filter *regexp.Regexp) bool {
	return (all || !f.IsDotfile()) &&
		(ignore == nil || !ignore.MatchString(f.Name())) &&
		(filter == nil || filter.MatchString(f.Name()))
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
	os.Stdout.WriteString(os.Args[0])
	os.Stdout.WriteString(" error: ")
	os.Stdout.WriteString(err)
	os.Stdout.WriteString("\n")
	os.Exit(1)
}

func fixSrc(src string) string {
	if osWindows {
		return strings.Replace(src, "~", os.Getenv("USERPROFILE"), -1)
	}
	return src
}
