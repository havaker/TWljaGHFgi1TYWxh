package middleware

import (
	"errors"
	"io"
	"net/http"
)

type limitedReadCloser struct {
	limit, read int
	reader      io.ReadCloser
}

// ErrBodyLimitExceeded is an error returned when request body is too large
var ErrBodyLimitExceeded = errors.New("middleware: body limit exceeded")

// BodyLimit is a middleware that discards requests with Content-Length larger
// then given limit, and wraps request.Body with reader that returns error
// when reading above limit
func BodyLimit(limit int, method string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			if r.Method == method {
				if r.ContentLength > int64(limit) {
					w.WriteHeader(http.StatusRequestEntityTooLarge)
					return
				}

				r.Body = &limitedReadCloser{limit: limit, reader: r.Body}
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func (l *limitedReadCloser) Read(b []byte) (n int, err error) {
	n, err = l.reader.Read(b)
	l.read += n
	if l.read > l.limit {
		return n, ErrBodyLimitExceeded
	}
	return
}

func (l *limitedReadCloser) Close() error {
	return l.reader.Close()
}
