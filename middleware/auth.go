package middleware

import (
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/flume-cloud-services/file-storage/controllers"
)

// Middleware allows to use multiple middleware on a single route
func Middleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}

// AuthMiddleware manage auth via JWT
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		admin := os.Getenv("FLUME_FILE_STORAGE_ADMIN")
		if len(admin) == 0 {
			admin = "admin"
		}

		jwtKey := []byte(os.Getenv("FLUME_FILE_STORAGE_SECRET"))
		if len(jwtKey) == 0 {
			jwtKey = []byte("this_is_a_secret_token")
		}

		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := c.Value
		claims := &controllers.Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
