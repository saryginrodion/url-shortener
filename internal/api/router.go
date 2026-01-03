package api

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"roadmap.restapi/internal/api/handlers"
	"roadmap.restapi/internal/api/middleware"
	"roadmap.restapi/internal/token"
	"roadmap.restapi/internal/url"
	"roadmap.restapi/internal/user"
)

func NewRouter(
	log *slog.Logger,
	users *user.UseCases,
	userRepo user.UserRepository,
	tokens *token.UseCases,
	tokenExtractor token.ClaimsExtractor,
	urlsRepo url.URLRepository,
	urls *url.UseCases,
) chi.Router {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Logging(log))
		r.Use(middleware.Recover())
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"*"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: false,
		}))

		r.Mount("/auth", handlers.AuthRouter(
			tokenExtractor,
			userRepo,
			users,
			tokens,
		))

		r.Mount("/urls", handlers.UrlsRouter(
			tokenExtractor,
			userRepo,
			urlsRepo,
			urls,
		))
	})

	return r
}
