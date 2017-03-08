package ipefmt

import (
	"bytes"

	"github.com/Nhanderu/ipe"
)

type treeFormatter struct {
	srcs []srcInfoBuffer
}

func newTreeFormatter(args ArgsInfo) *treeFormatter {
	f := &treeFormatter{make([]srcInfoBuffer, 0)}
	for _, src := range args.Sources {
		file, err := ipe.Read(fixInSrc(src))
		if err != nil {
			f.srcs = append(f.srcs, srcInfoBuffer{file, err, nil})
		} else {
			f.getDir(file, bytes.NewBuffer([]byte{}), args, []bool{})
		}
	}
	return f
}

func (f *treeFormatter) getDir(file ipe.File, buffer *bytes.Buffer, args ArgsInfo, corners []bool) {
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}

	if len(corners) == 0 {
		f.srcs = append(f.srcs, srcInfoBuffer{file, nil, buffer})
	}

	if args.Reverse {
		reverse(fs)
	}

	// First loop: preparation.
	for _, file := range fs {
		checkBiggestValues(file, args)
	}

	// Second loop: printing.
	for ii, file := range fs {
		f.getFile(file, buffer, args, append(corners, ii+1 == len(fs)))
	}
}

func (f *treeFormatter) getFile(file ipe.File, buffer *bytes.Buffer, args ArgsInfo, corners []bool) {
	if !shouldShow(file, args) {
		return
	}

	buffer.WriteString(makeTree(corners))
	if args.Classify {
		buffer.WriteString(file.ClassifiedName())
	} else {
		buffer.WriteString(file.Name())
	}
	buffer.WriteRune('\n')

	if args.Recursive && file.IsDir() && (args.Depth == 0 || int(args.Depth) >= len(corners)) {
		f.getDir(file, buffer, args, corners)
	}
}

func (f *treeFormatter) String() string {
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
			buffer.WriteString(src.buffer.String())
		}
		if writeNames {
			buffer.WriteString("\n")
		}
	}
	return buffer.String()
}
