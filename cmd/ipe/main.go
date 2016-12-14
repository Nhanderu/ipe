package main

import (
	"fmt"
	"time"

	"github.com/Nhanderu/ipe"
	"github.com/Nhanderu/tuyo/convert"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	kilobyte2 = 1024
	megabyte2 = kilobyte2 * 1024
	gigabyte2 = megabyte2 * 1024
	terabyte2 = gigabyte2 * 1024

	kilobyte10 = 1000
	megabyte10 = kilobyte10 * 1000
	gigabyte10 = megabyte10 * 1000
	terabyte10 = gigabyte10 * 1000
)

var (
	srcArg            = kingpin.Arg("src", "the directory to list contents").Default(".").String()
	separatorFlag     = kingpin.Flag("separator", "separator of the columns in long view").Default(" ").String()
	allFlag           = kingpin.Flag("all", "do not hide entries starting with .").Short('a').Bool()
	classifyFlag      = kingpin.Flag("classify", "append indicator (one of /=@|) to entries").Short('F').Bool()
	humanReadableFlag = kingpin.Flag("human-readable", "print sizes in human readable format (e.g., 1K 234M 2G)").Short('h').Bool()
	siFlag            = kingpin.Flag("si", "print sizes in human readable format, but use powers of 1000 not 1024").Bool()
	inodeFlag         = kingpin.Flag("inode", "print index number of each file").Short('i').Bool()
	ignoreFlag        = kingpin.Flag("ignore", "to not list implied entries matching shell PATTERN").Short('I').Regexp()
	longFlag          = kingpin.Flag("long", "use a long listing format").Short('l').Bool()
)

func main() {
	kingpin.Parse()

	fs, err := ipe.ReadDir(*srcArg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for i, f := range fs {

		n := f.Name()
		if (!*allFlag && n[0] == '.') ||
			(*ignoreFlag != nil && (*ignoreFlag).MatchString(n)) {
			continue
		}

		var name string
		if *classifyFlag {
			name = f.ClassifiedName()
		} else {
			name = n
		}

		var size string
		if *humanReadableFlag {
			size = fmt.Sprintf("%s%s", humanSize2(f.Size()), *separatorFlag)
		} else if *siFlag {
			size = fmt.Sprintf("%s%s", humanSize10(f.Size()), *separatorFlag)
		} else {
			size = fmt.Sprintf("%d%s", f.Size(), *separatorFlag)
		}

		var inode string
		if *inodeFlag {
			inode = fmt.Sprintf("%d%s", i+1, *separatorFlag)
		}

		if *longFlag {
			fmt.Printf("%s%s%s%s%s%s%s\n",
				inode,
				f.Mode().String(),
				*separatorFlag,
				size,
				fmtTime(f.ModTime()),
				*separatorFlag,
				name)
		} else {
			fmt.Printf("%s  ", name)
		}
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

func humanSize2(s int64) string {
	return humanSize(s, kilobyte2, megabyte2, gigabyte2, terabyte2)
}

func humanSize10(s int64) string {
	return humanSize(s, kilobyte10, megabyte10, gigabyte10, terabyte10)
}

func humanSize(s, kb, mb, gb, tb int64) string {
	if s < kb {
		return fmt.Sprintf("%6dB", s)
	} else if s < mb {
		return fmt.Sprintf("%5.1dKB", s/kb)
	} else if s < gb {
		return fmt.Sprintf("%5.1dMB", s/mb)
	} else if s < tb {
		return fmt.Sprintf("%5.1dGB", s/gb)
	} else {
		return fmt.Sprintf("%5.1dTB", s/tb)
	}
}
