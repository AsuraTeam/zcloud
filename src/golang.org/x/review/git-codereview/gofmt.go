// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

var gofmtList bool

func cmdGofmt(args []string) {
	flags.BoolVar(&gofmtList, "l", false, "list files that need to be formatted")
	flags.Parse(args)
	if len(flag.Args()) > 0 {
		fmt.Fprintf(stderr(), "Usage: %s gofmt %s [-l]\n", os.Args[0], globalFlags)
		os.Exit(2)
	}

	f := gofmtCommand
	if !gofmtList {
		f |= gofmtWrite
	}

	files, stderr := runGofmt(f)
	if gofmtList {
		w := stdout()
		for _, file := range files {
			fmt.Fprintf(w, "%s\n", file)
		}
	}
	if stderr != "" {
		dief("gofmt reported errors:\n\t%s", strings.Replace(strings.TrimSpace(stderr), "\n", "\n\t", -1))
	}
}

const (
	gofmtPreCommit = 1 << iota
	gofmtCommand
	gofmtWrite
)

// runGofmt runs the external gofmt command over modified files.
//
// The definition of "modified files" depends on the bit flags.
// If gofmtPreCommit is set, then runGofmt considers *.go files that
// differ between the index (staging area) and the branchpoint
// (the latest commit before the branch diverged from upstream).
// If gofmtCommand is set, then runGofmt considers all those files
// in addition to files with unstaged modifications.
// It never considers untracked files.
//
// As a special case for the main repo (but applied everywhere)
// *.go files under a top-level test directory are excluded from the
// formatting requirement, except run.go and those in test/bench/.
//
// If gofmtWrite is set (only with gofmtCommand, meaning this is 'git gofmt'),
// runGofmt replaces the original files with their formatted equivalents.
// Git makes this difficult. In general the file in the working tree
// (the local file system) can have unstaged changes that make it different
// from the equivalent file in the index. To help pass the precommit hook,
// 'git gofmt'  must make it easy to update the files in the index.
// One option is to run gofmt on all the files of the same name in the
// working tree and leave it to the user to sort out what should be staged
// back into the index. Another is to refuse to reformat files for which
// different versions exist in the index vs the working tree. Both of these
// options are unsatisfying: they foist busy work onto the user,
// and it's exactly the kind of busy work that a program is best for.
// Instead, when runGofmt finds files in the index that need
// reformatting, it reformats them there, bypassing the working tree.
// It also reformats files in the working tree that need reformatting.
// For both, only files modified since the branchpoint are considered.
// The result should be that both index and working tree get formatted
// correctly and diffs between the two remain meaningful (free of
// formatting distractions). Modifying files in the index directly may
// surprise Git users, but it seems the best of a set of bad choices, and
// of course those users can choose not to use 'git gofmt'.
// This is a bit more work than the other git commands do, which is
// a little worrying, but the choice being made has the nice property
// that if 'git gofmt' is interrupted, a second 'git gofmt' will put things into
// the same state the first would have.
//
// runGofmt returns a list of files that need (or needed) reformatting.
// If gofmtPreCommit is set, the names always refer to files in the index.
// If gofmtCommand is set, then a name without a suffix (see below)
// refers to both the copy in the index and the copy in the working tree
// and implies that the two copies are identical. Otherwise, in the case
// that the index and working tree differ, the file name will have an explicit
// " (staged)" or " (unstaged)" suffix saying which is meant.
//
// runGofmt also returns any standard error output from gofmt,
// usually indicating syntax errors in the Go source files.
// If gofmtCommand is set, syntax errors in index files that do not match
// the working tree show a " (staged)" suffix after the file name.
// The errors never use the " (unstaged)" suffix, in order to keep
// references to the local file system in the standard file:line form.
func runGofmt(flags int) (files []string, stderrText string) {
	pwd, err := os.Getwd()
	if err != nil {
		dief("%v", err)
	}
	pwd = filepath.Clean(pwd) // convert to host \ syntax
	if !strings.HasSuffix(pwd, string(filepath.Separator)) {
		pwd += string(filepath.Separator)
	}

	b := CurrentBranch()
	repo := repoRoot()
	if !strings.HasSuffix(repo, string(filepath.Separator)) {
		repo += string(filepath.Separator)
	}

	// Find files modified in the index compared to the branchpoint.
	branchpt := b.Branchpoint()
	if strings.Contains(cmdOutput("git", "branch", "-r", "--contains", b.FullName()), "origin/") {
		// This is a branch tag move, not an actual change.
		// Use HEAD as branch point, so nothing will appear changed.
		// We don't want to think about gofmt on already published
		// commits.
		branchpt = "HEAD"
	}
	indexFiles := addRoot(repo, filter(gofmtRequired, nonBlankLines(cmdOutput("git", "diff", "--name-only", "--diff-filter=ACM", "--cached", branchpt, "--"))))
	localFiles := addRoot(repo, filter(gofmtRequired, nonBlankLines(cmdOutput("git", "diff", "--name-only", "--diff-filter=ACM"))))
	localFilesMap := stringMap(localFiles)
	isUnstaged := func(file string) bool {
		return localFilesMap[file]
	}

	if len(indexFiles) == 0 && ((flags&gofmtCommand) == 0 || len(localFiles) == 0) {
		return
	}

	// Determine which files have unstaged changes and are therefore
	// different from their index versions. For those, the index version must
	// be copied into a temporary file in the local file system.
	needTemp := filter(isUnstaged, indexFiles)

	// Map between temporary file name and place in file tree where
	// file would be checked out (if not for the unstaged changes).
	tempToFile := map[string]string{}
	fileToTemp := map[string]string{}
	cleanup := func() {} // call before dying (defer won't run)
	if len(needTemp) > 0 {
		// Ask Git to copy the index versions into temporary files.
		// Git stores the temporary files, named .merge_*, in the repo root.
		// Unlike the Git commands above, the non-temp file names printed
		// here are relative to the current directory, not the repo root.

		// git checkout-index --temp is broken on windows. Running this command:
		//
		// git checkout-index --temp -- bad-bad-bad2.go bad-bad-broken.go bad-bad-good.go bad-bad2-bad.go bad-bad2-broken.go bad-bad2-good.go bad-broken-bad.go bad-broken-bad2.go bad-broken-good.go bad-good-bad.go bad-good-bad2.go bad-good-broken.go bad2-bad-bad2.go bad2-bad-broken.go bad2-bad-good.go bad2-bad2-bad.go bad2-bad2-broken.go bad2-bad2-good.go bad2-broken-bad.go bad2-broken-bad2.go bad2-broken-good.go bad2-good-bad.go bad2-good-bad2.go bad2-good-broken.go broken-bad-bad2.go broken-bad-broken.go broken-bad-good.go broken-bad2-bad.go broken-bad2-broken.go broken-bad2-good.go
		//
		// produces this output
		//
		// .merge_file_a05448      bad-bad-bad2.go
		// .merge_file_b05448      bad-bad-broken.go
		// .merge_file_c05448      bad-bad-good.go
		// .merge_file_d05448      bad-bad2-bad.go
		// .merge_file_e05448      bad-bad2-broken.go
		// .merge_file_f05448      bad-bad2-good.go
		// .merge_file_g05448      bad-broken-bad.go
		// .merge_file_h05448      bad-broken-bad2.go
		// .merge_file_i05448      bad-broken-good.go
		// .merge_file_j05448      bad-good-bad.go
		// .merge_file_k05448      bad-good-bad2.go
		// .merge_file_l05448      bad-good-broken.go
		// .merge_file_m05448      bad2-bad-bad2.go
		// .merge_file_n05448      bad2-bad-broken.go
		// .merge_file_o05448      bad2-bad-good.go
		// .merge_file_p05448      bad2-bad2-bad.go
		// .merge_file_q05448      bad2-bad2-broken.go
		// .merge_file_r05448      bad2-bad2-good.go
		// .merge_file_s05448      bad2-broken-bad.go
		// .merge_file_t05448      bad2-broken-bad2.go
		// .merge_file_u05448      bad2-broken-good.go
		// .merge_file_v05448      bad2-good-bad.go
		// .merge_file_w05448      bad2-good-bad2.go
		// .merge_file_x05448      bad2-good-broken.go
		// .merge_file_y05448      broken-bad-bad2.go
		// .merge_file_z05448      broken-bad-broken.go
		// error: unable to create file .merge_file_XXXXXX (No error)
		// .merge_file_XXXXXX      broken-bad-good.go
		// error: unable to create file .merge_file_XXXXXX (No error)
		// .merge_file_XXXXXX      broken-bad2-bad.go
		// error: unable to create file .merge_file_XXXXXX (No error)
		// .merge_file_XXXXXX      broken-bad2-broken.go
		// error: unable to create file .merge_file_XXXXXX (No error)
		// .merge_file_XXXXXX      broken-bad2-good.go
		//
		// so limit the number of file arguments to 25.
		for len(needTemp) > 0 {
			n := len(needTemp)
			if n > 25 {
				n = 25
			}
			args := []string{"checkout-index", "--temp", "--"}
			args = append(args, needTemp[:n]...)
			// Until Git 2.3.0, git checkout-index --temp is broken if not run in the repo root.
			// Work around by running in the repo root.
			// http://article.gmane.org/gmane.comp.version-control.git/261739
			// https://github.com/git/git/commit/74c4de5
			for _, line := range nonBlankLines(cmdOutputDir(repo, "git", args...)) {
				i := strings.Index(line, "\t")
				if i < 0 {
					continue
				}
				temp, file := line[:i], line[i+1:]
				temp = filepath.Join(repo, temp)
				file = filepath.Join(repo, file)
				tempToFile[temp] = file
				fileToTemp[file] = temp
			}
			needTemp = needTemp[n:]
		}
		cleanup = func() {
			for temp := range tempToFile {
				os.Remove(temp)
			}
			tempToFile = nil
		}
		defer cleanup()
	}
	dief := func(format string, args ...interface{}) {
		cleanup()
		dief(format, args...) // calling top-level dief function
	}

	// Run gofmt to find out which files need reformatting;
	// if gofmtWrite is set, reformat them in place.
	// For references to local files, remove leading pwd if present
	// to make relative to current directory.
	// Temp files and local-only files stay as absolute paths for easy matching in output.
	args := []string{"-l"}
	if flags&gofmtWrite != 0 {
		args = append(args, "-w")
	}
	for _, file := range indexFiles {
		if isUnstaged(file) {
			args = append(args, fileToTemp[file])
		} else {
			args = append(args, strings.TrimPrefix(file, pwd))
		}
	}
	if flags&gofmtCommand != 0 {
		for _, file := range localFiles {
			args = append(args, file)
		}
	}

	if *verbose > 1 {
		fmt.Fprintln(stderr(), commandString("gofmt", args))
	}
	cmd := exec.Command("gofmt", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	if stderr.Len() == 0 && err != nil {
		// Error but no stderr: usually can't find gofmt.
		dief("invoking gofmt: %v", err)
	}

	// Build file list.
	files = lines(stdout.String())

	// Restage files that need to be restaged.
	if flags&gofmtWrite != 0 {
		add := []string{"add"}
		write := []string{"hash-object", "-w", "--"}
		updateIndex := []string{}
		for _, file := range files {
			if real := tempToFile[file]; real != "" {
				write = append(write, file)
				updateIndex = append(updateIndex, strings.TrimPrefix(real, repo))
			} else if !isUnstaged(file) {
				add = append(add, file)
			}
		}
		if len(add) > 1 {
			run("git", add...)
		}
		if len(updateIndex) > 0 {
			hashes := nonBlankLines(cmdOutput("git", write...))
			if len(hashes) != len(write)-3 {
				dief("git hash-object -w did not write expected number of objects")
			}
			var buf bytes.Buffer
			for i, name := range updateIndex {
				fmt.Fprintf(&buf, "100644 %s\t%s\n", hashes[i], name)
			}
			verbosef("git update-index --index-info")
			cmd := exec.Command("git", "update-index", "--index-info")
			cmd.Stdin = &buf
			out, err := cmd.CombinedOutput()
			if err != nil {
				dief("git update-index: %v\n%s", err, out)
			}
		}
	}

	// Remap temp files back to original names for caller.
	for i, file := range files {
		if real := tempToFile[file]; real != "" {
			if flags&gofmtCommand != 0 {
				real += " (staged)"
			}
			files[i] = strings.TrimPrefix(real, pwd)
		} else if isUnstaged(file) {
			files[i] = strings.TrimPrefix(file+" (unstaged)", pwd)
		}
	}

	// Rewrite temp names in stderr, and shorten local file names.
	// No suffix added for local file names (see comment above).
	text := "\n" + stderr.String()
	for temp, file := range tempToFile {
		if flags&gofmtCommand != 0 {
			file += " (staged)"
		}
		text = strings.Replace(text, "\n"+temp+":", "\n"+strings.TrimPrefix(file, pwd)+":", -1)
	}
	for _, file := range localFiles {
		text = strings.Replace(text, "\n"+file+":", "\n"+strings.TrimPrefix(file, pwd)+":", -1)
	}
	text = text[1:]

	sort.Strings(files)
	return files, text
}

// gofmtRequired reports whether the specified file should be checked
// for gofmt'dness by the pre-commit hook.
// The file name is relative to the repo root.
func gofmtRequired(file string) bool {
	// TODO: Consider putting this policy into codereview.cfg.
	if !strings.HasSuffix(file, ".go") {
		return false
	}
	if !strings.HasPrefix(file, "test/") {
		return true
	}
	return strings.HasPrefix(file, "test/bench/") || file == "test/run.go"
}

// stringMap returns a map m such that m[s] == true if s was in the original list.
func stringMap(list []string) map[string]bool {
	m := map[string]bool{}
	for _, x := range list {
		m[x] = true
	}
	return m
}

// filter returns the elements in list satisfying f.
func filter(f func(string) bool, list []string) []string {
	var out []string
	for _, x := range list {
		if f(x) {
			out = append(out, x)
		}
	}
	return out
}

func addRoot(root string, list []string) []string {
	var out []string
	for _, x := range list {
		out = append(out, filepath.Join(root, x))
	}
	return out
}
