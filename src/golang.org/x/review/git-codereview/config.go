// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"
)

var (
	configRef    = "refs/remotes/origin/master:codereview.cfg"
	cachedConfig map[string]string
)

// Config returns the code review config.
// Configs consist of lines of the form "key: value".
// Lines beginning with # are comments.
// If there is no config, it returns an empty map.
// If the config is malformed, it dies.
func config() map[string]string {
	if cachedConfig != nil {
		return cachedConfig
	}
	raw, err := cmdOutputErr("git", "show", configRef)
	if err != nil {
		verbosef("%sfailed to load config from %q: %v", raw, configRef, err)
		cachedConfig = make(map[string]string)
		return cachedConfig
	}
	cachedConfig, err = parseConfig(raw)
	if err != nil {
		dief("%v", err)
	}
	return cachedConfig
}

func parseConfig(raw string) (map[string]string, error) {
	cfg := make(map[string]string)
	for _, line := range nonBlankLines(raw) {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			// comment line
			continue
		}
		fields := strings.SplitN(line, ":", 2)
		if len(fields) != 2 {
			return nil, fmt.Errorf("bad config line, expected 'key: value': %q", line)
		}
		cfg[strings.TrimSpace(fields[0])] = strings.TrimSpace(fields[1])
	}
	return cfg, nil
}
