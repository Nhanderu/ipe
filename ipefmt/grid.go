package ipefmt

import (
	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
)

type gridFormatter struct {
	*commonFormatter
	direction gridt.Direction
}

func newGridFormatter(args ArgsInfo) *gridFormatter {
	if args.Across {
		return &gridFormatter{&commonFormatter{args, make([]srcInfo, 0), 0}, gridt.LeftToRight}
	}
	return &gridFormatter{&commonFormatter{args, make([]srcInfo, 0), 0}, gridt.TopToBottom}
}

func (f *gridFormatter) getDir(file ipe.File, grid **gridt.Grid, corners []bool) {
	*grid = gridt.New(f.direction, f.args.Separator)
	f.srcs = append(f.srcs, srcInfo{file, nil, *grid})
}

func (f *gridFormatter) getFile(file ipe.File, grid *gridt.Grid, corners []bool) {
	grid.Add(f.getName(file))
}
