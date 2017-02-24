package ipe

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

// File represents a file.
type File struct {
	name    string
	dir     string
	size    int64
	accTime time.Time
	modTime time.Time
	crtTime time.Time
	mode    os.FileMode
	user    *user.User
	group   *user.Group
	inode   uint64
	sys     interface{}
}

// Name returns the base name of the file.
func (f File) Name() string { return f.name }

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

// FullName returns the full and absolute path of the file.
func (f File) FullName() string {
	return filepath.Join(f.dir, f.Name())
}

// Size returns the length in bytes for regular files.
func (f File) Size() int64 { return f.size }

// ModTime returns the last modification time.
func (f File) ModTime() time.Time { return f.modTime }

// AccTime returns the last access time.
func (f File) AccTime() time.Time { return f.accTime }

// CrtTime returns the creation time.
func (f File) CrtTime() time.Time { return f.crtTime }

// Mode returns the file mode bits.
func (f File) Mode() os.FileMode { return f.mode }

// IsDir reports whether `f` describes a directory.
func (f File) IsDir() bool { return f.mode&os.ModeDir != 0 }

// IsAppend reports whether `f` describes an append-only file.
func (f File) IsAppend() bool { return f.mode&os.ModeAppend != 0 }

// IsExclusive reports whether `f` describes an exclusive-use file.
func (f File) IsExclusive() bool { return f.mode&os.ModeExclusive != 0 }

// IsTemporary reports whether `f` describes a temporary file (not backed up).
func (f File) IsTemporary() bool { return f.mode&os.ModeTemporary != 0 }

// IsSymlink reports whether `f` describes a symbolic link.
func (f File) IsSymlink() bool { return f.mode&os.ModeSymlink != 0 }

// IsDevice reports whether `f` describes a device file.
func (f File) IsDevice() bool { return f.mode&os.ModeDevice != 0 }

// IsNamedPipe reports whether `f` describes a named pipe (FIFO).
func (f File) IsNamedPipe() bool { return f.mode&os.ModeNamedPipe != 0 }

// IsSocket reports whether `f` describes a socket.
func (f File) IsSocket() bool { return f.mode&os.ModeSocket != 0 }

// IsRegular reports whether `f` describes a regular file.
// That is, it tests that no mode type bits are set.
func (f File) IsRegular() bool { return f.mode&os.ModeType == 0 }

// IsDotfile reports whether `f` describes a dotfile.
func (f File) IsDotfile() bool { return f.name[0] == '.' }

// User returns the user of the file.
func (f File) User() *user.User { return f.user }

// Group returns the user group of the file.
func (f File) Group() *user.Group { return f.group }

// Sys represents the underlying data source of the file.
func (f File) Sys() interface{} { return f.sys }

// Children opens a directory and reads its contents.
func (f File) Children() []File {
	if !f.IsDir() {
		return nil
	}
	fs, err := ReadDir(f.FullName())
	if err != nil {
		return nil
	}
	return fs
}

// ReadDir opens and reads the directory path and return its contents.
func ReadDir(path string) ([]File, error) {
	path, f, err := read(path)
	if err != nil {
		return nil, err
	}
	fis, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	files := make([]File, len(fis))
	for i, fi := range fis {
		file, err := newFile(path, fi)
		if err != nil {
			return nil, err
		}
		files[i] = file
	}
	return files, nil
}

// Read opens and reads the path and return its content.
func Read(path string) (File, error) {
	path, f, err := read(path)
	if err != nil {
		return File{}, err
	}
	fi, err := f.Stat()
	if err != nil {
		return File{}, err
	}
	return newFile(filepath.Dir(path), fi)
}

func read(path string) (string, *os.File, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	f, err := os.Open(path)
	return path, f, err
}
