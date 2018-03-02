// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Branch describes a Git branch.
type Branch struct {
	Name          string    // branch name
	loadedPending bool      // following fields are valid
	originBranch  string    // upstream origin branch
	commitsAhead  int       // number of commits ahead of origin branch
	commitsBehind int       // number of commits behind origin branch
	branchpoint   string    // latest commit hash shared with origin branch
	pending       []*Commit // pending commits, newest first (children before parents)
}

// A Commit describes a single pending commit on a Git branch.
type Commit struct {
	Hash      string // commit hash
	ShortHash string // abbreviated commit hash
	Parent    string // parent hash
	Merge     string // for merges, hash of commit being merged into Parent
	Message   string // commit message
	Subject   string // first line of commit message
	ChangeID  string // Change-Id in commit message ("" if missing)

	// For use by pending command.
	g         *GerritChange // associated Gerrit change data
	gerr      error         // error loading Gerrit data
	committed []string      // list of files in this commit
}

// CurrentBranch returns the current branch.
func CurrentBranch() *Branch {
	name := strings.TrimPrefix(trim(cmdOutput("git", "rev-parse", "--abbrev-ref", "HEAD")), "heads/")
	return &Branch{Name: name}
}

// DetachedHead reports whether branch b corresponds to a detached HEAD
// (does not have a real branch name).
func (b *Branch) DetachedHead() bool {
	return b.Name == "HEAD"
}

// OriginBranch returns the name of the origin branch that branch b tracks.
// The returned name is like "origin/master" or "origin/dev.garbage" or
// "origin/release-branch.go1.4".
func (b *Branch) OriginBranch() string {
	if b.DetachedHead() {
		// Detached head mode.
		// "origin/HEAD" is clearly false, but it should be easy to find when it
		// appears in other commands. Really any caller of OriginBranch
		// should check for detached head mode.
		return "origin/HEAD"
	}

	if b.originBranch != "" {
		return b.originBranch
	}
	argv := []string{"git", "rev-parse", "--abbrev-ref", b.Name + "@{u}"}
	out, err := exec.Command(argv[0], argv[1:]...).CombinedOutput()
	if err == nil && len(out) > 0 {
		b.originBranch = string(bytes.TrimSpace(out))
		return b.originBranch
	}

	// Have seen both "No upstream configured" and "no upstream configured".
	if strings.Contains(string(out), "upstream configured") {
		// Assume branch was created before we set upstream correctly.
		b.originBranch = "origin/master"
		return b.originBranch
	}
	fmt.Fprintf(stderr(), "%v\n%s\n", commandString(argv[0], argv[1:]), out)
	dief("%v", err)
	panic("not reached")
}

func (b *Branch) FullName() string {
	if b.Name != "HEAD" {
		return "refs/heads/" + b.Name
	}
	return b.Name
}

// IsLocalOnly reports whether b is a local work branch (only local, not known to remote server).
func (b *Branch) IsLocalOnly() bool {
	return "origin/"+b.Name != b.OriginBranch()
}

// HasPendingCommit reports whether b has any pending commits.
func (b *Branch) HasPendingCommit() bool {
	b.loadPending()
	return b.commitsAhead > 0
}

// Pending returns b's pending commits, newest first (children before parents).
func (b *Branch) Pending() []*Commit {
	b.loadPending()
	return b.pending
}

// Branchpoint returns an identifier for the latest revision
// common to both this branch and its upstream branch.
func (b *Branch) Branchpoint() string {
	b.loadPending()
	return b.branchpoint
}

func (b *Branch) loadPending() {
	if b.loadedPending {
		return
	}
	b.loadedPending = true

	// In case of early return.
	b.branchpoint = trim(cmdOutput("git", "rev-parse", "HEAD"))

	if b.DetachedHead() {
		return
	}

	// Note: --topo-order means child first, then parent.
	origin := b.OriginBranch()
	const numField = 5
	all := trim(cmdOutput("git", "log", "--topo-order", "--format=format:%H%x00%h%x00%P%x00%B%x00%s%x00", origin+".."+b.FullName(), "--"))
	fields := strings.Split(all, "\x00")
	if len(fields) < numField {
		return // nothing pending
	}
	for i, field := range fields {
		fields[i] = strings.TrimLeft(field, "\r\n")
	}
	foundMergeBranchpoint := false
	for i := 0; i+numField <= len(fields); i += numField {
		c := &Commit{
			Hash:      fields[i],
			ShortHash: fields[i+1],
			Parent:    strings.TrimSpace(fields[i+2]), // %P starts with \n for some reason
			Message:   fields[i+3],
			Subject:   fields[i+4],
		}
		if j := strings.Index(c.Parent, " "); j >= 0 {
			c.Parent, c.Merge = c.Parent[:j], c.Parent[j+1:]
			// Found merge point.
			// Merges break the invariant that the last shared commit (the branchpoint)
			// is the parent of the final commit in the log output.
			// If c.Parent is on the origin branch, then since we are reading the log
			// in (reverse) topological order, we know that c.Parent is the actual branchpoint,
			// even if we later see additional commits on a different branch leading down to
			// a lower location on the same origin branch.
			// Check c.Merge (the second parent) too, so we don't depend on the parent order.
			if strings.Contains(cmdOutput("git", "branch", "-a", "--contains", c.Parent), " "+origin+"\n") {
				foundMergeBranchpoint = true
				b.branchpoint = c.Parent
			}
			if strings.Contains(cmdOutput("git", "branch", "-a", "--contains", c.Merge), " "+origin+"\n") {
				foundMergeBranchpoint = true
				b.branchpoint = c.Merge
			}
		}
		for _, line := range lines(c.Message) {
			// Note: Keep going even if we find one, so that
			// we take the last Change-Id line, just in case
			// there is a commit message quoting another
			// commit message.
			// I'm not sure this can come up at all, but just in case.
			if strings.HasPrefix(line, "Change-Id: ") {
				c.ChangeID = line[len("Change-Id: "):]
			}
		}

		b.pending = append(b.pending, c)
		if !foundMergeBranchpoint {
			b.branchpoint = c.Parent
		}
	}
	b.commitsAhead = len(b.pending)
	b.commitsBehind = len(lines(cmdOutput("git", "log", "--format=format:x", b.FullName()+".."+b.OriginBranch(), "--")))
}

// Submitted reports whether some form of b's pending commit
// has been cherry picked to origin.
func (b *Branch) Submitted(id string) bool {
	if id == "" {
		return false
	}
	line := "Change-Id: " + id
	out := cmdOutput("git", "log", "-n", "1", "-F", "--grep", line, b.Name+".."+b.OriginBranch(), "--")
	return strings.Contains(out, line)
}

var stagedRE = regexp.MustCompile(`^[ACDMR]  `)

// HasStagedChanges reports whether the working directory contains staged changes.
func HasStagedChanges() bool {
	for _, s := range nonBlankLines(cmdOutput("git", "status", "-b", "--porcelain")) {
		if stagedRE.MatchString(s) {
			return true
		}
	}
	return false
}

var unstagedRE = regexp.MustCompile(`^.[ACDMR]`)

// HasUnstagedChanges reports whether the working directory contains unstaged changes.
func HasUnstagedChanges() bool {
	for _, s := range nonBlankLines(cmdOutput("git", "status", "-b", "--porcelain")) {
		if unstagedRE.MatchString(s) {
			return true
		}
	}
	return false
}

// LocalChanges returns a list of files containing staged, unstaged, and untracked changes.
// The elements of the returned slices are typically file names, always relative to the root,
// but there are a few alternate forms. First, for renaming or copying, the element takes
// the form `from -> to`. Second, in the case of files with names that contain unusual characters,
// the files (or the from, to fields of a rename or copy) are quoted C strings.
// For now, we expect the caller only shows these to the user, so these exceptions are okay.
func LocalChanges() (staged, unstaged, untracked []string) {
	for _, s := range lines(cmdOutput("git", "status", "-b", "--porcelain")) {
		if len(s) < 4 || s[2] != ' ' {
			continue
		}
		switch s[0] {
		case 'A', 'C', 'D', 'M', 'R':
			staged = append(staged, s[3:])
		case '?':
			untracked = append(untracked, s[3:])
		}
		switch s[1] {
		case 'A', 'C', 'D', 'M', 'R':
			unstaged = append(unstaged, s[3:])
		}
	}
	return
}

// LocalBranches returns a list of all known local branches.
// If the current directory is in detached HEAD mode, one returned
// branch will have Name == "HEAD" and DetachedHead() == true.
func LocalBranches() []*Branch {
	var branches []*Branch
	current := CurrentBranch()
	for _, s := range nonBlankLines(cmdOutput("git", "branch", "-q")) {
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "* ") {
			// * marks current branch in output.
			// Normally the current branch has a name like any other,
			// but in detached HEAD mode the branch listing shows
			// a localized (translated) textual description instead of
			// a branch name. Avoid language-specific differences
			// by using CurrentBranch().Name for the current branch.
			// It detects detached HEAD mode in a more portable way.
			// (git rev-parse --abbrev-ref HEAD returns 'HEAD').
			s = current.Name
		}
		branches = append(branches, &Branch{Name: s})
	}
	return branches
}

func OriginBranches() []string {
	var branches []string
	for _, line := range nonBlankLines(cmdOutput("git", "branch", "-a", "-q")) {
		line = strings.TrimSpace(line)
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

// GerritChange returns the change metadata from the Gerrit server
// for the branch's pending change.
// The extra strings are passed to the Gerrit API request as o= parameters,
// to enable additional information. Typical values include "LABELS" and "CURRENT_REVISION".
// See https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html for details.
func (b *Branch) GerritChange(c *Commit, extra ...string) (*GerritChange, error) {
	if !b.HasPendingCommit() {
		return nil, fmt.Errorf("no changes pending")
	}
	id := fullChangeID(b, c)
	for i, x := range extra {
		if i == 0 {
			id += "?"
		} else {
			id += "&"
		}
		id += "o=" + x
	}
	return readGerritChange(id)
}

const minHashLen = 4 // git minimum hash length accepted on command line

// CommitByHash finds a unique pending commit by its hash prefix.
// It dies if the hash cannot be resolved to a pending commit,
// using the action ("mail", "submit") in the failure message.
func (b *Branch) CommitByHash(action, hash string) *Commit {
	if len(hash) < minHashLen {
		dief("cannot %s: commit hash %q must be at least %d digits long", action, hash, minHashLen)
	}
	var c *Commit
	for _, c1 := range b.Pending() {
		if strings.HasPrefix(c1.Hash, hash) {
			if c != nil {
				dief("cannot %s: commit hash %q is ambiguous in the current branch", action, hash)
			}
			c = c1
		}
	}
	if c == nil {
		dief("cannot %s: commit hash %q not found in the current branch", action, hash)
	}
	return c
}

// DefaultCommit returns the default pending commit for this branch.
// It dies if there is not exactly one pending commit,
// using the action ("mail", "submit") in the failure message.
func (b *Branch) DefaultCommit(action string) *Commit {
	work := b.Pending()
	if len(work) == 0 {
		dief("cannot %s: no changes pending", action)
	}
	if len(work) >= 2 {
		var buf bytes.Buffer
		for _, c := range work {
			fmt.Fprintf(&buf, "\n\t%s %s", c.ShortHash, c.Subject)
		}
		extra := ""
		if action == "submit" {
			extra = " or use submit -i"
		}
		dief("cannot %s: multiple changes pending; must specify commit hash on command line%s:%s", action, extra, buf.String())
	}
	return work[0]
}

func cmdBranchpoint(args []string) {
	expectZeroArgs(args, "sync")
	fmt.Fprintf(stdout(), "%s\n", CurrentBranch().Branchpoint())
}

func cmdRebaseWork(args []string) {
	expectZeroArgs(args, "rebase-work")
	b := CurrentBranch()
	if HasStagedChanges() || HasUnstagedChanges() {
		dief("cannot rebase with uncommitted work")
	}
	if len(b.Pending()) == 0 {
		dief("no pending work")
	}
	run("git", "rebase", "-i", b.Branchpoint())
}
