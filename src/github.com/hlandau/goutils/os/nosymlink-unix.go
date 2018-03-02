// +build linux freebsd openbsd netbsd dragonfly solaris darwin

package os

import (
	"os"
	"syscall"
)

func openFileNoSymlinks(path string, flags int, mode os.FileMode) (*os.File, error) {
	return os.OpenFile(path, flags|syscall.O_NOFOLLOW, mode)
}
