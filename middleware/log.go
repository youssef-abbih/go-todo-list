package middleware

import (
	"net"
	"time"
	"log/slog"
	"net/http"

)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytesWritten += n
	return n, err
}

func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriterWrapper{ResponseWriter: w}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		clientIP := r.RemoteAddr
		if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			clientIP = ip
		}

		slog.Info("HTTP request",
			"method", r.Method,
			"url", r.URL.String(),
			"proto", r.Proto,
			"remote_ip", clientIP,
			"user_agent", r.UserAgent(),
			"referer", r.Referer(),
			"status", wrapped.statusCode,
			"size", wrapped.bytesWritten,
			"duration", duration,
		)
	})
}
