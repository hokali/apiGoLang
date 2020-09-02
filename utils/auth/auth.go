package auth

import (
	"net/http"

	okta "../../controllers/okta"
)

// JwtVerify Middleware function
func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tk := okta.JwtMiddlewareChk(w, r)

		if tk != true {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
