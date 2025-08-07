package account

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nontypeable/financial-tracker/internal/auth"
	"github.com/nontypeable/financial-tracker/internal/delivery/account/dto"
	"github.com/nontypeable/financial-tracker/internal/domain/account"
	httpHelper "github.com/nontypeable/financial-tracker/internal/http"
)

type handler struct {
	service account.Service
}

func NewHandler(service account.Service) *handler {
	return &handler{service: service}
}

func (h *handler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/account", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Post("/", h.create)
		})
	})
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		if err := httpHelper.Error(w, http.StatusInternalServerError, "internal server error"); err != nil {
			log.Printf("httpHelper.Error: %v", err)
		}
		return
	}

	var payload dto.CreateRequest
	if err := httpHelper.DecodeAndValidate(r, &payload); err != nil {
		status, msg := httpHelper.MapAppErrorToHTTP(err)
		httpHelper.Error(w, status, msg)
		return
	}

	accountID, err := h.service.Create(r.Context(), userID, payload.Name, payload.Balance)
	if err != nil {
		status, msg := httpHelper.MapAppErrorToHTTP(err)
		httpHelper.Error(w, status, msg)
		return
	}

	if err := httpHelper.JSON(w, http.StatusCreated, &dto.CreateResponse{ID: accountID}); err != nil {
		log.Printf("httpHelper.JSON: %v", err)
	}
}
