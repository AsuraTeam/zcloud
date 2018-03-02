// +build !linux,!openbsd,!freebsd,!netbsd,!dragonfly,!solaris,!darwin

package os

import (
	"errors"
	"os"
)

var errNoSymlinksNotSupported = errors.New("opening files without following symlinks is not supported on this platform")

func openFileNoSymlinks(path string, flags int, mode os.FileMode) (*os.File, error) {
	return nil, errNoSymlinksNotSupported
}
