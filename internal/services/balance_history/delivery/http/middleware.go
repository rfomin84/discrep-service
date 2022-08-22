package balance_history

import (
	"github.com/rfomin84/discrep-service/config"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := config.GetConfig()
		bearerToken := r.Header.Get("Authorization")
		if len(bearerToken) == 0 {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		token := strings.Split(bearerToken, " ")
		if len(token) != 2 || token[0] != "Bearer" {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		if token[1] != cfg.GetString("API_TOKEN") {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
