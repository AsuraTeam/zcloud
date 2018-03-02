// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const lenChangeId = len("\n\nChange-Id: I") + 2*20

func TestHookCommitMsg(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()

	// Check that hook adds Change-Id.
	write(t, gt.client+"/msg.txt", "Test message.\n")
	testMain(t, "hook-invoke", "commit-msg", gt.client+"/msg.txt")
	data := read(t, gt.client+"/msg.txt")
	if !bytes.Contains(data, []byte("\n\nChange-Id: ")) {
		t.Fatalf("after hook-invoke commit-msg, missing Change-Id:\n%s", data)
	}

	// Check that hook is no-op when Change-Id is already present.
	testMain(t, "hook-invoke", "commit-msg", gt.client+"/msg.txt")
	data1 := read(t, gt.client+"/msg.txt")
	if !bytes.Equal(data, data1) {
		t.Fatalf("second hook-invoke commit-msg changed Change-Id:\nbefore:\n%s\n\nafter:\n%s", data, data1)
	}

	// Check that hook rejects multiple Change-Ids.
	write(t, gt.client+"/msgdouble.txt", string(data)+string(data))
	testMainDied(t, "hook-invoke", "commit-msg", gt.client+"/msgdouble.txt")
	const multiple = "git-codereview: multiple Change-Id lines\n"
	if got := testStderr.String(); got != multiple {
		t.Fatalf("unexpected output:\ngot: %q\nwant: %q", got, multiple)
	}

	// Check that hook fails when message is empty.
	write(t, gt.client+"/empty.txt", "\n\n# just a file with\n# comments\n")
	testMainDied(t, "hook-invoke", "commit-msg", gt.client+"/empty.txt")
	const empty = "git-codereview: empty commit message\n"
	if got := testStderr.String(); got != empty {
		t.Fatalf("unexpected output:\ngot: %q\nwant: %q", got, empty)
	}

	// Check that hook inserts a blank line after the first line as needed.
	rewrites := []struct {
		in   string
		want string
	}{
		{in: "all: gofmt", want: "all: gofmt"},
		{in: "all: gofmt\n", want: "all: gofmt\n"},
		{in: "all: gofmt\nahhh", want: "all: gofmt\n\nahhh"},
		{in: "all: gofmt\n\nahhh", want: "all: gofmt\n\nahhh"},
		{in: "all: gofmt\n\n\nahhh", want: "all: gofmt\n\n\nahhh"},
	}
	for _, tt := range rewrites {
		write(t, gt.client+"/in.txt", tt.in)
		testMain(t, "hook-invoke", "commit-msg", gt.client+"/in.txt")
		write(t, gt.client+"/want.txt", tt.want)
		testMain(t, "hook-invoke", "commit-msg", gt.client+"/want.txt")
		got, err := ioutil.ReadFile(gt.client + "/in.txt")
		if err != nil {
			t.Fatal(err)
		}
		want, err := ioutil.ReadFile(gt.client + "/want.txt")
		if err != nil {
			t.Fatal(err)
		}

		// pull off the Change-Id that got appended
		got = got[:len(got)-lenChangeId]
		want = want[:len(want)-lenChangeId]
		if !bytes.Equal(got, want) {
			t.Fatalf("failed to rewrite:\n%s\n\ngot:\n\n%s\n\nwant:\n\n%s\n", tt.in, got, want)
		}
	}
}

func TestHookCommitMsgIssueRepoRewrite(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()

	msgs := []string{
		// If there's no config, don't rewrite issue references.
		"math/big: catch all the rats\n\nFixes #99999, at least for now\n",
		// Fix the fix-message, even without config
		"math/big: catch all the rats\n\nFixes issue #99999, at least for now\n",
		"math/big: catch all the rats\n\nFixes issue 99999, at least for now\n",
		// Don't forget to write back if Change-Id already exists
		"math/big: catch all the rats\n\nFixes issue #99999, at least for now\n\nChange-Id: Ie77358867e38cf976a0688b6e2f80525dae3891e\n",
	}
	for _, msg := range msgs {
		write(t, gt.client+"/msg.txt", msg)
		testMain(t, "hook-invoke", "commit-msg", gt.client+"/msg.txt")
		got := read(t, gt.client+"/msg.txt")
		got = got[:len(got)-lenChangeId]
		const want = "math/big: catch all the rats\n\nFixes #99999, at least for now\n"
		if string(got) != want {
			t.Errorf("issue rewrite failed: got\n\n%s\nwant\n\n%s\nlen %d and %d", got, want, len(got), len(want))
		}
	}

	// Add issuerepo config.
	write(t, gt.client+"/codereview.cfg", "issuerepo: golang/go")
	trun(t, gt.client, "git", "add", "codereview.cfg")
	trun(t, gt.client, "git", "commit", "-m", "add issuerepo codereview config")

	// Look in master rather than origin/master for the config
	savedConfigRef := configRef
	configRef = "master:codereview.cfg"
	cachedConfig = nil

	// Check for the rewrite
	msgs = []string{
		"math/big: catch all the rats\n\nFixes #99999, at least for now\n",
		"math/big: catch all the rats\n\nFixes issue #99999, at least for now\n",
		"math/big: catch all the rats\n\nFixes issue 99999, at least for now\n",
		"math/big: catch all the rats\n\nFixes issue golang/go#99999, at least for now\n",
		"math/big: catch all the rats\n\nFixes issue #99999, at least for now\n\nChange-Id: Ie77358867e38cf976a0688b6e2f80525dae3891e\n",
	}
	for _, msg := range msgs {
		write(t, gt.client+"/msg.txt", msg)
		testMain(t, "hook-invoke", "commit-msg", gt.client+"/msg.txt")
		got := read(t, gt.client+"/msg.txt")
		got = got[:len(got)-lenChangeId]
		const want = "math/big: catch all the rats\n\nFixes golang/go#99999, at least for now\n"
		if string(got) != want {
			t.Errorf("issue rewrite failed: got\n\n%s\nwant\n\n%s", got, want)
		}
	}

	// Reset config state
	configRef = savedConfigRef
	cachedConfig = nil
}

func TestHookCommitMsgBranchPrefix(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()

	checkPrefix := func(prefix string) {
		write(t, gt.client+"/msg.txt", "Test message.\n")
		testMain(t, "hook-invoke", "commit-msg", gt.client+"/msg.txt")
		data, err := ioutil.ReadFile(gt.client + "/msg.txt")
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.HasPrefix(data, []byte(prefix)) {
			t.Errorf("after hook-invoke commit-msg on %s, want prefix %q:\n%s", CurrentBranch().Name, prefix, data)
		}

		if i := strings.Index(prefix, "]"); i >= 0 {
			prefix := prefix[:i+1]
			for _, magic := range []string{"fixup!", "squash!"} {
				write(t, gt.client+"/msg.txt", magic+" Test message.\n")
				testMain(t, "hook-invoke", "commit-msg", gt.client+"/msg.txt")
				data, err := ioutil.ReadFile(gt.client + "/msg.txt")
				if err != nil {
					t.Fatal(err)
				}
				if bytes.HasPrefix(data, []byte(prefix)) {
					t.Errorf("after hook-invoke commit-msg on %s with %s, found incorrect prefix %q:\n%s", CurrentBranch().Name, magic, prefix, data)
				}
			}
		}
	}

	// Create server branch and switch to server branch on client.
	// Test that commit hook adds prefix.
	trun(t, gt.server, "git", "checkout", "-b", "dev.cc")
	trun(t, gt.client, "git", "fetch", "-q")
	testMain(t, "change", "dev.cc")
	checkPrefix("[dev.cc] Test message.\n")

	// Work branch with server branch as upstream.
	testMain(t, "change", "ccwork")
	checkPrefix("[dev.cc] Test message.\n")

	// Master has no prefix.
	testMain(t, "change", "master")
	checkPrefix("Test message.\n")

	// Work branch from master has no prefix.
	testMain(t, "change", "work")
	checkPrefix("Test message.\n")
}

func TestHookPreCommit(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()

	// Write out a non-Go file.
	testMain(t, "change", "mybranch")
	write(t, gt.client+"/msg.txt", "A test message.")
	trun(t, gt.client, "git", "add", "msg.txt")
	testMain(t, "hook-invoke", "pre-commit") // should be no-op

	if err := os.MkdirAll(gt.client+"/test/bench", 0755); err != nil {
		t.Fatal(err)
	}
	write(t, gt.client+"/bad.go", badGo)
	write(t, gt.client+"/good.go", goodGo)
	write(t, gt.client+"/test/bad.go", badGo)
	write(t, gt.client+"/test/good.go", goodGo)
	write(t, gt.client+"/test/bench/bad.go", badGo)
	write(t, gt.client+"/test/bench/good.go", goodGo)
	trun(t, gt.client, "git", "add", ".")

	testMainDied(t, "hook-invoke", "pre-commit")
	testPrintedStderr(t, "gofmt needs to format these files (run 'git gofmt'):",
		"bad.go", "!good.go", fromSlash("!test/bad"), fromSlash("test/bench/bad.go"))

	write(t, gt.client+"/broken.go", brokenGo)
	trun(t, gt.client, "git", "add", "broken.go")
	testMainDied(t, "hook-invoke", "pre-commit")
	testPrintedStderr(t, "gofmt needs to format these files (run 'git gofmt'):",
		"bad.go", "!good.go", fromSlash("!test/bad"), fromSlash("test/bench/bad.go"),
		"gofmt reported errors:", "broken.go")
}

func TestHookChangeGofmt(t *testing.T) {
	// During git change, we run the gofmt check before invoking commit,
	// so we should not see the line about 'git commit' failing.
	// That is, the failure should come from git change, not from
	// the commit hook.
	gt := newGitTest(t)
	defer gt.done()
	gt.work(t)

	// Write out a non-Go file.
	write(t, gt.client+"/bad.go", badGo)
	trun(t, gt.client, "git", "add", ".")

	t.Logf("invoking commit hook explicitly")
	testMainDied(t, "hook-invoke", "pre-commit")
	testPrintedStderr(t, "gofmt needs to format these files (run 'git gofmt'):", "bad.go")

	t.Logf("change without hook installed")
	testCommitMsg = "foo: msg"
	testMainDied(t, "change")
	testPrintedStderr(t, "gofmt needs to format these files (run 'git gofmt'):", "bad.go", "!running: git")

	t.Logf("change with hook installed")
	restore := testInstallHook(t, gt)
	defer restore()
	testCommitMsg = "foo: msg"
	testMainDied(t, "change")
	testPrintedStderr(t, "gofmt needs to format these files (run 'git gofmt'):", "bad.go", "!running: git")
}

func TestHookPreCommitDetachedHead(t *testing.T) {
	// If we're in detached head mode, something special is going on,
	// like git rebase. We disable the gofmt-checking precommit hook,
	// since we expect it would just get in the way at that point.
	// (It also used to crash.)

	gt := newGitTest(t)
	defer gt.done()
	gt.work(t)

	write(t, gt.client+"/bad.go", badGo)
	trun(t, gt.client, "git", "add", ".")
	trun(t, gt.client, "git", "checkout", "HEAD^0")

	testMain(t, "hook-invoke", "pre-commit")
	testNoStdout(t)
	testNoStderr(t)
}

func TestHookPreCommitEnv(t *testing.T) {
	// If $GIT_GOFMT_HOOK == "off", gofmt hook should not complain.

	gt := newGitTest(t)
	defer gt.done()
	gt.work(t)

	write(t, gt.client+"/bad.go", badGo)
	trun(t, gt.client, "git", "add", ".")
	os.Setenv("GIT_GOFMT_HOOK", "off")
	defer os.Unsetenv("GIT_GOFMT_HOOK")

	testMain(t, "hook-invoke", "pre-commit")
	testNoStdout(t)
	testPrintedStderr(t, "git-gofmt-hook disabled by $GIT_GOFMT_HOOK=off")
}

func TestHookPreCommitUnstaged(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()
	gt.work(t)

	write(t, gt.client+"/bad.go", badGo)
	write(t, gt.client+"/good.go", goodGo)

	// The pre-commit hook is being asked about files in the index.
	// Make sure it is not looking at files in the working tree (current directory) instead.
	// There are three possible kinds of file: good, bad (misformatted), and broken (syntax error).
	// There are also three possible places files live: the most recent commit, the index,
	// and the working tree. We write a sequence of files that cover all possible
	// combination of kinds of file in the various places. For example,
	// good-bad-broken.go is a good file in the most recent commit,
	// a bad file in the index, and a broken file in the working tree.
	// After creating these files, we check that the gofmt hook reports
	// about the index only.

	const N = 3
	name := []string{"good", "bad", "broken"}
	content := []string{goodGo, badGo, brokenGo}
	var wantErr []string
	var allFiles []string
	writeFiles := func(n int) {
		allFiles = nil
		wantErr = nil
		for i := 0; i < N*N*N; i++ {
			// determine n'th digit of 3-digit base-N value i
			j := i
			for k := 0; k < (3 - 1 - n); k++ {
				j /= N
			}
			file := fmt.Sprintf("%s-%s-%s.go", name[i/N/N], name[(i/N)%N], name[i%N])
			allFiles = append(allFiles, file)
			write(t, gt.client+"/"+file, content[j%N])

			switch {
			case strings.Contains(file, "-bad-"):
				wantErr = append(wantErr, "\t"+file+"\n")
			case strings.Contains(file, "-broken-"):
				wantErr = append(wantErr, "\t"+file+":")
			default:
				wantErr = append(wantErr, "!"+file)
			}
		}
	}

	// committed files
	writeFiles(0)
	trun(t, gt.client, "git", "add", ".")
	trun(t, gt.client, "git", "commit", "-m", "msg")

	// staged files
	writeFiles(1)
	trun(t, gt.client, "git", "add", ".")

	// unstaged files
	writeFiles(2)

	wantErr = append(wantErr, "gofmt reported errors", "gofmt needs to format these files")

	testMainDied(t, "hook-invoke", "pre-commit")
	testPrintedStderr(t, wantErr...)
}

func TestHooks(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()

	gt.removeStubHooks()
	testMain(t, "hooks") // install hooks

	data, err := ioutil.ReadFile(gt.client + "/.git/hooks/commit-msg")
	if err != nil {
		t.Fatalf("hooks did not write commit-msg hook: %v", err)
	}
	if string(data) != "#!/bin/sh\nexec git-codereview hook-invoke commit-msg \"$@\"\n" {
		t.Fatalf("invalid commit-msg hook:\n%s", string(data))
	}
}

func TestHooksOverwriteOldCommitMsg(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()

	write(t, gt.client+"/.git/hooks/commit-msg", oldCommitMsgHook)
	testMain(t, "hooks") // install hooks
	data, err := ioutil.ReadFile(gt.client + "/.git/hooks/commit-msg")
	if err != nil {
		t.Fatalf("hooks did not write commit-msg hook: %v", err)
	}
	if string(data) == oldCommitMsgHook {
		t.Fatalf("hooks left old commit-msg hook in place")
	}
	if string(data) != "#!/bin/sh\nexec git-codereview hook-invoke commit-msg \"$@\"\n" {
		t.Fatalf("invalid commit-msg hook:\n%s", string(data))
	}
}

func testInstallHook(t *testing.T, gt *gitTest) (restore func()) {
	trun(t, gt.pwd, "go", "build", "-o", gt.client+"/git-codereview")
	path := os.Getenv("PATH")
	os.Setenv("PATH", gt.client+string(filepath.ListSeparator)+path)
	gt.removeStubHooks()
	testMain(t, "hooks") // install hooks

	return func() {
		os.Setenv("PATH", path)
	}
}

func TestHookCommitMsgFromGit(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()

	restore := testInstallHook(t, gt)
	defer restore()

	testMain(t, "change", "mybranch")
	write(t, gt.client+"/file", "more data")
	trun(t, gt.client, "git", "add", "file")
	trun(t, gt.client, "git", "commit", "-m", "mymsg")

	log := trun(t, gt.client, "git", "log", "-n", "1")
	if !strings.Contains(log, "mymsg") {
		t.Fatalf("did not find mymsg in git log output:\n%s", log)
	}
	// The 4 spaces are because git indents the commit message proper.
	if !strings.Contains(log, "\n    \n    Change-Id:") {
		t.Fatalf("did not find Change-Id in git log output:\n%s", log)
	}
}
