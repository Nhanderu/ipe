package ipefmt

import (
	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
)

type treeFormatter struct {
	commonFormatter
}

func newTreeFormatter(args ArgsInfo) *treeFormatter {
	f := &treeFormatter{commonFormatter{args, make([]srcInfo, 0), 1}}
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

func (f *treeFormatter) getDir(file ipe.File, grid *gridt.Grid, corners []bool) {
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}
	if len(corners) == 0 {
		f.srcs = append(f.srcs, srcInfo{file, nil, grid})
	}
	if f.args.Sort != ArgSortNone {
		sortFiles(fs, f.args.Sort)
	}
	if f.args.Reverse {
		reverseFiles(fs)
	}
	for ii, file := range fs {
		f.getFile(file, grid, append(corners, ii+1 == len(fs)))
	}
}

func (f *treeFormatter) getFile(file ipe.File, grid *gridt.Grid, corners []bool) {
	if !shouldShow(file, f.args) {
		return
	}

	grid.Add(makeTree(corners) + f.getName(file))

	if f.args.Recursive && file.IsDir() && (f.args.Depth == 0 || int(f.args.Depth) >= len(corners)) {
		f.getDir(file, grid, corners)
	}
}
