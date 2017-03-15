package ipefmt

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Nhanderu/ipe"
)

const (
	kilobyte = 1024
	megabyte = kilobyte * 1024
	gigabyte = megabyte * 1024
	terabyte = gigabyte * 1024

	osWindows = runtime.GOOS == "windows"
)

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

func fmtBlocks(f ipe.File) string {
	if f.IsDir() {
		return "-"
	}
	return strconv.FormatInt(f.Blocks(), 10)
}

func fmtTime(t time.Time) string {
	year, month, day := t.Date()
	str := fmt.Sprintf("%2d %s ", day, month.String()[:3])
	if year == time.Now().Year() {
		return fmt.Sprintf("%s%2d:%02d", str, t.Hour(), t.Minute())
	}
	return fmt.Sprintf("%s%d ", str, year)
}

func fixInSrc(src string) string {
	if osWindows {
		return strings.Replace(src, "~", os.Getenv("USERPROFILE"), -1)
	}
	return src
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
