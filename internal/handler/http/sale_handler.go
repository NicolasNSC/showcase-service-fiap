package handler

import (
	"encoding/json"
	"net/http"

	"github.com/NicolasNSC/showcase-service-fiap/internal/dto"
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

// CreateListing lida com a requisição interna para criar uma nova listagem.
// @Summary      Create a new sale listing
// @Description  Creates a new sale listing when notified by the catalog-service. This is an internal endpoint.
// @Tags         Internal
// @Accept       json
// @Produce      json
// @Param        listing  body      dto.InputCreateListingDTO  true  "Listing Data"
// @Success      201      {object}  dto.OutputCreateListingDTO
// @Failure      400      {string}  string "Invalid request body"
// @Failure      500      {string}  string "Internal server error"
// @Router       /listings [post]
func (h *SaleHandler) CreateListing(w http.ResponseWriter, r *http.Request) {
	var input dto.InputCreateListingDTO
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

// ListAvailable lida com a requisição para listar veículos à venda.
// @Summary      List available vehicles
// @Description  Get a list of all vehicles available for sale, sorted by price.
// @Tags         Sales
// @Accept       json
// @Produce      json
// @Success      200 {array}   dto.OutputSaleItemDTO
// @Failure      500 {string}  string "Internal server error"
// @Router       /sales/available [get]
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

// ListSold lida com a requisição para listar veículos vendidos.
// @Summary      List sold vehicles
// @Description  Get a list of all vehicles that have been sold, sorted by price.
// @Tags         Sales
// @Accept       json
// @Produce      json
// @Success      200 {array}   dto.OutputSaleItemDTO
// @Failure      500 {string}  string "Internal server error"
// @Router       /sales/sold [get]
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

// UpdateListing lida com a requisição interna para atualizar uma listagem.
// @Summary      Update a sale listing
// @Description  Updates a sale listing's data when notified by the catalog-service. This is an internal endpoint.
// @Tags         Internal
// @Accept       json
// @Produce      json
// @Param        vehicle_id  path      string                         true  "Vehicle ID"
// @Param        listing     body      dto.InputUpdateListingDTO  true  "Listing Data to Update"
// @Success      200         {string}  string "OK"
// @Failure      404         {string}  string "Listing not found"
// @Failure      500         {string}  string "Internal server error"
// @Router       /listings/vehicle/{vehicle_id} [put]
func (h *SaleHandler) UpdateListing(w http.ResponseWriter, r *http.Request) {
	vehicleID := chi.URLParam(r, "vehicle_id")
	if vehicleID == "" {
		http.Error(w, "Vehicle ID is required", http.StatusBadRequest)
		return
	}

	var input *dto.InputUpdateListingDTO
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

// Purchase lida com a requisição para iniciar a compra de um veículo.
// @Summary      Purchase a vehicle
// @Description  Initiates the purchase process for a specific sale listing.
// @Tags         Sales
// @Accept       json
// @Produce      json
// @Param        id        path      string                   true  "Sale ID"
// @Param        purchase  body      dto.InputPurchaseDTO true  "Buyer's CPF"
// @Success      202       {object}  dto.OutputPurchaseDTO
// @Failure      400       {string}  string "Invalid request body or ID"
// @Failure      404       {string}  string "Sale not found"
// @Failure      409       {string}  string "Sale is not available for purchase"
// @Failure      500       {string}  string "Internal server error"
// @Router       /sales/{id}/purchase [post]
func (h *SaleHandler) Purchase(w http.ResponseWriter, r *http.Request) {
	saleID := chi.URLParam(r, "id")
	if saleID == "" {
		http.Error(w, "Sale ID is required", http.StatusBadRequest)
		return
	}

	var input dto.InputPurchaseDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	output, err := h.useCase.Purchase(r.Context(), saleID, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(output)
}

// HandlePaymentWebhook lida com a notificação de pagamento do sistema externo.
// @Summary      Handle a payment webhook
// @Description  Receives payment status notifications from an external payment gateway.
// @Tags         Webhooks
// @Accept       json
// @Produce      json
// @Param        notification  body      dto.InputWebhookDTO  true  "Payment Notification Payload"
// @Success      204           {string}  string "No Content"
// @Failure      400           {string}  string "Invalid request body"
// @Failure      500           {string}  string "Failed to process webhook"
// @Router       /webhooks/payments [post]
func (h *SaleHandler) HandlePaymentWebhook(w http.ResponseWriter, r *http.Request) {
	var input *dto.InputWebhookDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.useCase.HandlePaymentWebhook(r.Context(), input)
	if err != nil {
		http.Error(w, "Failed to process webhook", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
