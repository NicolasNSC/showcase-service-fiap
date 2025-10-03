package repository

import (
	"context"

	"github.com/NicolasNSC/showcase-service-fiap/internal/domain"
)

//go:generate mockgen -source=sale_repository.go -destination=./mocks/sale_repository_mock.go -package=mocks
type SaleRepository interface {
	Save(ctx context.Context, sale *domain.Sale) error
	Update(ctx context.Context, sale *domain.Sale) error
	GetByID(ctx context.Context, id string) (*domain.Sale, error)
	GetByVehicleID(ctx context.Context, vehicleID string) (*domain.Sale, error)
	GetByPaymentID(ctx context.Context, paymentID string) (*domain.Sale, error)
	GetAvailableByPrice(ctx context.Context) ([]*domain.Sale, error)
	GetSoldByPrice(ctx context.Context) ([]*domain.Sale, error)
}
