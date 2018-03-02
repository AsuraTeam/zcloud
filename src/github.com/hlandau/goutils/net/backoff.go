package net

import "math"
import "math/rand"
import "time"

var randr *rand.Rand

func init() {
	t := time.Now()
	s := rand.NewSource(t.Unix() ^ t.UnixNano())
	randr = rand.New(s)
}

// Expresses a backoff and retry specification.
//
// The nil value of this structure results in sensible defaults being used.
type Backoff struct {
	// The maximum number of attempts which may be made.
	// If this is 0, the number of attempts is unlimited.
	MaxTries int

	// The initial delay, in milliseconds. This is the delay used after the first
	// failed attempt.
	InitialDelay time.Duration // ms

	// The maximum delay, in milliseconds. This is the maximum delay between
	// attempts.
	MaxDelay time.Duration // ms

	// Determines when the maximum delay should be reached. If this is 5, the
	// maximum delay will be reached after 5 attempts have been made.
	MaxDelayAfterTries int

	// Positive float expressing the maximum factor by which the delay value may
	// be randomly inflated. e.g. specify 0.05 for a 5% variation. Set to zero to
	// disable jitter.
	Jitter float64

	// The current try. You should not need to set this yourself.
	CurrentTry int
}

// Initialises any nil field in Backoff with sensible defaults. You
// normally do not need to call this method yourself, as it will be called
// automatically.
func (rc *Backoff) InitDefaults() {
	if rc.InitialDelay == 0 {
		rc.InitialDelay = 5 * time.Second
	}
	if rc.MaxDelay == 0 {
		rc.MaxDelay = 120 * time.Second
	}
	if rc.MaxDelayAfterTries == 0 {
		rc.MaxDelayAfterTries = 10
	}
}

// Gets the next delay in milliseconds and increments the internal try counter.
func (rc *Backoff) NextDelay() time.Duration {
	rc.InitDefaults()

	if rc.MaxTries != 0 && rc.CurrentTry >= rc.MaxTries {
		return time.Duration(0)
	}

	initialDelay := float64(rc.InitialDelay)
	maxDelay := float64(rc.MaxDelay)
	maxDelayAfterTries := float64(rc.MaxDelayAfterTries)
	currentTry := float64(rc.CurrentTry)

	// [from backoff.c]
	k := math.Log2(maxDelay/initialDelay) / maxDelayAfterTries
	d := time.Duration(initialDelay * math.Exp2(currentTry*k))
	rc.CurrentTry++

	if d > rc.MaxDelay {
		d = rc.MaxDelay
	}

	if rc.Jitter != 0 {
		f := (randr.Float64() - 0.5) * 2 // random value in range [-1,1)
		d = time.Duration(float64(d) * (1 + rc.Jitter*f))
	}

	return d
}

// Sleep for the duration returned by NextDelay().
func (rc *Backoff) Sleep() bool {
	d := rc.NextDelay()
	if d != 0 {
		time.Sleep(d)
		return true
	}
	return false
}

// Sets the internal try counter to zero; the next delay returned will be
// InitialDelay again.
func (rc *Backoff) Reset() {
	rc.CurrentTry = 0
}
