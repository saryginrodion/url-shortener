package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"roadmap.restapi/internal/api/middleware"
	"roadmap.restapi/internal/api/request"
	"roadmap.restapi/internal/api/response"
	"roadmap.restapi/internal/ctxlogging"
	"roadmap.restapi/internal/url"
)

func generateUrlID() string {
	newUrlUUID := strings.ReplaceAll(uuid.New().String(), "-", "")
	return string([]rune(newUrlUUID)[:10])
}

func urlCreate(urls url.URLRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log := ctxlogging.Get(r.Context())
		uid := r.Context().Value(middleware.CTX_USER_ID).(uuid.UUID)
		body, err := request.ParseAndValidateJson(validate, r.Body, UrlCreateRequest{})
		if err != nil {
			response.WriteJsonErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		if body.ID == "" {
			body.ID = generateUrlID()
		}

		newUrl := url.URL{
			ID:       body.ID,
			URL:      body.URL,
			Name:     body.Name,
			AuthorID: uid,
		}

		err = urls.Create(r.Context(), &newUrl)
		if err != nil {
			if errors.Is(err, url.ErrURLAlreadyExists) {
				response.WriteJsonErrorResponse(w, err, http.StatusConflict)
			} else {
				log.Error("unhandled error", "err", err)
				response.WriteJsonErrorResponse(w, err, http.StatusInternalServerError)
			}
			return
		}

		log.Debug("new url created", "url", newUrl)
		response.WriteJsonResponse(
			w,
			response.NewResponse(UrlDTO{
				ID:        newUrl.ID,
				Name:      newUrl.Name,
				URL:       newUrl.URL,
				CreatedAt: newUrl.CreatedAt,
			}),
			http.StatusCreated,
		)
	}
}

func urlRedirect(urls url.URLRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log := ctxlogging.Get(r.Context())
		urlID := chi.URLParam(r, "url-id")

		u, err := urls.ByID(r.Context(), urlID)
		if err != nil {
			if errors.Is(err, url.ErrURLNotFound) {
				response.WriteJsonErrorResponse(w, err, http.StatusNotFound)
			} else {
				log.Error("unhandled error", "err", err)
				response.WriteJsonErrorResponse(w, err, http.StatusInternalServerError)
			}
			return
		}

		log.Debug("url redirect", "url", u)
		http.Redirect(w, r, u.URL, http.StatusTemporaryRedirect)
	}
}

func urlDelete(urls *url.UseCases) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log := ctxlogging.Get(r.Context())
		uid := r.Context().Value(middleware.CTX_USER_ID).(uuid.UUID)
		urlID := chi.URLParam(r, "url-id")

		err := urls.Delete(r.Context(), uid, urlID)
		if err != nil {
			if errors.Is(err, url.ErrUserIsNotAuthor) {
				response.WriteJsonErrorResponse(w, err, http.StatusForbidden)
			} else if errors.Is(err, url.ErrURLNotFound) {
				response.WriteJsonErrorResponse(w, err, http.StatusNotFound)
			} else {
				log.Error("unhandled error", "err", err)
				response.WriteJsonErrorResponse(w, err, http.StatusInternalServerError)
			}
			return
		}

		log.Debug("url deleted", "urlID", urlID)
		response.WriteJsonResponse(
			w,
			struct{}{},
			http.StatusNoContent,
		)
	}
}
