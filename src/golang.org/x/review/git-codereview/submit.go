// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

// TODO(rsc): Add -tbr, along with standard exceptions (doc/go1.5.txt)

func cmdSubmit(args []string) {
	var interactive bool
	flags.BoolVar(&interactive, "i", false, "interactively select commits to submit")
	flags.Usage = func() {
		fmt.Fprintf(stderr(), "Usage: %s submit %s [-i | commit-hash...]\n", os.Args[0], globalFlags)
	}
	flags.Parse(args)
	if interactive && flags.NArg() > 0 {
		flags.Usage()
		os.Exit(2)
	}

	b := CurrentBranch()
	var cs []*Commit
	if interactive {
		hashes := submitHashes(b)
		if len(hashes) == 0 {
			printf("nothing to submit")
			return
		}
		for _, hash := range hashes {
			cs = append(cs, b.CommitByHash("submit", hash))
		}
	} else if args := flags.Args(); len(args) >= 1 {
		for _, arg := range args {
			cs = append(cs, b.CommitByHash("submit", arg))
		}
	} else {
		cs = append(cs, b.DefaultCommit("submit"))
	}

	// No staged changes.
	// Also, no unstaged changes, at least for now.
	// This makes sure the sync at the end will work well.
	// We can relax this later if there is a good reason.
	checkStaged("submit")
	checkUnstaged("submit")

	// Submit the changes.
	var g *GerritChange
	for _, c := range cs {
		printf("submitting %s %s", c.ShortHash, c.Subject)
		g = submit(b, c)
	}

	// Sync client to revision that Gerrit committed, but only if we can do it cleanly.
	// Otherwise require user to run 'git sync' themselves (if they care).
	run("git", "fetch", "-q")
	if len(cs) == 1 && len(b.Pending()) == 1 {
		if err := runErr("git", "checkout", "-q", "-B", b.Name, g.CurrentRevision, "--"); err != nil {
			dief("submit succeeded, but cannot sync local branch\n"+
				"\trun 'git sync' to sync, or\n"+
				"\trun 'git branch -D %s; git change master; git sync' to discard local branch", b.Name)
		}
	} else {
		printf("submit succeeded; run 'git sync' to sync")
	}

	// Done! Change is submitted, branch is up to date, ready for new work.
}

// submit submits a single commit c on branch b and returns the
// GerritChange for the submitted change. It dies if the submit fails.
func submit(b *Branch, c *Commit) *GerritChange {
	// Fetch Gerrit information about this change.
	g, err := b.GerritChange(c, "LABELS", "CURRENT_REVISION")
	if err != nil {
		dief("%v", err)
	}

	// Pre-check that this change appears submittable.
	// The final submit will check this too, but it is better to fail now.
	if err = submitCheck(g); err != nil {
		dief("cannot submit: %v", err)
	}

	// Upload most recent revision if not already on server.

	if c.Hash != g.CurrentRevision {
		run("git", "push", "-q", "origin", b.PushSpec(c))

		// Refetch change information, especially mergeable.
		g, err = b.GerritChange(c, "LABELS", "CURRENT_REVISION")
		if err != nil {
			dief("%v", err)
		}
	}

	// Don't bother if the server can't merge the changes.
	if !g.Mergeable {
		// Server cannot merge; explicit sync is needed.
		dief("cannot submit: conflicting changes submitted, run 'git sync'")
	}

	if *noRun {
		printf("stopped before submit")
		return g
	}

	// Otherwise, try the submit. Sends back updated GerritChange,
	// but we need extended information and the reply is in the
	// "SUBMITTED" state anyway, so ignore the GerritChange
	// in the response and fetch a new one below.
	if err := gerritAPI("/a/changes/"+fullChangeID(b, c)+"/submit", []byte(`{"wait_for_merge": true}`), nil); err != nil {
		dief("cannot submit: %v", err)
	}

	// It is common to get back "SUBMITTED" for a split second after the
	// request is made. That indicates that the change has been queued for submit,
	// but the first merge (the one wait_for_merge waited for)
	// failed, possibly due to a spurious condition. We see this often, and the
	// status usually changes to MERGED shortly thereafter.
	// Wait a little while to see if we can get to a different state.
	const steps = 6
	const max = 2 * time.Second
	for i := 0; i < steps; i++ {
		time.Sleep(max * (1 << uint(i+1)) / (1 << steps))
		g, err = b.GerritChange(c, "LABELS", "CURRENT_REVISION")
		if err != nil {
			dief("waiting for merge: %v", err)
		}
		if g.Status != "SUBMITTED" {
			break
		}
	}

	switch g.Status {
	default:
		dief("submit error: unexpected post-submit Gerrit change status %q", g.Status)

	case "MERGED":
		// good

	case "SUBMITTED":
		// see above
		dief("cannot submit: timed out waiting for change to be submitted by Gerrit")
	}

	return g
}

// submitCheck checks that g should be submittable. This is
// necessarily a best-effort check.
//
// g must have the "LABELS" option.
func submitCheck(g *GerritChange) error {
	// Check Gerrit change status.
	switch g.Status {
	default:
		return fmt.Errorf("unexpected Gerrit change status %q", g.Status)

	case "NEW", "SUBMITTED":
		// Not yet "MERGED", so try the submit.
		// "SUBMITTED" is a weird state. It means that Submit has been clicked once,
		// but it hasn't happened yet, usually because of a merge failure.
		// The user may have done git sync and may now have a mergable
		// copy waiting to be uploaded, so continue on as if it were "NEW".

	case "MERGED":
		// Can happen if moving between different clients.
		return fmt.Errorf("change already submitted, run 'git sync'")

	case "ABANDONED":
		return fmt.Errorf("change abandoned")
	}

	// Check for label approvals (like CodeReview+2).
	for _, name := range g.LabelNames() {
		label := g.Labels[name]
		if label.Optional {
			continue
		}
		if label.Rejected != nil {
			return fmt.Errorf("change has %s rejection", name)
		}
		if label.Approved == nil {
			return fmt.Errorf("change missing %s approval", name)
		}
	}

	return nil
}

// submitHashes interactively prompts for commits to submit.
func submitHashes(b *Branch) []string {
	// Get pending commits on b.
	pending := b.Pending()
	for _, c := range pending {
		// Note that DETAILED_LABELS does not imply LABELS.
		c.g, c.gerr = b.GerritChange(c, "CURRENT_REVISION", "LABELS", "DETAILED_LABELS")
		if c.g == nil {
			c.g = new(GerritChange)
		}
	}

	// Construct submit script.
	var script bytes.Buffer
	for i := len(pending) - 1; i >= 0; i-- {
		c := pending[i]

		if c.g.ID == "" {
			fmt.Fprintf(&script, "# change not on Gerrit:\n#")
		} else if err := submitCheck(c.g); err != nil {
			fmt.Fprintf(&script, "# %v:\n#", err)
		}

		formatCommit(&script, c, true)
	}

	fmt.Fprintf(&script, `
# The above commits will be submitted in order from top to bottom
# when you exit the editor.
#
# These lines can be re-ordered, removed, and commented out.
#
# If you remove all lines, the submit will be aborted.
`)

	// Edit the script.
	final := editor(script.String())

	// Parse the final script.
	var hashes []string
	for _, line := range lines(final) {
		line := strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		if i := strings.Index(line, " "); i >= 0 {
			line = line[:i]
		}
		hashes = append(hashes, line)
	}

	return hashes
}
