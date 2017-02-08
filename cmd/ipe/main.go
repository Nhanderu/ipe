package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Nhanderu/ipe"
	"github.com/Nhanderu/trena"
	"github.com/Nhanderu/tuyo/convert"
	"github.com/fatih/color"
	isatty "github.com/mattn/go-isatty"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	kilobyte = 1024
	megabyte = kilobyte * 1024
	gigabyte = megabyte * 1024
	terabyte = gigabyte * 1024
)

var (
	srcArg            = kingpin.Arg("src", "the directory to list contents").Default(".").String()
	separatorFlag     = kingpin.Flag("separator", "separator of the columns in long view").Default(" ").String()
	allFlag           = kingpin.Flag("all", "do not hide entries starting with .").Short('a').Bool()
	colorFlag         = kingpin.Flag("color", "control  whether  color is used to distinguish file types").Enum("never", "always", "auto")
	classifyFlag      = kingpin.Flag("classify", "append indicator (one of /=@|) to entries").Short('F').Bool()
	humanReadableFlag = kingpin.Flag("human-readable", "print sizes in human readable format (e.g., 1K 234M 2G)").Short('h').Bool()
	inodeFlag         = kingpin.Flag("inode", "print index number of each file").Short('i').Bool()
	ignoreFlag        = kingpin.Flag("ignore", "to not list implied entries matching shell PATTERN").Short('I').Regexp()
	longFlag          = kingpin.Flag("long", "use a long listing format").Short('l').Bool()
	reverseFlag       = kingpin.Flag("reverse", "").Short('r').Bool()
	recursiveFlag     = kingpin.Flag("recursive", "").Short('R').Bool()

	width, biggestMode, biggestSize, biggestTime int
)

func main() {
	kingpin.Parse()

	// Gets the necessary info.
	fs, err := ipe.ReadDir(*srcArg)
	if err != nil {
		endWithErr(err)
	}
	width, _, err = trena.Size()
	if err != nil {
		endWithErr(err)
	}
	if *reverseFlag {
		reverse(fs)
	}

	// First loop: preparation.
	for _, f := range fs {
		checkBiggestValues(f)
	}

	// Second loop: printing.
	for i, f := range fs {
		printFile(i, f, 0)
	}

	fmt.Println()
}

func checkBiggestValues(f ipe.File) {
	if m := len(fmtMode(f.Mode().String(), "")); m > biggestMode {
		biggestMode = m
	}
	if s := len(fmtSize(f.Size(), "")); s > biggestSize {
		biggestSize = s
	}
	if t := len(fmtTime(f.ModTime(), "")); t > biggestTime {
		biggestTime = t
	}
	if *recursiveFlag {
		for _, ff := range f.Children() {
			checkBiggestValues(ff)
		}
	}
}

func printFile(i int, f ipe.File, t int) {
	n := f.Name()
	if (!*allFlag && f.IsDotfile()) ||
		(*ignoreFlag != nil && (*ignoreFlag).MatchString(n)) {
		return
	}

	if *colorFlag == "auto" {
		color.NoColor = !isatty.IsTerminal(os.Stdout.Fd()) || os.Getenv("TERM") == "dumb"
	} else {
		color.NoColor = *colorFlag == "never"
	}

	var name string
	if *classifyFlag {
		name = f.ClassifiedName()
	} else {
		name = n
	}

	var inode string
	if *inodeFlag {
		inode = fmt.Sprintf("%d%s", i+1, *separatorFlag)
	}

	if *longFlag {
		fmt.Printf("%s%s%s%s%s%s\n",
			inode,
			getMode(f, *separatorFlag),
			getSize(f, *separatorFlag),
			getTime(f, *separatorFlag),
			strings.Repeat("--> ", t),
			name)
	} else {
		fmt.Printf("%s  ", name)
	}

	if *recursiveFlag {
		children := f.Children()
		if *reverseFlag {
			reverse(children)
		}
		for _, c := range children {
			printFile(i, c, t+1)
		}
	}
}

func getMode(f ipe.File, sep string) string {
	return padLeft(fmtMode(f.Mode().String(), sep), " ", biggestMode+len(sep))
}

func fmtMode(m, sep string) string {
	return fmt.Sprintf("%s%s", m, sep)
}

func getSize(f ipe.File, sep string) string {
	return padLeft(fmtSize(f.Size(), sep), " ", biggestSize+len(sep))
}

func fmtSize(s int64, sep string) string {
	if *humanReadableFlag {
		return fmt.Sprintf("%s%s", humanSize(s), sep)
	}
	return fmt.Sprintf("%s%s", convert.ToString(s), sep)
}

func humanSize(s int64) string {
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
	return padRight(fmtTime(f.ModTime(), sep), " ", biggestTime+len(sep))
}

func fmtTime(t time.Time, sep string) string {
	year, month, day := t.Date()
	str := fmt.Sprintf("%2d %s ", day, month.String()[:3])
	if year == time.Now().Year() {
		return fmt.Sprintf("%s%2d:%02d%s", str, t.Hour(), t.Minute(), sep)
	}
	return fmt.Sprintf("%s%d%s", str, year, sep)
}

func reverse(a []ipe.File) {
	for l, r := 0, len(a)-1; l < r; l, r = l+1, r-1 {
		a[l], a[r] = a[r], a[l]
	}
}

func endWithErr(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

func padLeft(a, b string, l int) string {
	if l <= len(a) {
		return a
	}
	return fmt.Sprintf("%s%s", strings.Repeat(b, l-len(a)), a)
}

func padRight(a, b string, l int) string {
	if l <= len(a) {
		return a
	}
	return fmt.Sprintf("%s%s", a, strings.Repeat(b, l-len(a)))
}
