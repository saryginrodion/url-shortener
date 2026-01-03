package handlers

import (
	"errors"
	"net/http"
	"time"
	"github.com/google/uuid"
	"roadmap.restapi/internal/api/middleware"
	"roadmap.restapi/internal/api/request"
	"roadmap.restapi/internal/api/response"
	"roadmap.restapi/internal/config"
	"roadmap.restapi/internal/ctxlogging"
	"roadmap.restapi/internal/token"
	"roadmap.restapi/internal/user"
)

const API_COOKIE_PATH = "/api/v1"

func registration(users *user.UseCases, tokens *token.UseCases) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log := ctxlogging.Get(r.Context())
		body, err := request.ParseAndValidateJson(validate, r.Body, RegistrationRequest{})
		if err != nil {
			response.WriteJsonErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		newUser, err := users.NewUser(r.Context(), body.Email, body.Password)
		if err != nil {
			if errors.Is(err, user.ErrUserAlreadyExists) {
				response.WriteJsonErrorResponse(w, err, http.StatusConflict)
			} else {
				log.Error("unknown error on user creation", "err", err)
				response.WriteJsonErrorResponse(w, err, http.StatusInternalServerError)
			}
			return
		}

		tokenPair, err := tokens.NewPair(r.Context(), newUser.ID)
		if err != nil {
			log.Error("unknown error on token generation", "err", err)
			response.WriteJsonErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		accessCookie := &http.Cookie{
			Name:     middleware.COOKIE_ACCESS,
			Value:    tokenPair.Access,
			Secure:   config.Cfg().IsProd(),
			HttpOnly: config.Cfg().IsProd(),
			Path: API_COOKIE_PATH,
		}

		refreshCookie := &http.Cookie{
			Name:     middleware.COOKIE_REFRESH,
			Value:    tokenPair.Refresh,
			Secure:   config.Cfg().IsProd(),
			HttpOnly: config.Cfg().IsProd(),
			Path: API_COOKIE_PATH,
		}

		http.SetCookie(w, accessCookie)
		http.SetCookie(w, refreshCookie)

		response.WriteJsonResponse(w, response.NewResponse(UserDTO{
			Email:     newUser.Email,
			CreatedAt: newUser.CreatedAt,
		}), http.StatusOK)
	}
}

func logout(tokens *token.UseCases) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		refTokenCookie, err := r.Cookie(middleware.COOKIE_REFRESH)
		if err != nil {
			response.WriteJsonErrorResponse(w, err, http.StatusUnprocessableEntity)
			return
		}

		if err = tokens.RevokeToken(r.Context(), refTokenCookie.Value); err != nil {
			response.WriteJsonErrorResponse(w, err, http.StatusUnprocessableEntity)
			return
		}

		accessCookie := &http.Cookie{
			Name:     middleware.COOKIE_ACCESS,
			Value:    "",
			Expires:  time.Now(),
			Secure:   config.Cfg().IsProd(),
			HttpOnly: config.Cfg().IsProd(),
			Path: API_COOKIE_PATH,
		}

		refreshCookie := &http.Cookie{
			Name:     middleware.COOKIE_REFRESH,
			Value:    "",
			Expires:  time.Now(),
			Secure:   config.Cfg().IsProd(),
			HttpOnly: config.Cfg().IsProd(),
			Path: API_COOKIE_PATH,
		}

		http.SetCookie(w, accessCookie)
		http.SetCookie(w, refreshCookie)

		response.WriteJsonResponse(w, response.NewResponse(struct{}{}), http.StatusNoContent)
	}
}

func refresh(tokens *token.UseCases) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		refTokenCookie, err := r.Cookie(middleware.COOKIE_REFRESH)
		if err != nil {
			response.WriteJsonErrorResponse(w, err, http.StatusUnprocessableEntity)
			return
		}

		newPair, err := tokens.RefreshTokenPair(r.Context(), refTokenCookie.Value)
		if err != nil {
			response.WriteJsonErrorResponse(w, err, http.StatusUnprocessableEntity)
			return
		}

		accessCookie := &http.Cookie{
			Name:     middleware.COOKIE_ACCESS,
			Value:    newPair.Access,
			Secure:   config.Cfg().IsProd(),
			HttpOnly: config.Cfg().IsProd(),
			Path: API_COOKIE_PATH,
		}

		refreshCookie := &http.Cookie{
			Name:     middleware.COOKIE_REFRESH,
			Value:    newPair.Refresh,
			Secure:   config.Cfg().IsProd(),
			HttpOnly: config.Cfg().IsProd(),
			Path: API_COOKIE_PATH,
		}

		http.SetCookie(w, accessCookie)
		http.SetCookie(w, refreshCookie)

		response.WriteJsonResponse(w, response.NewResponse(struct{}{}), http.StatusNoContent)
	}
}

func login(users *user.UseCases, tokens *token.UseCases) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log := ctxlogging.Get(r.Context())
		body, err := request.ParseAndValidateJson(validate, r.Body, LoginRequest{})
		if err != nil {
			response.WriteJsonErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		u, err := users.ByEmailAndPassword(r.Context(), body.Email, body.Password)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, user.ErrPasswordIncorrect) {
				response.WriteJsonErrorResponse(w, errors.New("incorrect email or password"), http.StatusNotFound)
			} else {
				log.Error("unknown error on user creation", "err", err)
				response.WriteJsonErrorResponse(w, err, http.StatusInternalServerError)
			}
			return
		}

		tokenPair, err := tokens.NewPair(r.Context(), u.ID)
		if err != nil {
			log.Error("unknown error on token generation", "err", err)
			response.WriteJsonErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		accessCookie := &http.Cookie{
			Name:     middleware.COOKIE_ACCESS,
			Value:    tokenPair.Access,
			Secure:   config.Cfg().IsProd(),
			HttpOnly: config.Cfg().IsProd(),
			Path: API_COOKIE_PATH,
		}

		refreshCookie := &http.Cookie{
			Name:     middleware.COOKIE_REFRESH,
			Value:    tokenPair.Refresh,
			Secure:   config.Cfg().IsProd(),
			HttpOnly: config.Cfg().IsProd(),
			Path: API_COOKIE_PATH,
		}

		http.SetCookie(w, accessCookie)
		http.SetCookie(w, refreshCookie)

		response.WriteJsonResponse(w, response.NewResponse(UserDTO{
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		}), http.StatusOK)
	}
}

func me(userRepo user.UserRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log := ctxlogging.Get(r.Context())
		uid := r.Context().Value(middleware.CTX_USER_ID).(uuid.UUID)

		me, err := userRepo.ByID(r.Context(), uid)
		if err != nil {
			log.Error("unknown error on user get", "err", err)
			response.WriteJsonErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		response.WriteJsonResponse(w, response.NewResponse(UserDTO{
			Email:     me.Email,
			CreatedAt: me.CreatedAt,
		}), http.StatusOK)
	}
}
