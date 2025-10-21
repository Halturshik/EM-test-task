package api

import (
	"log"
	"net/http"
	"time"
)

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *logResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &logResponseWriter{ResponseWriter: w, statusCode: 200}

		log.Printf("→ %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		if lrw.statusCode >= 400 {
			log.Printf("← %s %s завершился с ошибкой %d (заняло %s)", r.Method, r.URL.Path, lrw.statusCode, duration)
		} else {
			log.Printf("← %s %s (заняло %s)", r.Method, r.URL.Path, duration)
		}
	})
}
