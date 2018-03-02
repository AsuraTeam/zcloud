// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const (
	goodGo      = "package good\n"
	badGo       = " package bad1 "
	badGoFixed  = "package bad1\n"
	bad2Go      = " package bad2 "
	bad2GoFixed = "package bad2\n"
	brokenGo    = "package B R O K E N"
)

func TestGofmt(t *testing.T) {
	// Test of basic operations.
	gt := newGitTest(t)
	defer gt.done()

	gt.work(t)

	if err := os.MkdirAll(gt.client+"/test/bench", 0755); err != nil {
		t.Fatal(err)
	}
	write(t, gt.client+"/bad.go", badGo)
	write(t, gt.client+"/good.go", goodGo)
	write(t, gt.client+"/test/bad.go", badGo)
	write(t, gt.client+"/test/good.go", goodGo)
	write(t, gt.client+"/test/bench/bad.go", badGo)
	write(t, gt.client+"/test/bench/good.go", goodGo)
	trun(t, gt.client, "git", "add", ".") // make files tracked

	testMain(t, "gofmt", "-l")
	testPrintedStdout(t, "bad.go\n", "!good.go", fromSlash("!test/bad"), fromSlash("test/bench/bad.go"))

	testMain(t, "gofmt", "-l")
	testPrintedStdout(t, "bad.go\n", "!good.go", fromSlash("!test/bad"), fromSlash("test/bench/bad.go"))

	testMain(t, "gofmt")
	testNoStdout(t)

	testMain(t, "gofmt", "-l")
	testNoStdout(t)

	write(t, gt.client+"/bad.go", badGo)
	write(t, gt.client+"/broken.go", brokenGo)
	trun(t, gt.client, "git", "add", ".")
	testMainDied(t, "gofmt", "-l")
	testPrintedStdout(t, "bad.go")
	testPrintedStderr(t, "gofmt reported errors", "broken.go")
}

func TestGofmtSubdir(t *testing.T) {
	// Check that gofmt prints relative paths for files in or below the current directory.
	gt := newGitTest(t)
	defer gt.done()

	gt.work(t)

	mkdir(t, gt.client+"/dir1")
	mkdir(t, gt.client+"/longnamedir2")
	write(t, gt.client+"/dir1/bad1.go", badGo)
	write(t, gt.client+"/longnamedir2/bad2.go", badGo)
	trun(t, gt.client, "git", "add", ".") // make files tracked

	chdir(t, gt.client)
	testMain(t, "gofmt", "-l")
	testPrintedStdout(t, fromSlash("dir1/bad1.go"), fromSlash("longnamedir2/bad2.go"))

	chdir(t, gt.client+"/dir1")
	testMain(t, "gofmt", "-l")
	testPrintedStdout(t, "bad1.go", fromSlash("!/bad1.go"), fromSlash("longnamedir2/bad2.go"))

	chdir(t, gt.client+"/longnamedir2")
	testMain(t, "gofmt", "-l")
	testPrintedStdout(t, "bad2.go", fromSlash("!/bad2.go"), fromSlash("dir1/bad1.go"))

	mkdir(t, gt.client+"/z")
	chdir(t, gt.client+"/z")
	testMain(t, "gofmt", "-l")
	testPrintedStdout(t, fromSlash("longnamedir2/bad2.go"), fromSlash("dir1/bad1.go"))
}

func TestGofmtSubdirIndexCheckout(t *testing.T) {
	// Like TestGofmtSubdir but bad Go files are only in index, not working copy.
	// Check also that prints a correct path (relative or absolute) for files outside the
	// current directory, even when running with Git before 2.3.0 which doesn't
	// handle those right in git checkout-index --temp.

	gt := newGitTest(t)
	defer gt.done()

	gt.work(t)

	mkdir(t, gt.client+"/dir1")
	mkdir(t, gt.client+"/longnamedir2")
	write(t, gt.client+"/dir1/bad1.go", badGo)
	write(t, gt.client+"/longnamedir2/bad2.go", badGo)
	trun(t, gt.client, "git", "add", ".") // put files in index
	write(t, gt.client+"/dir1/bad1.go", goodGo)
	write(t, gt.client+"/longnamedir2/bad2.go", goodGo)

	chdir(t, gt.client)
	testMain(t, "gofmt", "-l")
	testPrintedStdout(t, fromSlash("dir1/bad1.go (staged)"), fromSlash("longnamedir2/bad2.go (staged)"))

	chdir(t, gt.client+"/dir1")
	testMain(t, "gofmt", "-l")
	testPrintedStdout(t, "bad1.go (staged)", fromSlash("!/bad1.go"), fromSlash("longnamedir2/bad2.go (staged)"))

	chdir(t, gt.client+"/longnamedir2")
	testMain(t, "gofmt", "-l")
	testPrintedStdout(t, "bad2.go (staged)", fromSlash("!/bad2.go"), fromSlash("dir1/bad1.go (staged)"))

	mkdir(t, gt.client+"/z")
	chdir(t, gt.client+"/z")
	testMain(t, "gofmt", "-l")
	testPrintedStdout(t, fromSlash("longnamedir2/bad2.go (staged)"), fromSlash("dir1/bad1.go (staged)"))
}

func TestGofmtUnstaged(t *testing.T) {
	// Test when unstaged files are different from staged ones.
	// See TestHookPreCommitUnstaged for an explanation.
	// In this test we use two different kinds of bad files, so that
	// we can test having a bad file in the index and a different
	// bad file in the working directory.

	gt := newGitTest(t)
	defer gt.done()
	gt.work(t)

	name := []string{"good", "bad", "bad2", "broken"}
	orig := []string{goodGo, badGo, bad2Go, brokenGo}
	fixed := []string{goodGo, badGoFixed, bad2GoFixed, brokenGo}
	const N = 4

	var allFiles, wantOut, wantErr []string
	writeFiles := func(n int) {
		allFiles = nil
		wantOut = nil
		wantErr = nil
		for i := 0; i < N*N*N; i++ {
			// determine n'th digit of 3-digit base-N value i
			j := i
			for k := 0; k < (3 - 1 - n); k++ {
				j /= N
			}
			text := orig[j%N]
			file := fmt.Sprintf("%s-%s-%s.go", name[i/N/N], name[(i/N)%N], name[i%N])
			allFiles = append(allFiles, file)
			write(t, gt.client+"/"+file, text)

			if (i/N)%N != i%N {
				staged := file + " (staged)"
				switch {
				case strings.Contains(file, "-bad-"), strings.Contains(file, "-bad2-"):
					wantOut = append(wantOut, staged)
					wantErr = append(wantErr, "!"+staged)
				case strings.Contains(file, "-broken-"):
					wantOut = append(wantOut, "!"+staged)
					wantErr = append(wantErr, staged)
				default:
					wantOut = append(wantOut, "!"+staged)
					wantErr = append(wantErr, "!"+staged)
				}
			}
			switch {
			case strings.Contains(file, "-bad.go"), strings.Contains(file, "-bad2.go"):
				if (i/N)%N != i%N {
					file += " (unstaged)"
				}
				wantOut = append(wantOut, file+"\n")
				wantErr = append(wantErr, "!"+file+":", "!"+file+" (unstaged)")
			case strings.Contains(file, "-broken.go"):
				wantOut = append(wantOut, "!"+file+"\n", "!"+file+" (unstaged)")
				wantErr = append(wantErr, file+":")
			default:
				wantOut = append(wantOut, "!"+file+"\n", "!"+file+":", "!"+file+" (unstaged)")
				wantErr = append(wantErr, "!"+file+"\n", "!"+file+":", "!"+file+" (unstaged)")
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

	// Check that gofmt -l shows the right output and errors.
	testMainDied(t, "gofmt", "-l")
	testPrintedStdout(t, wantOut...)
	testPrintedStderr(t, wantErr...)

	// Again (last command should not have written anything).
	testMainDied(t, "gofmt", "-l")
	testPrintedStdout(t, wantOut...)
	testPrintedStderr(t, wantErr...)

	// Reformat in place.
	testMainDied(t, "gofmt")
	testNoStdout(t)
	testPrintedStderr(t, wantErr...)

	// Read files to make sure unstaged did not bleed into staged.
	for i, file := range allFiles {
		if data, err := ioutil.ReadFile(gt.client + "/" + file); err != nil {
			t.Errorf("%v", err)
		} else if want := fixed[i%N]; string(data) != want {
			t.Errorf("%s: working tree = %q, want %q", file, string(data), want)
		}
		if data, want := trun(t, gt.client, "git", "show", ":"+file), fixed[i/N%N]; data != want {
			t.Errorf("%s: index = %q, want %q", file, data, want)
		}
		if data, want := trun(t, gt.client, "git", "show", "HEAD:"+file), orig[i/N/N]; data != want {
			t.Errorf("%s: commit = %q, want %q", file, data, want)
		}
	}

	// Check that gofmt -l still shows the errors.
	testMainDied(t, "gofmt", "-l")
	testNoStdout(t)
	testPrintedStderr(t, wantErr...)
}

func TestGofmtAmbiguousRevision(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()

	t.Logf("creating file that conflicts with revision parameter")
	write(t, gt.client+"/HEAD", "foo")

	testMain(t, "gofmt")
}

func TestGofmtFastForwardMerge(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()

	// merge dev.branch into master
	write(t, gt.server+"/file", "more work")
	trun(t, gt.server, "git", "commit", "-m", "work", "file")
	trun(t, gt.server, "git", "merge", "-m", "merge", "dev.branch")

	// add bad go file on master
	write(t, gt.server+"/bad.go", "package {\n")
	trun(t, gt.server, "git", "add", "bad.go")
	trun(t, gt.server, "git", "commit", "-m", "bad go")

	// update client
	trun(t, gt.client, "git", "checkout", "master")
	trun(t, gt.client, "git", "pull")
	testMain(t, "change", "dev.branch")
	trun(t, gt.client, "git", "pull")

	// merge master into dev.branch, fast forward merge
	trun(t, gt.client, "git", "merge", "--ff-only", "master")

	// verify that now client is in a state where just the tag is changing; there's no new commit.
	masterHash := strings.TrimSpace(trun(t, gt.server, "git", "rev-parse", "master"))
	devHash := strings.TrimSpace(trun(t, gt.client, "git", "rev-parse", "HEAD"))

	if masterHash != devHash {
		t.Logf("branches:\n%s", trun(t, gt.client, "git", "branch", "-a", "-v"))
		t.Logf("log:\n%s", trun(t, gt.client, "git", "log", "--graph", "--decorate"))
		t.Fatalf("setup wrong - got different commit hashes on master and dev branch")
	}

	// check that gofmt finds nothing to do, ignoring the bad (but committed) file1.go.
	testMain(t, "gofmt")
	testNoStdout(t)
	testNoStderr(t)
}
