package net_test

import "testing"
import "github.com/hlandau/goutils/net"
import "time"

func TestBackoff(t *testing.T) {
	b := net.Backoff{}

	eq := func(d time.Duration, ms int) {
		a := int(d / time.Millisecond)
		if a != ms {
			t.Errorf("Backoff #%d:  %v should be %v", b.CurrentTry, a, ms)
		}
	}

	eq(b.NextDelay(), 5000)
	eq(b.NextDelay(), 6870)
	eq(b.NextDelay(), 9440)
	eq(b.NextDelay(), 12972)
	eq(b.NextDelay(), 17826)
	eq(b.NextDelay(), 24494)
	eq(b.NextDelay(), 33658)
	eq(b.NextDelay(), 46250)
	eq(b.NextDelay(), 63553)
	eq(b.NextDelay(), 87329)
	eq(b.NextDelay(), 120000)
}
