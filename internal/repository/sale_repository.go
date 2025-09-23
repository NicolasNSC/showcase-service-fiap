package repository

import (
	"context"

	"github.com/NicolasNSC/showcase-service-fiap/internal/domain"
)

type SaleRepository interface {
	Save(ctx context.Context, sale *domain.Sale) error
	GetAvailableByPrice(ctx context.Context) ([]*domain.Sale, error)
	GetSoldByPrice(ctx context.Context) ([]*domain.Sale, error)
}
