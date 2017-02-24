package main

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
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
	filterFlag    = kingpin.Flag("filter", "only show entries that matches the pattern").Short('f').Regexp()
	ignoreFlag    = kingpin.Flag("ignore", "do not show entries that matches the pattern").Short('I').Regexp()
	longFlag      = kingpin.Flag("long", "show entries in the \"long view\"").Short('l').Bool()
	reverseFlag   = kingpin.Flag("reverse", "reverse order of entries").Short('r').Bool()
	recursiveFlag = kingpin.Flag("recursive", "list subdirectories recursively").Short('R').Bool()
	treeFlag      = kingpin.Flag("tree", "shows the entries in the tree view").Short('t').Bool()

	gridView                                           bool
	biggestMode, biggestSize, biggestTime, biggestUser int
	grids                                              []dirGrid
	direction                                          gridt.Direction

	osWindows = runtime.GOOS == "windows"
)

type dirGrid struct {
	file ipe.File
	grid *gridt.Grid
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

	grids = make([]dirGrid, 0)
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
		printDir(src, f, 0, []bool{})
	}

	if gridView {
		for _, grid := range grids {
			if len(grids) > 1 {
				os.Stdout.WriteString(grid.file.FullName())
				os.Stdout.WriteString("\n")
			}
			g, ok := grid.grid.FitIntoWidth(width)
			if !ok {
				for _, cell := range grid.grid.Cells() {
					os.Stdout.WriteString(cell)
					os.Stdout.WriteString("\n")
				}
			} else {
				os.Stdout.WriteString(g.String())
			}
			os.Stdout.WriteString("\n")
		}
	} else {
		os.Stdout.WriteString("\n")
	}
}

func printDir(src string, file ipe.File, depth int, corners []bool) {
	if !file.IsDir() {
		fmt.Println(src, "is not a directory")
		return
	}
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}

	grid := gridt.New(direction, *separatorFlag)
	grids = append(grids, dirGrid{file, grid})

	if *reverseFlag {
		reverse(fs)
	}

	// First loop: preparation.
	for _, f := range fs {
		checkBiggestValues(f, *allFlag, *ignoreFlag, *filterFlag)
	}

	// Second loop: printing.
	for ii, f := range fs {
		printFile(
			f,
			grid,
			depth,
			append(corners, ii+1 == len(fs)))
	}
}

func printFile(file ipe.File, grid *gridt.Grid, depth int, corners []bool) {
	if !show(file, *allFlag, *ignoreFlag, *filterFlag) {
		return
	}

	var name string
	if *classifyFlag {
		name = file.ClassifiedName()
	} else {
		name = file.Name()
	}

	if !gridView {
		if *longFlag {
			os.Stdout.WriteString(getMode(file, *separatorFlag))
			os.Stdout.WriteString(getSize(file, *separatorFlag))
			os.Stdout.WriteString(getTime(file, *separatorFlag))
			os.Stdout.WriteString(getUser(file, *separatorFlag))
		}
		if *treeFlag {
			os.Stdout.WriteString(makeTree(corners))
		}
		os.Stdout.WriteString(name)
		os.Stdout.WriteString("\n")
	} else {
		grid.Add(name)
	}

	if *recursiveFlag && file.IsDir() {
		printDir(file.Name(), file, depth+1, corners)
	}
}

func checkBiggestValues(f ipe.File, all bool, ignore, filter *regexp.Regexp) {
	if !show(f, all, ignore, filter) {
		return
	}
	if m := len(fmtMode(f)); m > biggestMode {
		biggestMode = m
	}
	if s := len(fmtSize(f)); s > biggestSize {
		biggestSize = s
	}
	if t := len(fmtTime(f)); t > biggestTime {
		biggestTime = t
	}
	if u := len(fmtUser(f)); u > biggestUser {
		biggestUser = u
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

func getMode(f ipe.File, sep string) string {
	return addSep(text.PadLeft(fmtMode(f), " ", biggestMode), sep)
}

func fmtMode(f ipe.File) string {
	return f.Mode().String()
}

func getSize(f ipe.File, sep string) string {
	return addSep(text.PadLeft(fmtSize(f), " ", biggestSize), sep)
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

func getTime(f ipe.File, sep string) string {
	return addSep(text.PadRight(fmtTime(f), " ", biggestTime), sep)
}

func fmtTime(f ipe.File) string {
	t := f.ModTime()
	year, month, day := t.Date()
	str := fmt.Sprintf("%2d %s ", day, month.String()[:3])
	if year == time.Now().Year() {
		return fmt.Sprintf("%s%2d:%02d", str, t.Hour(), t.Minute())
	}
	return fmt.Sprintf("%s%d", str, year)
}

func getUser(f ipe.File, sep string) string {
	if osWindows {
		return ""
	}
	return addSep(text.PadLeft(fmtUser(f), " ", biggestUser), sep)
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
	os.Stdout.WriteString(err)
	os.Stdout.WriteString("\n")
	os.Exit(1)
}

func fixSrc(src string) string {
	if runtime.GOOS == "windows" {
		return strings.Replace(src, "~", os.Getenv("USERPROFILE"), -1)
	}
	return src
}
