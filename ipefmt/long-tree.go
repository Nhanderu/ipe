package ipefmt

import (
	"strconv"

	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
)

type longTreeFormatter struct {
	commonFormatter
	long *longFormatter
}

func newLongTreeFormatter(args ArgsInfo) *longTreeFormatter {
	f := &longTreeFormatter{commonFormatter{args, make([]srcInfo, 0), 3}, newLongFormatter(args)}
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

	f.long.write(
		grid,
		strconv.FormatUint(file.Inode(), 10),
		file.Mode().String(),
		fmtSize(file),
		fmtTime(file.AccTime()),
		fmtTime(file.ModTime()),
		fmtTime(file.CrtTime()),
		file.User().Username,
		makeTree(corners)+f.getName(file),
	)

	if f.args.Recursive && file.IsDir() && (f.args.Depth == 0 || int(f.args.Depth) >= len(corners)) {
		f.getDir(file, grid, corners)
	}
}
