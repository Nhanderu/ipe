package ipe

import "os"

type File struct {
	fileInfo os.FileInfo
}

type RecurFile struct {
	File
	Children []RecurFile
}

func (f File) Name() string { return f.fileInfo.Name() }
func (f File) Size() int64  { return f.fileInfo.Size() }

func (f File) IsDir() bool        { return f.fileInfo.Mode()&os.ModeDir != 0 }
func (f File) IsAppend() bool     { return f.fileInfo.Mode()&os.ModeAppend != 0 }
func (f File) IsExclusive() bool  { return f.fileInfo.Mode()&os.ModeExclusive != 0 }
func (f File) IsTemporary() bool  { return f.fileInfo.Mode()&os.ModeTemporary != 0 }
func (f File) IsSymlink() bool    { return f.fileInfo.Mode()&os.ModeSymlink != 0 }
func (f File) IsDevice() bool     { return f.fileInfo.Mode()&os.ModeDevice != 0 }
func (f File) IsNamedPipe() bool  { return f.fileInfo.Mode()&os.ModeNamedPipe != 0 }
func (f File) IsSocket() bool     { return f.fileInfo.Mode()&os.ModeSocket != 0 }
func (f File) IsSetuid() bool     { return f.fileInfo.Mode()&os.ModeSetuid != 0 }
func (f File) IsSetgid() bool     { return f.fileInfo.Mode()&os.ModeSetgid != 0 }
func (f File) IsCharDevice() bool { return f.fileInfo.Mode()&os.ModeCharDevice != 0 }
func (f File) IsSticky() bool     { return f.fileInfo.Mode()&os.ModeSticky != 0 }

func Open(path string) ([]File, error) {
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
