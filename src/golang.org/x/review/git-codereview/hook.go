// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var hookPath = ".git/hooks/"
var hookFiles = []string{
	"commit-msg",
	"pre-commit",
}

func installHook() {
	for _, hookFile := range hookFiles {
		filename := filepath.Join(repoRoot(), hookPath+hookFile)
		hookContent := fmt.Sprintf(hookScript, hookFile)

		if data, err := ioutil.ReadFile(filename); err == nil {
			// Special case: remove old hooks that use 'git-review'
			oldHookContent := fmt.Sprintf(oldHookScript, hookFile)
			if string(data) == oldHookContent {
				verbosef("removing old %v hook", hookFile)
				os.Remove(filename)
			}
			// Special case: remove old commit-msg shell script
			// in favor of invoking the git-codereview hook
			// implementation, which will be easier to change in
			// the future.
			if hookFile == "commit-msg" && string(data) == oldCommitMsgHook {
				verbosef("removing old commit-msg hook")
				os.Remove(filename)
			}
		}

		// If hook file exists, assume it is okay.
		_, err := os.Stat(filename)
		if err == nil {
			if *verbose > 0 {
				data, err := ioutil.ReadFile(filename)
				if err != nil {
					verbosef("reading hook: %v", err)
				} else if string(data) != hookContent {
					verbosef("unexpected hook content in %s", filename)
				}
			}
			continue
		}

		if !os.IsNotExist(err) {
			dief("checking hook: %v", err)
		}
		verbosef("installing %s hook", hookFile)
		if err := ioutil.WriteFile(filename, []byte(hookContent), 0700); err != nil {
			dief("writing hook: %v", err)
		}
	}
}

func repoRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		dief("could not get current directory: %v", err)
	}
	rootlen := 1
	if runtime.GOOS == "windows" {
		rootlen += len(filepath.VolumeName(dir))
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		if len(dir) == rootlen && dir[rootlen-1] == filepath.Separator {
			dief("git root not found. Rerun from within the Git tree.")
		}
		dir = filepath.Dir(dir)
	}
}

var hookScript = `#!/bin/sh
exec git-codereview hook-invoke %s "$@"
`

var oldHookScript = `#!/bin/sh
exec git-review hook-invoke %s "$@"
`

func cmdHookInvoke(args []string) {
	flags.Parse(args)
	args = flags.Args()
	if len(args) == 0 {
		dief("usage: git-codereview hook-invoke <hook-name> [args...]")
	}
	switch args[0] {
	case "commit-msg":
		hookCommitMsg(args[1:])
	case "pre-commit":
		hookPreCommit(args[1:])
	}
}

var (
	issueRefRE         = regexp.MustCompile(`(?P<space>\s)(?P<ref>#\d+\w)`)
	oldFixesRETemplate = `Fixes +(issue +(%s)?#?)?(?P<issueNum>[0-9]+)`
)

// hookCommitMsg is installed as the git commit-msg hook.
// It adds a Change-Id line to the bottom of the commit message
// if there is not one already.
func hookCommitMsg(args []string) {
	if len(args) != 1 {
		dief("usage: git-codereview hook-invoke commit-msg message.txt\n")
	}

	b := CurrentBranch()
	if b.DetachedHead() {
		// Likely executing rebase or some other internal operation.
		// Probably a mistake to make commit message changes.
		return
	}

	file := args[0]
	oldData, err := ioutil.ReadFile(file)
	if err != nil {
		dief("%v", err)
	}
	data := append([]byte{}, oldData...)
	data = stripComments(data)

	// Empty message not allowed.
	if len(bytes.TrimSpace(data)) == 0 {
		dief("empty commit message")
	}

	// Insert a blank line between first line and subsequent lines if not present.
	eol := bytes.IndexByte(data, '\n')
	if eol != -1 && len(data) > eol+1 && data[eol+1] != '\n' {
		data = append(data, 0)
		copy(data[eol+1:], data[eol:])
		data[eol+1] = '\n'
	}

	issueRepo := config()["issuerepo"]
	// Update issue references to point to issue repo, if set.
	if issueRepo != "" {
		data = issueRefRE.ReplaceAll(data, []byte("${space}"+issueRepo+"${ref}"))
	}
	// TestHookCommitMsgIssueRepoRewrite makes sure the regex is valid
	oldFixesRE := regexp.MustCompile(fmt.Sprintf(oldFixesRETemplate, regexp.QuoteMeta(issueRepo)))
	data = oldFixesRE.ReplaceAll(data, []byte("Fixes "+issueRepo+"#${issueNum}"))

	// Complain if two Change-Ids are present.
	// This can happen during an interactive rebase;
	// it is easy to forget to remove one of them.
	nChangeId := bytes.Count(data, []byte("\nChange-Id: "))
	if nChangeId > 1 {
		dief("multiple Change-Id lines")
	}

	// Add Change-Id to commit message if not present.
	if nChangeId == 0 {
		n := len(data)
		for n > 0 && data[n-1] == '\n' {
			n--
		}
		var id [20]byte
		if _, err := io.ReadFull(rand.Reader, id[:]); err != nil {
			dief("generating Change-Id: %v", err)
		}
		data = append(data[:n], fmt.Sprintf("\n\nChange-Id: I%x\n", id[:])...)
	}

	// Add branch prefix to commit message if not present and not on master
	// and not a special Git fixup! or squash! commit message.
	branch := strings.TrimPrefix(b.OriginBranch(), "origin/")
	if branch != "master" {
		prefix := "[" + branch + "] "
		if !bytes.HasPrefix(data, []byte(prefix)) && !isFixup(data) {
			data = []byte(prefix + string(data))
		}
	}

	// Write back.
	if !bytes.Equal(data, oldData) {
		if err := ioutil.WriteFile(file, data, 0666); err != nil {
			dief("%v", err)
		}
	}
}

var (
	fixupBang  = []byte("fixup!")
	squashBang = []byte("squash!")
)

// isFixup reports whether text is a Git fixup! or squash! commit,
// which must not have a prefix.
func isFixup(text []byte) bool {
	return bytes.HasPrefix(text, fixupBang) || bytes.HasPrefix(text, squashBang)
}

// stripComments strips lines that begin with "#".
func stripComments(in []byte) []byte {
	return regexp.MustCompile(`(?m)^#.*\n`).ReplaceAll(in, nil)
}

// hookPreCommit is installed as the git pre-commit hook.
// It prevents commits to the master branch.
// It checks that the Go files added, copied, or modified by
// the change are gofmt'd, and if not it prints gofmt instructions
// and exits with nonzero status.
func hookPreCommit(args []string) {
	// Prevent commits to master branches.
	b := CurrentBranch()
	if b.DetachedHead() {
		// This is an internal commit such as during git rebase.
		// Don't die, and don't force gofmt.
		return
	}
	if !b.IsLocalOnly() {
		dief("cannot commit on %s branch", b.Name)
	}

	hookGofmt()
}

func hookGofmt() {
	if os.Getenv("GIT_GOFMT_HOOK") == "off" {
		fmt.Fprintf(stderr(), "git-gofmt-hook disabled by $GIT_GOFMT_HOOK=off\n")
		return
	}

	files, stderr := runGofmt(gofmtPreCommit)

	if stderr != "" {
		msgf := printf
		if len(files) == 0 {
			msgf = dief
		}
		msgf("gofmt reported errors:\n\t%s", strings.Replace(strings.TrimSpace(stderr), "\n", "\n\t", -1))
	}

	if len(files) == 0 {
		return
	}

	dief("gofmt needs to format these files (run 'git gofmt'):\n\t%s",
		strings.Join(files, "\n\t"))
}

// This is NOT USED ANYMORE.
// It is here only for comparing against old commit-hook files.
var oldCommitMsgHook = `#!/bin/sh
# From Gerrit Code Review 2.2.1
#
# Part of Gerrit Code Review (http://code.google.com/p/gerrit/)
#
# Copyright (C) 2009 The Android Open Source Project
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

CHANGE_ID_AFTER="Bug|Issue"
MSG="$1"

# Check for, and add if missing, a unique Change-Id
#
add_ChangeId() {
	clean_message=` + "`" + `sed -e '
		/^diff --git a\/.*/{
			s///
			q
		}
		/^Signed-off-by:/d
		/^#/d
	' "$MSG" | git stripspace` + "`" + `
	if test -z "$clean_message"
	then
		return
	fi

	if grep -i '^Change-Id:' "$MSG" >/dev/null
	then
		return
	fi

	id=` + "`" + `_gen_ChangeId` + "`" + `
	perl -e '
		$MSG = shift;
		$id = shift;
		$CHANGE_ID_AFTER = shift;

		undef $/;
		open(I, $MSG); $_ = <I>; close I;
		s|^diff --git a/.*||ms;
		s|^#.*$||mg;
		exit unless $_;

		@message = split /\n/;
		$haveFooter = 0;
		$startFooter = @message;
		for($line = @message - 1; $line >= 0; $line--) {
			$_ = $message[$line];

			if (/^[a-zA-Z0-9-]+:/ && !m,^[a-z0-9-]+://,) {
				$haveFooter++;
				next;
			}
			next if /^[ []/;
			$startFooter = $line if ($haveFooter && /^\r?$/);
			last;
		}

		@footer = @message[$startFooter+1..@message];
		@message = @message[0..$startFooter];
		push(@footer, "") unless @footer;

		for ($line = 0; $line < @footer; $line++) {
			$_ = $footer[$line];
			next if /^($CHANGE_ID_AFTER):/i;
			last;
		}
		splice(@footer, $line, 0, "Change-Id: I$id");

		$_ = join("\n", @message, @footer);
		open(O, ">$MSG"); print O; close O;
	' "$MSG" "$id" "$CHANGE_ID_AFTER"
}
_gen_ChangeIdInput() {
	echo "tree ` + "`" + `git write-tree` + "`" + `"
	if parent=` + "`" + `git rev-parse HEAD^0 2>/dev/null` + "`" + `
	then
		echo "parent $parent"
	fi
	echo "author ` + "`" + `git var GIT_AUTHOR_IDENT` + "`" + `"
	echo "committer ` + "`" + `git var GIT_COMMITTER_IDENT` + "`" + `"
	echo
	printf '%s' "$clean_message"
}
_gen_ChangeId() {
	_gen_ChangeIdInput |
	git hash-object -t commit --stdin
}


add_ChangeId
`
