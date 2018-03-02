// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"reflect"
	"testing"
)

func TestParseConfig(t *testing.T) {
	cases := []struct {
		raw     string
		want    map[string]string
		wanterr bool
	}{
		{raw: "", want: map[string]string{}},
		{raw: "issuerepo: golang/go", want: map[string]string{"issuerepo": "golang/go"}},
		{raw: "# comment", want: map[string]string{}},
		{raw: "# comment\n  k  :   v   \n# comment 2\n\n k2:v2\n", want: map[string]string{"k": "v", "k2": "v2"}},
	}

	for _, tt := range cases {
		cfg, err := parseConfig(tt.raw)
		if err != nil != tt.wanterr {
			t.Errorf("parse(%q) error: %v", tt.raw, err)
			continue
		}
		if !reflect.DeepEqual(cfg, tt.want) {
			t.Errorf("parse(%q)=%v want %v", tt.raw, cfg, tt.want)
		}
	}
}
