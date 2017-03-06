// +build windows

package ipe

import (
	"errors"
	"os"
	"os/user"
	"syscall"
	"time"
)

func newFile(dir string, fi os.FileInfo) (File, error) {
	sys := fi.Sys().(*syscall.Win32FileAttributeData)
	if sys == nil {
		return File{}, errors.New("invalid file attributes")
	}
	var u *user.User
	var g *user.Group
	lib := syscall.NewLazyDLL("advapi32.lib")
	if lib != nil {
		// lib.NewProc("GetSecurityInfo").Call()
	}
	//syscall.LookupAccountName
	return File{
		fi.Name(),
		dir,
		fi.Size(),
		time.Unix(0, sys.LastAccessTime.Nanoseconds()),
		time.Unix(0, sys.LastWriteTime.Nanoseconds()),
		time.Unix(0, sys.CreationTime.Nanoseconds()),
		fi.Mode(),
		u, // That's a problem.
		g, // That's a problem.
		0, // That's a problem.
		sys,
	}, nil
}
