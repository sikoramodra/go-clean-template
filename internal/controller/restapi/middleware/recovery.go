package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/sikoramodra/go-clean-template/pkg/logger"
)

func buildPanicMessage(r *http.Request, err any) string {
	var result strings.Builder

	result.WriteString(r.RemoteAddr)
	result.WriteString(" - ")
	result.WriteString(r.Method)
	result.WriteString(" ")
	result.WriteString(r.URL.RequestURI())
	result.WriteString(" PANIC DETECTED: ")
	fmt.Fprintf(&result, "%v\n%s\n", err, debug.Stack())

	return result.String()
}

func Recovery(l logger.Interface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					l.Error(buildPanicMessage(r, err))
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
