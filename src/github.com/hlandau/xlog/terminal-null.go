// +build !darwin,!freebsd,!openbsd,!netbsd,!linux,!windows

package xlog

import "io"

func isTerminal(w io.Writer) bool {
	return false
}
