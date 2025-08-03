package user

import (
	"net/http"
	"strings"
	"time"

	"github.com/nontypeable/financial-tracker/internal/auth"
	"github.com/nontypeable/financial-tracker/internal/delivery/user/dto"
	httpHelper "github.com/nontypeable/financial-tracker/internal/http"

	"github.com/go-chi/chi/v5"
	"github.com/nontypeable/financial-tracker/internal/domain/user"
)

type handler struct {
	service user.Service
}

func NewHandler(service user.Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/sign-up", h.signUp)
		r.Post("/sign-in", h.signIn)
		r.Post("/refresh", h.refresh)
		r.With(authMiddleware).Post("/update", h.update)
	})
}

func (h *handler) signUp(w http.ResponseWriter, r *http.Request) {
	var payload dto.SignUpRequest
	err := httpHelper.DecodeAndValidate(r, &payload)
	if err != nil {
		httpHelper.Error(w, http.StatusBadRequest, "invalid json payload")
		return
	}

	accessToken, refreshToken, err := h.service.SignUp(
		r.Context(),
		payload.Email,
		payload.Password,
		payload.FirstName,
		payload.LastName,
	)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			httpHelper.Error(w, http.StatusConflict, "user already registered")
			return
		}
		httpHelper.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/user/refresh",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})

	httpHelper.JSON(w, http.StatusCreated, &dto.SignUpResponse{AccessToken: accessToken})
}

func (h *handler) signIn(w http.ResponseWriter, r *http.Request) {
	var payload dto.SignInRequest
	err := httpHelper.DecodeAndValidate(r, &payload)
	if err != nil {
		httpHelper.Error(w, http.StatusBadRequest, "invalid json payload")
		return
	}

	accessToken, refreshToken, err := h.service.SignIn(
		r.Context(),
		payload.Email,
		payload.Password,
	)
	if err != nil {
		httpHelper.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/user/refresh",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})

	httpHelper.JSON(w, http.StatusOK, &dto.SignInResponse{AccessToken: accessToken})
}

func (h *handler) refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		httpHelper.Error(w, http.StatusUnauthorized, "missing refresh token")
		return
	}

	accessToken, refreshToken, err := h.service.Refresh(r.Context(), cookie.Value)
	if err != nil {
		httpHelper.Error(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/user/auth/refresh",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})

	httpHelper.JSON(w, http.StatusOK, &dto.RefreshResponse{AccessToken: accessToken})
}

func (h *handler) update(w http.ResponseWriter, r *http.Request) {
	var payload dto.UpdateRequest
	err := httpHelper.DecodeAndValidate(r, &payload)
	if err != nil {
		httpHelper.Error(w, http.StatusBadRequest, "invalid json payload")
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpHelper.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	err = h.service.Update(r.Context(), userID, payload.FirstName, payload.LastName)
	if err != nil {
		httpHelper.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	httpHelper.JSON(w, http.StatusOK, nil)
}
