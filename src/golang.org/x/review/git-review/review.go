// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO(adg): recognize non-master remote branches
// TODO(adg): accept -a flag on 'commit' (like git commit -a)
// TODO(adg): check style of commit message
// TOOD(adg): print gerrit votes on 'pending'
// TODO(adg): add gofmt commit hook
// TODO(adg): print changed files on review sync
// TODO(adg): translate email addresses without @ by looking up somewhere

// Command git-review provides a simple command-line user interface for
// working with git repositories and the Gerrit code review system.
// See "git-review help" for details.
package main // import "golang.org/x/review/git-review"

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	flags   *flag.FlagSet
	verbose = new(count) // installed as -v below
	noRun   = new(bool)
)

func initFlags() {
	flags = flag.NewFlagSet("", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0], os.Args[0])
	}
	flags.Var(verbose, "v", "report git commands")
	flags.BoolVar(noRun, "n", false, "print but do not run commands")
}

const globalFlags = "[-n] [-v]"

const usage = `Usage: %s <command> ` + globalFlags + `
Type "%s help" for more information.
`

const help = `Usage: %s <command> ` + globalFlags + `

The review command is a wrapper for the git command that provides a simple
interface to the "single-commit feature branch" development model.

Available commands:

	change [name]
		Create a change commit, or amend an existing change commit,
		with the staged changes. If a branch name is provided, check
		out that branch (creating it if it does not exist).
		(Does not amend the existing commit when switching branches.)

	gofmt
		TBD

	help
		Show this help text.

	hooks
		Install Git commit hooks for Gerrit and gofmt.
		Every other operation except help also does this, if they are not
		already installed.

	mail [-f] [-r reviewer,...] [-cc mail,...]
		Upload change commit to the code review server and send mail
		requesting a code review.
		If -f is specified, upload even if there are staged changes.

	mail -diff
		Show the changes but do not send mail or upload.

	pending [-r]
		Show local branches and their head commits.
		If -r is specified, show additional information from Gerrit.

	sync
		Fetch changes from the remote repository and merge them into
		the current branch, rebasing the change commit on top of them.


`

func main() {
	initFlags()

	if len(os.Args) < 2 {
		flags.Usage()
		if dieTrap != nil {
			dieTrap()
		}
		os.Exit(2)
	}
	command, args := os.Args[1], os.Args[2:]

	if command == "help" {
		fmt.Fprintf(os.Stdout, help, os.Args[0])
		return
	}

	installHook()

	switch command {
	case "change", "c":
		change(args)
	case "gofmt":
		dief("gofmt not implemented")
	case "hooks":
		// done - installHook already ran
	case "mail", "m":
		mail(args)
	case "pending", "p":
		pending(args)
	case "sync", "s":
		doSync(args)
	default:
		flags.Usage()
	}
}

func expectZeroArgs(args []string, command string) {
	flags.Parse(args)
	if len(flags.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s %s %s\n", os.Args[0], command, globalFlags)
		os.Exit(2)
	}
}

func run(command string, args ...string) {
	if err := runErr(command, args...); err != nil {
		if *verbose == 0 {
			// If we're not in verbose mode, print the command
			// before dying to give context to the failure.
			fmt.Fprintln(os.Stderr, commandString(command, args))
		}
		dief("%v", err)
	}
}

var runLog []string

func runErr(command string, args ...string) error {
	if *verbose > 0 || *noRun {
		fmt.Fprintln(os.Stderr, commandString(command, args))
	}
	if *noRun {
		return nil
	}
	if runLog != nil {
		runLog = append(runLog, strings.TrimSpace(command+" "+strings.Join(args, " ")))
	}
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// getOutput runs the specified command and returns its combined standard
// output and standard error outputs.
// NOTE: It should only be used to run commands that return information,
// **not** commands that make any actual changes.
func getOutput(command string, args ...string) string {
	// NOTE: We only show these non-state-modifying commands with -v -v.
	// Otherwise things like 'git sync -v' show all our internal "find out about
	// the git repo" commands, which is confusing if you are just trying to find
	// out what git sync means.
	if *verbose > 1 {
		fmt.Fprintln(os.Stderr, commandString(command, args))
	}
	b, err := exec.Command(command, args...).CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n%s\n", commandString(command, args), b)
		dief("%v", err)
	}
	return string(bytes.TrimSpace(b))
}

// getLines is like getOutput but it returns non-empty output lines.
// NOTE: It should only be used to run commands that return information,
// **not** commands that make any actual changes.
func getLines(command string, args ...string) []string {
	var s []string
	for _, l := range strings.Split(getOutput(command, args...), "\n") {
		s = append(s, strings.TrimSpace(l))
	}
	return s
}

func commandString(command string, args []string) string {
	return strings.Join(append([]string{command}, args...), " ")
}

var dieTrap func()

func dief(format string, args ...interface{}) {
	printf(format, args...)
	if dieTrap != nil {
		dieTrap()
	}
	os.Exit(1)
}

func verbosef(format string, args ...interface{}) {
	if *verbose > 0 {
		printf(format, args...)
	}
}

func printf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], fmt.Sprintf(format, args...))
}

// count is a flag.Value that is like a flag.Bool and a flag.Int.
// If used as -name, it increments the count, but -name=x sets the count.
// Used for verbose flag -v.
type count int

func (c *count) String() string {
	return fmt.Sprint(int(*c))
}

func (c *count) Set(s string) error {
	switch s {
	case "true":
		*c++
	case "false":
		*c = 0
	default:
		n, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("invalid count %q", s)
		}
		*c = count(n)
	}
	return nil
}

func (c *count) IsBoolFlag() bool {
	return true
}
