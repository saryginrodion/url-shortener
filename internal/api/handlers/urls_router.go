package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"roadmap.restapi/internal/api/middleware"
	"roadmap.restapi/internal/token"
	"roadmap.restapi/internal/url"
	"roadmap.restapi/internal/user"
)

func UrlsRouter(
	extractor token.ClaimsExtractor,
	userRepo user.UserRepository,
	urlsRepo url.URLRepository,
	urls *url.UseCases,
) chi.Router {
	r := chi.NewRouter()
	authMW := middleware.Auth(extractor, userRepo)

	r.Get("/{url-id}", http.HandlerFunc(urlRedirect(urlsRepo)))

	r.With(authMW).Post("/", http.HandlerFunc(urlCreate(urlsRepo)))
	r.With(authMW).Delete("/{url-id}", http.HandlerFunc(urlDelete(urls)))

	return r
}
