// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"testing"
)

var authTests = []struct {
	netrc       string
	cookiefile  string
	user        string
	password    string
	cookieName  string
	cookieValue string
	died        bool
}{
	{
		died: true,
	},
	{
		netrc:    "machine go.googlesource.com login u1 password pw\n",
		user:     "u1",
		password: "pw",
	},
	{
		cookiefile: "go.googlesource.com	TRUE	/	TRUE	2147483647	o2	git-u2=pw\n",
		cookieName:  "o2",
		cookieValue: "git-u2=pw",
	},
	{
		cookiefile: ".googlesource.com	TRUE	/	TRUE	2147483647	o3	git-u3=pw\n",
		cookieName:  "o3",
		cookieValue: "git-u3=pw",
	},
	{
		cookiefile: ".googlesource.com	TRUE	/	TRUE	2147483647	o4	WRONG\n" +
			"go.googlesource.com	TRUE	/	TRUE	2147483647	o4	git-u4=pw\n",
		cookieName:  "o4",
		cookieValue: "git-u4=pw",
	},
	{
		cookiefile: "go.googlesource.com	TRUE	/	TRUE	2147483647	o5	git-u5=pw\n" +
			".googlesource.com	TRUE	/	TRUE	2147483647	o5	WRONG\n",
		cookieName:  "o5",
		cookieValue: "git-u5=pw",
	},
	{
		netrc:      "machine go.googlesource.com login u6 password pw\n",
		cookiefile: "BOGUS",
		user:       "u6",
		password:   "pw",
	},
	{
		netrc: "BOGUS",
		cookiefile: "go.googlesource.com	TRUE	/	TRUE	2147483647	o7	git-u7=pw\n",
		cookieName:  "o7",
		cookieValue: "git-u7=pw",
	},
	{
		netrc:      "machine go.googlesource.com login u8 password pw\n",
		cookiefile: "MISSING",
		user:       "u8",
		password:   "pw",
	},
}

func TestLoadAuth(t *testing.T) {
	gt := newGitTest(t)
	defer gt.done()
	gt.work(t)

	defer os.Setenv("HOME", os.Getenv("HOME"))
	os.Setenv("HOME", gt.client)
	trun(t, gt.client, "git", "config", "remote.origin.url", "https://go.googlesource.com/go")

	for i, tt := range authTests {
		t.Logf("#%d", i)
		auth.user = ""
		auth.password = ""
		auth.cookieName = ""
		auth.cookieValue = ""
		trun(t, gt.client, "git", "config", "http.cookiefile", "XXX")
		trun(t, gt.client, "git", "config", "--unset", "http.cookiefile")

		remove(t, gt.client+"/.netrc")
		remove(t, gt.client+"/.cookies")
		if tt.netrc != "" {
			write(t, gt.client+"/.netrc", tt.netrc)
		}
		if tt.cookiefile != "" {
			if tt.cookiefile != "MISSING" {
				write(t, gt.client+"/.cookies", tt.cookiefile)
			}
			trun(t, gt.client, "git", "config", "http.cookiefile", gt.client+"/.cookies")
		}

		// Run command via testMain to trap stdout, stderr, death.
		// mail -n will load auth info for us.
		if tt.died {
			testMainDied(t, "test-loadAuth")
		} else {
			testMain(t, "test-loadAuth")
		}

		if auth.user != tt.user || auth.password != tt.password {
			t.Errorf("#%d: have user, password = %q, %q, want %q, %q", i, auth.user, auth.password, tt.user, tt.password)
		}
		if auth.cookieName != tt.cookieName || auth.cookieValue != tt.cookieValue {
			t.Errorf("#%d: have cookie name, value = %q, %q, want %q, %q", i, auth.cookieName, auth.cookieValue, tt.cookieName, tt.cookieValue)
		}
	}
}
