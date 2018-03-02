// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"testing"
)

func TestMail(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()
	gt.work(t)

	h := CurrentBranch().Pending()[0].ShortHash

	// fake auth information to avoid Gerrit error
	auth.host = "gerrit.fake"
	auth.user = "not-a-user"
	defer func() {
		auth.host = ""
		auth.user = ""
	}()

	testMain(t, "mail")
	testRan(t,
		"git push -q origin HEAD:refs/for/master",
		"git tag -f work.mailed "+h)
}

func TestMailGitHub(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()
	gt.work(t)

	trun(t, gt.client, "git", "config", "remote.origin.url", "https://github.com/golang/go")

	testMainDied(t, "mail")
	testPrintedStderr(t, "git origin must be a Gerrit host, not GitHub: https://github.com/golang/go")
}

func TestMailAmbiguousRevision(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()
	gt.work(t)

	t.Logf("creating file that conflicts with revision parameter")
	b := CurrentBranch()
	mkdir(t, gt.client+"/origin")
	write(t, gt.client+"/"+b.Branchpoint()+"..HEAD", "foo")

	testMain(t, "mail", "-diff")
}

var reviewerLog = []string{
	"Fake 1 <r1@fake.com>",
	"Fake 1 <r1@fake.com>",
	"Fake 1 <r1@fake.com>",
	"Reviewer 1 <r1@golang.org>",
	"Reviewer 1 <r1@golang.org>",
	"Reviewer 1 <r1@golang.org>",
	"Reviewer 1 <r1@golang.org>",
	"Reviewer 1 <r1@golang.org>",
	"Other <other@golang.org>",
	"<anon@golang.org>",
}

func TestMailShort(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()

	// fake auth information to avoid Gerrit error
	auth.host = "gerrit.fake"
	auth.user = "not-a-user"
	defer func() {
		auth.host = ""
		auth.user = ""
	}()

	// Seed commit history with reviewers.
	for i, addr := range reviewerLog {
		write(t, gt.server+"/file", fmt.Sprintf("v%d", i))
		trun(t, gt.server, "git", "commit", "-a", "-m", "msg\n\nReviewed-by: "+addr+"\n")
	}
	trun(t, gt.client, "git", "pull")

	// Do some work.
	gt.work(t)

	h := CurrentBranch().Pending()[0].ShortHash

	testMain(t, "mail")
	testRan(t,
		"git push -q origin HEAD:refs/for/master",
		"git tag -f work.mailed "+h)

	testMain(t, "mail", "-r", "r1")
	testRan(t,
		"git push -q origin HEAD:refs/for/master%r=r1@golang.org",
		"git tag -f work.mailed "+h)

	testMain(t, "mail", "-r", "other,anon", "-cc", "r1,full@email.com")
	testRan(t,
		"git push -q origin HEAD:refs/for/master%r=other@golang.org,r=anon@golang.org,cc=r1@golang.org,cc=full@email.com",
		"git tag -f work.mailed "+h)

	testMainDied(t, "mail", "-r", "other", "-r", "anon,r1,missing")
	testPrintedStderr(t, "unknown reviewer: missing")
}

func TestMailTopic(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()
	gt.work(t)

	h := CurrentBranch().Pending()[0].ShortHash

	// fake auth information to avoid Gerrit error
	auth.host = "gerrit.fake"
	auth.user = "not-a-user"
	defer func() {
		auth.host = ""
		auth.user = ""
	}()

	testMainDied(t, "mail", "-topic", "contains,comma")
	testPrintedStderr(t, "topic may not contain a comma")

	testMain(t, "mail", "-topic", "test-topic")
	testRan(t,
		"git push -q origin HEAD:refs/for/master%topic=test-topic",
		"git tag -f work.mailed "+h)
}
