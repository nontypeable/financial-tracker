package user

import (
	"net/http"
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
		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", h.signUp)
			r.Post("/sign-in", h.signIn)
			r.Post("/refresh", h.refresh)
		})

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Get("/me", h.getUserInfo)
			r.Patch("/me", h.update)
			r.Patch("/me/email", h.updateEmail)
			r.Patch("/me/password", h.updatePassword)
		})
	})
}

func (h *handler) signUp(w http.ResponseWriter, r *http.Request) {
	var payload dto.SignUpRequest
	if err := httpHelper.DecodeAndValidate(r, &payload); err != nil {
		status, msg := httpHelper.MapAppErrorToHTTP(err)
		httpHelper.Error(w, status, msg)
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
		status, msg := httpHelper.MapAppErrorToHTTP(err)
		httpHelper.Error(w, status, msg)
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

	httpHelper.JSON(w, http.StatusCreated, &dto.SignUpResponse{AccessToken: accessToken})
}

func (h *handler) signIn(w http.ResponseWriter, r *http.Request) {
	var payload dto.SignInRequest
	if err := httpHelper.DecodeAndValidate(r, &payload); err != nil {
		status, msg := httpHelper.MapAppErrorToHTTP(err)
		httpHelper.Error(w, status, msg)
		return
	}

	accessToken, refreshToken, err := h.service.SignIn(r.Context(), payload.Email, payload.Password)
	if err != nil {
		status, msg := httpHelper.MapAppErrorToHTTP(err)
		httpHelper.Error(w, status, msg)
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

	httpHelper.JSON(w, http.StatusOK, &dto.SignInResponse{AccessToken: accessToken})
}

func (h *handler) refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		status, msg := httpHelper.MapAppErrorToHTTP(err)
		httpHelper.Error(w, status, msg)
		return
	}

	accessToken, refreshToken, err := h.service.Refresh(r.Context(), cookie.Value)
	if err != nil {
		status, msg := httpHelper.MapAppErrorToHTTP(err)
		httpHelper.Error(w, status, msg)
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

func (h *handler) getUserInfo(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpHelper.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	user, err := h.service.GetUserInfo(r.Context(), userID)
	if err != nil {
		status, msg := httpHelper.MapAppErrorToHTTP(err)
		httpHelper.Error(w, status, msg)
		return
	}

	httpHelper.JSON(w, http.StatusOK, &dto.GetUserInfoResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})
}

func (h *handler) update(w http.ResponseWriter, r *http.Request) {
	var payload dto.UpdateRequest
	if err := httpHelper.DecodeAndValidate(r, &payload); err != nil {
		status, msg := httpHelper.MapAppErrorToHTTP(err)
		httpHelper.Error(w, status, msg)
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpHelper.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := h.service.Update(r.Context(), userID, payload.FirstName, payload.LastName); err != nil {
		status, msg := httpHelper.MapAppErrorToHTTP(err)
		httpHelper.Error(w, status, msg)
		return
	}

	httpHelper.JSON(w, http.StatusOK, nil)
}

func (h *handler) updateEmail(w http.ResponseWriter, r *http.Request) {
	var payload dto.UpdateEmailRequest
	if err := httpHelper.DecodeAndValidate(r, &payload); err != nil {
		httpHelper.Error(w, http.StatusBadRequest, "invalid json payload")
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpHelper.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := h.service.ChangeEmail(r.Context(), userID, payload.Email, payload.Password); err != nil {
		status, msg := httpHelper.MapAppErrorToHTTP(err)
		httpHelper.Error(w, status, msg)
		return
	}

	httpHelper.JSON(w, http.StatusOK, nil)
}

func (h *handler) updatePassword(w http.ResponseWriter, r *http.Request) {
	var payload dto.UpdatePasswordRequest
	if err := httpHelper.DecodeAndValidate(r, &payload); err != nil {
		status, msg := httpHelper.MapAppErrorToHTTP(err)
		httpHelper.Error(w, status, msg)
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpHelper.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := h.service.ChangePassword(r.Context(), userID, payload.NewPassword, payload.CurrentPassword); err != nil {
		status, msg := httpHelper.MapAppErrorToHTTP(err)
		httpHelper.Error(w, status, msg)
		return
	}

	httpHelper.JSON(w, http.StatusOK, nil)
}
