package ipefmt

import (
	"bytes"

	"github.com/Nhanderu/ipe"
)

type longFormatter struct {
	srcs []srcInfoBuffer
}

func newLongFormatter(args ArgsInfo) *longFormatter {
	f := &longFormatter{make([]srcInfoBuffer, 0)}
	for _, src := range args.Sources {
		file, err := ipe.Read(fixInSrc(src))
		if err != nil {
			f.srcs = append(f.srcs, srcInfoBuffer{file, err, nil})
		} else {
			f.getDir(file, args, 0)
		}
	}
	return f
}

func (f *longFormatter) getDir(file ipe.File, args ArgsInfo, depth uint8) {
	fs := file.Children()
	if fs == nil || len(fs) == 0 {
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	f.srcs = append(f.srcs, srcInfoBuffer{file, nil, buffer})

	if args.Reverse {
		reverse(fs)
	}

	// First loop: preparation.
	for _, file := range fs {
		checkBiggestValues(file, args)
	}

	// Second loop: printing.
	for _, file := range fs {
		f.getFile(file, buffer, args, depth+1)
	}
}

func (f *longFormatter) getFile(file ipe.File, buffer *bytes.Buffer, args ArgsInfo, depth uint8) {
	if !shouldShow(file, args) {
		return
	}

	acc, mod, crt := timesToShow(args)
	if args.Inode && !osWindows {
		buffer.WriteString(fmtColumn(fmtInode(file), args.Separator, bgstInode))
	}
	buffer.WriteString(fmtColumn(fmtMode(file), args.Separator, bgstMode))
	buffer.WriteString(fmtColumn(fmtSize(file), args.Separator, bgstSize))
	if acc {
		buffer.WriteString(fmtColumn(fmtAccTime(file), args.Separator, bgstAccTime))
	}
	if mod {
		buffer.WriteString(fmtColumn(fmtModTime(file), args.Separator, bgstModTime))
	}
	if crt {
		buffer.WriteString(fmtColumn(fmtCrtTime(file), args.Separator, bgstCrtTime))
	}
	if !osWindows {
		buffer.WriteString(fmtColumn(fmtUser(file), args.Separator, bgstUser))
	}
	if args.Classify {
		buffer.WriteString(file.ClassifiedName())
	} else {
		buffer.WriteString(file.Name())
	}
	buffer.WriteRune('\n')

	if args.Recursive && file.IsDir() && (args.Depth == 0 || args.Depth >= depth) {
		f.getDir(file, args, depth)
	}
}

func (f *longFormatter) String() string {
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
