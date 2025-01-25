package api

import (
	"net/http"

	"github.com/dotcreep/mailapi/internal/utils"
	"golang.org/x/time/rate"
)

func Midleware(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(1, 5)
	Json := utils.Json{}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			Json.NewResponse(false, w, nil, "Too many requests", http.StatusTooManyRequests, nil)
			return
		}
		cfg, err := utils.OpenYAML()
		if err != nil {
			Json.NewResponse(false, w, nil, "internal server error", http.StatusInternalServerError, err.Error())
			return
		}
		secret := cfg.Config.Token
		apiKey := r.Header.Get("X-Auth-Key")
		if apiKey != secret {
			Json.NewResponse(false, w, nil, "Unauthorized", http.StatusUnauthorized, nil)
			return
		}

		next.ServeHTTP(w, r)
	})

}
