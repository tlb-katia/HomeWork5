package middleware

import (
	"log/slog"
	"net/http"
)

func LoggingMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("HTTP request received",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("remote_addr", r.RemoteAddr))
			next.ServeHTTP(w, r)
		})
	}
}
