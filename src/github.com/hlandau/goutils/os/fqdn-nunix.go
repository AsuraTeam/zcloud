// +build windows plan9

package os

func machineFQDN() (string, error) {
	return os.Hostname()
}
