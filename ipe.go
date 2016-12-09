package ipe

import (
	"errors"
	"os"
)

var errIsNotDir = errors.New("entry is not a dir")

type Entry struct {
	absPath  string
	fileInfo os.FileInfo
}

func (e Entry) Name() string {
	return e.fileInfo.Name()
}

func (e Entry) Size() int64 {
	return e.fileInfo.Size()
}

func (e Entry) IsDir() bool {
	return e.fileInfo.IsDir()
}

func (e Entry) AbsPath() string {
	return e.absPath
}
