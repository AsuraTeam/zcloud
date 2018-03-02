// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO(rsc): Document multi-change branch behavior.

// Command git-codereview provides a simple command-line user interface for
// working with git repositories and the Gerrit code review system.
// See "git-codereview help" for details.
package main // import "golang.org/x/review/git-codereview"

import (
	"bytes"
	"flag"
	"fmt"
	"io"
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
		fmt.Fprintf(stderr(), usage, os.Args[0], os.Args[0])
	}
	flags.Var(verbose, "v", "report commands")
	flags.BoolVar(noRun, "n", false, "print but do not run commands")
}

const globalFlags = "[-n] [-v]"

const usage = `Usage: %s <command> ` + globalFlags + `
Type "%s help" for more information.
`

const help = `Usage: %s <command> ` + globalFlags + `

The git-codereview command is a wrapper for the git command that provides a
simple interface to the "single-commit feature branch" development model.

See the docs for details: https://godoc.org/golang.org/x/review/git-codereview

The -v flag prints all commands that make changes.
The -n flag prints all commands that would be run, but does not run them.

Available commands:

	change [name]
		Create a change commit, or amend an existing change commit,
		with the staged changes. If a branch name is provided, check
		out that branch (creating it if it does not exist).
		Does not amend the existing commit when switching branches.
		If -q is specified, skip the editing of an extant pending
		change's commit message.
		If -a is specified, automatically add any unstaged changes in
		tracked files during commit.

	gofmt [-l]
		Run gofmt on all tracked files in the staging area and the
		working tree.
		If -l is specified, list files that need formatting.
		Otherwise, reformat files in place.

	help
		Show this help text.

	hooks
		Install Git commit hooks for Gerrit and gofmt.
		Every other operation except help also does this,
		if they are not already installed.

	mail [-f] [-r reviewer,...] [-cc mail,...]
		Upload change commit to the code review server and send mail
		requesting a code review.
		If -f is specified, upload even if there are staged changes.
		The -r and -cc flags identify the email addresses of people to
		do the code review and to be CC'ed about the code review.
		Multiple addresses are given as a comma-separated list.

	mail -diff
		Show the changes but do not send mail or upload.

	pending [-c] [-l] [-s]
		Show the status of all pending changes and staged, unstaged,
		and untracked files in the local repository.
		If -c is specified, show only changes on the current branch.
		If -l is specified, only use locally available information.
		If -s is specified, show short output.

	submit [-i | commit-hash...]
		Push the pending change to the Gerrit server and tell Gerrit to
		submit it to the master branch.

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
		fmt.Fprintf(stdout(), help, os.Args[0])
		return
	}

	installHook()

	switch command {
	case "branchpoint":
		cmdBranchpoint(args)
	case "change":
		cmdChange(args)
	case "gofmt":
		cmdGofmt(args)
	case "hook-invoke":
		cmdHookInvoke(args)
	case "hooks":
		// done - installHook already ran
	case "mail", "m":
		cmdMail(args)
	case "pending":
		cmdPending(args)
	case "rebase-work":
		cmdRebaseWork(args)
	case "submit":
		cmdSubmit(args)
	case "sync":
		cmdSync(args)
	case "test-loadAuth": // for testing only
		loadAuth()
	default:
		flags.Usage()
	}
}

func expectZeroArgs(args []string, command string) {
	flags.Parse(args)
	if len(flags.Args()) > 0 {
		fmt.Fprintf(stderr(), "Usage: %s %s %s\n", os.Args[0], command, globalFlags)
		os.Exit(2)
	}
}

func run(command string, args ...string) {
	if err := runErr(command, args...); err != nil {
		if *verbose == 0 {
			// If we're not in verbose mode, print the command
			// before dying to give context to the failure.
			fmt.Fprintf(stderr(), "(running: %s)\n", commandString(command, args))
		}
		dief("%v", err)
	}
}

func runErr(command string, args ...string) error {
	return runDirErr("", command, args...)
}

var runLogTrap []string

func runDirErr(dir, command string, args ...string) error {
	if *verbose > 0 || *noRun {
		fmt.Fprintln(stderr(), commandString(command, args))
	}
	if *noRun {
		return nil
	}
	if runLogTrap != nil {
		runLogTrap = append(runLogTrap, strings.TrimSpace(command+" "+strings.Join(args, " ")))
	}
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout()
	cmd.Stderr = stderr()
	return cmd.Run()
}

// cmdOutput runs the command line, returning its output.
// If the command cannot be run or does not exit successfully,
// cmdOutput dies.
//
// NOTE: cmdOutput must be used only to run commands that read state,
// not for commands that make changes. Commands that make changes
// should be run using runDirErr so that the -v and -n flags apply to them.
func cmdOutput(command string, args ...string) string {
	return cmdOutputDir(".", command, args...)
}

// cmdOutputDir runs the command line in dir, returning its output.
// If the command cannot be run or does not exit successfully,
// cmdOutput dies.
//
// NOTE: cmdOutput must be used only to run commands that read state,
// not for commands that make changes. Commands that make changes
// should be run using runDirErr so that the -v and -n flags apply to them.
func cmdOutputDir(dir, command string, args ...string) string {
	s, err := cmdOutputDirErr(dir, command, args...)
	if err != nil {
		fmt.Fprintf(stderr(), "%v\n%s\n", commandString(command, args), s)
		dief("%v", err)
	}
	return s
}

// cmdOutputErr runs the command line in dir, returning its output
// and any error results.
//
// NOTE: cmdOutputErr must be used only to run commands that read state,
// not for commands that make changes. Commands that make changes
// should be run using runDirErr so that the -v and -n flags apply to them.
func cmdOutputErr(command string, args ...string) (string, error) {
	return cmdOutputDirErr(".", command, args...)
}

// cmdOutputDirErr runs the command line in dir, returning its output
// and any error results.
//
// NOTE: cmdOutputDirErr must be used only to run commands that read state,
// not for commands that make changes. Commands that make changes
// should be run using runDirErr so that the -v and -n flags apply to them.
func cmdOutputDirErr(dir, command string, args ...string) (string, error) {
	// NOTE: We only show these non-state-modifying commands with -v -v.
	// Otherwise things like 'git sync -v' show all our internal "find out about
	// the git repo" commands, which is confusing if you are just trying to find
	// out what git sync means.
	if *verbose > 1 {
		fmt.Fprintln(stderr(), commandString(command, args))
	}
	cmd := exec.Command(command, args...)
	if dir != "." {
		cmd.Dir = dir
	}
	b, err := cmd.CombinedOutput()
	return string(b), err
}

// trim is shorthand for strings.TrimSpace.
func trim(text string) string {
	return strings.TrimSpace(text)
}

// trimErr applies strings.TrimSpace to the result of cmdOutput(Dir)Err,
// passing the error along unmodified.
func trimErr(text string, err error) (string, error) {
	return strings.TrimSpace(text), err
}

// lines returns the lines in text.
func lines(text string) []string {
	out := strings.Split(text, "\n")
	// Split will include a "" after the last line. Remove it.
	if n := len(out) - 1; n >= 0 && out[n] == "" {
		out = out[:n]
	}
	return out
}

// nonBlankLines returns the non-blank lines in text.
func nonBlankLines(text string) []string {
	var out []string
	for _, s := range lines(text) {
		if strings.TrimSpace(s) != "" {
			out = append(out, s)
		}
	}
	return out
}

func commandString(command string, args []string) string {
	return strings.Join(append([]string{command}, args...), " ")
}

var dieTrap func()

func dief(format string, args ...interface{}) {
	printf(format, args...)
	die()
}

func die() {
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

var stdoutTrap, stderrTrap *bytes.Buffer

func stdout() io.Writer {
	if stdoutTrap != nil {
		return stdoutTrap
	}
	return os.Stdout
}

func stderr() io.Writer {
	if stderrTrap != nil {
		return stderrTrap
	}
	return os.Stderr
}

func printf(format string, args ...interface{}) {
	fmt.Fprintf(stderr(), "%s: %s\n", os.Args[0], fmt.Sprintf(format, args...))
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
