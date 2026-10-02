package middleware

import (
	"io"
	"net/http"
)

func httpMaxBytesReader(w http.ResponseWriter, r io.ReadCloser, n int64) io.ReadCloser {
	return http.MaxBytesReader(w, r, n)
}
