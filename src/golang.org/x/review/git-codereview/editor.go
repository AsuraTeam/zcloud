// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

// editor invokes an interactive editor on a temporary file containing
// initial, blocks until the editor exits, and returns the (possibly
// edited) contents of the temporary file. It follows the conventions
// of git for selecting and invoking the editor (see git-var(1)).
func editor(initial string) string {
	// Query the git editor command.
	gitEditor := trim(cmdOutput("git", "var", "GIT_EDITOR"))

	// Create temporary file.
	temp, err := ioutil.TempFile("", "git-codereview")
	if err != nil {
		dief("creating temp file: %v", err)
	}
	tempName := temp.Name()
	defer os.Remove(tempName)
	if _, err := io.WriteString(temp, initial); err != nil {
		dief("%v", err)
	}
	if err := temp.Close(); err != nil {
		dief("%v", err)
	}

	// Invoke the editor. See git's prepare_shell_cmd.
	cmd := exec.Command("sh", "-c", gitEditor+" \"$@\"", gitEditor, tempName)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tempName)
		dief("editor exited with: %v", err)
	}

	// Read the edited file.
	b, err := ioutil.ReadFile(tempName)
	if err != nil {
		dief("%v", err)
	}
	return string(b)
}
