package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/Nhanderu/tuyo/convert"
)

const (
	kilobyte = 1024
	megabyte = kilobyte * 1024
	gigabyte = megabyte * 1024
	terabyte = gigabyte * 1024
)

func main() {
	var src string
	if len(os.Args) > 1 {
		src = os.Args[1]
	} else {
		src = "."
	}
	fs, err := ioutil.ReadDir(src)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, f := range fs {
		sep := " "
		name := f.Name()
		if f.IsDir() {
			name += "/"
		}
		fmt.Printf("%s%s%s%s%s%s%s\n",
			f.Mode().String(),
			sep,
			humanSize(f.Size()),
			sep,
			fmtTime(f.ModTime()),
			sep,
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
