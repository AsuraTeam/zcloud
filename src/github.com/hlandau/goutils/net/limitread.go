package net

import (
	"errors"
	"io"
)

type limitReader struct {
	r         io.Reader
	remaining int
}

// Error returned if the number of bytes read from a LimitReader exceeds the
// limit.
var ErrLimitExceeded = errors.New("reader size limit exceeded")

func (lr *limitReader) Read(b []byte) (int, error) {
	n, err := lr.r.Read(b)
	lr.remaining = lr.remaining - n
	if lr.remaining < 0 {
		err = ErrLimitExceeded
	}

	return n, err
}

// Creates a limit reader. This is similar to io.LimitReader, except that if
// more than limit bytes are read from the stream, an error is returned rather
// than EOF. Thus, accidental truncation of oversized byte streams is avoided
// in favour of a hard error.
//
// Returns ErrLimitExceeded once more than limit bytes are read. The read
// operation still occurs.
func LimitReader(r io.Reader, limit int) io.Reader {
	return &limitReader{
		r:         r,
		remaining: limit,
	}
}
