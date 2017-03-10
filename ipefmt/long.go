package ipefmt

import (
	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
)

type longFormatter struct {
	commonFormatter
	showAcc bool
	showMod bool
	showCrt bool
}

func newLongFormatter(args ArgsInfo) *longFormatter {
	acc, mod, crt := timesToShow(args)
	f := &longFormatter{commonFormatter{args, make([]srcInfo, 0), 3}, acc, mod, crt}
	if f.showAcc {
		f.cols++
	}
	if f.showMod {
		f.cols++
	}
	if f.showCrt {
		f.cols++
	}
	if !osWindows {
		f.cols++
		if f.args.Inode {
			f.cols++
		}
	}
	for _, src := range args.Sources {
		file, err := ipe.Read(fixInSrc(src))
		if err != nil {
			f.srcs = append(f.srcs, srcInfo{file, err, nil})
		} else {
			f.getDir(file, 0)
		}
	}
	return f
}

func (f *longFormatter) getDir(file ipe.File, depth uint8) {
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}
	grid := gridt.New(gridt.LeftToRight, f.args.Separator)
	f.srcs = append(f.srcs, srcInfo{file, nil, grid})
	if f.args.Reverse {
		reverse(fs)
	}
	for _, file := range fs {
		f.getFile(file, grid, depth+1)
	}
}

func (f *longFormatter) getFile(file ipe.File, grid *gridt.Grid, depth uint8) {
	if !shouldShow(file, f.args) {
		return
	}

	if f.args.Inode && !osWindows {
		grid.Add(fmtInode(file))
	}
	grid.Add(fmtMode(file))
	grid.Add(fmtSize(file))
	if f.showAcc {
		grid.Add(fmtAccTime(file))
	}
	if f.showMod {
		grid.Add(fmtModTime(file))
	}
	if f.showCrt {
		grid.Add(fmtCrtTime(file))
	}
	if !osWindows {
		grid.Add(fmtUser(file))
	}
	grid.Add(f.getName(file))

	if f.args.Recursive && file.IsDir() && (f.args.Depth == 0 || f.args.Depth >= depth) {
		f.getDir(file, depth)
	}
}
