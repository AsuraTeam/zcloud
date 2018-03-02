// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "strings"

func doSync(args []string) {
	expectZeroArgs(args, "sync")

	// Get current branch and commit ID for fixup after pull.
	b := CurrentBranch()
	id := b.ChangeID()

	// Pull remote changes into local branch.
	// We do this in one command so that people following along with 'git sync -v'
	// see fewer commands to understand.
	// We want to pull in the remote changes from the upstream branch
	// and rebase the current pending commit (if any) on top of them.
	// If there is no pending commit, the pull will do a fast-forward merge.
	run("git", "pull", "-q", "-r", "origin", strings.TrimPrefix(b.OriginBranch(), "origin/"))

	// If the change commit has been submitted,
	// roll back change leaving any changes unstaged.
	// Pull should have done this for us, but check just in case.
	if b.Submitted(id) && b.HasPendingCommit() {
		run("git", "reset", "HEAD^")
	}
}
