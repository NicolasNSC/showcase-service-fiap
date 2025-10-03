package domain_test

import (
	"testing"

	"github.com/NicolasNSC/showcase-service-fiap/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewSale_AllScenarios(t *testing.T) {
	t.Run("should create a new sale successfully", func(t *testing.T) {
		sale, err := domain.NewSale("vehicle-uuid", "Fiat", "Toro", 150000)
		assert.NoError(t, err)
		assert.NotNil(t, sale)
		assert.Equal(t, "vehicle-uuid", sale.VehicleID)
		assert.Equal(t, "Fiat", sale.Brand)
		assert.Equal(t, "Toro", sale.Model)
		assert.Equal(t, 150000.0, sale.Price)
		assert.Equal(t, domain.StatusAvailable, sale.Status)
		assert.NotEmpty(t, sale.ID)
		assert.NotZero(t, sale.CreatedAt)
		assert.NotZero(t, sale.UpdatedAt)
	})

	t.Run("should return error for empty brand", func(t *testing.T) {
		sale, err := domain.NewSale("vehicle-uuid", "", "Toro", 150000)
		assert.Error(t, err)
		assert.Nil(t, sale)
		assert.Equal(t, "brand and model are required for listing", err.Error())
	})

	t.Run("should return error for empty model", func(t *testing.T) {
		sale, err := domain.NewSale("vehicle-uuid", "Fiat", "", 150000)
		assert.Error(t, err)
		assert.Nil(t, sale)
		assert.Equal(t, "brand and model are required for listing", err.Error())
	})

	t.Run("should return error for empty brand and model", func(t *testing.T) {
		sale, err := domain.NewSale("vehicle-uuid", "", "", 150000)
		assert.Error(t, err)
		assert.Nil(t, sale)
		assert.Equal(t, "brand and model are required for listing", err.Error())
	})

	t.Run("should return error for negative price", func(t *testing.T) {
		sale, err := domain.NewSale("vehicle-uuid", "Fiat", "Toro", -100)
		assert.Error(t, err)
		assert.Nil(t, sale)
		assert.Equal(t, "price must be greater than zero", err.Error())
	})

	t.Run("should return error for zero price and empty brand/model", func(t *testing.T) {
		sale, err := domain.NewSale("vehicle-uuid", "", "", 0)
		assert.Error(t, err)
		assert.Nil(t, sale)
		assert.Equal(t, "price must be greater than zero", err.Error())
	})

	t.Run("should return error when vehicle_id is empty along with other invalid fields", func(t *testing.T) {
		sale, err := domain.NewSale("", "", "", -50)
		assert.Error(t, err)
		assert.Nil(t, sale)
		assert.Equal(t, "vehicle_id cannot be empty", err.Error())
	})
}
