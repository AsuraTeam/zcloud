// +build !windows

package os

import "os"
import "syscall"

func isNotEmpty(err error) bool {
	switch pe := err.(type) {
	default:
		return false
	case *os.PathError:
		err = pe.Err
	case *os.LinkError:
		err = pe.Err
	}

	return err == syscall.ENOTEMPTY || err == ErrNotEmpty
}
