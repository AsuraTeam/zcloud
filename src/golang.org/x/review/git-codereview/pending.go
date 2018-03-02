// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

var (
	pendingLocal       bool // -l flag, use only local operations (no network)
	pendingCurrentOnly bool // -c flag, show only current branch
	pendingShort       bool // -s flag, short display
)

// A pendingBranch collects information about a single pending branch.
// We overlap the reading of this information for each branch.
type pendingBranch struct {
	*Branch            // standard Branch functionality
	current   bool     // is this the current branch?
	staged    []string // files in staging area, only if current==true
	unstaged  []string // files unstaged in local directory, only if current==true
	untracked []string // files untracked in local directory, only if current==true
}

// load populates b with information about the branch.
func (b *pendingBranch) load() {
	b.loadPending()
	if !b.current && b.commitsAhead == 0 {
		// Won't be displayed, don't bother looking any closer.
		return
	}
	b.OriginBranch() // cache result
	if b.current {
		b.staged, b.unstaged, b.untracked = LocalChanges()
	}
	for _, c := range b.Pending() {
		c.committed = nonBlankLines(cmdOutput("git", "diff", "--name-only", c.Parent, c.Hash, "--"))
		if !pendingLocal {
			c.g, c.gerr = b.GerritChange(c, "DETAILED_LABELS", "CURRENT_REVISION", "MESSAGES", "DETAILED_ACCOUNTS")
		}
		if c.g == nil {
			c.g = new(GerritChange) // easier for formatting code
		}
	}
}

func cmdPending(args []string) {
	flags.BoolVar(&pendingCurrentOnly, "c", false, "show only current branch")
	flags.BoolVar(&pendingLocal, "l", false, "use only local information - no network operations")
	flags.BoolVar(&pendingShort, "s", false, "show short listing")
	flags.Parse(args)
	if len(flags.Args()) > 0 {
		fmt.Fprintf(stderr(), "Usage: %s pending %s [-c] [-l] [-s]\n", os.Args[0], globalFlags)
		os.Exit(2)
	}

	// Fetch info about remote changes, so that we can say which branches need sync.
	if !pendingLocal {
		run("git", "fetch", "-q")
		http.DefaultClient.Timeout = 60 * time.Second
	}

	// Build list of pendingBranch structs to be filled in.
	// The current branch is always first.
	var branches []*pendingBranch
	branches = []*pendingBranch{{Branch: CurrentBranch(), current: true}}
	if !pendingCurrentOnly {
		current := CurrentBranch().Name
		for _, b := range LocalBranches() {
			if b.Name != current {
				branches = append(branches, &pendingBranch{Branch: b})
			}
		}
	}

	// The various data gathering is a little slow,
	// especially run in serial with a lot of branches.
	// Overlap inspection of multiple branches.
	// Each branch is only accessed by a single worker.

	// Build work queue.
	work := make(chan *pendingBranch, len(branches))
	done := make(chan bool, len(branches))
	for _, b := range branches {
		work <- b
	}
	close(work)

	// Kick off goroutines to do work.
	n := len(branches)
	if n > 10 {
		n = 10
	}
	for i := 0; i < n; i++ {
		go func() {
			for b := range work {
				b.load()
				done <- true
			}
		}()
	}

	// Wait for goroutines to finish.
	// Note: Counting work items, not goroutines (there may be fewer goroutines).
	for range branches {
		<-done
	}

	// Print output.
	// If there are multiple changes in the current branch, the output splits them out into separate sections,
	// in reverse commit order, to match git log output.
	//
	//	wbshadow 7a524a1..a496c1e (current branch, all mailed, 23 behind, tracking master)
	//	+ uncommitted changes
	//		Files unstaged:
	//			src/runtime/proc1.go
	//
	//	+ a496c1e https://go-review.googlesource.com/2064 (mailed)
	//		runtime: add missing write barriers in append's copy of slice data
	//
	//		Found with GODEBUG=wbshadow=1 mode.
	//		Eventually that will run automatically, but right now
	//		it still detects other missing write barriers.
	//
	//		Change-Id: Ic8624401d7c8225a935f719f96f2675c6f5c0d7c
	//
	//		Code-Review:
	//			+0 Austin Clements, Rick Hudson
	//		Files in this change:
	//			src/runtime/slice.go
	//
	//	+ 95390c7 https://go-review.googlesource.com/2061 (mailed)
	//		runtime: add GODEBUG wbshadow for finding missing write barriers
	//
	//		This is the detection code. It works well enough that I know of
	//		a handful of missing write barriers. However, those are subtle
	//		enough that I'll address them in separate followup CLs.
	//
	//		Change-Id: If863837308e7c50d96b5bdc7d65af4969bf53a6e
	//
	//		Code-Review:
	//			+0 Austin Clements, Rick Hudson
	//		Files in this change:
	//			src/runtime/extern.go
	//			src/runtime/malloc1.go
	//			src/runtime/malloc2.go
	//			src/runtime/mgc.go
	//			src/runtime/mgc0.go
	//			src/runtime/proc1.go
	//			src/runtime/runtime1.go
	//			src/runtime/runtime2.go
	//			src/runtime/stack1.go
	//
	// The first line only gives information that applies to the entire branch:
	// the name, the commit range, whether this is the current branch, whether
	// all the commits are mailed/submitted, how far behind, what remote branch
	// it is tracking.
	// The individual change sections have per-change information: the hash of that
	// commit, the URL on the Gerrit server, whether it is mailed/submitted, the list of
	// files in that commit. The uncommitted file modifications are shown as a separate
	// section, at the beginning, to fit better into the reverse commit order.
	//
	// The short view compresses the listing down to two lines per commit:
	//	wbshadow 7a524a1..a496c1e (current branch, all mailed, 23 behind, tracking master)
	//	+ uncommitted changes
	//		Files unstaged:
	//			src/runtime/proc1.go
	//	+ a496c1e runtime: add missing write barriers in append's copy of slice data (CL 2064, mailed)
	//	+ 95390c7 runtime: add GODEBUG wbshadow for finding missing write barriers (CL 2061, mailed)

	var buf bytes.Buffer
	printFileList := func(name string, list []string) {
		if len(list) == 0 {
			return
		}
		fmt.Fprintf(&buf, "\tFiles %s:\n", name)
		for _, file := range list {
			fmt.Fprintf(&buf, "\t\t%s\n", file)
		}
	}

	for _, b := range branches {
		if !b.current && b.commitsAhead == 0 {
			// Hide branches with no work on them.
			continue
		}

		fmt.Fprintf(&buf, "%s", b.Name)
		work := b.Pending()
		if len(work) > 0 {
			fmt.Fprintf(&buf, " %.7s..%s", b.branchpoint, work[0].ShortHash)
		}
		var tags []string
		if b.current {
			tags = append(tags, "current branch")
		}
		if allMailed(work) && len(work) > 0 {
			tags = append(tags, "all mailed")
		}
		if allSubmitted(work) && len(work) > 0 {
			tags = append(tags, "all submitted")
		}
		if b.commitsBehind > 0 {
			tags = append(tags, fmt.Sprintf("%d behind", b.commitsBehind))
		}
		if b.OriginBranch() != "origin/master" {
			tags = append(tags, "tracking "+strings.TrimPrefix(b.OriginBranch(), "origin/"))
		}
		if len(tags) > 0 {
			fmt.Fprintf(&buf, " (%s)", strings.Join(tags, ", "))
		}
		fmt.Fprintf(&buf, "\n")
		printed := false
		if text := b.errors(); text != "" {
			fmt.Fprintf(&buf, "\tERROR: %s\n", strings.Replace(strings.TrimSpace(text), "\n", "\n\t", -1))
			if !pendingShort {
				printed = true
				fmt.Fprintf(&buf, "\n")
			}
		}

		if b.current && len(b.staged)+len(b.unstaged)+len(b.untracked) > 0 {
			printed = true
			fmt.Fprintf(&buf, "+ uncommitted changes\n")
			printFileList("untracked", b.untracked)
			printFileList("unstaged", b.unstaged)
			printFileList("staged", b.staged)
			if !pendingShort {
				fmt.Fprintf(&buf, "\n")
			}
		}

		for _, c := range work {
			printed = true
			fmt.Fprintf(&buf, "+ ")
			formatCommit(&buf, c, pendingShort)
			if !pendingShort {
				printFileList("in this change", c.committed)
				fmt.Fprintf(&buf, "\n")
			}
		}
		if pendingShort || !printed {
			fmt.Fprintf(&buf, "\n")
		}
	}

	stdout().Write(buf.Bytes())
}

// formatCommit writes detailed information about c to w. c.g must
// have the "CURRENT_REVISION" (or "ALL_REVISIONS") and
// "DETAILED_LABELS" options set.
//
// If short is true, this writes a single line overview.
//
// If short is false, this writes detailed information about the
// commit and its Gerrit state.
func formatCommit(w io.Writer, c *Commit, short bool) {
	g := c.g
	msg := strings.TrimRight(c.Message, "\r\n")
	fmt.Fprintf(w, "%s", c.ShortHash)
	var tags []string
	if short {
		if i := strings.Index(msg, "\n"); i >= 0 {
			msg = msg[:i]
		}
		fmt.Fprintf(w, " %s", msg)
		if g.Number != 0 {
			tags = append(tags, fmt.Sprintf("CL %d%s", g.Number, codeReviewScores(g)))
		}
	} else {
		if g.Number != 0 {
			fmt.Fprintf(w, " %s/%d", auth.url, g.Number)
		}
	}
	if g.CurrentRevision == c.Hash {
		tags = append(tags, "mailed")
	}
	if g.Status == "MERGED" {
		tags = append(tags, "submitted")
	}
	if len(tags) > 0 {
		fmt.Fprintf(w, " (%s)", strings.Join(tags, ", "))
	}
	fmt.Fprintf(w, "\n")
	if short {
		return
	}

	fmt.Fprintf(w, "\t%s\n", strings.Replace(msg, "\n", "\n\t", -1))
	fmt.Fprintf(w, "\n")

	for _, name := range g.LabelNames() {
		label := g.Labels[name]
		minValue := 10000
		maxValue := -10000
		byScore := map[int][]string{}
		for _, x := range label.All {
			// Hide CL owner unless owner score is nonzero.
			if g.Owner != nil && x.ID == g.Owner.ID && x.Value == 0 {
				continue
			}
			byScore[x.Value] = append(byScore[x.Value], x.Name)
			if minValue > x.Value {
				minValue = x.Value
			}
			if maxValue < x.Value {
				maxValue = x.Value
			}
		}
		// Unless there are scores to report, do not show labels other than Code-Review.
		// This hides Run-TryBot and TryBot-Result.
		if minValue >= 0 && maxValue <= 0 && name != "Code-Review" {
			continue
		}
		fmt.Fprintf(w, "\t%s:\n", name)
		for score := maxValue; score >= minValue; score-- {
			who := byScore[score]
			if len(who) == 0 || score == 0 && name != "Code-Review" {
				continue
			}
			sort.Strings(who)
			fmt.Fprintf(w, "\t\t%+d %s\n", score, strings.Join(who, ", "))
		}
	}
}

// codeReviewScores reports the code review scores as tags for the short output.
//
// g must have the "DETAILED_LABELS" option set.
func codeReviewScores(g *GerritChange) string {
	label := g.Labels["Code-Review"]
	if label == nil {
		return ""
	}
	minValue := 10000
	maxValue := -10000
	for _, x := range label.All {
		if minValue > x.Value {
			minValue = x.Value
		}
		if maxValue < x.Value {
			maxValue = x.Value
		}
	}
	var scores string
	if minValue < 0 {
		scores += fmt.Sprintf(" %d", minValue)
	}
	if maxValue > 0 {
		scores += fmt.Sprintf(" %+d", maxValue)
	}
	return scores
}

// allMailed reports whether all commits in work have been posted to Gerrit.
func allMailed(work []*Commit) bool {
	for _, c := range work {
		if c.Hash != c.g.CurrentRevision {
			return false
		}
	}
	return true
}

// allSubmitted reports whether all commits in work have been submitted to the origin branch.
func allSubmitted(work []*Commit) bool {
	for _, c := range work {
		if c.g.Status != "MERGED" {
			return false
		}
	}
	return true
}

// errors returns any errors that should be displayed
// about the state of the current branch, diagnosing common mistakes.
func (b *Branch) errors() string {
	b.loadPending()
	var buf bytes.Buffer
	if !b.IsLocalOnly() && b.commitsAhead > 0 {
		fmt.Fprintf(&buf, "Branch contains %d commit%s not on origin/%s.\n", b.commitsAhead, suffix(b.commitsAhead, "s"), b.Name)
		fmt.Fprintf(&buf, "\tDo not commit directly to %s branch.\n", b.Name)
	}
	return buf.String()
}

// suffix returns an empty string if n == 1, s otherwise.
func suffix(n int, s string) string {
	if n == 1 {
		return ""
	}
	return s
}
