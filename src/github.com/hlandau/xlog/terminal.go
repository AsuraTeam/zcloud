// +build darwin freebsd openbsd netbsd linux windows

package xlog

import "io"
import "github.com/mattn/go-isatty"

func isTerminal(w io.Writer) bool {
	wf, ok := w.(interface {
		Fd() uintptr
	})
	if !ok {
		return false
	}

	return isatty.IsTerminal(wf.Fd())
}
