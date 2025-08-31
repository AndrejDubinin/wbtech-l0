package http

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	logger interface {
		Info(msg string, fields ...zap.Field)
		Error(msg string, fields ...zap.Field)
	}
)

func AccessLogMiddleware(next http.Handler, logger logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Info("response",
			zap.String("method", r.Method),
			zap.String("remote_address", r.RemoteAddr),
			zap.String("path", r.URL.Path),
			zap.String("response_time", time.Since(start).String()))
	})
}
