package handler

import (
	"encoding/json"
	"net/http"

	"github.com/NicolasNSC/showcase-service-fiap/internal/usecase"
	"github.com/go-chi/chi"
)

type SaleHandler struct {
	useCase usecase.SaleUseCaseInterface
}

func NewSaleHandler(useCase usecase.SaleUseCaseInterface) *SaleHandler {
	return &SaleHandler{
		useCase: useCase,
	}
}

func (h *SaleHandler) CreateListing(w http.ResponseWriter, r *http.Request) {
	var input usecase.InputCreateListingDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	output, err := h.useCase.CreateListing(r.Context(), &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *SaleHandler) ListAvailable(w http.ResponseWriter, r *http.Request) {
	output, err := h.useCase.ListAvailable(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

func (h *SaleHandler) ListSold(w http.ResponseWriter, r *http.Request) {
	output, err := h.useCase.ListSold(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

func (h *SaleHandler) UpdateListing(w http.ResponseWriter, r *http.Request) {
	vehicleID := chi.URLParam(r, "vehicle_id")
	if vehicleID == "" {
		http.Error(w, "Vehicle ID is required", http.StatusBadRequest)
		return
	}

	var input usecase.InputUpdateListingDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.useCase.UpdateListing(r.Context(), vehicleID, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
