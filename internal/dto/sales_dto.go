package dto

import "time"

type OutputSaleItemDTO struct {
	SaleID    string  `json:"sale_id"`
	VehicleID string  `json:"vehicle_id"`
	Brand     string  `json:"brand"`
	Model     string  `json:"model"`
	Price     float64 `json:"price"`
}

type InputCreateListingDTO struct {
	VehicleID string  `json:"vehicle_id"`
	Brand     string  `json:"brand"`
	Model     string  `json:"model"`
	Price     float64 `json:"price"`
}

type OutputCreateListingDTO struct {
	SaleID    string    `json:"sale_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type InputUpdateListingDTO struct {
	Brand string  `json:"brand"`
	Model string  `json:"model"`
	Price float64 `json:"price"`
}

type InputPurchaseDTO struct {
	BuyerCPF string `json:"buyer_cpf"`
}

type OutputPurchaseDTO struct {
	PaymentID string `json:"payment_id"`
}

type InputWebhookDTO struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
}
