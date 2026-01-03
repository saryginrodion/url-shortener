package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"roadmap.restapi/internal/api/middleware"
	"roadmap.restapi/internal/token"
	"roadmap.restapi/internal/user"
)

func AuthRouter(
	extractor token.ClaimsExtractor,
	userRepo user.UserRepository,
	users *user.UseCases,
	tokens *token.UseCases,
) chi.Router {
	r := chi.NewRouter()
	r.Post("/registration", registration(users, tokens))
	r.Post("/login", login(users, tokens))
	r.Post("/logout", logout(tokens))
	r.Post("/refresh", refresh(tokens))
	r.With(middleware.Auth(extractor, userRepo)).Get("/me", http.HandlerFunc(me(userRepo)))
	return r
}

