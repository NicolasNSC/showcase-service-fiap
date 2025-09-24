package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type SaleStatus string

const (
	StatusAvailable      SaleStatus = "AVAILABLE"
	StatusPendingPayment SaleStatus = "PENDING_PAYMENT"
	StatusSold           SaleStatus = "SOLD"
	StatusCanceled       SaleStatus = "CANCELED"
)

type Sale struct {
	ID        string     `json:"id"`
	VehicleID string     `json:"vehicle_id"`
	Brand     string     `json:"brand"`
	Model     string     `json:"model"`
	Price     float64    `json:"price"`
	Status    SaleStatus `json:"status"`
	PaymentID string     `json:"payment_id,omitempty"`
	BuyerCPF  *string    `json:"buyer_cpf,omitempty"`
	SaleDate  *time.Time `json:"sale_date,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewSale(vehicleID, brand, model string, price float64) (*Sale, error) {
	if vehicleID == "" {
		return nil, errors.New("vehicle_id cannot be empty")
	}
	if price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}
	if brand == "" || model == "" {
		return nil, errors.New("brand and model are required for listing")
	}

	return &Sale{
		ID:        uuid.New().String(),
		VehicleID: vehicleID,
		Brand:     brand,
		Model:     model,
		Price:     price,
		Status:    StatusAvailable,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
