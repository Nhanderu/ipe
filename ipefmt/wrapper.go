package ipefmt

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
	"github.com/fatih/color"
)

type formatter interface {
	fmt.Stringer

	getDir(file ipe.File, grid **gridt.Grid, corners []bool)
	getFile(file ipe.File, grid *gridt.Grid, corners []bool)

	appendSource(src srcInfo)
}

type formatterWrapper struct {
	formatter formatter
	args      ArgsInfo
}

func (f *formatterWrapper) read(formatter formatter, args ArgsInfo) {
	f.formatter = formatter
	f.args = args
	if args.Color != ArgColorAuto {
		color.NoColor = args.Color == ArgColorNever
	}
	for _, src := range args.Sources {
		file, err := ipe.Read(fixInSrc(src))
		if err != nil {
			f.formatter.appendSource(srcInfo{file, err, nil})
		} else {
			f.getDir(file, gridt.New(gridt.LeftToRight, f.args.Separator), []bool{})
		}
	}
}

func (f *formatterWrapper) getDir(file ipe.File, grid *gridt.Grid, corners []bool) {
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}
	f.formatter.getDir(file, &grid, corners)
	if f.args.Sort != ArgSortNone {
		sort.Slice(fs, func(i, j int) bool {
			switch f.args.Sort {
			case ArgSortInode:
				return fs[i].Inode() < fs[j].Inode()
			case ArgSortMode:
				r := strings.NewReplacer("-", "")
				return r.Replace(fs[i].Mode().String()) < r.Replace(fs[j].Mode().String())
			case ArgSortSize:
				return fs[i].Size() < fs[j].Size()
			case ArgSortAccessed:
				return fs[i].AccTime().Unix() < fs[j].AccTime().Unix()
			case ArgSortModified:
				return fs[i].ModTime().Unix() < fs[j].ModTime().Unix()
			case ArgSortCreated:
				return fs[i].CrtTime().Unix() < fs[j].CrtTime().Unix()
			case ArgSortUser:
				return fs[i].User().Uid < fs[j].User().Uid
			case ArgSortGroup:
				return fs[i].Group().Gid < fs[j].Group().Gid
			case ArgSortName:
				return fs[i].Name() < fs[j].Name()
			default:
				return true
			}
		})
	}
	if f.args.DirsFirst {
		sort.Slice(fs, func(i, j int) bool {
			return fs[i].IsDir() || !fs[j].IsDir()
		})
	}
	if f.args.Reverse {
		for l, r := 0, len(fs)-1; l < r; l, r = l+1, r-1 {
			fs[l], fs[r] = fs[r], fs[l]
		}
	}
	for i, child := range fs {
		f.getFile(child, grid, append(corners, i+1 == len(fs)))
	}
}

func (f *formatterWrapper) getFile(file ipe.File, grid *gridt.Grid, corners []bool) {
	for _, f := range f.args.Filter {
		if !f.MatchString(file.Name()) && !f.MatchString(file.FullName()) {
			return
		}
	}
	for _, i := range f.args.Ignore {
		if i.MatchString(file.Name()) || i.MatchString(file.FullName()) {
			return
		}
	}
	if !f.args.All && file.IsDotfile() {
		return
	}

	f.formatter.getFile(file, grid, corners)

	if f.args.Recursive && file.IsDir() && (f.args.Depth == 0 || int(f.args.Depth) >= len(corners)) {
		f.getDir(file, grid, corners)
	}
}

func (f formatterWrapper) String() string { return f.formatter.String() }
