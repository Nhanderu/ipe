package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/Nhanderu/tuyo/convert"
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
		fmt.Printf("%s%s%7d%s%s%s%s\n",
			f.Mode().String(), sep,
			f.Size(), sep,
			fmtTime(f.ModTime()), sep,
			name)
	}
}

func fmtTime(t time.Time) string {
	year, month, day := t.Date()
	str := fmt.Sprintf("%2d %s ", day, month.String()[:3])
	if year == time.Now().Year() {
		str += fmt.Sprintf("%2d:%2d", t.Hour(), t.Minute())
	} else {
		str += convert.ToString(year)
	}
	return str
}
