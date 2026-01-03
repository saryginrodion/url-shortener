package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"roadmap.restapi/internal/ctxlogging"
)

func Logging(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			startTime := time.Now()

			traceID := uuid.New().String()[0:13]
			log := log.With(
				"traceID", traceID,
			)
			newCtx := req.WithContext(ctxlogging.Add(req.Context(), log))
			wrapWriter := middleware.NewWrapResponseWriter(w, req.ProtoMajor)

			next.ServeHTTP(wrapWriter, newCtx)

			defer func() {
				log.Info("served request",
					"method", req.Method,
					"url", req.URL.Path,
					"status", wrapWriter.Status(),
					"bytes", wrapWriter.BytesWritten(),
					"elapsed", time.Since(startTime),
				)
			}()
		})
	}
}
