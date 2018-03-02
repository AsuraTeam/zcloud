package os

import "os"

// Opens a file but does not follow symlinks.
func OpenFileNoSymlinks(path string, flags int, mode os.FileMode) (*os.File, error) {
	return openFileNoSymlinks(path, flags, mode)
}

// See OpenFileNoSymlinks.
func OpenNoSymlinks(path string) (*os.File, error) {
	return OpenFileNoSymlinks(path, os.O_RDONLY, 0)
}
