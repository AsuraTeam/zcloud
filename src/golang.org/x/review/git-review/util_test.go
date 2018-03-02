// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type gitTest struct {
	pwd    string // current directory before test
	tmpdir string // temporary directory holding repos
	server string // server repo root
	client string // client repo root
}

func (gt *gitTest) done() {
	os.RemoveAll(gt.tmpdir)
	os.Chdir(gt.pwd)
}

func newGitTest(t *testing.T) *gitTest {
	tmpdir, err := ioutil.TempDir("", "git-review-test")
	if err != nil {
		t.Fatal(err)
	}

	server := tmpdir + "/git-origin"

	mkdir(t, server)
	write(t, server+"/file", "this is master")
	trun(t, server, "git", "init", ".")
	trun(t, server, "git", "add", "file")
	trun(t, server, "git", "commit", "-m", "on master")

	for _, name := range []string{"dev.branch", "release.branch"} {
		trun(t, server, "git", "checkout", "master")
		trun(t, server, "git", "branch", name)
		write(t, server+"/file", "this is "+name)
		trun(t, server, "git", "commit", "-a", "-m", "on "+name)
	}

	client := tmpdir + "/git-client"
	mkdir(t, client)
	trun(t, client, "git", "clone", server, ".")
	trun(t, client, "git", "config", "core.editor", "false")
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(client); err != nil {
		t.Fatal(err)
	}

	gt := &gitTest{
		pwd:    pwd,
		tmpdir: tmpdir,
		server: server,
		client: client,
	}

	return gt
}

func mkdir(t *testing.T, dir string) {
	if err := os.Mkdir(dir, 0777); err != nil {
		t.Fatal(err)
	}
}

func write(t *testing.T, file, data string) {
	if err := ioutil.WriteFile(file, []byte(data), 0666); err != nil {
		t.Fatal(err)
	}
}

func trun(t *testing.T, dir string, cmdline ...string) {
	cmd := exec.Command(cmdline[0], cmdline[1:]...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("in %s/, ran %s: %v\n%s", filepath.Base(dir), cmdline, err, out)
	}
}

func testMain(t *testing.T, args ...string) {
	t.Logf("git-review %s", strings.Join(args, " "))
	runLog = []string{} // non-nil, to trigger saving of commands

	defer func() {
		if err := recover(); err != nil {
			runLog = nil
			dieTrap = nil
			t.Fatalf("panic: %v", err)
		}
	}()

	dieTrap = func() {
		panic("died")
	}

	os.Args = append([]string{"git-review"}, args...)
	main()

	dieTrap = nil
}

func testRan(t *testing.T, cmds ...string) {
	if !reflect.DeepEqual(runLog, cmds) {
		t.Errorf("ran:\n%s", strings.Join(runLog, "\n"))
		t.Errorf("wanted:\n%s", strings.Join(cmds, "\n"))
	}
}
