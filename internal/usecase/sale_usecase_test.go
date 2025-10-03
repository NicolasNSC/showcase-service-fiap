package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/NicolasNSC/showcase-service-fiap/internal/domain"
	"github.com/NicolasNSC/showcase-service-fiap/internal/dto"
	"github.com/NicolasNSC/showcase-service-fiap/internal/repository/mocks"
	"github.com/NicolasNSC/showcase-service-fiap/internal/usecase"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type SaleUseCaseSuite struct {
	suite.Suite

	ctx        context.Context
	repository *mocks.MockSaleRepository
}

func (suite *SaleUseCaseSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.ctx = context.Background()
	suite.repository = mocks.NewMockSaleRepository(ctrl)
}

func Test_SaleUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(SaleUseCaseSuite))
}

func (suite *SaleUseCaseSuite) Test_CreateListing() {
	input := &dto.InputCreateListingDTO{
		VehicleID: "fc338f17-9fe8-40d1-8232-461fb1ecd080",
		Brand:     "Toyota",
		Model:     "Corolla",
		Price:     50000,
	}

	suite.T().Run("should create listing successfully", func(t *testing.T) {
		usecase := usecase.NewSaleUseCase(suite.repository)

		suite.repository.EXPECT().Save(suite.ctx, gomock.Any()).Return(nil)

		output, err := usecase.CreateListing(suite.ctx, input)
		suite.NoError(err)
		suite.NotNil(output)
		suite.NotEmpty(output.SaleID)
		suite.Equal(string(domain.StatusAvailable), output.Status)
		suite.WithinDuration(time.Now(), output.CreatedAt, time.Second)
	})

	suite.T().Run("should return error when domain.NewSale fails", func(t *testing.T) {
		usecase := usecase.NewSaleUseCase(suite.repository)

		input := &dto.InputCreateListingDTO{
			VehicleID: "",
			Brand:     "Toyota",
			Model:     "Corolla",
			Price:     50000,
		}

		output, err := usecase.CreateListing(suite.ctx, input)
		suite.Error(err)
		suite.Nil(output)
	})

	suite.T().Run("should return error when repo.Save fails", func(t *testing.T) {
		usecase := usecase.NewSaleUseCase(suite.repository)

		suite.repository.EXPECT().Save(suite.ctx, gomock.Any()).Return(errors.New("db error"))

		output, err := usecase.CreateListing(suite.ctx, input)
		suite.Error(err)
		suite.Nil(output)
	})
}

func (suite *SaleUseCaseSuite) Test_UpdateListing() {
	vehicleID := "fc338f17-9fe8-40d1-8232-461fb1ecd080"
	existingSale := &domain.Sale{
		ID:        "sale-123",
		VehicleID: vehicleID,
		Brand:     "Toyota",
		Model:     "Corolla",
		Price:     50000,
		Status:    domain.StatusAvailable,
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Hour),
	}
	input := &dto.InputUpdateListingDTO{
		Brand: "Honda",
		Model: "Civic",
		Price: 60000,
	}

	suite.T().Run("should update listing successfully", func(t *testing.T) {
		usecase := usecase.NewSaleUseCase(suite.repository)

		suite.repository.EXPECT().GetByVehicleID(suite.ctx, vehicleID).Return(existingSale, nil)
		suite.repository.EXPECT().Update(suite.ctx, gomock.Any()).Return(nil)

		err := usecase.UpdateListing(suite.ctx, vehicleID, input)
		suite.NoError(err)
	})

	suite.T().Run("should return error when repo.GetByVehicleID fails", func(t *testing.T) {
		usecase := usecase.NewSaleUseCase(suite.repository)

		suite.repository.EXPECT().GetByVehicleID(suite.ctx, vehicleID).Return(nil, errors.New("not found"))

		err := usecase.UpdateListing(suite.ctx, vehicleID, input)
		suite.Error(err)
	})

	suite.T().Run("should return error when repo.Update fails", func(t *testing.T) {
		usecase := usecase.NewSaleUseCase(suite.repository)

		suite.repository.EXPECT().GetByVehicleID(suite.ctx, vehicleID).Return(existingSale, nil)
		suite.repository.EXPECT().Update(suite.ctx, gomock.Any()).Return(errors.New("db error"))

		err := usecase.UpdateListing(suite.ctx, vehicleID, input)
		suite.Error(err)
	})
}

func (suite *SaleUseCaseSuite) Test_Purchase() {
	saleID := "sale-123"
	buyerCPF := "12345678900"
	input := dto.InputPurchaseDTO{
		BuyerCPF: buyerCPF,
	}

	suite.T().Run("should purchase successfully", func(t *testing.T) {
		now := time.Now()
		existingSale := &domain.Sale{
			ID:        saleID,
			VehicleID: "vehicle-1",
			Brand:     "Toyota",
			Model:     "Corolla",
			Price:     50000,
			Status:    domain.StatusAvailable,
			CreatedAt: now.Add(-time.Hour),
			UpdatedAt: now.Add(-time.Hour),
		}

		usecase := usecase.NewSaleUseCase(suite.repository)
		suite.repository.EXPECT().GetByID(suite.ctx, saleID).Return(existingSale, nil)
		suite.repository.EXPECT().Update(suite.ctx, gomock.Any()).Return(nil)

		output, err := usecase.Purchase(suite.ctx, saleID, input)
		suite.NoError(err)
		suite.NotNil(output)
		suite.NotEmpty(output.PaymentID)
	})

	suite.T().Run("should return error if repo.GetByID fails", func(t *testing.T) {
		usecase := usecase.NewSaleUseCase(suite.repository)
		suite.repository.EXPECT().GetByID(suite.ctx, saleID).Return(nil, errors.New("not found"))

		output, err := usecase.Purchase(suite.ctx, saleID, input)
		suite.Error(err)
		suite.Nil(output)
	})

	suite.T().Run("should return error if sale is not available", func(t *testing.T) {
		now := time.Now()
		notAvailableSale := &domain.Sale{
			ID:        saleID,
			VehicleID: "vehicle-1",
			Brand:     "Toyota",
			Model:     "Corolla",
			Price:     50000,
			Status:    domain.StatusSold,
			CreatedAt: now.Add(-time.Hour),
			UpdatedAt: now.Add(-time.Hour),
		}

		usecase := usecase.NewSaleUseCase(suite.repository)
		suite.repository.EXPECT().GetByID(suite.ctx, saleID).Return(notAvailableSale, nil)

		output, err := usecase.Purchase(suite.ctx, saleID, input)
		suite.Error(err)
		suite.Nil(output)
	})

	suite.T().Run("should return error if repo.Update fails", func(t *testing.T) {
		now := time.Now()
		existingSale := &domain.Sale{
			ID:        saleID,
			VehicleID: "vehicle-1",
			Brand:     "Toyota",
			Model:     "Corolla",
			Price:     50000,
			Status:    domain.StatusAvailable,
			CreatedAt: now.Add(-time.Hour),
			UpdatedAt: now.Add(-time.Hour),
		}

		usecase := usecase.NewSaleUseCase(suite.repository)
		suite.repository.EXPECT().GetByID(suite.ctx, saleID).Return(existingSale, nil)
		suite.repository.EXPECT().Update(suite.ctx, gomock.Any()).Return(errors.New("update failed"))

		output, err := usecase.Purchase(suite.ctx, saleID, input)
		suite.Error(err)
		suite.Nil(output)
	})
}

func (suite *SaleUseCaseSuite) Test_HandlePaymentWebhook() {
	paymentID := "payment-123"

	suite.T().Run("should update status to sold on APPROVED", func(t *testing.T) {
		now := time.Now()
		sale := &domain.Sale{
			ID:        "sale-1",
			VehicleID: "vehicle-1",
			Brand:     "Toyota",
			Model:     "Corolla",
			Price:     50000,
			Status:    domain.StatusPendingPayment,
			PaymentID: paymentID,
			CreatedAt: now.Add(-time.Hour),
			UpdatedAt: now.Add(-time.Hour),
		}

		usecase := usecase.NewSaleUseCase(suite.repository)
		input := &dto.InputWebhookDTO{
			PaymentID: paymentID,
			Status:    "APPROVED",
		}
		suite.repository.EXPECT().GetByPaymentID(suite.ctx, paymentID).Return(sale, nil)
		suite.repository.EXPECT().Update(suite.ctx, gomock.Any()).Return(nil)

		err := usecase.HandlePaymentWebhook(suite.ctx, input)
		suite.NoError(err)
	})

	suite.T().Run("should update status to sold on EFETUADO", func(t *testing.T) {
		now := time.Now()
		sale := &domain.Sale{
			ID:        "sale-1",
			VehicleID: "vehicle-1",
			Brand:     "Toyota",
			Model:     "Corolla",
			Price:     50000,
			Status:    domain.StatusPendingPayment,
			PaymentID: paymentID,
			CreatedAt: now.Add(-time.Hour),
			UpdatedAt: now.Add(-time.Hour),
		}

		usecase := usecase.NewSaleUseCase(suite.repository)
		input := &dto.InputWebhookDTO{
			PaymentID: paymentID,
			Status:    "EFETUADO",
		}
		suite.repository.EXPECT().GetByPaymentID(suite.ctx, paymentID).Return(sale, nil)
		suite.repository.EXPECT().Update(suite.ctx, gomock.Any()).Return(nil)

		err := usecase.HandlePaymentWebhook(suite.ctx, input)
		suite.NoError(err)
	})

	suite.T().Run("should update status to canceled on CANCELED", func(t *testing.T) {
		now := time.Now()
		sale := &domain.Sale{
			ID:        "sale-1",
			VehicleID: "vehicle-1",
			Brand:     "Toyota",
			Model:     "Corolla",
			Price:     50000,
			Status:    domain.StatusPendingPayment,
			PaymentID: paymentID,
			CreatedAt: now.Add(-time.Hour),
			UpdatedAt: now.Add(-time.Hour),
		}

		usecase := usecase.NewSaleUseCase(suite.repository)
		input := &dto.InputWebhookDTO{
			PaymentID: paymentID,
			Status:    "CANCELED",
		}
		suite.repository.EXPECT().GetByPaymentID(suite.ctx, paymentID).Return(sale, nil)
		suite.repository.EXPECT().Update(suite.ctx, gomock.Any()).Return(nil)

		err := usecase.HandlePaymentWebhook(suite.ctx, input)
		suite.NoError(err)
	})

	suite.T().Run("should update status to canceled on CANCELADO", func(t *testing.T) {
		now := time.Now()
		sale := &domain.Sale{
			ID:        "sale-1",
			VehicleID: "vehicle-1",
			Brand:     "Toyota",
			Model:     "Corolla",
			Price:     50000,
			Status:    domain.StatusPendingPayment,
			PaymentID: paymentID,
			CreatedAt: now.Add(-time.Hour),
			UpdatedAt: now.Add(-time.Hour),
		}

		usecase := usecase.NewSaleUseCase(suite.repository)
		input := &dto.InputWebhookDTO{
			PaymentID: paymentID,
			Status:    "CANCELADO",
		}
		suite.repository.EXPECT().GetByPaymentID(suite.ctx, paymentID).Return(sale, nil)
		suite.repository.EXPECT().Update(suite.ctx, gomock.Any()).Return(nil)

		err := usecase.HandlePaymentWebhook(suite.ctx, input)
		suite.NoError(err)
	})

	suite.T().Run("should return error for invalid status", func(t *testing.T) {
		now := time.Now()
		sale := &domain.Sale{
			ID:        "sale-1",
			VehicleID: "vehicle-1",
			Brand:     "Toyota",
			Model:     "Corolla",
			Price:     50000,
			Status:    domain.StatusPendingPayment,
			PaymentID: paymentID,
			CreatedAt: now.Add(-time.Hour),
			UpdatedAt: now.Add(-time.Hour),
		}

		usecase := usecase.NewSaleUseCase(suite.repository)
		input := &dto.InputWebhookDTO{
			PaymentID: paymentID,
			Status:    "INVALID_STATUS",
		}
		suite.repository.EXPECT().GetByPaymentID(suite.ctx, paymentID).Return(sale, nil)

		err := usecase.HandlePaymentWebhook(suite.ctx, input)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid payment status")
	})

	suite.T().Run("should return error if GetByPaymentID fails", func(t *testing.T) {
		usecase := usecase.NewSaleUseCase(suite.repository)
		input := &dto.InputWebhookDTO{
			PaymentID: paymentID,
			Status:    "APPROVED",
		}
		suite.repository.EXPECT().GetByPaymentID(suite.ctx, paymentID).Return(nil, errors.New("not found"))

		err := usecase.HandlePaymentWebhook(suite.ctx, input)
		suite.Error(err)
		suite.Contains(err.Error(), "not found")
	})

	suite.T().Run("should return error if sale is not pending payment", func(t *testing.T) {
		now := time.Now()
		sale := &domain.Sale{
			ID:        "sale-1",
			VehicleID: "vehicle-1",
			Brand:     "Toyota",
			Model:     "Corolla",
			Price:     50000,
			Status:    domain.StatusAvailable,
			PaymentID: paymentID,
			CreatedAt: now.Add(-time.Hour),
			UpdatedAt: now.Add(-time.Hour),
		}

		usecase := usecase.NewSaleUseCase(suite.repository)
		input := &dto.InputWebhookDTO{
			PaymentID: paymentID,
			Status:    "APPROVED",
		}
		suite.repository.EXPECT().GetByPaymentID(suite.ctx, paymentID).Return(sale, nil)

		err := usecase.HandlePaymentWebhook(suite.ctx, input)
		suite.Error(err)
		suite.Contains(err.Error(), "sale is not in pending payment status")
	})

	suite.T().Run("should return error if repo.Update fails", func(t *testing.T) {
		now := time.Now()
		sale := &domain.Sale{
			ID:        "sale-1",
			VehicleID: "vehicle-1",
			Brand:     "Toyota",
			Model:     "Corolla",
			Price:     50000,
			Status:    domain.StatusPendingPayment,
			PaymentID: paymentID,
			CreatedAt: now.Add(-time.Hour),
			UpdatedAt: now.Add(-time.Hour),
		}

		usecase := usecase.NewSaleUseCase(suite.repository)
		input := &dto.InputWebhookDTO{
			PaymentID: paymentID,
			Status:    "APPROVED",
		}
		suite.repository.EXPECT().GetByPaymentID(suite.ctx, paymentID).Return(sale, nil)
		suite.repository.EXPECT().Update(suite.ctx, gomock.Any()).Return(errors.New("update error"))

		err := usecase.HandlePaymentWebhook(suite.ctx, input)
		suite.Error(err)
		suite.Contains(err.Error(), "update error")
	})
}

func (suite *SaleUseCaseSuite) Test_ListAvailable() {
	suite.T().Run("should return available listings ordered by price", func(t *testing.T) {
		usecase := usecase.NewSaleUseCase(suite.repository)
		sales := []*domain.Sale{
			{
				ID:        "sale-1",
				VehicleID: "vehicle-1",
				Brand:     "Toyota",
				Model:     "Corolla",
				Price:     50000,
			},
			{
				ID:        "sale-2",
				VehicleID: "vehicle-2",
				Brand:     "Honda",
				Model:     "Civic",
				Price:     60000,
			},
		}
		suite.repository.EXPECT().GetAvailableByPrice(suite.ctx).Return(sales, nil)

		output, err := usecase.ListAvailable(suite.ctx)
		suite.NoError(err)
		suite.Len(output, 2)
		suite.Equal("sale-1", output[0].SaleID)
		suite.Equal("vehicle-1", output[0].VehicleID)
		suite.Equal("Toyota", output[0].Brand)
		suite.Equal("Corolla", output[0].Model)
		suite.Equal(float64(50000), output[0].Price)
		suite.Equal("sale-2", output[1].SaleID)
	})

	suite.T().Run("should return error if repo.GetAvailableByPrice fails", func(t *testing.T) {
		usecase := usecase.NewSaleUseCase(suite.repository)
		suite.repository.EXPECT().GetAvailableByPrice(suite.ctx).Return(nil, errors.New("db error"))

		output, err := usecase.ListAvailable(suite.ctx)
		suite.Error(err)
		suite.Nil(output)
	})
}

func (suite *SaleUseCaseSuite) Test_ListSold() {
	suite.T().Run("should return sold listings ordered by price", func(t *testing.T) {
		usecase := usecase.NewSaleUseCase(suite.repository)
		sales := []*domain.Sale{
			{
				ID:        "sale-1",
				VehicleID: "vehicle-1",
				Brand:     "Toyota",
				Model:     "Corolla",
				Price:     50000,
			},
			{
				ID:        "sale-2",
				VehicleID: "vehicle-2",
				Brand:     "Honda",
				Model:     "Civic",
				Price:     60000,
			},
		}
		suite.repository.EXPECT().GetSoldByPrice(suite.ctx).Return(sales, nil)

		output, err := usecase.ListSold(suite.ctx)
		suite.NoError(err)
		suite.Len(output, 2)
		suite.Equal("sale-1", output[0].SaleID)
		suite.Equal("vehicle-1", output[0].VehicleID)
		suite.Equal("Toyota", output[0].Brand)
		suite.Equal("Corolla", output[0].Model)
		suite.Equal(float64(50000), output[0].Price)
		suite.Equal("sale-2", output[1].SaleID)
	})

	suite.T().Run("should return error if repo.GetSoldByPrice fails", func(t *testing.T) {
		usecase := usecase.NewSaleUseCase(suite.repository)
		suite.repository.EXPECT().GetSoldByPrice(suite.ctx).Return(nil, errors.New("db error"))

		output, err := usecase.ListSold(suite.ctx)
		suite.Error(err)
		suite.Nil(output)
	})
}
