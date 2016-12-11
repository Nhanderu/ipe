package ipe

import "os"
import "fmt"
import "path/filepath"

type File struct {
	os.FileInfo
}

type RecurFile struct {
	File
	Children []RecurFile
}

func (f File) ClassifiedName() string {
	switch {
	case f.IsDir():
		return fmt.Sprintf("%s%s", f.Name(), "/")
	case f.IsSymlink():
		return fmt.Sprintf("%s%s", f.Name(), "@")
	case f.IsNamedPipe():
		return fmt.Sprintf("%s%s", f.Name(), "|")
	case f.IsSocket():
		return fmt.Sprintf("%s%s", f.Name(), "=")
	default:
		return f.Name()
	}
}

func (f File) IsDir() bool        { return f.Mode()&os.ModeDir != 0 }
func (f File) IsAppend() bool     { return f.Mode()&os.ModeAppend != 0 }
func (f File) IsExclusive() bool  { return f.Mode()&os.ModeExclusive != 0 }
func (f File) IsTemporary() bool  { return f.Mode()&os.ModeTemporary != 0 }
func (f File) IsSymlink() bool    { return f.Mode()&os.ModeSymlink != 0 }
func (f File) IsDevice() bool     { return f.Mode()&os.ModeDevice != 0 }
func (f File) IsNamedPipe() bool  { return f.Mode()&os.ModeNamedPipe != 0 }
func (f File) IsSocket() bool     { return f.Mode()&os.ModeSocket != 0 }
func (f File) IsSetuid() bool     { return f.Mode()&os.ModeSetuid != 0 }
func (f File) IsSetgid() bool     { return f.Mode()&os.ModeSetgid != 0 }
func (f File) IsCharDevice() bool { return f.Mode()&os.ModeCharDevice != 0 }
func (f File) IsSticky() bool     { return f.Mode()&os.ModeSticky != 0 }

func ReadDir(path string) ([]File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	rlist := make([]File, len(list))
	for _, file := range list {
		rlist = append(rlist, File{file})
	}
	return rlist, nil
}

func ReadRecurDir(path string) ([]RecurFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	rlist := make([]RecurFile, len(list))
	for _, file := range list {
		var innerFiles []RecurFile
		if file.IsDir() {
			innerFiles, err = ReadRecurDir(filepath.Join(path, file.Name()))
			if err != nil {
				return nil, err
			}
		}
		rlist = append(rlist, RecurFile{File{file}, innerFiles})
	}
	return rlist, nil
}
