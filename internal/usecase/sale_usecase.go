package usecase

import (
	"context"

	"github.com/NicolasNSC/showcase-service-fiap/internal/repository"
)

type OutputSaleItemDTO struct {
	VehicleID string  `json:"vehicle_id"`
	Brand     string  `json:"brand"`
	Model     string  `json:"model"`
	Price     float64 `json:"price"`
}

type SaleUseCaseInterface interface {
	ListAvailable(ctx context.Context) ([]*OutputSaleItemDTO, error)
	ListSold(ctx context.Context) ([]*OutputSaleItemDTO, error)
}

type saleUseCase struct {
	repo repository.SaleRepository
}

func NewSaleUseCase(repo repository.SaleRepository) SaleUseCaseInterface {
	return &saleUseCase{
		repo: repo,
	}
}

func (uc *saleUseCase) ListAvailable(ctx context.Context) ([]*OutputSaleItemDTO, error) {
	sales, err := uc.repo.GetAvailableByPrice(ctx)
	if err != nil {
		return nil, err
	}

	var output []*OutputSaleItemDTO
	for _, sale := range sales {
		output = append(output, &OutputSaleItemDTO{
			VehicleID: sale.VehicleID,
			Brand:     sale.Brand,
			Model:     sale.Model,
			Price:     sale.Price,
		})
	}

	return output, nil
}

func (uc *saleUseCase) ListSold(ctx context.Context) ([]*OutputSaleItemDTO, error) {
	sales, err := uc.repo.GetSoldByPrice(ctx)
	if err != nil {
		return nil, err
	}

	var output []*OutputSaleItemDTO
	for _, sale := range sales {
		output = append(output, &OutputSaleItemDTO{
			VehicleID: sale.VehicleID,
			Brand:     sale.Brand,
			Model:     sale.Model,
			Price:     sale.Price,
		})
	}

	return output, nil
}
