package ipefmt

import (
	"bytes"

	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
)

// srcInfo represents the common infomation for an output node.
type srcInfo struct {
	file ipe.File
	err  error
	grid *gridt.Grid
}

// commonFormatter represents the common infomation and methods for the formatters.
type commonFormatter struct {
	args ArgsInfo
	srcs []srcInfo
	cols int
}

// Format outputs the formatter into a correct string.
func (f commonFormatter) String() string {
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
			var d gridt.Dimensions
			var ok bool
			if f.cols > 0 {
				d, ok = src.grid.FitIntoColumns(f.cols)
			} else {
				d, ok = src.grid.FitIntoWidth(f.args.Width)
			}
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

// getName returns the name of the file, based on the arguments.
func (f commonFormatter) getName(file ipe.File) string {
	if f.args.Classify {
		return file.ClassifiedName()
	}
	return file.Name()
}

// appendSource appends another `srcInfo` to its list.
func (f *commonFormatter) appendSource(src srcInfo) {
	f.srcs = append(f.srcs, src)
}

// Format formats the arguments into the correct output.
func Format(args ArgsInfo) string {
	var f formatter
	if args.Long && args.Tree {
		f = newLongTreeFormatter(args)
	} else if args.Long {
		f = newLongFormatter(args)
	} else if args.Tree {
		f = newTreeFormatter(args)
	} else {
		f = newGridFormatter(args)
	}
	var w formatterWrapper
	w.read(f, args)
	return w.String()
}
