package ipe

import (
	"fmt"
	"os"
	"path/filepath"
)

// File represents a file.
type File struct {
	os.FileInfo
	dir string
}

// ClassifiedName returns the name with an appended type indicator.
func (f File) ClassifiedName() string {
	switch {
	case f.IsDir():
		return fmt.Sprint(f.Name(), "/")
	case f.IsSymlink():
		return fmt.Sprint(f.Name(), "@")
	case f.IsNamedPipe():
		return fmt.Sprint(f.Name(), "|")
	case f.IsSocket():
		return fmt.Sprint(f.Name(), "=")
	default:
		return f.Name()
	}
}

// IsDir reports whether `f` describes a directory.
func (f File) IsDir() bool { return f.Mode()&os.ModeDir != 0 }

// IsAppend reports whether `f` describes an append-only file.
func (f File) IsAppend() bool { return f.Mode()&os.ModeAppend != 0 }

// IsExclusive reports whether `f` describes an exclusive-use file.
func (f File) IsExclusive() bool { return f.Mode()&os.ModeExclusive != 0 }

// IsTemporary reports whether `f` describes a temporary file (not backed up).
func (f File) IsTemporary() bool { return f.Mode()&os.ModeTemporary != 0 }

// IsSymlink reports whether `f` describes a symbolic link.
func (f File) IsSymlink() bool { return f.Mode()&os.ModeSymlink != 0 }

// IsDevice reports whether `f` describes a device file.
func (f File) IsDevice() bool { return f.Mode()&os.ModeDevice != 0 }

// IsNamedPipe reports whether `f` describes a named pipe (FIFO).
func (f File) IsNamedPipe() bool { return f.Mode()&os.ModeNamedPipe != 0 }

// IsSocket reports whether `f` describes a socket.
func (f File) IsSocket() bool { return f.Mode()&os.ModeSocket != 0 }

// IsDotfile reports whether `f` describes a dotfile.
func (f File) IsDotfile() bool { return f.Name()[0] == '.' }

// Children opens a directory and reads its contents.
func (f File) Children() []File {
	if !f.IsDir() {
		return nil
	}
	fs, err := ReadDir(filepath.Join(f.dir, f.Name()))
	if err != nil {
		return nil
	}
	return fs
}

// ReadDir opens and reads the contents of a given directory.
func ReadDir(dirpath string) ([]File, error) {
	f, err := os.Open(dirpath)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	rlist := make([]File, len(list))
	for i, file := range list {
		rlist[i] = File{file, dirpath}
	}
	return rlist, nil
}
