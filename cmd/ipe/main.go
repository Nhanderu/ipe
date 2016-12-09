package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Nhanderu/tuyo/convert"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	kilobyte = 1024
	megabyte = kilobyte * 1024
	gigabyte = megabyte * 1024
	terabyte = gigabyte * 1024
)

var (
	inodeFlag     = kingpin.Flag("inode", "print index number of each file").Short('i').Bool()
	srcArg        = kingpin.Arg("src", "the directory to list contents").Default(".").String()
	separatorFlag = kingpin.Flag("separator", "").Default(" ").String()
)

func main() {
	kingpin.Parse()

	fs, err := ioutil.ReadDir(*srcArg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for i, f := range fs {
		name := f.Name()
		if f.IsDir() {
			name += "/"
		}

		var inode string
		if *inodeFlag {
			inode = fmt.Sprintf("%d%s", i+1, *separatorFlag)
		} else {
			inode = ""
		}

		fmt.Printf("%s%s%s%s%s%s%s%s\n",
			inode,
			f.Mode().String(),
			*separatorFlag,
			humanSize(f.Size()),
			*separatorFlag,
			fmtTime(f.ModTime()),
			*separatorFlag,
			name)
	}
}

func fmtTime(t time.Time) string {
	year, month, day := t.Date()
	str := fmt.Sprintf("%2d %s ", day, month.String()[:3])
	if year == time.Now().Year() {
		str += fmt.Sprintf("%2d:%02d", t.Hour(), t.Minute())
	} else {
		str += convert.ToString(year)
	}
	return str
}

func humanSize(s int64) string {
	if s < kilobyte {
		return fmt.Sprintf("%6dB", s)
	} else if s < megabyte {
		return fmt.Sprintf("%5.1dKB", s/kilobyte)
	} else if s < gigabyte {
		return fmt.Sprintf("%5.1dMB", s/megabyte)
	} else if s < terabyte {
		return fmt.Sprintf("%5.1dGB", s/gigabyte)
	} else {
		return fmt.Sprintf("%5.1dTB", s/terabyte)
	}
}