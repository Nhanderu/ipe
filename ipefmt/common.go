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
	bgstMode, bgstSize, bgstUser, bgstAccTime, bgstModTime, bgstCrtTime, bgstInode int
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

// NewFormatter returns the correct formatter, based on the arguments.
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

func fmtAccTime(f ipe.File) string {
	return fmtTime(f.AccTime())
}

func fmtModTime(f ipe.File) string {
	return fmtTime(f.ModTime())
}

func fmtCrtTime(f ipe.File) string {
	return fmtTime(f.CrtTime())
}

func fmtTime(t time.Time) string {
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

func timesToShow(args ArgsInfo) (bool, bool, bool) {
	var acc, mod, crt bool
	for _, t := range args.Time {
		if t == ArgTimeAcc {
			acc = true
		}
		if t == ArgTimeMod {
			mod = true
		}
		if t == ArgTimeCrt {
			crt = true
		}
	}
	return acc, mod, crt
}
