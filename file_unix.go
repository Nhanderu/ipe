// +build !windows

package ipe

import (
	"errors"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

func fileno(name string) int {
	var fd int
	for {
		var err error
		fd, err = syscall.Open(name, syscall.O_RDONLY|syscall.O_CLOEXEC, 0)
		if err != nil {
			if runtime.GOOS == "darwin" && err == syscall.EINTR {
				continue
			}
			return -1
		}
		break
	}

	// TODO
	// if !supportsCloseOnExec {
	// 	syscall.CloseOnExec(fd)
	// }

	return fd
}

func newFile(dir string, fi os.FileInfo, fd int) (File, error) {
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
		fd,
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
