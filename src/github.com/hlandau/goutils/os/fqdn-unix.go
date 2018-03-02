// +build !windows,!plan9

package os

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var reValidHostname = regexp.MustCompile(`^([0-9a-zA-Z-]+\.)[0-9a-zA-Z-]+$`)

func machineFQDN() (string, error) {
	hn, err := os.Hostname()
	if err != nil {
		return "", err
	}

	hn = strings.TrimSuffix(hn, ".")
	n := strings.Count(hn, ".")
	if n > 0 {
		return hn, nil
	}

	buf := bytes.Buffer{}
	cmd := exec.Command("hostname", "-f")
	cmd.Stdout = &buf

	err = cmd.Run()
	if err != nil {
		return hn, err
	}

	fqdn := strings.TrimSuffix(strings.TrimSpace(buf.String()), ".")
	if !reValidHostname.MatchString(fqdn) || strings.Count(fqdn, ".") < n {
		return hn, nil
	}

	return fqdn, nil
}
