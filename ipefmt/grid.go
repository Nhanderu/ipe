package ipefmt

import (
	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
)

type gridFormatter struct {
	commonFormatter
	direction gridt.Direction
}

func newGridFormatter(args ArgsInfo) *gridFormatter {
	var direction gridt.Direction
	if args.Across {
		direction = gridt.LeftToRight
	} else {
		direction = gridt.TopToBottom
	}
	f := &gridFormatter{commonFormatter{args, make([]srcInfo, 0), 0}, direction}
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

func (f *gridFormatter) getDir(file ipe.File, depth uint8) {
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}
	grid := gridt.New(f.direction, f.args.Separator)
	f.srcs = append(f.srcs, srcInfo{file, nil, grid})
	if f.args.Sort != ArgSortNone {
		sortFiles(fs, f.args.Sort)
	}
	if f.args.Reverse {
		reverseFiles(fs)
	}
	for _, file := range fs {
		f.getFile(file, grid, depth+1)
	}
}

func (f *gridFormatter) getFile(file ipe.File, grid *gridt.Grid, depth uint8) {
	if !shouldShow(file, f.args) {
		return
	}

	grid.Add(f.getName(file))

	if f.args.Recursive && file.IsDir() && (f.args.Depth == 0 || f.args.Depth >= depth) {
		f.getDir(file, depth)
	}
}
