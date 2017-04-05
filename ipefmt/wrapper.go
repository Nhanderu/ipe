package ipefmt

import (
	"sort"
	"strings"

	"path/filepath"

	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
	"github.com/fatih/color"
)

type formatterWrapper struct {
	Formatter
	args       ArgsInfo
	filterGlob []string
	ignoreGlob []string
}

func (f *formatterWrapper) getDir(file ipe.File, grid **gridt.Grid, corners []bool) {
	// Gets all the files inside the directory.
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}
	f.Formatter.getDir(file, grid, corners)

	// Sorts the files, based on the flags.
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

	// Formats every file.
	for i, child := range fs {
		f.getFile(child, *grid, append(corners, i+1 == len(fs)))
	}
}

func (f *formatterWrapper) getFile(file ipe.File, grid *gridt.Grid, corners []bool) {
	// Validates, if the file should really appear, based on the flags.
	for _, re := range f.args.FilterRegex {
		if !re.MatchString(file.Name()) && !re.MatchString(file.FullName()) {
			return
		}
	}
	for _, g := range f.filterGlob {
		if g != file.FullName() {
			return
		}
	}
	for _, re := range f.args.IgnoreRegex {
		if re.MatchString(file.Name()) || re.MatchString(file.FullName()) {
			return
		}
	}
	for _, g := range f.ignoreGlob {
		if g == file.FullName() {
			return
		}
	}
	if !f.args.All && file.IsDotfile() {
		return
	}

	// Adds the files to the specific formatter.
	f.Formatter.getFile(file, grid, corners)

	// Recurses.
	if f.args.Recursive && file.IsDir() && (f.args.Depth == 0 || int(f.args.Depth) >= len(corners)) {
		f.getDir(file, &grid, corners)
	}
}

func wrap(formatter Formatter, args ArgsInfo) *formatterWrapper {
	var f formatterWrapper
	f.Formatter = formatter
	f.args = args
	f.filterGlob = make([]string, 0)
	for _, glob := range args.FilterGlob {
		matches, err := filepath.Glob(glob)
		if err == nil && matches != nil {
			for _, match := range matches {
				abs, err := filepath.Abs(match)
				if err != nil {
					f.filterGlob = append(f.filterGlob, abs)
				}
			}
		}
	}
	f.ignoreGlob = make([]string, 0)
	for _, glob := range args.IgnoreGlob {
		matches, err := filepath.Glob(glob)
		if err == nil && matches != nil {
			for _, match := range matches {
				abs, err := filepath.Abs(match)
				if err != nil {
					f.ignoreGlob = append(f.ignoreGlob, abs)
				}
			}
		}
	}
	if f.args.Color != ArgColorAuto {
		color.NoColor = f.args.Color == ArgColorNever
	}
	for _, src := range f.args.Sources {
		file, err := ipe.Read(fixInSrc(src))
		if err != nil {
			f.Formatter.appendSource(srcInfo{file, err, nil})
		} else {
			g := gridt.New(gridt.LeftToRight, f.args.Separator)
			f.getDir(file, &g, []bool{})
		}
	}
	return &f
}
