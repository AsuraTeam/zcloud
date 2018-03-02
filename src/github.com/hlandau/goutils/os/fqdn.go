package os

// Returns the machine hostname as a fully-qualified domain name.
func MachineFQDN() (string, error) {
	return machineFQDN()
}
