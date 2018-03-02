// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "testing"

func TestSync(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()

	testMain(t, "change", "work")

	// check for error with unstaged changes
	write(t, gt.client+"/file1", "")
	trun(t, gt.client, "git", "add", "file1")
	write(t, gt.client+"/file1", "actual content")
	testMainDied(t, "sync")
	testPrintedStderr(t, "cannot sync: unstaged changes exist",
		"git status", "git stash", "git add", "git-codereview change")
	testNoStdout(t)

	// check for error with staged changes
	trun(t, gt.client, "git", "add", "file1")
	testMainDied(t, "sync")
	testPrintedStderr(t, "cannot sync: staged changes exist",
		"git status", "!git stash", "!git add", "git-codereview change")
	testNoStdout(t)

	// check for success after stash
	trun(t, gt.client, "git", "stash")
	testMain(t, "sync")
	testNoStdout(t)
	testNoStderr(t)

	// make server 1 step ahead of client
	write(t, gt.server+"/file", "new content")
	trun(t, gt.server, "git", "add", "file")
	trun(t, gt.server, "git", "commit", "-m", "msg")

	// check for success
	testMain(t, "sync")
	testNoStdout(t)
	testNoStderr(t)
}

func TestSyncRebase(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()

	// client 3 ahead
	gt.work(t)
	gt.work(t)
	gt.work(t)

	b := CurrentBranch()
	if len(b.Pending()) != 3 {
		t.Fatalf("have %d pending CLs, want 3", len(b.Pending()))
	}
	top := b.Pending()[0].Hash

	// check for success for sync no-op
	testMain(t, "sync")
	testNoStdout(t)
	testNoStderr(t)

	b = CurrentBranch()
	if len(b.Pending()) != 3 {
		t.Fatalf("have %d pending CLs after no-op sync, want 3", len(b.Pending()))
	}
	if b.Pending()[0].Hash != top {
		t.Fatalf("CL hashes changed during no-op sync")
	}

	// submit first two CLs - gt.serverWork does same thing gt.work does, but on server

	gt.serverWork(t)
	gt.serverWorkUnrelated(t) // wedge in unrelated work to get different hashes
	gt.serverWork(t)

	testMain(t, "sync")
	testNoStdout(t)
	testNoStderr(t)

	// there should be one left, and it should be a different hash
	b = CurrentBranch()
	if len(b.Pending()) != 1 {
		t.Fatalf("have %d pending CLs after submitting two, want 1", len(b.Pending()))
	}
	if b.Pending()[0].Hash == top {
		t.Fatalf("CL hashes DID NOT change during sync after submit")
	}

	// submit final change
	gt.serverWork(t)

	testMain(t, "sync")
	testNoStdout(t)
	testNoStderr(t)

	// there should be none left
	b = CurrentBranch()
	if len(b.Pending()) != 0 {
		t.Fatalf("have %d pending CLs after final sync, want 0", len(b.Pending()))
	}
}
