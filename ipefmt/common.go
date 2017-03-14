package ipefmt

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
)

const (
	kilobyte = 1024
	megabyte = kilobyte * 1024
	gigabyte = megabyte * 1024
	terabyte = gigabyte * 1024

	osWindows = runtime.GOOS == "windows"
)

type srcInfo struct {
	file ipe.File
	err  error
	grid *gridt.Grid
}

type commonFormatter struct {
	args ArgsInfo
	srcs []srcInfo
	cols int
}

func (f commonFormatter) String() string {
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
			var d gridt.Dimensions
			var ok bool
			if f.cols > 0 {
				d, ok = src.grid.FitIntoColumns(f.cols)
			} else {
				d, ok = src.grid.FitIntoWidth(f.args.Width)
			}
			if !ok || f.args.OneLine {
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

func (f commonFormatter) getName(file ipe.File) string {
	if f.args.Classify {
		return file.ClassifiedName()
	}
	return file.Name()
}

func (f *commonFormatter) appendSource(src srcInfo) {
	f.srcs = append(f.srcs, src)
}

// NewFormatter returns the correct formatter, based on the arguments.
func NewFormatter(args ArgsInfo) fmt.Stringer {
	var f formatter
	if args.Long && args.Tree {
		f = newLongTreeFormatter(args)
	} else if args.Long {
		f = newLongFormatter(args)
	} else if args.Tree {
		f = newTreeFormatter(args)
	} else {
		f = newGridFormatter(args)
	}
	var w formatterWrapper
	w.read(f, args)
	return w
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
