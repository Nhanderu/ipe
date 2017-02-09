package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Nhanderu/ipe"
	"github.com/Nhanderu/trena"
	"github.com/Nhanderu/tuyo/convert"
	"github.com/Nhanderu/tuyo/text"
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
	treeFlag          = kingpin.Flag("tree", "shows the entries in the tree view").Short('t').Bool()

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
		printFile(i, f, 0, []bool{i+1 == len(fs)})
	}

	fmt.Println()
}

func checkBiggestValues(f ipe.File) {
	if !show(f) {
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
	if *recursiveFlag {
		for _, ff := range f.Children() {
			checkBiggestValues(ff)
		}
	}
}

func show(f ipe.File) bool {
	return (*allFlag || !f.IsDotfile()) && (*ignoreFlag == nil || !(*ignoreFlag).MatchString(f.Name()))
}

func printFile(i int, f ipe.File, t int, corners []bool) {
	if !show(f) {
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
		name = f.Name()
	}

	var inode string
	if *inodeFlag {
		inode = fmt.Sprintf("%d%s", i+1, *separatorFlag)
	}

	if *longFlag {
		if *treeFlag {
			name = fmt.Sprint(makeTree(corners), name)
		}
		fmt.Printf("%s%s%s%s%s\n",
			inode,
			getMode(f, *separatorFlag),
			getSize(f, *separatorFlag),
			getTime(f, *separatorFlag),
			name)
	} else {
		fmt.Print(name, "  ")
	}

	if *recursiveFlag {
		children := f.Children()
		if *reverseFlag {
			reverse(children)
		}
		for ii, c := range children {
			printFile(i, c, t+1, append(corners, ii+1 == len(children)))
		}
	}
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
	s := f.Size()
	if !*humanReadableFlag {
		return convert.ToString(s)
	}
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

func endWithErr(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
