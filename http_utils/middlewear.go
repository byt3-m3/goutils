package http_utils

import (
	"log/slog"
	"net/http"
)

// LoggingMiddleware returns an HTTP middleware that logs requests using slog.
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			logger.Info("received request",
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
				slog.String("remoteAddr", req.RemoteAddr),
				slog.Int64("contentLen", req.ContentLength),
			)

			next.ServeHTTP(w, req)
		})
	}
}

func RecoveryMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered",
						slog.Any("error", rec),
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
						slog.String("remoteAddr", r.RemoteAddr),
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
