package ipefmt

import (
	"strconv"

	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
)

type longFormatter struct {
	*commonFormatter

	showLinks  bool
	showBlocks bool
	showAcc    bool
	showMod    bool
	showCrt    bool
	showInode  bool
	showUser   bool
	showGroup  bool
}

func newLongFormatter(args ArgsInfo) *longFormatter {
	f := &longFormatter{
		&commonFormatter{args, make([]srcInfo, 0), 0},
		!osWindows,
		!osWindows,
		false,
		false,
		false,
		args.Inode && !osWindows,
		!osWindows,
		args.Group && !osWindows,
	}
	f.showAcc, f.showMod, f.showCrt = timesToShow(args)
	f.cols = f.calculateCols()
	return f
}

func (f *longFormatter) getDir(file ipe.File, grid **gridt.Grid, corners []bool) {
	*grid = gridt.New(gridt.LeftToRight, f.args.Separator)
	f.appendSource(srcInfo{file, nil, *grid})
	f.writeHeader(*grid)
}

func (f *longFormatter) getFile(file ipe.File, grid *gridt.Grid, corners []bool) {
	f.writeAllButName(grid, file, f.getName(file))
}

func (f longFormatter) calculateCols() int {
	cols := 3
	if f.showInode {
		cols++
	}
	if f.showAcc {
		cols++
	}
	if f.showMod {
		cols++
	}
	if f.showCrt {
		cols++
	}
	if f.showUser {
		cols++
	}
	if f.showGroup {
		cols++
	}
	return cols
}

func (f *longFormatter) writeHeader(grid *gridt.Grid) {
	if f.args.Header {
		f.write(
			grid,
			ArgSortInode,
			ArgSortMode,
			ArgSortSize,
			ArgSortLinks,
			ArgSortBlocks,
			ArgSortAccessed,
			ArgSortModified,
			ArgSortCreated,
			ArgSortUser,
			ArgSortGroup,
			ArgSortName,
		)
	}
}

func (f *longFormatter) writeAllButName(grid *gridt.Grid, file ipe.File, name string) {
	f.write(
		grid,
		strconv.FormatUint(file.Inode(), 10),
		file.Mode().String(),
		fmtSize(file),
		strconv.FormatUint(file.Links(), 10),
		strconv.FormatInt(file.Blocks(), 10),
		fmtTime(file.AccTime()),
		fmtTime(file.ModTime()),
		fmtTime(file.CrtTime()),
		file.User().Username,
		file.Group().Name,
		name,
	)
}

func (f *longFormatter) write(grid *gridt.Grid, inode, mode, size, link, block, acc, mod, crt, user, group, name string) {
	if f.showInode {
		grid.Add(inode)
	}
	grid.Add(mode)
	grid.Add(size)
	if f.showLinks {
		grid.Add(link)
	}
	if f.showBlocks {
		grid.Add(block)
	}
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
	if f.showGroup {
		grid.Add(group)
	}
	grid.Add(name)
}
