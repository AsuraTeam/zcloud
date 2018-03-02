package os

import (
	"fmt"
	"os/exec"
	"syscall"
)

// Given an error, determines if that error is an os/exec ExitError. If it is,
// and the platform is supported, returns the exitCode and nil error.
// Otherwise, passes through the error, and exitCode is undefined.
func GetExitCode(exitErr error) (exitCode int, err error) {
	if e, ok := exitErr.(*exec.ExitError); ok {
		return getExitCodeSys(e.ProcessState.Sys())
	}

	return 0, err
}

func getExitCodeSys(sys interface{}) (exitCode int, err error) {
	ws, ok := sys.(syscall.WaitStatus)
	if !ok {
		// Should never happen.
		panic(fmt.Sprintf("unknown system exit code type: %T", sys))
	}

	return ws.ExitStatus(), nil
}
