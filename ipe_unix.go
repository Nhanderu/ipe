// +build !windows

package ipe

import (
	"errors"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"
)

func newFile(dir string, fi os.FileInfo) (File, error) {
	sys := fi.Sys().(*syscall.Stat_t)
	if sys == nil {
		return File{}, errors.New("invalid file attributes")
	}
	u, err := user.LookupId(strconv.FormatUint(uint64(sys.Uid), 10))
	if err != nil {
		return File{}, err
	}
	g, err := user.LookupGroupId(strconv.FormatUint(uint64(sys.Gid), 10))
	if err != nil {
		return File{}, err
	}
	return File{
		fi.Name(),
		dir,
		fi.Size(),
		time.Unix(sys.Atim.Sec, sys.Atim.Nsec),
		time.Unix(sys.Mtim.Sec, sys.Mtim.Nsec),
		time.Unix(sys.Ctim.Sec, sys.Ctim.Nsec),
		fi.Mode(),
		u,
		g,
		sys.Ino,
		sys,
	}, nil
}
