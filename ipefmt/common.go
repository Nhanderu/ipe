package ipefmt

import (
	"bytes"
	"fmt"
	"io"

	"github.com/Nhanderu/gridt"
	"github.com/Nhanderu/ipe"
)

type Formatter interface {
	fmt.Stringer
	io.WriterTo

	getDir(file ipe.File, grid **gridt.Grid, corners []bool)
	getFile(file ipe.File, grid *gridt.Grid, corners []bool)
	appendSource(src srcInfo)
}

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

// String outputs the formatter into a correct string.
func (f commonFormatter) String() string {
	var buffer bytes.Buffer
	f.WriteTo(&buffer)
	return buffer.String()
}

// WriteTo writes the values of the formatter into a writer.
func (f commonFormatter) WriteTo(w io.Writer) (int64, error) {
	writeNames := len(f.srcs) > 1
	var total int
	var err error
	for _, src := range f.srcs {
		var n int
		if writeNames {
			n, err = w.Write([]byte(src.file.FullName()))
			if total += n; err != nil {
				break
			}
			n, err = w.Write([]byte("\n"))
			if total += n; err != nil {
				break
			}
		}
		if src.err != nil {
			n, err = w.Write([]byte("Error: "))
			if total += n; err != nil {
				break
			}
			n, err = w.Write([]byte(src.err.Error()))
			if total += n; err != nil {
				break
			}
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
					n, err = w.Write([]byte(cell))
					if total += n; err != nil {
						break
					}
					n, err = w.Write([]byte("\n"))
					if total += n; err != nil {
						break
					}
				}
			} else {
				n, err = w.Write([]byte(d.String()))
				if total += n; err != nil {
					break
				}
			}
		}
		if writeNames {
			n, err = w.Write([]byte("\n"))
			if total += n; err != nil {
				break
			}
		}
	}
	return int64(total), err
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

// NewFormatter returns the correct formatter based on the arguments.
func NewFormatter(args ArgsInfo) Formatter {
	if args.Long && args.Tree {
		return wrap(newLongTreeFormatter(args), args)
	}
	if args.Long {
		return wrap(newLongFormatter(args), args)
	}
	if args.Tree {
		return wrap(newTreeFormatter(args), args)
	}
	return wrap(newGridFormatter(args), args)
}
