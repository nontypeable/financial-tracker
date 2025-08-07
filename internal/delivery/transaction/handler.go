package transaction

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nontypeable/financial-tracker/internal/delivery/transaction/dto"
	"github.com/nontypeable/financial-tracker/internal/domain/transaction"
	httpHelper "github.com/nontypeable/financial-tracker/internal/http"
)

type handler struct {
	service transaction.Service
}

func NewHandler(service transaction.Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/transaction", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Post("/", h.create)
		})
	})
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) {
	var payload dto.CreateRequest
	if err := httpHelper.DecodeAndValidate(r, &payload); err != nil {
		log.Print(err.Error())
		if err := httpHelper.Error(w, http.StatusBadRequest, "invalid json payload"); err != nil {
			log.Printf("httpHelper.Error: %v", err)
		}
		return
	}

	id, err := h.service.Create(r.Context(), payload.AccountID, payload.Amount, payload.Type, payload.Description)
	if err != nil {
		log.Print(err.Error())
		status, msg := httpHelper.MapAppErrorToHTTP(err)
		if err := httpHelper.Error(w, status, msg); err != nil {
			log.Printf("httpHelper.Error: %v", err)
		}
		return
	}

	if err := httpHelper.JSON(w, http.StatusCreated, &dto.CreateResponse{ID: id}); err != nil {
		log.Printf("httpHelper.JSON: %v", err)
	}
}
