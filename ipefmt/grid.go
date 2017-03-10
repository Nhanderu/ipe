package ipefmt

import (
	"bytes"

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

	if f.args.Reverse {
		reverse(fs)
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

func (f gridFormatter) String() string {
	var buffer bytes.Buffer
	writeNames := len(f.srcs) > 1
	for _, src := range f.srcs {
		if writeNames {
			buffer.WriteString(src.file.FullName())
			buffer.WriteString("\n")
		}
		if src.err != nil {
			buffer.WriteString("Error: ")
			buffer.WriteString(src.err.Error())
		} else {
			d, ok := src.grid.FitIntoWidth(f.args.Width)
			if !ok || f.args.OneLine {
				for _, cell := range src.grid.Cells() {
					buffer.WriteString(cell)
					buffer.WriteString("\n")
				}
			} else {
				buffer.WriteString(d.String())
			}
		}
		if writeNames {
			buffer.WriteString("\n")
		}
	}
	return buffer.String()
}
