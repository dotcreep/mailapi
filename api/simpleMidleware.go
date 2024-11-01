package api

import (
	"mailapi/env"
	"mailapi/utils"
	"net/http"
	"os"

	"golang.org/x/time/rate"
)

func Midleware(next http.HandlerFunc) http.HandlerFunc {
	limiter := rate.NewLimiter(1, 5)
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			utils.RespondJSON(w, http.StatusTooManyRequests, map[string]interface{}{
				"success": false,
				"message": "Too many requests",
			})
			return
		}
		env.Load(".env")
		secret := os.Getenv("SECRET_ACCESS")
		apiKey := r.Header.Get("X-Auth-Key")
		if apiKey != secret {
			utils.RespondJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"message": "Unauthorized",
			})
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}

}
