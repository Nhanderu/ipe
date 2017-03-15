// +build windows

package ipe

import (
	"errors"
	"os"
	"os/user"
	"syscall"
	"time"
	"unsafe"
)

func fileno(name string) int {
	fd, _ := syscall.Open(name, syscall.O_RDONLY, 0)
	handle := int(fd)
	value, _, err := syscall.NewLazyDLL("msvcrt.dll").NewProc("_get_osfhandle").Call(uintptr(fd))
	if err.(syscall.Errno) == 0 && value != ^uintptr(0) {
		handle = int(value)
	}
	return handle
}

func newFile(dir string, fi os.FileInfo, fd int) (File, error) {
	sys := fi.Sys().(*syscall.Win32FileAttributeData)
	if sys == nil {
		return File{}, errors.New("invalid file attributes")
	}
	// u, g := getUserAndGroup(fd)
	return File{
		fd,
		fi.Name(),
		dir,
		fi.Size(),
		time.Unix(0, sys.LastAccessTime.Nanoseconds()),
		time.Unix(0, sys.LastWriteTime.Nanoseconds()),
		time.Unix(0, sys.CreationTime.Nanoseconds()),
		fi.Mode(),
		&user.User{},  // That's a problem.
		&user.Group{}, // That's a problem.
		0,             // That's a problem.
		0,             // That's not a problem.
		0,             // That's not a problem.
		sys,
	}, nil
}

func getUserAndGroup(fd int) (*user.User, *user.Group) {
	dll := syscall.NewLazyDLL("advapi32.dll")
	if dll.Load() != nil {
		return nil, nil
	}
	const ownerFlag, groupFlag uintptr = 0x1, 0x2
	var uid, gid int
	dll.NewProc("GetSecurityInfo").Call(
		uintptr(syscall.Handle(fd)),
		1,
		ownerFlag|groupFlag,
		uintptr(unsafe.Pointer(&uid)),
		uintptr(unsafe.Pointer(&gid)),
		0,
		0,
		0)
	return nil, nil
}
