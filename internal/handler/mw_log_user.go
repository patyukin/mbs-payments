package handler

import (
	"auth-telegram/internal/metrics"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h *Handler) LogUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.IncomingTraffic.Inc() // TODO use prometheus

		log.Error().Msgf(
			"Got request, method: %s, url: %s, user-agent: %s, x-user-id: %s, x-user-role: %s",
			r.Method,
			r.URL.String(),
			r.UserAgent(),
			r.Header.Get("X-User-ID"),
			r.Header.Get("X-User-Role"),
		)

		next.ServeHTTP(w, r)
	})
}
