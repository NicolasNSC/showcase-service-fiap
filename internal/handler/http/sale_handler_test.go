package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NicolasNSC/showcase-service-fiap/internal/dto"
	h "github.com/NicolasNSC/showcase-service-fiap/internal/handler/http"
	"github.com/NicolasNSC/showcase-service-fiap/internal/usecase/mocks"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type SaleHandlerSuite struct {
	suite.Suite

	ctx     context.Context
	useCase *mocks.MockSaleUseCaseInterface
	handler *h.SaleHandler
}

func (suite *SaleHandlerSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.ctx = context.Background()
	suite.useCase = mocks.NewMockSaleUseCaseInterface(ctrl)
	suite.handler = h.NewSaleHandler(suite.useCase)
}

func Test_SaleHandlerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(SaleHandlerSuite))
}

func (suite *SaleHandlerSuite) Test_CreateListing() {
	input := &dto.InputCreateListingDTO{
		VehicleID: "vehicle-id",
		Brand:     "Toyota",
		Model:     "Corolla",
		Price:     50000,
	}

	suite.T().Run("Create Listing - Success", func(t *testing.T) {
		expectedOutput := &dto.OutputCreateListingDTO{
			SaleID:    "sale-id",
			Status:    "AVAILABLE",
			CreatedAt: time.Now(),
		}

		suite.useCase.EXPECT().CreateListing(suite.ctx, input).Return(expectedOutput, nil)

		body, _ := json.Marshal(input)
		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodPost, "/listings", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		suite.handler.CreateListing(rr, req)

		suite.Equal(http.StatusCreated, rr.Code)
		suite.Equal("application/json", rr.Header().Get("Content-Type"))

		var resp dto.OutputCreateListingDTO
		err := json.NewDecoder(rr.Body).Decode(&resp)
		suite.NoError(err)
	})

	suite.T().Run("Create Listing - Invalid Body", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodPost, "/listings", bytes.NewReader([]byte("invalid json")))
		rr := httptest.NewRecorder()

		suite.handler.CreateListing(rr, req)

		suite.Equal(http.StatusBadRequest, rr.Code)
		suite.Contains(rr.Body.String(), "Invalid request body")
	})

	suite.T().Run("Create Listing - Use Case Error", func(t *testing.T) {
		expectedErr := errors.New("usecase error")

		suite.useCase.EXPECT().CreateListing(suite.ctx, input).Return(nil, expectedErr)

		body, _ := json.Marshal(input)
		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodPost, "/listings", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		suite.handler.CreateListing(rr, req)

		suite.Equal(http.StatusInternalServerError, rr.Code)
		suite.Contains(rr.Body.String(), expectedErr.Error())
	})
}

func (suite *SaleHandlerSuite) Test_ListAvailable() {
	suite.T().Run("List Available - Success", func(t *testing.T) {
		expectedOutput := []*dto.OutputSaleItemDTO{
			{
				SaleID:    "sale-id-1",
				VehicleID: "vehicle-id-1",
				Brand:     "Toyota",
				Model:     "Corolla",
				Price:     50000,
			},
			{
				SaleID:    "sale-id-2",
				VehicleID: "vehicle-id-2",
				Brand:     "Honda",
				Model:     "Civic",
				Price:     60000,
			},
		}

		suite.useCase.EXPECT().ListAvailable(suite.ctx).Return(expectedOutput, nil)

		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodGet, "/listings/available", nil)
		rr := httptest.NewRecorder()

		suite.handler.ListAvailable(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
		suite.Equal("application/json", rr.Header().Get("Content-Type"))

		var resp []*dto.OutputSaleItemDTO
		err := json.NewDecoder(rr.Body).Decode(&resp)
		suite.NoError(err)
		suite.Len(resp, 2)
		suite.Equal("sale-id-1", resp[0].SaleID)
		suite.Equal("sale-id-2", resp[1].SaleID)
	})

	suite.T().Run("List Available - Use Case Error", func(t *testing.T) {
		expectedErr := errors.New("usecase error")

		suite.useCase.EXPECT().ListAvailable(suite.ctx).Return(([]*dto.OutputSaleItemDTO)(nil), expectedErr)

		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodGet, "/listings/available", nil)
		rr := httptest.NewRecorder()

		suite.handler.ListAvailable(rr, req)

		suite.Equal(http.StatusInternalServerError, rr.Code)
		suite.Contains(rr.Body.String(), expectedErr.Error())
	})
}

func (suite *SaleHandlerSuite) Test_ListSold() {
	suite.T().Run("List Sold - Success", func(t *testing.T) {
		expectedOutput := []*dto.OutputSaleItemDTO{
			{
				SaleID:    "sale-id-1",
				VehicleID: "vehicle-id-1",
				Brand:     "Toyota",
				Model:     "Corolla",
				Price:     50000,
			},
			{
				SaleID:    "sale-id-2",
				VehicleID: "vehicle-id-2",
				Brand:     "Honda",
				Model:     "Civic",
				Price:     60000,
			},
		}

		suite.useCase.EXPECT().ListSold(suite.ctx).Return(expectedOutput, nil)

		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodGet, "/listings/sold", nil)
		rr := httptest.NewRecorder()

		suite.handler.ListSold(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
		suite.Equal("application/json", rr.Header().Get("Content-Type"))

		var resp []*dto.OutputSaleItemDTO
		err := json.NewDecoder(rr.Body).Decode(&resp)
		suite.NoError(err)
		suite.Len(resp, 2)
		suite.Equal("sale-id-1", resp[0].SaleID)
		suite.Equal("sale-id-2", resp[1].SaleID)
	})

	suite.T().Run("List Sold - Use Case Error", func(t *testing.T) {
		expectedErr := errors.New("usecase error")

		suite.useCase.EXPECT().ListSold(suite.ctx).Return(([]*dto.OutputSaleItemDTO)(nil), expectedErr)

		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodGet, "/listings/sold", nil)
		rr := httptest.NewRecorder()

		suite.handler.ListSold(rr, req)

		suite.Equal(http.StatusInternalServerError, rr.Code)
		suite.Contains(rr.Body.String(), "usecase error")
	})
}

func (suite *SaleHandlerSuite) Test_UpdateListing() {
	vehicleID := "vehicle-123"
	input := &dto.InputUpdateListingDTO{
		Brand: "Ford",
		Model: "Focus",
		Price: 45000,
	}

	suite.T().Run("Update Listing - Success", func(t *testing.T) {
		suite.useCase.EXPECT().UpdateListing(gomock.Any(), vehicleID, input).Return(nil)

		body, _ := json.Marshal(input)
		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodPut, "/listings/"+vehicleID, bytes.NewReader(body))
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"vehicle_id"},
				Values: []string{vehicleID},
			},
		}))
		rr := httptest.NewRecorder()

		suite.handler.UpdateListing(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
	})

	suite.T().Run("Update Listing - Missing Vehicle ID", func(t *testing.T) {
		body, _ := json.Marshal(input)
		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodPut, "/listings/", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		suite.handler.UpdateListing(rr, req)

		suite.Equal(http.StatusBadRequest, rr.Code)
		suite.Contains(rr.Body.String(), "Vehicle ID is required")
	})

	suite.T().Run("Update Listing - Invalid Body", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodPut, "/listings/"+vehicleID, bytes.NewReader([]byte("invalid json")))
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"vehicle_id"},
				Values: []string{vehicleID},
			},
		}))
		rr := httptest.NewRecorder()

		suite.handler.UpdateListing(rr, req)

		suite.Equal(http.StatusBadRequest, rr.Code)
		suite.Contains(rr.Body.String(), "Invalid request body")
	})

	suite.T().Run("Update Listing - Use Case Error", func(t *testing.T) {
		expectedErr := errors.New("update error")
		suite.useCase.EXPECT().UpdateListing(gomock.Any(), vehicleID, input).Return(expectedErr)

		body, _ := json.Marshal(input)
		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodPut, "/listings/"+vehicleID, bytes.NewReader(body))
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"vehicle_id"},
				Values: []string{vehicleID},
			},
		}))
		rr := httptest.NewRecorder()

		suite.handler.UpdateListing(rr, req)

		suite.Equal(http.StatusInternalServerError, rr.Code)
		suite.Contains(rr.Body.String(), expectedErr.Error())
	})
}

func (suite *SaleHandlerSuite) Test_Purchase() {
	saleID := "sale-123"
	input := dto.InputPurchaseDTO{
		BuyerCPF: "buyer-456",
	}

	suite.T().Run("Purchase - Success", func(t *testing.T) {
		expectedOutput := &dto.OutputPurchaseDTO{
			PaymentID: "payment-789",
		}
		suite.useCase.EXPECT().Purchase(gomock.Any(), saleID, input).Return(expectedOutput, nil)

		body, _ := json.Marshal(input)
		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodPost, "/listings/"+saleID+"/purchase", bytes.NewReader(body))
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{saleID},
			},
		}))
		rr := httptest.NewRecorder()

		suite.handler.Purchase(rr, req)

		suite.Equal(http.StatusAccepted, rr.Code)
		suite.Equal("application/json", rr.Header().Get("Content-Type"))

		var resp dto.OutputPurchaseDTO
		err := json.NewDecoder(rr.Body).Decode(&resp)
		suite.NoError(err)
		suite.Equal(expectedOutput.PaymentID, resp.PaymentID)
	})

	suite.T().Run("Purchase - Missing Sale ID", func(t *testing.T) {
		body, _ := json.Marshal(input)
		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodPost, "/listings//purchase", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		suite.handler.Purchase(rr, req)

		suite.Equal(http.StatusBadRequest, rr.Code)
		suite.Contains(rr.Body.String(), "Sale ID is required")
	})

	suite.T().Run("Purchase - Invalid Body", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodPost, "/listings/"+saleID+"/purchase", bytes.NewReader([]byte("invalid json")))
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{saleID},
			},
		}))
		rr := httptest.NewRecorder()

		suite.handler.Purchase(rr, req)

		suite.Equal(http.StatusBadRequest, rr.Code)
		suite.Contains(rr.Body.String(), "Invalid request body")
	})

	suite.T().Run("Purchase - Use Case Error", func(t *testing.T) {
		expectedErr := errors.New("purchase error")
		suite.useCase.EXPECT().Purchase(gomock.Any(), saleID, input).Return(nil, expectedErr)

		body, _ := json.Marshal(input)
		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodPost, "/listings/"+saleID+"/purchase", bytes.NewReader(body))
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{saleID},
			},
		}))
		rr := httptest.NewRecorder()

		suite.handler.Purchase(rr, req)

		suite.Equal(http.StatusInternalServerError, rr.Code)
		suite.Contains(rr.Body.String(), expectedErr.Error())
	})
}

func (suite *SaleHandlerSuite) Test_HandlePaymentWebhook() {
	suite.T().Run("HandlePaymentWebhook - Success", func(t *testing.T) {
		input := &dto.InputWebhookDTO{
			PaymentID: "payment-123",
			Status:    "PAID",
		}
		suite.useCase.EXPECT().HandlePaymentWebhook(gomock.Any(), input).Return(nil)

		body, _ := json.Marshal(input)
		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodPost, "/webhook/payment", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		suite.handler.HandlePaymentWebhook(rr, req)

		suite.Equal(http.StatusNoContent, rr.Code)
	})

	suite.T().Run("HandlePaymentWebhook - Invalid Body", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodPost, "/webhook/payment", bytes.NewReader([]byte("invalid json")))
		rr := httptest.NewRecorder()

		suite.handler.HandlePaymentWebhook(rr, req)

		suite.Equal(http.StatusBadRequest, rr.Code)
		suite.Contains(rr.Body.String(), "Invalid request body")
	})

	suite.T().Run("HandlePaymentWebhook - Use Case Error", func(t *testing.T) {
		input := &dto.InputWebhookDTO{
			PaymentID: "payment-456",
			Status:    "FAILED",
		}
		suite.useCase.EXPECT().HandlePaymentWebhook(gomock.Any(), input).Return(errors.New("webhook error"))

		body, _ := json.Marshal(input)
		req, _ := http.NewRequestWithContext(suite.ctx, http.MethodPost, "/webhook/payment", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		suite.handler.HandlePaymentWebhook(rr, req)

		suite.Equal(http.StatusInternalServerError, rr.Code)
		suite.Contains(rr.Body.String(), "Failed to process webhook")
	})
}
