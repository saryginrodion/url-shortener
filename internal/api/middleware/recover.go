package middleware

import (
	"errors"
	"net/http"

	"roadmap.restapi/internal/api/response"
	"roadmap.restapi/internal/ctxlogging"
)

func Recover() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := ctxlogging.Get(r.Context())

			defer func() {
				if rec := recover(); rec != nil {
					log.Error("server panicked!", "recover", rec, "url", r.URL, "ip", r.RemoteAddr)
					response.WriteJsonErrorResponse(w, errors.New("internal server error"), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
