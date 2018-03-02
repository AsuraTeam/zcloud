// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	sourceHost   = ".googlesource.com"
	reviewSuffix = "-review"
)

var notFound = errors.New("not found")

func getChange(origin, id string) (*Change, error) {
	u, err := url.Parse(origin)
	if err != nil {
		return nil, fmt.Errorf("parsing origin URL: %v", err)
	}
	if !strings.HasSuffix(u.Host, sourceHost) {
		return nil, fmt.Errorf("origin URL not on %v", sourceHost)
	}
	u.Host = strings.TrimSuffix(u.Host, sourceHost) + reviewSuffix + sourceHost
	u.Path = "/changes/" + id
	r, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.StatusCode == http.StatusNotFound {
		return nil, notFound
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %v", r.Status)
	}
	br := bufio.NewReader(r.Body)
	br.ReadSlice('\n') // throw away first line
	var c Change
	if err := json.NewDecoder(br).Decode(&c); err != nil {
		return nil, err
	}
	u.Path = fmt.Sprintf("/%v", c.Number)
	c.URL = u.String()
	return &c, nil
}

type Change struct {
	ChangeId string `json:"change_id"`
	Subject  string
	Status   string
	Number   int `json:"_number"`
	Messages []*Message
	URL      string
}

type Message struct {
	Author struct {
		Name string
	}
	Message string
}
