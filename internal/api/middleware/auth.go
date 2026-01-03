package middleware

import (
	"context"
	"errors"
	"net/http"

	"roadmap.restapi/internal/api/response"
	"roadmap.restapi/internal/ctxlogging"
	"roadmap.restapi/internal/token"
	"roadmap.restapi/internal/user"
)

const (
	COOKIE_ACCESS  = "access-token"
	COOKIE_REFRESH = "refresh-token"

	CTX_USER_ID = "user_id"
)

var ErrNotAuthenticated = errors.New("Unauthenticated")

func Auth(extractor token.ClaimsExtractor, userRepo user.UserRepository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := ctxlogging.Get(r.Context())
			cookie, err := r.Cookie(COOKIE_ACCESS)
			if err != nil {
				response.WriteJsonResponse(w, response.NewErrorResponse(err.Error()), http.StatusUnauthorized)
				return
			}

			claims, err := extractor.ParseAndValidate(r.Context(), cookie.Value)
			if err != nil {
				response.WriteJsonResponse(w, response.NewErrorResponse(err.Error()), http.StatusUnauthorized)
				return
			}

			log.Info("claims", "claims", claims)

			_, err = userRepo.ByID(r.Context(), claims.UID)
			if err != nil {
				response.WriteJsonErrorResponse(w, err, http.StatusUnauthorized)
				return
			}

			newCtx := context.WithValue(r.Context(), CTX_USER_ID, claims.UID)

			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}
