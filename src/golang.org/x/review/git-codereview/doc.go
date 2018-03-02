// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Git-codereview manages the code review process for Git changes using a Gerrit
server.

The git-codereview tool manages "change branches" in the local git repository.
Each such branch tracks a single commit, or "pending change",
that is reviewed using a Gerrit server; the Gerrit remote must be
named 'origin' in the local git repo.

Modifications to the pending change are applied by amending the commit.
This process implements the "single-commit feature branch" model.

Once installed as git-codereview, the tool's commands are available through git
either by running

	git codereview <command>

or, if aliases are installed, as

	git <command>

The review tool's command names do not conflict with any extant git commands.
This document uses the first form for clarity but most users install these
aliases in their .gitconfig file:

	[alias]
		change = codereview change
		gofmt = codereview gofmt
		mail = codereview mail
		pending = codereview pending
		rebase-work = codereview rebase-work
		submit = codereview submit
		sync = codereview sync

All commands accept these global flags:

The -v flag prints all commands that make changes.

The -n flag prints all commands that would be run, but does not run them.

Descriptions of each command follow.

Branchpoint

	git codereview branchpoint

The branchpoint command prints the commit hash of the most recent change
on the current branch that is shared with the Gerrit server. This is the point
where local work branched from the published tree. The command is intended
mainly for use in scripts. For example, "git diff $(git codereview branchpoint)"
or "git log $(git codereview branchpoint)..HEAD".

Change

The change command creates and moves between Git branches and maintains the
pending commits on work branches.

	git codereview change [-a] [-q] [branchname]

Given a branch name as an argument, the change command switches to the named
branch, creating it if necessary. If the branch is created and there are staged
changes, it will commit the changes to the branch, creating a new pending
change.

With no argument, the change command creates a new pending change from the
staged changes in the current branch or, if there is already a pending change,
amends that change.

The -q option skips the editing of an extant pending change's commit message.

The -a option automatically adds any unstaged changes in tracked files during
commit; it is equivalent to the 'git commit' -a option.

Gofmt

The gofmt command applies the gofmt program to all files modified in the
current work branch, both in the staging area (index) and the working tree
(local directory).

	git codereview gofmt [-l]

The -l option causes the command to list the files that need reformatting but
not reformat them. Otherwise, the gofmt command reformats modified files in
place. That is, files in the staging area are reformatted in the staging area,
and files in the working tree are reformatted in the working tree.

Help

The help command displays basic usage instructions.

	git codereview help

Hooks

The hooks command installs the Git hooks to enforce code review conventions.

	git codereview hooks

The pre-commit hook checks that all Go code is formatted with gofmt and that
the commit is not being made directly to the master branch.

The commit-msg hook adds the Gerrit "Change-Id" line to the commit message if
not present. It also checks that the message uses the convention established by
the Go project that the first line has the form, pkg/path: summary.

The hooks command will not overwrite an existing hook.
If it is not installing hooks, use 'git codereview hooks -v' for details.
This hook installation is also done at startup by all other git codereview
commands, except 'help'.

Hook-Invoke

The hook-invoke command is an internal command that invokes the named Git hook.

	git codereview hook-invoke <hook> [args]

It is run by the shell scripts installed by the "git codereview hooks" command.

Mail

The mail command starts the code review process for the pending change.

	git codereview mail [-f] [-r email] [-cc email]

It pushes the pending change commit in the current branch to the Gerrit code
review server and prints the URL for the change on the server.
If the change already exists on the server, the mail command updates that
change with a new changeset.

The -r and -cc flags identify the email addresses of people to do the code
review and to be CC'ed about the code review.
Multiple addresses are given as a comma-separated list.

An email address passed to -r or -cc can be shortened from name@domain to name.
The mail command resolves such shortenings by reading the list of past reviewers
from the git repository log to find email addresses of the form name@somedomain
and then, in case of ambiguity, using the reviewer who appears most often.

The mail command fails if there are staged changes that are not committed.
The -f flag overrides this behavior.

The mail command updates the tag <branchname>.mailed to refer to the
commit that was most recently mailed, so running 'git diff <branchname>.mailed'
shows diffs between what is on the Gerrit server and the current directory.

Pending

The pending command prints to standard output the status of all pending changes
and staged, unstaged, and untracked files in the local repository.

	git codereview pending [-l]

The -l flag causes the command to use only locally available information.
By default, it fetches recent commits and code review information from the
Gerrit server.

Rebase-work

The rebase-work command runs git rebase in interactive mode over pending changes.
It is shorthand for "git rebase -i $(git codereview branchpoint)".
It differs from plain "git rebase -i" in that the latter will try to incorporate
new commits from the origin branch during the rebase, and git rebase-work
does not.

Submit

The submit command pushes the pending change to the Gerrit server and tells
Gerrit to submit it to the master branch.

	git codereview submit

The command fails if there are modified files (staged or unstaged) that are not
part of the pending change.

After submitting the change, the change command tries to synchronize the
current branch to the submitted commit, if it can do so cleanly.
If not, it will prompt the user to run 'git codereview sync' manually.

After a successful sync, the branch can be used to prepare a new change.

Sync

The sync command updates the local repository.

	git codereview sync

It fetches changes from the remote repository and merges changes from the
upstream branch to the current branch, rebasing the pending change, if any,
onto those changes.

*/
package main
