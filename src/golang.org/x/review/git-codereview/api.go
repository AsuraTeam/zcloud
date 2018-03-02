// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

// auth holds cached data about authentication to Gerrit.
var auth struct {
	host    string // "go.googlesource.com"
	url     string // "https://go-review.googlesource.com"
	project string // "go", "tools", "crypto", etc

	// Authentication information.
	// Either cookie name + value from git cookie file
	// or username and password from .netrc.
	cookieName  string
	cookieValue string
	user        string
	password    string
}

// loadGerritOrigin loads the Gerrit host name from the origin remote.
// If the origin remote does not appear to be a Gerrit server
// (is missing, is GitHub, is not https, has too many path elements),
// loadGerritOrigin dies.
func loadGerritOrigin() {
	if auth.host != "" {
		return
	}

	// Gerrit must be set, either explicitly via the code review config or
	// implicitly as Git's origin remote.
	origin := config()["gerrit"]
	if origin == "" {
		origin = trim(cmdOutput("git", "config", "remote.origin.url"))
	}

	if strings.Contains(origin, "github.com") {
		dief("git origin must be a Gerrit host, not GitHub: %s", origin)
	}

	if !strings.HasPrefix(origin, "https://") {
		dief("git origin must be an https:// URL: %s", origin)
	}
	// https:// prefix and then one slash between host and top-level name
	if strings.Count(origin, "/") != 3 {
		dief("git origin is malformed: %s", origin)
	}
	host := origin[len("https://"):strings.LastIndex(origin, "/")]

	// In the case of Google's Gerrit, host is go.googlesource.com
	// and apiURL uses go-review.googlesource.com, but the Gerrit
	// setup instructions do not write down a cookie explicitly for
	// go-review.googlesource.com, so we look for the non-review
	// host name instead.
	url := origin
	if i := strings.Index(url, ".googlesource.com"); i >= 0 {
		url = url[:i] + "-review" + url[i:]
	}
	i := strings.LastIndex(url, "/")
	url, project := url[:i], url[i+1:]

	auth.host = host
	auth.url = url
	auth.project = project
}

// loadAuth loads the authentication tokens for making API calls to
// the Gerrit origin host.
func loadAuth() {
	if auth.user != "" || auth.cookieName != "" {
		return
	}

	loadGerritOrigin()

	// First look in Git's http.cookiefile, which is where Gerrit
	// now tells users to store this information.
	if cookieFile, _ := trimErr(cmdOutputErr("git", "config", "http.cookiefile")); cookieFile != "" {
		data, _ := ioutil.ReadFile(cookieFile)
		maxMatch := -1
		for _, line := range lines(string(data)) {
			f := strings.Split(line, "\t")
			if len(f) >= 7 && (f[0] == auth.host || strings.HasPrefix(f[0], ".") && strings.HasSuffix(auth.host, f[0])) {
				if len(f[0]) > maxMatch {
					auth.cookieName = f[5]
					auth.cookieValue = f[6]
					maxMatch = len(f[0])
				}
			}
		}
		if maxMatch > 0 {
			return
		}
	}

	// If not there, then look in $HOME/.netrc, which is where Gerrit
	// used to tell users to store the information, until the passwords
	// got so long that old versions of curl couldn't handle them.
	data, _ := ioutil.ReadFile(os.Getenv("HOME") + "/.netrc")
	for _, line := range lines(string(data)) {
		if i := strings.Index(line, "#"); i >= 0 {
			line = line[:i]
		}
		f := strings.Fields(line)
		if len(f) >= 6 && f[0] == "machine" && f[1] == auth.host && f[2] == "login" && f[4] == "password" {
			auth.user = f[3]
			auth.password = f[5]
			return
		}
	}

	dief("cannot find authentication info for %s", auth.host)
}

// gerritError is an HTTP error response served by Gerrit.
type gerritError struct {
	url        string
	statusCode int
	status     string
	body       string
}

func (e *gerritError) Error() string {
	if e.statusCode == http.StatusNotFound {
		return "change not found on Gerrit server"
	}

	extra := strings.TrimSpace(e.body)
	if extra != "" {
		extra = ": " + extra
	}
	return fmt.Sprintf("%s%s", e.status, extra)
}

// gerritAPI executes a GET or POST request to a Gerrit API endpoint.
// It uses GET when requestBody is nil, otherwise POST. If target != nil,
// gerritAPI expects to get a 200 response with a body consisting of an
// anti-xss line (]})' or some such) followed by JSON.
// If requestBody != nil, gerritAPI sets the Content-Type to application/json.
func gerritAPI(path string, requestBody []byte, target interface{}) error {
	// Strictly speaking, we might be able to use unauthenticated
	// access, by removing the /a/ from the URL, but that assumes
	// that all the information we care about is publicly visible.
	// Using authentication makes it possible for this to work with
	// non-public CLs or Gerrit hosts too.
	loadAuth()

	if !strings.HasPrefix(path, "/") {
		dief("internal error: gerritAPI called with malformed path")
	}

	url := auth.url + path
	method := "GET"
	var reader io.Reader
	if requestBody != nil {
		method = "POST"
		reader = bytes.NewReader(requestBody)
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return err
	}
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if auth.cookieName != "" {
		req.AddCookie(&http.Cookie{
			Name:  auth.cookieName,
			Value: auth.cookieValue,
		})
	} else {
		req.SetBasicAuth(auth.user, auth.password)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return &gerritError{url, resp.StatusCode, resp.Status, string(body)}
	}

	if target != nil {
		i := bytes.IndexByte(body, '\n')
		if i < 0 {
			return fmt.Errorf("%s: malformed json response", url)
		}
		body = body[i:]
		if err := json.Unmarshal(body, target); err != nil {
			return fmt.Errorf("%s: malformed json response", url)
		}
	}
	return nil
}

// fullChangeID returns the unambigous Gerrit change ID for the commit c on branch b.
// The retruned ID has the form project~originbranch~Ihexhexhexhexhex.
// See https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#change-id for details.
func fullChangeID(b *Branch, c *Commit) string {
	loadGerritOrigin()
	return auth.project + "~" + strings.TrimPrefix(b.OriginBranch(), "origin/") + "~" + c.ChangeID
}

// readGerritChange reads the metadata about a change from the Gerrit server.
// The changeID should use the syntax project~originbranch~Ihexhexhexhexhex returned
// by fullChangeID. Using only Ihexhexhexhexhex will work provided it uniquely identifies
// a single change on the server.
// The changeID can have additional query parameters appended to it, as in "normalid?o=LABELS".
// See https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#change-id for details.
func readGerritChange(changeID string) (*GerritChange, error) {
	var c GerritChange
	err := gerritAPI("/a/changes/"+changeID, nil, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// GerritChange is the JSON struct returned by a Gerrit CL query.
type GerritChange struct {
	ID              string
	Project         string
	Branch          string
	ChangeId        string `json:"change_id"`
	Subject         string
	Status          string
	Created         string
	Updated         string
	Mergeable       bool
	Insertions      int
	Deletions       int
	Number          int `json:"_number"`
	Owner           *GerritAccount
	Labels          map[string]*GerritLabel
	CurrentRevision string `json:"current_revision"`
	Revisions       map[string]*GerritRevision
	Messages        []*GerritMessage
}

// LabelNames returns the label names for the change, in lexicographic order.
func (g *GerritChange) LabelNames() []string {
	var names []string
	for name := range g.Labels {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GerritMessage is the JSON struct for a Gerrit MessageInfo.
type GerritMessage struct {
	Author struct {
		Name string
	}
	Message string
}

// GerritLabel is the JSON struct for a Gerrit LabelInfo.
type GerritLabel struct {
	Optional bool
	Blocking bool
	Approved *GerritAccount
	Rejected *GerritAccount
	All      []*GerritApproval
}

// GerritAccount is the JSON struct for a Gerrit AccountInfo.
type GerritAccount struct {
	ID       int `json:"_account_id"`
	Name     string
	Email    string
	Username string
}

// GerritApproval is the JSON struct for a Gerrit ApprovalInfo.
type GerritApproval struct {
	GerritAccount
	Value int
	Date  string
}

// GerritRevision is the JSON struct for a Gerrit RevisionInfo.
type GerritRevision struct {
	Number int `json:"_number"`
	Ref    string
	Fetch  map[string]*GerritFetch
}

// GerritFetch is the JSON struct for a Gerrit FetchInfo
type GerritFetch struct {
	URL string
	Ref string
}
