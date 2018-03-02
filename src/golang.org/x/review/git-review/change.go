// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var changeQuick bool

func change(args []string) {
	flags.BoolVar(&changeQuick, "q", false, "do not edit commit msg when updating commit")
	flags.Parse(args)
	if len(flags.Args()) > 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s change %s [branch]\n", os.Args[0], globalFlags)
		os.Exit(2)

	}

	// Checkout or create branch, if specified.
	target := flags.Arg(0)
	if target != "" {
		checkoutOrCreate(target)
	}

	// Create or amend change commit.
	b := CurrentBranch()
	if !b.IsLocalOnly() {
		if target != "" {
			// Permit "review change master".
			return
		}
		dief("can't commit to %s branch (use '%s change branchname').", b.Name, os.Args[0])
	}

	if b.ChangeID() == "" {
		// No change commit on this branch, create one.
		commitChanges(false)
		return
	}
	if target != "" {
		// If we switched to an existing branch, don't amend the
		// commit. (The user can run 'review change' to do that.)
		return
	}
	// Amend the commit.
	commitChanges(true)
}

var testCommitMsg string

func commitChanges(amend bool) {
	if !HasStagedChanges() {
		printf("no staged changes. Did you forget to 'git add'?")
	}
	commit := func(amend bool) {
		args := []string{"commit", "-q", "--allow-empty"}
		if amend {
			args = append(args, "--amend")
			if changeQuick {
				args = append(args, "--no-edit")
			}
		}
		if testCommitMsg != "" {
			args = append(args, "-m", testCommitMsg)
		}
		run("git", args...)
	}
	commit(amend)
	for !commitMessageOK() {
		fmt.Print("re-edit commit message (y/n)? ")
		if !scanYes() {
			break
		}
		commit(true)
	}
	printf("change updated.")
}

func checkoutOrCreate(target string) {
	// If local branch exists, check it out.
	for _, b := range LocalBranches() {
		if b.Name == target {
			run("git", "checkout", "-q", target)
			printf("changed to branch %v.", target)
			return
		}
	}

	// If origin branch exists, create local branch tracking it.
	for _, name := range OriginBranches() {
		if name == "origin/"+target {
			run("git", "checkout", "-q", "-t", "-b", target, name)
			printf("created branch %v tracking %s.", target, name)
			return
		}
	}

	// Otherwise, this is a request to create a local work branch.
	// Check for reserved names. We take everything with a dot.
	if strings.Contains(target, ".") {
		dief("invalid branch name %v: branch names with dots are reserved for git-review.", target)
	}

	// If the current branch has a pending commit, building
	// on top of it will not help. Don't allow that.
	// Otherwise, inherit HEAD and upstream from the current branch.
	b := CurrentBranch()
	if b.HasPendingCommit() {
		dief("cannot branch from work branch; change back to %v first.", strings.TrimPrefix(b.OriginBranch(), "origin/"))
	}

	origin := b.OriginBranch()

	// NOTE: This is different from git checkout -q -t -b branch. It does not move HEAD.
	run("git", "checkout", "-q", "-b", target)
	run("git", "branch", "-q", "--set-upstream-to", origin)
	printf("created branch %v tracking %s.", target, origin)
}

var (
	messageRE  = regexp.MustCompile(`^(\[[a-zA-Z0-9.-]+\] )?[a-zA-Z0-9-/, ]+: `)
	oldFixesRE = regexp.MustCompile(`Fixes +(issue +#?)?([0-9]+)`)
)

func commitMessageOK() bool {
	body := getOutput("git", "log", "--format=format:%B", "-n", "1")
	ok := true
	if !messageRE.MatchString(body) {
		fmt.Print(commitMessageWarning)
		ok = false
	}
	if m := oldFixesRE.FindStringSubmatch(body); m != nil {
		fmt.Printf(fixesIssueWarning, m[0], m[2])
		ok = false
	}
	return ok
}

const commitMessageWarning = `
Your CL description appears not to use the standard form.

The first line of your change description is conventionally a one-line summary
of the change, prefixed by the primary affected package, and is used as the
subject for code review mail; the rest of the description elaborates.

Examples:

	encoding/rot13: new package

	math: add IsInf, IsNaN

	net: fix cname in LookupHost

	unicode: update to Unicode 5.0.2

`

const fixesIssueWarning = `
Your CL description contains the string %q, which is
the old Google Code way of linking commits to issues.

You should rewrite it to use the GitHub convention: "Fixes #%v".

`

func scanYes() bool {
	var s string
	fmt.Scan(&s)
	return strings.HasPrefix(strings.ToLower(s), "y")
}
