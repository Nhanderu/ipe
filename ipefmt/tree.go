package ipefmt

import (
	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
)

type treeFormatter struct {
	*commonFormatter
}

func newTreeFormatter(args ArgsInfo) *treeFormatter {
	return &treeFormatter{&commonFormatter{args, make([]srcInfo, 0), 1}}
}

func (f *treeFormatter) getDir(file ipe.File, grid **gridt.Grid, corners []bool) {
	if len(corners) == 0 {
		f.appendSource(srcInfo{file, nil, *grid})
	}
}

func (f *treeFormatter) getFile(file ipe.File, grid *gridt.Grid, corners []bool) {
	grid.Add(makeTree(corners) + f.getName(file))
}
