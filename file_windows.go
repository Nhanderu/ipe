// +build windows

package ipe

import (
	"errors"
	"os"
	"syscall"
	"time"
)

func newFile(dir string, fi os.FileInfo) (File, error) {
	sys := fi.Sys().(*syscall.Win32FileAttributeData)
	if sys == nil {
		return File{}, errors.New("invalid file attributes")
	}
	return File{
		fi.Name(),
		dir,
		fi.Size(),
		time.Unix(0, sys.LastAccessTime.Nanoseconds()),
		time.Unix(0, sys.LastWriteTime.Nanoseconds()),
		time.Unix(0, sys.CreationTime.Nanoseconds()),
		fi.Mode(),
		nil, // That's a problem.
		nil, // That's a problem.
		0,   // That's a problem.
		sys,
	}, nil
}
