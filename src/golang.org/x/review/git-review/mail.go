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

func mail(args []string) {
	var (
		diff   = flags.Bool("diff", false, "show change commit diff and don't upload or mail")
		force  = flags.Bool("f", false, "mail even if there are staged changes")
		rList  = new(stringList) // installed below
		ccList = new(stringList) // installed below
	)
	flags.Var(rList, "r", "comma-separated list of reviewers")
	flags.Var(ccList, "cc", "comma-separated list of people to CC:")

	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s mail %s [-r reviewer,...] [-cc mail,...]\n", os.Args[0], globalFlags)
	}
	flags.Parse(args)
	if len(flags.Args()) != 0 {
		flags.Usage()
		os.Exit(2)
	}

	b := CurrentBranch()
	if b.ChangeID() == "" {
		dief("no pending change; can't mail.")
	}

	if *diff {
		run("git", "diff", "HEAD^..HEAD")
		return
	}

	if !*force && HasStagedChanges() {
		dief("there are staged changes; aborting.\n" +
			"Use 'review change' to include them or 'review mail -f' to force it.")
	}

	refSpec := "HEAD:refs/for/master"
	start := "%"
	if *rList != "" {
		refSpec += mailList(start, "r", string(*rList))
		start = ","
	}
	if *ccList != "" {
		refSpec += mailList(start, "cc", string(*ccList))
	}
	run("git", "push", "-q", "origin", refSpec)
}

// mailAddressRE matches the mail addresses we admit. It's restrictive but admits
// all the addresses in the Go CONTRIBUTORS file at time of writing (tested separately).
var mailAddressRE = regexp.MustCompile(`^[a-zA-Z0-9][-_.a-zA-Z0-9]*@[-_.a-zA-Z0-9]+$`)

// mailList turns the list of mail addresses from the flag value into the format
// expected by gerrit. The start argument is a % or , depending on where we
// are in the processing sequence.
func mailList(start, tag string, flagList string) string {
	spec := start
	for i, addr := range strings.Split(flagList, ",") {
		if !mailAddressRE.MatchString(addr) {
			dief("%q is not a valid reviewer mail address", addr)
		}
		if i > 0 {
			spec += ","
		}
		spec += tag + "=" + addr
	}
	return spec
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
