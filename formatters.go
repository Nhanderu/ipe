package ipe

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Nhanderu/gridt"
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

type srcInfoBuffer struct {
	file   File
	buffer *bytes.Buffer
}

type srcInfoGrid struct {
	file File
	grid *gridt.Grid
}

func NewFormatter(args ArgsInfo) (fmt.Stringer, error) {
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

type GridFormatter struct {
	direction gridt.Direction
}

func newGridFormatter(args ArgsInfo) (GridFormatter, error) {
	var direction gridt.Direction
	if args.Across {
		direction = gridt.LeftToRight
	} else {
		direction = gridt.TopToBottom
	}
	return GridFormatter{direction}, nil
}

func (f GridFormatter) String() string { return "" }

type TreeFormatter struct {
	buffer bytes.Buffer
}

func newTreeFormatter(args ArgsInfo) (TreeFormatter, error) {
	return TreeFormatter{}, nil
}

func (f TreeFormatter) String() string {
	return f.buffer.String()
}

type LongFormatter struct {
	srcs []srcInfoBuffer
}

func newLongFormatter(args ArgsInfo) (LongFormatter, error) {
	f := LongFormatter{}
	f.srcs = make([]srcInfoBuffer, 0)
	for _, src := range args.Source {
		src = fixInSrc(src)
		file, err := Read(src)
		if err != nil {
			return LongFormatter{}, err
		}
		f.getDir(src, file, []bool{}, args)
	}
	return f, nil
}

func (f LongFormatter) getDir(src string, file File, corners []bool, args ArgsInfo) {
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	f.srcs = append(f.srcs, srcInfoBuffer{file, buffer})

	if args.Reverse {
		reverse(fs)
	}

	// First loop: preparation.
	// for _, file := range fs {
	// 	checkBiggestValues(f, args)
	// }

	// Second loop: printing.
	for ii, file := range fs {
		f.getFile(srcInfoBuffer{file, buffer}, append(corners, ii+1 == len(fs)), args)
	}
}

func (f LongFormatter) getFile(srcInfo srcInfoBuffer, corners []bool, args ArgsInfo) {

}

func (f LongFormatter) String() string { return "" }

type LongTreeFormatter struct{}

func newLongTreeFormatter(args ArgsInfo) (LongTreeFormatter, error) {
	return LongTreeFormatter{}, nil
}

func (f LongTreeFormatter) String() string { return "" }

func fmtColumn(column, sep string, size int) string {
	var buf bytes.Buffer
	buf.WriteString(text.PadLeft(column, " ", size))
	buf.WriteString(sep)
	return buf.String()
}

func fmtInode(f File) string {
	return strconv.FormatUint(f.Inode(), 10)
}

func fmtMode(f File) string {
	return f.Mode().String()
}

func fmtSize(f File) string {
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

func fmtTime(f File) string {
	t := f.ModTime()
	year, month, day := t.Date()
	str := fmt.Sprintf("%2d %s ", day, month.String()[:3])
	if year == time.Now().Year() {
		return fmt.Sprintf("%s%2d:%02d", str, t.Hour(), t.Minute())
	}
	return fmt.Sprintf("%s%d ", str, year)
}

func fmtUser(f File) string {
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

func fixOutSrc(src string) string {
	if osWindows {
		return strings.Replace(src, os.Getenv("USERPROFILE"), "~", -1)
	}
	return src
}

func reverse(a []File) {
	for l, r := 0, len(a)-1; l < r; l, r = l+1, r-1 {
		a[l], a[r] = a[r], a[l]
	}
}
