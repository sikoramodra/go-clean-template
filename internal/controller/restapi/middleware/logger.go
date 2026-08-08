package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/evrone/go-clean-template/pkg/logger"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (sw *statusWriter) WriteHeader(status int) {
	sw.status = status
	sw.ResponseWriter.WriteHeader(status)
}

func (sw *statusWriter) Write(b []byte) (int, error) {
	size, err := sw.ResponseWriter.Write(b)
	sw.size += size

	return size, err
}

func buildRequestMessage(r *http.Request, status, size int) string {
	var result strings.Builder

	result.WriteString(r.RemoteAddr)
	result.WriteString(" - ")
	result.WriteString(r.Method)
	result.WriteString(" ")
	result.WriteString(r.URL.RequestURI())
	result.WriteString(" - ")
	result.WriteString(strconv.Itoa(status))
	result.WriteString(" ")
	result.WriteString(strconv.Itoa(size))

	return result.String()
}

func Logger(l logger.Interface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := &statusWriter{ResponseWriter: w}

			next.ServeHTTP(sw, r)

			l.Info("%s", buildRequestMessage(r, sw.status, sw.size))
		})
	}
}
