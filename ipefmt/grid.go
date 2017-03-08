package ipefmt

import (
	"bytes"

	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
)

type gridFormatter struct {
	srcs      []srcInfoGrid
	direction gridt.Direction
	width     int
	oneLine   bool
}

func newGridFormatter(args ArgsInfo) *gridFormatter {
	var direction gridt.Direction
	if args.Across {
		direction = gridt.LeftToRight
	} else {
		direction = gridt.TopToBottom
	}
	f := &gridFormatter{make([]srcInfoGrid, 0), direction, args.Width, args.OneLine}
	for _, src := range args.Sources {
		file, err := ipe.Read(fixInSrc(src))
		if err != nil {
			f.srcs = append(f.srcs, srcInfoGrid{file, err, nil})
		} else {
			f.getDir(file, args, 0)
		}
	}
	return f
}

func (f *gridFormatter) getDir(file ipe.File, args ArgsInfo, depth uint8) {
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}

	grid := gridt.New(f.direction, args.Separator)
	f.srcs = append(f.srcs, srcInfoGrid{file, nil, grid})

	if args.Reverse {
		reverse(fs)
	}

	// First loop: preparation.
	for _, file := range fs {
		checkBiggestValues(file, args)
	}

	// Second loop: printing.
	for _, file := range fs {
		f.getFile(file, grid, args, depth+1)
	}
}

func (f *gridFormatter) getFile(file ipe.File, grid *gridt.Grid, args ArgsInfo, depth uint8) {
	if !shouldShow(file, args) {
		return
	}

	if args.Classify {
		grid.Add(file.ClassifiedName())
	} else {
		grid.Add(file.Name())
	}

	if args.Recursive && file.IsDir() && (args.Depth == 0 || args.Depth >= depth) {
		f.getDir(file, args, depth)
	}
}

func (f *gridFormatter) String() string {
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
			d, ok := src.grid.FitIntoWidth(f.width)
			if !ok || f.oneLine {
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
