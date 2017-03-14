package ipefmt

import (
	"strconv"

	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
)

type longTreeFormatter struct {
	commonFormatter
	showAcc   bool
	showMod   bool
	showCrt   bool
	showInode bool
	showUser  bool
}

func newLongTreeFormatter(args ArgsInfo) *longTreeFormatter {
	f := &longTreeFormatter{commonFormatter{args, make([]srcInfo, 0), 3}, false, false, false, args.Inode && !osWindows, !osWindows}
	f.showAcc, f.showMod, f.showCrt = timesToShow(args)
	if f.showInode {
		f.cols++
	}
	if f.showAcc {
		f.cols++
	}
	if f.showMod {
		f.cols++
	}
	if f.showCrt {
		f.cols++
	}
	if f.showUser {
		f.cols++
	}
	for _, src := range args.Sources {
		file, err := ipe.Read(fixInSrc(src))
		if err != nil {
			f.srcs = append(f.srcs, srcInfo{file, err, nil})
		} else {
			f.getDir(file, gridt.New(gridt.LeftToRight, f.args.Separator), []bool{})
		}
	}
	return f
}

func (f *longTreeFormatter) getDir(file ipe.File, grid *gridt.Grid, corners []bool) {
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}
	if len(corners) == 0 {
		f.srcs = append(f.srcs, srcInfo{file, nil, grid})
	}
	if f.args.Reverse {
		reverse(fs)
	}
	for ii, file := range fs {
		f.getFile(file, grid, append(corners, ii+1 == len(fs)))
	}
}

func (f *longTreeFormatter) getFile(file ipe.File, grid *gridt.Grid, corners []bool) {
	if !shouldShow(file, f.args) {
		return
	}

	if f.showInode {
		grid.Add(strconv.FormatUint(file.Inode(), 10))
	}
	grid.Add(file.Mode().String())
	grid.Add(fmtSize(file))
	if f.showAcc {
		grid.Add(fmtTime(file.AccTime()))
	}
	if f.showMod {
		grid.Add(fmtTime(file.ModTime()))
	}
	if f.showCrt {
		grid.Add(fmtTime(file.CrtTime()))
	}
	if f.showUser {
		grid.Add(file.User().Username)
	}
	grid.Add(makeTree(corners) + f.getName(file))

	if f.args.Recursive && file.IsDir() && (f.args.Depth == 0 || int(f.args.Depth) >= len(corners)) {
		f.getDir(file, grid, corners)
	}
}
