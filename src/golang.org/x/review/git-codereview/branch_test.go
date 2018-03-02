// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestCurrentBranch(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()

	t.Logf("on master")
	checkCurrentBranch(t, "master", "origin/master", false, false, "", "")

	t.Logf("on newbranch")
	trun(t, gt.client, "git", "checkout", "-b", "newbranch")
	checkCurrentBranch(t, "newbranch", "origin/master", true, false, "", "")

	t.Logf("making change")
	write(t, gt.client+"/file", "i made a change")
	trun(t, gt.client, "git", "commit", "-a", "-m", "My change line.\n\nChange-Id: I0123456789abcdef0123456789abcdef\n")
	checkCurrentBranch(t, "newbranch", "origin/master", true, true, "I0123456789abcdef0123456789abcdef", "My change line.")

	t.Logf("on dev.branch")
	trun(t, gt.client, "git", "checkout", "-t", "-b", "dev.branch", "origin/dev.branch")
	checkCurrentBranch(t, "dev.branch", "origin/dev.branch", false, false, "", "")

	t.Logf("on newdev")
	trun(t, gt.client, "git", "checkout", "-t", "-b", "newdev", "origin/dev.branch")
	checkCurrentBranch(t, "newdev", "origin/dev.branch", true, false, "", "")

	t.Logf("making change")
	write(t, gt.client+"/file", "i made another change")
	trun(t, gt.client, "git", "commit", "-a", "-m", "My other change line.\n\nChange-Id: I1123456789abcdef0123456789abcdef\n")
	checkCurrentBranch(t, "newdev", "origin/dev.branch", true, true, "I1123456789abcdef0123456789abcdef", "My other change line.")

	t.Logf("detached head mode")
	trun(t, gt.client, "git", "checkout", "HEAD^0")
	checkCurrentBranch(t, "HEAD", "origin/HEAD", false, false, "", "")
}

func checkCurrentBranch(t *testing.T, name, origin string, isLocal, hasPending bool, changeID, subject string) {
	b := CurrentBranch()
	if b.Name != name {
		t.Errorf("b.Name = %q, want %q", b.Name, name)
	}
	if x := b.OriginBranch(); x != origin {
		t.Errorf("b.OriginBranch() = %q, want %q", x, origin)
	}
	if x := b.IsLocalOnly(); x != isLocal {
		t.Errorf("b.IsLocalOnly() = %v, want %v", x, isLocal)
	}
	if x := b.HasPendingCommit(); x != hasPending {
		t.Errorf("b.HasPendingCommit() = %v, want %v", x, isLocal)
	}
	if work := b.Pending(); len(work) > 0 {
		c := work[0]
		if x := c.ChangeID; x != changeID {
			t.Errorf("b.Pending()[0].ChangeID = %q, want %q", x, changeID)
		}
		if x := c.Subject; x != subject {
			t.Errorf("b.Pending()[0].Subject = %q, want %q", x, subject)
		}
	}
}

func TestLocalBranches(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()

	t.Logf("on master")
	checkLocalBranches(t, "master")

	t.Logf("on dev branch")
	trun(t, gt.client, "git", "checkout", "-b", "newbranch")
	checkLocalBranches(t, "master", "newbranch")

	t.Logf("detached head mode")
	trun(t, gt.client, "git", "checkout", "HEAD^0")
	checkLocalBranches(t, "HEAD", "master", "newbranch")
}

func checkLocalBranches(t *testing.T, want ...string) {
	var names []string
	branches := LocalBranches()
	for _, b := range branches {
		names = append(names, b.Name)
	}
	if !reflect.DeepEqual(names, want) {
		t.Errorf("LocalBranches() = %v, want %v", names, want)
	}
}

func TestAmbiguousRevision(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()
	gt.work(t)

	t.Logf("creating file paths that conflict with revision parameters")
	mkdir(t, gt.client+"/origin")
	write(t, gt.client+"/origin/master..work", "Uh-Oh! SpaghettiOs")
	mkdir(t, gt.client+"/work..origin")
	write(t, gt.client+"/work..origin/master", "Be sure to drink your Ovaltine")

	b := CurrentBranch()
	b.Submitted("I123456789")
}

func TestBranchpoint(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()

	// Get hash corresponding to checkout (known to server).
	hash := strings.TrimSpace(trun(t, gt.client, "git", "rev-parse", "HEAD"))

	// Any work we do after this point should find hash as branchpoint.
	for i := 0; i < 4; i++ {
		testMain(t, "branchpoint")
		t.Logf("numCommits=%d", i)
		testPrintedStdout(t, hash)
		testNoStderr(t)

		gt.work(t)
	}
}

func TestRebaseWork(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()

	// Get hash corresponding to checkout (known to server).
	// Any work we do after this point should find hash as branchpoint.
	hash := strings.TrimSpace(trun(t, gt.client, "git", "rev-parse", "HEAD"))

	testMainDied(t, "rebase-work", "-n")
	testPrintedStderr(t, "no pending work")

	write(t, gt.client+"/file", "uncommitted")
	testMainDied(t, "rebase-work", "-n")
	testPrintedStderr(t, "cannot rebase with uncommitted work")

	gt.work(t)

	for i := 0; i < 4; i++ {
		testMain(t, "rebase-work", "-n")
		t.Logf("numCommits=%d", i)
		testPrintedStderr(t, "git rebase -i "+hash)

		gt.work(t)
	}
}

func TestBranchpointMerge(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()

	// commit more work on master
	write(t, gt.server+"/file", "more work")
	trun(t, gt.server, "git", "commit", "-m", "work", "file")

	// update client
	trun(t, gt.client, "git", "checkout", "master")
	trun(t, gt.client, "git", "pull")

	hash := strings.TrimSpace(trun(t, gt.client, "git", "rev-parse", "HEAD"))

	// merge dev.branch
	testMain(t, "change", "work")
	trun(t, gt.client, "git", "merge", "-m", "merge", "origin/dev.branch")

	// check branchpoint is old head (despite this commit having two parents)
	bp := CurrentBranch().Branchpoint()
	if bp != hash {
		t.Logf("branches:\n%s", trun(t, gt.client, "git", "branch", "-a", "-v"))
		t.Logf("log:\n%s", trun(t, gt.client, "git", "log", "--graph", "--decorate"))
		t.Logf("log origin/master..HEAD:\n%s", trun(t, gt.client, "git", "log", "origin/master..HEAD"))
		t.Fatalf("branchpoint=%q, want %q", bp, hash)
	}
}
