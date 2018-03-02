package os

import "errors"

var ErrNotEmpty = errors.New("directory not empty")

// Returns true if the error is ErrNotEmpty or the underlying POSIX error code
// error is ENOTEMPTY. Currently always returns false on Windows.
func IsNotEmpty(err error) bool {
	return isNotEmpty(err)
}
