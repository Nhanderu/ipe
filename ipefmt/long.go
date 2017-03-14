package ipefmt

import (
	"strconv"

	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
)

type longFormatter struct {
	commonFormatter
	showAcc   bool
	showMod   bool
	showCrt   bool
	showInode bool
	showUser  bool
}

func newLongFormatter(args ArgsInfo) *longFormatter {
	f := &longFormatter{commonFormatter{args, make([]srcInfo, 0), 3}, false, false, false, args.Inode && !osWindows, !osWindows}
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
	if f.args.Header {
		f.write(
			grid,
			"inode",
			"mode",
			"size",
			"accessed",
			"modified",
			"created",
			"user",
			"file name",
		)
	}
	for _, file := range fs {
		f.getFile(file, grid, depth+1)
	}
}

func (f *longFormatter) getFile(file ipe.File, grid *gridt.Grid, depth uint8) {
	if !shouldShow(file, f.args) {
		return
	}

	f.write(
		grid,
		strconv.FormatUint(file.Inode(), 10),
		file.Mode().String(),
		fmtSize(file),
		fmtTime(file.AccTime()),
		fmtTime(file.ModTime()),
		fmtTime(file.CrtTime()),
		file.User().Username,
		f.getName(file),
	)

	if f.args.Recursive && file.IsDir() && (f.args.Depth == 0 || f.args.Depth >= depth) {
		f.getDir(file, depth)
	}
}

func (f longFormatter) write(grid *gridt.Grid, inode, mode, size, acc, mod, crt, user, name string) {
	if f.showInode {
		grid.Add(inode)
	}
	grid.Add(mode)
	grid.Add(size)
	if f.showAcc {
		grid.Add(acc)
	}
	if f.showMod {
		grid.Add(mod)
	}
	if f.showCrt {
		grid.Add(crt)
	}
	if f.showUser {
		grid.Add(user)
	}
	grid.Add(name)
}
