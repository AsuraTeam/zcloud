// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "strings"

func cmdSync(args []string) {
	expectZeroArgs(args, "sync")

	// Get current branch and commit ID for fixup after pull.
	b := CurrentBranch()
	var id string
	if work := b.Pending(); len(work) > 0 {
		id = work[0].ChangeID
	}

	// Don't sync with staged or unstaged changes.
	// rebase is going to complain if we don't, and we can give a nicer error.
	checkStaged("sync")
	checkUnstaged("sync")

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
	b = CurrentBranch() // discard any cached information
	if len(b.Pending()) == 1 && b.Submitted(id) {
		run("git", "reset", b.Branchpoint())
	}
}

func checkStaged(cmd string) {
	if HasStagedChanges() {
		dief("cannot %s: staged changes exist\n"+
			"\trun 'git status' to see changes\n"+
			"\trun 'git-codereview change' to commit staged changes", cmd)
	}
}

func checkUnstaged(cmd string) {
	if HasUnstagedChanges() {
		dief("cannot %s: unstaged changes exist\n"+
			"\trun 'git status' to see changes\n"+
			"\trun 'git stash' to save unstaged changes\n"+
			"\trun 'git add' and 'git-codereview change' to commit staged changes", cmd)
	}
}
