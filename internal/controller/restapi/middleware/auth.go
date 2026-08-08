package middleware

import (
	"context"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session"
)

type ctxKey string

const userIDKey ctxKey = "userID"

func Session() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return session.VerifySession(nil, func(w http.ResponseWriter, r *http.Request) {
			sess := session.GetSessionFromRequestContext(r.Context())
			if sess == nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)

				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, sess.GetUserID())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)

	return id, ok
}
