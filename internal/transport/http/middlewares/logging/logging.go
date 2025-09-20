package middlewares

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		parts := strings.Split(fwd, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

func routePattern(r *http.Request) string {
	if rc := chi.RouteContext(r.Context()); rc != nil {
		if pat := rc.RoutePattern(); pat != "" {
			return pat
		}
	}
	return r.URL.Path
}

func LoggingMiddleware(base *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			reqID := middleware.GetReqID(r.Context())
			log := base.With(
				"request_id", reqID,
				"method", r.Method,
				"route", routePattern(r),
				"remote_ip", clientIP(r),
				"user_agent", r.UserAgent(),
			)
			r = r.WithContext(WithLogger(r.Context(), log))

			next.ServeHTTP(ww, r)

			status := ww.Status()
			if status == 0 {
				status = http.StatusOK
			}
			latency := time.Since(start).Milliseconds()

			attrs := []any{
				"status", status,
				"size", ww.BytesWritten(),
				"duration_ms", latency,
			}

			route := routePattern(r)
			if route == "/healthz" || strings.HasPrefix(route, "/metrics") {
				log.Debug("http_request", attrs...)
				return
			}

			switch {
			case status >= 500:
				log.Error("http_request", attrs...)
			case status >= 400:
				log.Warn("http_request", attrs...)
			default:
				log.Info("http_request", attrs...)
			}
		})
	}
}
