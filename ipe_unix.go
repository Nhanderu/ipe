// +build !windows

package ipe

import (
	"errors"
	"os"
	"os/user"
	"syscall"
	"time"
)

func newFile(dir string, fi os.FileInfo) (File, error) {
	sys := fi.Sys().(*syscall.Stat_t)
	if sys == nil {
		return File{}, errors.New("invalid file attributes")
	}
	u, err := user.LookupId(sys.Uid)
	if err != nil {
		return File{}, err
	}
	g, err := user.LookupGroupId(sys.Gui)
	if err != nil {
		return File{}, err
	}
	println(time.Unix(sys.Mtime, sys.MtimeNsec) == fi.ModTime())
	return File{
		fi.Name(),
		dir,
		fi.Size(),
		time.Unix(sys.Atime, sys.AtimeNsec),
		time.Unix(sys.Mtime, sys.MtimeNsec),
		time.Unix(sys.Ctime, sys.CtimeNsec),
		fi.Mode(),
		u,
		g,
		sys.Ino,
		sys,
	}, nil
}
