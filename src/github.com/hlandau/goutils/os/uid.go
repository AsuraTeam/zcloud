// +build !windows

package os

import (
	"fmt"
	"os"
	"syscall"
)

// Returns the UID for a file.
func GetFileUID(fi os.FileInfo) (int, error) {
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, fmt.Errorf("unknown file info type: %T", fi.Sys())
	}

	return int(st.Uid), nil
}

// Returns the GID for a file.
func GetFileGID(fi os.FileInfo) (int, error) {
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, fmt.Errorf("unknown file info type: %T", fi.Sys())
	}

	return int(st.Gid), nil
}
