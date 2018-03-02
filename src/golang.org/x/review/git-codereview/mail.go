// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

func cmdMail(args []string) {
	var (
		diff   = flags.Bool("diff", false, "show change commit diff and don't upload or mail")
		force  = flags.Bool("f", false, "mail even if there are staged changes")
		topic  = flags.String("topic", "", "set Gerrit topic")
		trybot = flags.Bool("trybot", false, "run trybots on the uploaded CLs")
		rList  = new(stringList) // installed below
		ccList = new(stringList) // installed below
	)
	flags.Var(rList, "r", "comma-separated list of reviewers")
	flags.Var(ccList, "cc", "comma-separated list of people to CC:")

	flags.Usage = func() {
		fmt.Fprintf(stderr(), "Usage: %s mail %s [-r reviewer,...] [-cc mail,...] [-topic topic] [-trybot] [commit-hash]\n", os.Args[0], globalFlags)
	}
	flags.Parse(args)
	if len(flags.Args()) > 1 {
		flags.Usage()
		os.Exit(2)
	}

	b := CurrentBranch()

	var c *Commit
	if len(flags.Args()) == 1 {
		c = b.CommitByHash("mail", flags.Arg(0))
	} else {
		c = b.DefaultCommit("mail")
	}

	if *diff {
		run("git", "diff", b.Branchpoint()[:7]+".."+c.ShortHash, "--")
		return
	}

	if !*force && HasStagedChanges() {
		dief("there are staged changes; aborting.\n"+
			"Use '%s change' to include them or '%s mail -f' to force it.", os.Args[0], os.Args[0])
	}

	// for side effect of dying with a good message if origin is GitHub
	loadGerritOrigin()

	refSpec := b.PushSpec(c)
	start := "%"
	if *rList != "" {
		refSpec += mailList(start, "r", string(*rList))
		start = ","
	}
	if *ccList != "" {
		refSpec += mailList(start, "cc", string(*ccList))
		start = ","
	}
	if *topic != "" {
		// There's no way to escape the topic, but the only
		// ambiguous character is ',' (though other characters
		// like ' ' will be rejected outright by git).
		if strings.Contains(*topic, ",") {
			dief("topic may not contain a comma")
		}
		refSpec += start + "topic=" + *topic
		start = ","
	}
	if *trybot {
		refSpec += start + "l=Run-TryBot"
		start = ","
	}
	run("git", "push", "-q", "origin", refSpec)

	// Create local tag for mailed change.
	// If in the 'work' branch, this creates or updates work.mailed.
	// Older mailings are in the reflog, so work.mailed is newest,
	// work.mailed@{1} is the one before that, work.mailed@{2} before that,
	// and so on.
	// Git doesn't actually have a concept of a local tag,
	// but Gerrit won't let people push tags to it, so the tag
	// can't propagate out of the local client into the official repo.
	// There is no conflict with the branch names people are using
	// for work, because git change rejects any name containing a dot.
	// The space of names with dots is ours (the Go team's) to define.
	run("git", "tag", "-f", b.Name+".mailed", c.ShortHash)
}

// PushSpec returns the spec for a Gerrit push command to publish the change c in b.
// If c is nil, PushSpec returns a spec for pushing all changes in b.
func (b *Branch) PushSpec(c *Commit) string {
	local := "HEAD"
	if c != nil && (len(b.Pending()) == 0 || b.Pending()[0].Hash != c.Hash) {
		local = c.ShortHash
	}
	return local + ":refs/for/" + strings.TrimPrefix(b.OriginBranch(), "origin/")
}

// mailAddressRE matches the mail addresses we admit. It's restrictive but admits
// all the addresses in the Go CONTRIBUTORS file at time of writing (tested separately).
var mailAddressRE = regexp.MustCompile(`^([a-zA-Z0-9][-_.a-zA-Z0-9]*)(@[-_.a-zA-Z0-9]+)?$`)

// mailList turns the list of mail addresses from the flag value into the format
// expected by gerrit. The start argument is a % or , depending on where we
// are in the processing sequence.
func mailList(start, tag string, flagList string) string {
	errors := false
	spec := start
	short := ""
	long := ""
	for i, addr := range strings.Split(flagList, ",") {
		m := mailAddressRE.FindStringSubmatch(addr)
		if m == nil {
			printf("invalid reviewer mail address: %s", addr)
			errors = true
			continue
		}
		if m[2] == "" {
			email := mailLookup(addr)
			if email == "" {
				printf("unknown reviewer: %s", addr)
				errors = true
				continue
			}
			short += "," + addr
			long += "," + email
			addr = email
		}
		if i > 0 {
			spec += ","
		}
		spec += tag + "=" + addr
	}
	if short != "" {
		verbosef("expanded %s to %s", short[1:], long[1:])
	}
	if errors {
		die()
	}
	return spec
}

// reviewers is the list of reviewers for the current repository,
// sorted by how many reviews each has done.
var reviewers []reviewer

type reviewer struct {
	addr  string
	count int
}

// mailLookup translates the short name (like adg) into a full
// email address (like adg@golang.org).
// It returns "" if no translation is found.
// The algorithm for expanding short user names is as follows:
// Look at the git commit log for the current repository,
// extracting all the email addresses in Reviewed-By lines
// and sorting by how many times each address appears.
// For each short user name, walk the list, most common
// address first, and use the first address found that has
// the short user name on the left side of the @.
func mailLookup(short string) string {
	loadReviewers()

	short += "@"
	for _, r := range reviewers {
		if strings.HasPrefix(r.addr, short) {
			return r.addr
		}
	}
	return ""
}

// loadReviewers reads the reviewer list from the current git repo
// and leaves it in the global variable reviewers.
// See the comment on mailLookup for a description of how the
// list is generated and used.
func loadReviewers() {
	if reviewers != nil {
		return
	}
	countByAddr := map[string]int{}
	for _, line := range nonBlankLines(cmdOutput("git", "log", "--format=format:%B")) {
		if strings.HasPrefix(line, "Reviewed-by:") {
			f := strings.Fields(line)
			addr := f[len(f)-1]
			if strings.HasPrefix(addr, "<") && strings.Contains(addr, "@") && strings.HasSuffix(addr, ">") {
				countByAddr[addr[1:len(addr)-1]]++
			}
		}
	}

	reviewers = []reviewer{}
	for addr, count := range countByAddr {
		reviewers = append(reviewers, reviewer{addr, count})
	}
	sort.Sort(reviewersByCount(reviewers))
}

type reviewersByCount []reviewer

func (x reviewersByCount) Len() int      { return len(x) }
func (x reviewersByCount) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x reviewersByCount) Less(i, j int) bool {
	if x[i].count != x[j].count {
		return x[i].count > x[j].count
	}
	return x[i].addr < x[j].addr
}

// stringList is a flag.Value that is like flag.String, but if repeated
// keeps appending to the old value, inserting commas as separators.
// This allows people to write -r rsc,adg (like the old hg command)
// but also -r rsc -r adg (like standard git commands).
// This does change the meaning of -r rsc -r adg (it used to mean just adg).
type stringList string

func (x *stringList) String() string {
	return string(*x)
}

func (x *stringList) Set(s string) error {
	if *x != "" && s != "" {
		*x += ","
	}
	*x += stringList(s)
	return nil
}
