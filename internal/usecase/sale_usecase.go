package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/NicolasNSC/showcase-service-fiap/internal/domain"
	"github.com/NicolasNSC/showcase-service-fiap/internal/dto"
	"github.com/NicolasNSC/showcase-service-fiap/internal/repository"
	"github.com/google/uuid"
)

//go:generate mockgen -source=sale_usecase.go -destination=./mocks/sale_usecase_mock.go -package=mocks
type SaleUseCaseInterface interface {
	CreateListing(ctx context.Context, input *dto.InputCreateListingDTO) (*dto.OutputCreateListingDTO, error)
	UpdateListing(ctx context.Context, vehicleID string, input *dto.InputUpdateListingDTO) error
	Purchase(ctx context.Context, saleID string, input dto.InputPurchaseDTO) (*dto.OutputPurchaseDTO, error)
	HandlePaymentWebhook(ctx context.Context, input *dto.InputWebhookDTO) error
	ListAvailable(ctx context.Context) ([]*dto.OutputSaleItemDTO, error)
	ListSold(ctx context.Context) ([]*dto.OutputSaleItemDTO, error)
}

type saleUseCase struct {
	repo repository.SaleRepository
}

func NewSaleUseCase(repo repository.SaleRepository) SaleUseCaseInterface {
	return &saleUseCase{
		repo: repo,
	}
}

func (uc *saleUseCase) CreateListing(ctx context.Context, input *dto.InputCreateListingDTO) (*dto.OutputCreateListingDTO, error) {
	sale, err := domain.NewSale(input.VehicleID, input.Brand, input.Model, input.Price)
	if err != nil {
		return nil, err
	}

	err = uc.repo.Save(ctx, sale)
	if err != nil {
		return nil, err
	}

	output := &dto.OutputCreateListingDTO{
		SaleID:    sale.ID,
		Status:    string(sale.Status),
		CreatedAt: sale.CreatedAt,
	}

	return output, nil
}

func (uc *saleUseCase) UpdateListing(ctx context.Context, vehicleID string, input *dto.InputUpdateListingDTO) error {
	sale, err := uc.repo.GetByVehicleID(ctx, vehicleID)
	if err != nil {
		return err
	}

	sale.Brand = input.Brand
	sale.Model = input.Model
	sale.Price = input.Price
	sale.UpdatedAt = time.Now()

	return uc.repo.Update(ctx, sale)
}

func (uc *saleUseCase) Purchase(ctx context.Context, saleID string, input dto.InputPurchaseDTO) (*dto.OutputPurchaseDTO, error) {
	sale, err := uc.repo.GetByID(ctx, saleID)
	if err != nil {
		return nil, err
	}

	if sale.Status != domain.StatusAvailable {
		return nil, errors.New("sale is not available for purchase")
	}

	now := time.Now()
	sale.Status = domain.StatusPendingPayment
	sale.BuyerCPF = &input.BuyerCPF
	sale.SaleDate = &now
	sale.PaymentID = uuid.New().String()
	sale.UpdatedAt = now

	err = uc.repo.Update(ctx, sale)
	if err != nil {
		return nil, err
	}

	output := &dto.OutputPurchaseDTO{
		PaymentID: sale.PaymentID,
	}

	return output, nil
}

func (uc *saleUseCase) HandlePaymentWebhook(ctx context.Context, input *dto.InputWebhookDTO) error {
	sale, err := uc.repo.GetByPaymentID(ctx, input.PaymentID)
	if err != nil {
		return err
	}

	if sale.Status != domain.StatusPendingPayment {
		return errors.New("sale is not in pending payment status")
	}

	switch strings.ToUpper(input.Status) {
	case "APPROVED", "EFETUADO":
		sale.Status = domain.StatusSold
	case "CANCELED", "CANCELADO":
		sale.Status = domain.StatusCanceled
	default:
		return errors.New("invalid payment status received from webhook")
	}

	sale.UpdatedAt = time.Now()

	return uc.repo.Update(ctx, sale)
}

func (uc *saleUseCase) ListAvailable(ctx context.Context) ([]*dto.OutputSaleItemDTO, error) {
	sales, err := uc.repo.GetAvailableByPrice(ctx)
	if err != nil {
		return nil, err
	}

	var output []*dto.OutputSaleItemDTO
	for _, sale := range sales {
		output = append(output, &dto.OutputSaleItemDTO{
			SaleID:    sale.ID,
			VehicleID: sale.VehicleID,
			Brand:     sale.Brand,
			Model:     sale.Model,
			Price:     sale.Price,
		})
	}

	return output, nil
}

func (uc *saleUseCase) ListSold(ctx context.Context) ([]*dto.OutputSaleItemDTO, error) {
	sales, err := uc.repo.GetSoldByPrice(ctx)
	if err != nil {
		return nil, err
	}

	var output []*dto.OutputSaleItemDTO
	for _, sale := range sales {
		output = append(output, &dto.OutputSaleItemDTO{
			SaleID:    sale.ID,
			VehicleID: sale.VehicleID,
			Brand:     sale.Brand,
			Model:     sale.Model,
			Price:     sale.Price,
		})
	}

	return output, nil
}
