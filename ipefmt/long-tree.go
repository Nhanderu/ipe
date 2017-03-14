package ipefmt

import (
	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
)

type longTreeFormatter struct {
	*commonFormatter
	long *longFormatter
	tree *treeFormatter
}

func newLongTreeFormatter(args ArgsInfo) *longTreeFormatter {
	f := &longTreeFormatter{
		&commonFormatter{args, make([]srcInfo, 0), 0},
		newLongFormatter(args),
		newTreeFormatter(args),
	}
	f.cols = f.long.calculateCols()
	return f
}

func (f *longTreeFormatter) getDir(file ipe.File, grid **gridt.Grid, corners []bool) {
	if len(corners) == 0 {
		f.long.writeHeader(*grid)
	}
	f.tree.getDir(file, grid, corners)
	f.srcs = f.tree.srcs
}

func (f *longTreeFormatter) getFile(file ipe.File, grid *gridt.Grid, corners []bool) {
	f.long.writeAllButName(grid, file, makeTree(corners)+f.getName(file))
}
