// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "testing"

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
	if x := b.ChangeID(); x != changeID {
		t.Errorf("b.ChangeID() = %q, want %q", x, changeID)
	}
	if x := b.Subject(); x != subject {
		t.Errorf("b.Subject() = %q, want %q", x, subject)
	}
}
