// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Branch describes a Git branch.
type Branch struct {
	Name     string // branch name
	changeID string // Change-Id of pending commit ("" if nothing pending)
	subject  string // first line of pending commit ("" if nothing pending)
}

// CurrentBranch returns the current branch.
func CurrentBranch() *Branch {
	name := getOutput("git", "rev-parse", "--abbrev-ref", "HEAD")
	return &Branch{Name: name}
}

func (b *Branch) OriginBranch() string {
	argv := []string{"git", "rev-parse", "--abbrev-ref", "@{u}"}
	out, err := exec.Command(argv[0], argv[1:]...).CombinedOutput()
	if err == nil && len(out) > 0 {
		return string(bytes.TrimSpace(out))
	}
	if strings.Contains(string(out), "No upstream configured") {
		// Assume branch was created before we set upstream correctly.
		return "origin/master"
	}
	fmt.Fprintf(os.Stderr, "%v\n%s\n", commandString(argv[0], argv[1:]), out)
	dief("%v", err)
	panic("not reached")
}

func (b *Branch) IsLocalOnly() bool {
	return "origin/"+b.Name != b.OriginBranch()
}

func (b *Branch) HasPendingCommit() bool {
	head := getOutput("git", "rev-parse", b.Name)
	base := getOutput("git", "merge-base", b.OriginBranch(), b.Name)
	return base != head
}

func (b *Branch) ChangeID() string {
	if b.changeID == "" {
		if b.HasPendingCommit() {
			b.changeID = headChangeID(b.Name)
			b.subject = commitSubject(b.Name)
		}
	}
	return b.changeID
}

func (b *Branch) Subject() string {
	b.ChangeID() // page in subject
	return b.subject
}

func commitSubject(ref string) string {
	const f = "--format=format:%s"
	return getOutput("git", "log", "-n", "1", f, ref, "--")
}

func headChangeID(branch string) string {
	const (
		p = "Change-Id: "
		f = "--format=format:%b"
	)
	for _, s := range getLines("git", "log", "-n", "1", f, branch, "--") {
		if strings.HasPrefix(s, p) {
			return strings.TrimSpace(strings.TrimPrefix(s, p))
		}
	}
	dief("no Change-Id line found in HEAD commit on branch %s.", branch)
	panic("unreachable")
}

// Submitted reports whether some form of b's pending commit
// has been cherry picked to origin.
func (b *Branch) Submitted(id string) bool {
	return len(getOutput("git", "log", "--grep", "Change-Id: "+id, b.OriginBranch())) > 0
}

var stagedRe = regexp.MustCompile(`^[ACDMR]  `)

func HasStagedChanges() bool {
	for _, s := range getLines("git", "status", "-b", "--porcelain") {
		if stagedRe.MatchString(s) {
			return true
		}
	}
	return false
}

func LocalBranches() []*Branch {
	var branches []*Branch
	for _, s := range getLines("git", "branch", "-q") {
		branches = append(branches, &Branch{Name: strings.TrimPrefix(s, "* ")})
	}
	return branches
}

func OriginBranches() []string {
	var branches []string
	for _, line := range getLines("git", "branch", "-a", "-q") {
		if i := strings.Index(line, " -> "); i >= 0 {
			line = line[:i]
		}
		name := strings.TrimSpace(strings.TrimPrefix(line, "* "))
		if strings.HasPrefix(name, "remotes/origin/") {
			branches = append(branches, strings.TrimPrefix(name, "remotes/"))
		}
	}
	return branches
}
