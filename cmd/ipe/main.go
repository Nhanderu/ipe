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
	"github.com/alecthomas/kingpin"
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

	if args.Color != ipe.ArgColorAuto {
		color.NoColor = args.Color == ipe.ArgColorNever
	}

	if args.Across {
		direction = gridt.LeftToRight
	} else {
		direction = gridt.TopToBottom
	}

	gridView = !args.Long && !args.Tree
	treeView = !args.Long && args.Tree
	longView = args.Long && !args.Tree
	longTreeView = args.Long && args.Tree

	width, _, err := trena.Size()
	if err != nil {
		endWithErr(err.Error())
	}

	srcs = make([]srcInfo, 0)
	for _, src := range args.Source {
		src = fixSrc(src)
		f, err := ipe.Read(src)
		if err != nil {
			endWithErr(err.Error())
		}
		printDir(src, f, []bool{}, args)
	}

	if args.Tree {
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
				if !ok || args.OneLine {
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

func parseArgs() ipe.ArgsInfo {
	var args ipe.ArgsInfo
	kingpin.Arg("source", "the directory to list contents").Default(".").StringsVar(&args.Source)
	kingpin.Flag("separator", "separator of the columns").Short('S').Default("  ").StringVar(&args.Separator)
	kingpin.Flag("across", "writes the entries by lines instead of by columns").Short('x').BoolVar(&args.Across)
	kingpin.Flag("all", "do not hide entries starting with .").Short('a').BoolVar(&args.All)
	kingpin.Flag("color", "control whether color is used to distinguish file types").Default(ipe.ArgColorAuto).EnumVar(&args.Color, ipe.ArgColorNever, ipe.ArgColorAlways, ipe.ArgColorAuto)
	kingpin.Flag("classify", "append indicator to the entries").Short('F').BoolVar(&args.Classify)
	kingpin.Flag("depth", "maximum depth of recursion").Short('D').IntVar(&args.Depth)
	kingpin.Flag("filter", "only show entries that matches the pattern").Short('f').RegexpVar(&args.Filter)
	kingpin.Flag("ignore", "do not show entries that matches the pattern").Short('I').RegexpVar(&args.Ignore)
	kingpin.Flag("inode", "show entry inode").Short('i').BoolVar(&args.Inode)
	kingpin.Flag("long", "show entries in the \"long view\"").Short('l').BoolVar(&args.Long)
	kingpin.Flag("one-line", "show one entry per line").Short('1').BoolVar(&args.OneLine)
	kingpin.Flag("reverse", "reverse order of entries").Short('r').BoolVar(&args.Reverse)
	kingpin.Flag("recursive", "list subdirectories recursively").Short('R').BoolVar(&args.Recursive)
	kingpin.Flag("tree", "shows the entries in the tree view").Short('t').BoolVar(&args.Tree)
	kingpin.Parse()
	return args
}

func printDir(src string, file ipe.File, corners []bool, args ipe.ArgsInfo) {
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}

	var buffer *bytes.Buffer
	if args.Tree {
		buffer = &outBuffer
	} else {
		buffer = bytes.NewBuffer([]byte{})
	}
	grid := gridt.New(direction, args.Separator)
	srcs = append(srcs, srcInfo{file, grid, buffer})

	if args.Reverse {
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

func checkBiggestValues(f ipe.File, args ipe.ArgsInfo) {
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

func printFile(src srcInfo, corners []bool, args ipe.ArgsInfo) {
	if !shouldShow(src.file, args) {
		return
	}

	var name string
	if args.Classify {
		name = src.file.ClassifiedName()
	} else {
		name = src.file.Name()
	}

	if !gridView {
		src.grid.Add(name)
	} else {
		if args.Long {
			src.writec(fmtInode, args.Separator, bgstInode, args.Inode, !osWindows)
			src.writec(fmtMode, args.Separator, bgstMode)
			src.writec(fmtSize, args.Separator, bgstSize)
			src.writec(fmtTime, args.Separator, 0)
			src.writec(fmtUser, args.Separator, bgstUser, !osWindows)
		}
		if args.Tree {
			src.write(makeTree(corners))
		}
		src.write(name)
		src.write("\n")
	}

	if args.Recursive && src.file.IsDir() && (args.Depth == 0 || args.Depth >= len(corners)) {
		printDir(src.file.Name(), src.file, corners, args)
	}
}

func shouldShow(f ipe.File, args ipe.ArgsInfo) bool {
	return (args.All || !f.IsDotfile()) &&
		(args.Ignore == nil || !args.Ignore.MatchString(f.Name())) &&
		(args.Filter == nil || args.Filter.MatchString(f.Name()))
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
