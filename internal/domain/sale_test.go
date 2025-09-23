package domain_test

import (
	"testing"

	"github.com/NicolasNSC/showcase-service-fiap/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewSale(t *testing.T) {
	t.Run("should create a new sale successfully", func(t *testing.T) {
		sale, err := domain.NewSale("vehicle-uuid", "Fiat", "Toro", 150000)

		assert.NoError(t, err)
		assert.NotNil(t, sale)
		assert.NotEmpty(t, sale.ID)
		assert.Equal(t, "vehicle-uuid", sale.VehicleID)
		assert.Equal(t, domain.StatusAvailable, sale.Status)
		assert.Equal(t, float64(150000), sale.Price)
		assert.NotZero(t, sale.CreatedAt)
	})

	t.Run("should return error for zero price", func(t *testing.T) {
		sale, err := domain.NewSale("vehicle-uuid", "Fiat", "Toro", 0)
		assert.Error(t, err)
		assert.Nil(t, sale)
		assert.Equal(t, "price must be greater than zero", err.Error())
	})

	t.Run("should return error for empty vehicle_id", func(t *testing.T) {
		sale, err := domain.NewSale("", "Fiat", "Toro", 150000)
		assert.Error(t, err)
		assert.Nil(t, sale)
		assert.Equal(t, "vehicle_id cannot be empty", err.Error())
	})
}
