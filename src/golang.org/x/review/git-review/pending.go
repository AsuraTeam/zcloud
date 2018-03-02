// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

func pending(args []string) {
	expectZeroArgs(args, "pending")
	// TODO(adg): implement -r

	current := CurrentBranch().Name
	for _, branch := range LocalBranches() {
		p := "  "
		if branch.Name == current {
			p = "* "
		}
		pending := branch.HasPendingCommit()
		if pending {
			fmt.Printf("%v%v: %v\n", p, branch.Name, branch.Subject())
		} else if branch.Name == current {
			// Nothing pending but print the line to show where we are.
			fmt.Printf("%v%v: (no pending change)\n", p, branch.Name)
		}
	}
}
