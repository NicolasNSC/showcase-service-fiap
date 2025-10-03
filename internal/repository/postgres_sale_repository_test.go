package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/NicolasNSC/showcase-service-fiap/internal/domain"
	"github.com/NicolasNSC/showcase-service-fiap/internal/repository"
	"github.com/stretchr/testify/suite"
)

type PostgresSaleRepositoryTestSuite struct {
	suite.Suite
}

func Test_PostgresSaleRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(PostgresSaleRepositoryTestSuite))
}

func (suite *PostgresSaleRepositoryTestSuite) Test_Save() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.T().Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	repo := repository.NewPostgresSaleRepository(db)

	sale := &domain.Sale{
		ID:        "sale-id",
		VehicleID: "vehicle-id",
		Brand:     "BrandX",
		Model:     "ModelY",
		Price:     10000.0,
		Status:    "available",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	suite.T().Run("should save sale successfully", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO sales`).
			WithArgs(
				sale.ID,
				sale.VehicleID,
				sale.Brand,
				sale.Model,
				sale.Price,
				sale.Status,
				sale.CreatedAt,
				sale.UpdatedAt,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.Save(context.Background(), sale)
		suite.NoError(err)
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should return error when db fails", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO sales`).
			WithArgs(
				sale.ID,
				sale.VehicleID,
				sale.Brand,
				sale.Model,
				sale.Price,
				sale.Status,
				sale.CreatedAt,
				sale.UpdatedAt,
			).
			WillReturnError(errors.New("db error"))

		err = repo.Save(context.Background(), sale)
		suite.Error(err)
		suite.EqualError(err, "db error")
		suite.NoError(mock.ExpectationsWereMet())
	})
}

func (suite *PostgresSaleRepositoryTestSuite) Test_Update() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.T().Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	repo := repository.NewPostgresSaleRepository(db)

	now := time.Now()
	buyerCPF := "12345678900"
	saleDate := now.Add(-time.Hour)
	sale := &domain.Sale{
		ID:        "sale-id",
		VehicleID: "vehicle-id",
		Brand:     "BrandX",
		Model:     "ModelY",
		Price:     10000.0,
		Status:    "sold",
		PaymentID: "payment-id",
		BuyerCPF:  &buyerCPF,
		SaleDate:  &saleDate,
		UpdatedAt: now,
	}

	suite.T().Run("should update sale successfully", func(t *testing.T) {
		mock.ExpectExec(`UPDATE sales`).
			WithArgs(
				sale.VehicleID,
				sale.Brand,
				sale.Model,
				sale.Price,
				sale.Status,
				sql.NullString{String: sale.PaymentID, Valid: true},
				sql.NullString{String: buyerCPF, Valid: true},
				sql.NullTime{Time: saleDate, Valid: true},
				sale.UpdatedAt,
				sale.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Update(context.Background(), sale)
		suite.NoError(err)
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should update sale with nil BuyerCPF and SaleDate", func(t *testing.T) {
		saleNoCPF := *sale
		saleNoCPF.BuyerCPF = nil
		saleNoCPF.SaleDate = nil

		mock.ExpectExec(`UPDATE sales`).
			WithArgs(
				saleNoCPF.VehicleID,
				saleNoCPF.Brand,
				saleNoCPF.Model,
				saleNoCPF.Price,
				saleNoCPF.Status,
				sql.NullString{String: saleNoCPF.PaymentID, Valid: true},
				sql.NullString{Valid: false},
				sql.NullTime{Valid: false},
				saleNoCPF.UpdatedAt,
				saleNoCPF.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Update(context.Background(), &saleNoCPF)
		suite.NoError(err)
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should return error when db fails", func(t *testing.T) {
		mock.ExpectExec(`UPDATE sales`).
			WithArgs(
				sale.VehicleID,
				sale.Brand,
				sale.Model,
				sale.Price,
				sale.Status,
				sql.NullString{String: sale.PaymentID, Valid: true},
				sql.NullString{String: buyerCPF, Valid: true},
				sql.NullTime{Time: saleDate, Valid: true},
				sale.UpdatedAt,
				sale.ID,
			).
			WillReturnError(errors.New("db error"))

		err := repo.Update(context.Background(), sale)
		suite.Error(err)
		suite.EqualError(err, "db error")
		suite.NoError(mock.ExpectationsWereMet())
	})
}

func (suite *PostgresSaleRepositoryTestSuite) Test_GetByID() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.T().Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	repo := repository.NewPostgresSaleRepository(db)

	now := time.Now()
	buyerCPF := "12345678900"
	saleDate := now.Add(-time.Hour)

	suite.T().Run("should get sale by id successfully with all fields", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "vehicle_id", "brand", "model", "price", "status",
			"payment_id", "buyer_cpf", "sale_date", "created_at", "updated_at",
		}).
			AddRow(
				"sale-id", "vehicle-id", "BrandX", "ModelY", 10000.0, "SOLD",
				"payment-id", buyerCPF, saleDate, now, now,
			)

		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, payment_id, buyer_cpf, sale_date, created_at, updated_at FROM sales WHERE id = \$1`).
			WithArgs("sale-id").
			WillReturnRows(rows)

		sale, err := repo.GetByID(context.Background(), "sale-id")
		suite.NoError(err)
		suite.NotNil(sale)
		suite.Equal("sale-id", sale.ID)
		suite.Equal("vehicle-id", sale.VehicleID)
		suite.Equal("BrandX", sale.Brand)
		suite.Equal("ModelY", sale.Model)
		suite.Equal(10000.0, sale.Price)
		suite.Equal(domain.StatusSold, sale.Status)
		suite.Equal("payment-id", sale.PaymentID)
		suite.NotNil(sale.BuyerCPF)
		suite.Equal(buyerCPF, *sale.BuyerCPF)
		suite.NotNil(sale.SaleDate)
		suite.WithinDuration(saleDate, *sale.SaleDate, time.Second)
		suite.WithinDuration(now, sale.CreatedAt, time.Second)
		suite.WithinDuration(now, sale.UpdatedAt, time.Second)
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should get sale by id with nil BuyerCPF and SaleDate", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "vehicle_id", "brand", "model", "price", "status",
			"payment_id", "buyer_cpf", "sale_date", "created_at", "updated_at",
		}).
			AddRow(
				"sale-id", "vehicle-id", "BrandX", "ModelY", 10000.0, "sold",
				"payment-id", nil, nil, now, now,
			)

		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, payment_id, buyer_cpf, sale_date, created_at, updated_at FROM sales WHERE id = \$1`).
			WithArgs("sale-id").
			WillReturnRows(rows)

		sale, err := repo.GetByID(context.Background(), "sale-id")
		suite.NoError(err)
		suite.NotNil(sale)
		suite.Equal("payment-id", sale.PaymentID)
		suite.Nil(sale.BuyerCPF)
		suite.Nil(sale.SaleDate)
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should return error when sale not found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "vehicle_id", "brand", "model", "price", "status",
			"payment_id", "buyer_cpf", "sale_date", "created_at", "updated_at",
		})

		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, payment_id, buyer_cpf, sale_date, created_at, updated_at FROM sales WHERE id = \$1`).
			WithArgs("not-found-id").
			WillReturnRows(rows)

		sale, err := repo.GetByID(context.Background(), "not-found-id")
		suite.Error(err)
		suite.Nil(sale)
		suite.EqualError(err, "sale not found")
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should return error when db fails", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, payment_id, buyer_cpf, sale_date, created_at, updated_at FROM sales WHERE id = \$1`).
			WithArgs("sale-id").
			WillReturnError(errors.New("db error"))

		sale, err := repo.GetByID(context.Background(), "sale-id")
		suite.Error(err)
		suite.Nil(sale)
		suite.EqualError(err, "db error")
		suite.NoError(mock.ExpectationsWereMet())
	})
}

func (suite *PostgresSaleRepositoryTestSuite) Test_GetByVehicleID() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.T().Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	repo := repository.NewPostgresSaleRepository(db)

	now := time.Now()
	buyerCPF := "12345678900"
	saleDate := now.Add(-time.Hour)

	suite.T().Run("should get sale by vehicle_id successfully with all fields", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "vehicle_id", "brand", "model", "price", "status",
			"payment_id", "buyer_cpf", "sale_date", "created_at", "updated_at",
		}).
			AddRow(
				"sale-id", "vehicle-id", "BrandX", "ModelY", 10000.0, "SOLD",
				"payment-id", buyerCPF, saleDate, now, now,
			)

		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, payment_id, buyer_cpf, sale_date, created_at, updated_at FROM sales WHERE vehicle_id = \$1`).
			WithArgs("vehicle-id").
			WillReturnRows(rows)

		sale, err := repo.GetByVehicleID(context.Background(), "vehicle-id")
		suite.NoError(err)
		suite.NotNil(sale)
		suite.Equal("sale-id", sale.ID)
		suite.Equal("vehicle-id", sale.VehicleID)
		suite.Equal("BrandX", sale.Brand)
		suite.Equal("ModelY", sale.Model)
		suite.Equal(10000.0, sale.Price)
		suite.Equal(domain.StatusSold, sale.Status)
		suite.Equal("payment-id", sale.PaymentID)
		suite.NotNil(sale.BuyerCPF)
		suite.Equal(buyerCPF, *sale.BuyerCPF)
		suite.NotNil(sale.SaleDate)
		suite.WithinDuration(saleDate, *sale.SaleDate, time.Second)
		suite.WithinDuration(now, sale.CreatedAt, time.Second)
		suite.WithinDuration(now, sale.UpdatedAt, time.Second)
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should get sale by vehicle_id with nil BuyerCPF and SaleDate", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "vehicle_id", "brand", "model", "price", "status",
			"payment_id", "buyer_cpf", "sale_date", "created_at", "updated_at",
		}).
			AddRow(
				"sale-id", "vehicle-id", "BrandX", "ModelY", 10000.0, "sold",
				"payment-id", nil, nil, now, now,
			)

		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, payment_id, buyer_cpf, sale_date, created_at, updated_at FROM sales WHERE vehicle_id = \$1`).
			WithArgs("vehicle-id").
			WillReturnRows(rows)

		sale, err := repo.GetByVehicleID(context.Background(), "vehicle-id")
		suite.NoError(err)
		suite.NotNil(sale)
		suite.Equal("payment-id", sale.PaymentID)
		suite.Nil(sale.BuyerCPF)
		suite.Nil(sale.SaleDate)
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should return error when sale not found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "vehicle_id", "brand", "model", "price", "status",
			"payment_id", "buyer_cpf", "sale_date", "created_at", "updated_at",
		})

		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, payment_id, buyer_cpf, sale_date, created_at, updated_at FROM sales WHERE vehicle_id = \$1`).
			WithArgs("not-found-vehicle-id").
			WillReturnRows(rows)

		sale, err := repo.GetByVehicleID(context.Background(), "not-found-vehicle-id")
		suite.Error(err)
		suite.Nil(sale)
		suite.EqualError(err, "sale listing for the given vehicle_id not found")
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should return error when db fails", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, payment_id, buyer_cpf, sale_date, created_at, updated_at FROM sales WHERE vehicle_id = \$1`).
			WithArgs("vehicle-id").
			WillReturnError(errors.New("db error"))

		sale, err := repo.GetByVehicleID(context.Background(), "vehicle-id")
		suite.Error(err)
		suite.Nil(sale)
		suite.EqualError(err, "db error")
		suite.NoError(mock.ExpectationsWereMet())
	})
}

func (suite *PostgresSaleRepositoryTestSuite) Test_GetByPaymentID() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.T().Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	repo := repository.NewPostgresSaleRepository(db)

	now := time.Now()
	buyerCPF := "12345678900"
	saleDate := now.Add(-time.Hour)

	suite.T().Run("should get sale by payment_id successfully with all fields", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "vehicle_id", "brand", "model", "price", "status",
			"payment_id", "buyer_cpf", "sale_date", "created_at", "updated_at",
		}).
			AddRow(
				"sale-id", "vehicle-id", "BrandX", "ModelY", 10000.0, "SOLD",
				"payment-id", buyerCPF, saleDate, now, now,
			)

		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, payment_id, buyer_cpf, sale_date, created_at, updated_at FROM sales WHERE payment_id = \$1`).
			WithArgs("payment-id").
			WillReturnRows(rows)

		sale, err := repo.GetByPaymentID(context.Background(), "payment-id")
		suite.NoError(err)
		suite.NotNil(sale)
		suite.Equal("sale-id", sale.ID)
		suite.Equal("vehicle-id", sale.VehicleID)
		suite.Equal("BrandX", sale.Brand)
		suite.Equal("ModelY", sale.Model)
		suite.Equal(10000.0, sale.Price)
		suite.Equal(domain.StatusSold, sale.Status)
		suite.Equal("payment-id", sale.PaymentID)
		suite.NotNil(sale.BuyerCPF)
		suite.Equal(buyerCPF, *sale.BuyerCPF)
		suite.NotNil(sale.SaleDate)
		suite.WithinDuration(saleDate, *sale.SaleDate, time.Second)
		suite.WithinDuration(now, sale.CreatedAt, time.Second)
		suite.WithinDuration(now, sale.UpdatedAt, time.Second)
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should get sale by payment_id with nil BuyerCPF and SaleDate", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "vehicle_id", "brand", "model", "price", "status",
			"payment_id", "buyer_cpf", "sale_date", "created_at", "updated_at",
		}).
			AddRow(
				"sale-id", "vehicle-id", "BrandX", "ModelY", 10000.0, "sold",
				"payment-id", nil, nil, now, now,
			)

		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, payment_id, buyer_cpf, sale_date, created_at, updated_at FROM sales WHERE payment_id = \$1`).
			WithArgs("payment-id").
			WillReturnRows(rows)

		sale, err := repo.GetByPaymentID(context.Background(), "payment-id")
		suite.NoError(err)
		suite.NotNil(sale)
		suite.Equal("payment-id", sale.PaymentID)
		suite.Nil(sale.BuyerCPF)
		suite.Nil(sale.SaleDate)
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should return error when sale not found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "vehicle_id", "brand", "model", "price", "status",
			"payment_id", "buyer_cpf", "sale_date", "created_at", "updated_at",
		})

		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, payment_id, buyer_cpf, sale_date, created_at, updated_at FROM sales WHERE payment_id = \$1`).
			WithArgs("not-found-payment-id").
			WillReturnRows(rows)

		sale, err := repo.GetByPaymentID(context.Background(), "not-found-payment-id")
		suite.Error(err)
		suite.Nil(sale)
		suite.EqualError(err, "sale not found for the given payment_id")
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should return error when db fails", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, payment_id, buyer_cpf, sale_date, created_at, updated_at FROM sales WHERE payment_id = \$1`).
			WithArgs("payment-id").
			WillReturnError(errors.New("db error"))

		sale, err := repo.GetByPaymentID(context.Background(), "payment-id")
		suite.Error(err)
		suite.Nil(sale)
		suite.EqualError(err, "db error")
		suite.NoError(mock.ExpectationsWereMet())
	})
}

func (suite *PostgresSaleRepositoryTestSuite) Test_GetAvailableByPrice() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.T().Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	repo := repository.NewPostgresSaleRepository(db)

	now := time.Now()

	suite.T().Run("should return available sales ordered by price", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "vehicle_id", "brand", "model", "price", "status", "created_at", "updated_at",
		}).
			AddRow("sale-1", "vehicle-1", "BrandA", "ModelA", 5000.0, "AVAILABLE", now, now).
			AddRow("sale-2", "vehicle-2", "BrandB", "ModelB", 7000.0, "AVAILABLE", now, now)

		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, created_at, updated_at FROM sales WHERE status = \$1 ORDER BY price ASC`).
			WithArgs("AVAILABLE").
			WillReturnRows(rows)

		sales, err := repo.GetAvailableByPrice(context.Background())
		suite.NoError(err)
		suite.Len(sales, 2)
		suite.Equal("sale-1", sales[0].ID)
		suite.Equal(5000.0, sales[0].Price)
		suite.Equal("sale-2", sales[1].ID)
		suite.Equal(7000.0, sales[1].Price)
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should return empty slice if no available sales", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "vehicle_id", "brand", "model", "price", "status", "created_at", "updated_at",
		})

		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, created_at, updated_at FROM sales WHERE status = \$1 ORDER BY price ASC`).
			WithArgs("AVAILABLE").
			WillReturnRows(rows)

		sales, err := repo.GetAvailableByPrice(context.Background())
		suite.NoError(err)
		suite.Len(sales, 0)
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should return error when db fails", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, created_at, updated_at FROM sales WHERE status = \$1 ORDER BY price ASC`).
			WithArgs("AVAILABLE").
			WillReturnError(errors.New("db error"))

		sales, err := repo.GetAvailableByPrice(context.Background())
		suite.Error(err)
		suite.Nil(sales)
		suite.EqualError(err, "db error")
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should return error when scan fails", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "vehicle_id", "brand", "model", "price", "status", "created_at", "updated_at",
		}).
			AddRow("sale-1", "vehicle-1", "BrandA", "ModelA", "invalid-price", "AVAILABLE", now, now)

		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, created_at, updated_at FROM sales WHERE status = \$1 ORDER BY price ASC`).
			WithArgs("AVAILABLE").
			WillReturnRows(rows)

		sales, err := repo.GetAvailableByPrice(context.Background())
		suite.Error(err)
		suite.Nil(sales)
		suite.Contains(err.Error(), "Scan error")
		suite.NoError(mock.ExpectationsWereMet())
	})
}

func (suite *PostgresSaleRepositoryTestSuite) Test_GetSoldByPrice() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.T().Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	repo := repository.NewPostgresSaleRepository(db)

	now := time.Now()

	suite.T().Run("should return sold sales ordered by price", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "vehicle_id", "brand", "model", "price", "status", "created_at", "updated_at",
		}).
			AddRow("sale-1", "vehicle-1", "BrandA", "ModelA", 8000.0, "SOLD", now, now).
			AddRow("sale-2", "vehicle-2", "BrandB", "ModelB", 12000.0, "SOLD", now, now)

		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, created_at, updated_at FROM sales WHERE status = \$1 ORDER BY price ASC`).
			WithArgs("SOLD").
			WillReturnRows(rows)

		sales, err := repo.GetSoldByPrice(context.Background())
		suite.NoError(err)
		suite.Len(sales, 2)
		suite.Equal("sale-1", sales[0].ID)
		suite.Equal(8000.0, sales[0].Price)
		suite.Equal("sale-2", sales[1].ID)
		suite.Equal(12000.0, sales[1].Price)
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should return empty slice if no sold sales", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "vehicle_id", "brand", "model", "price", "status", "created_at", "updated_at",
		})

		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, created_at, updated_at FROM sales WHERE status = \$1 ORDER BY price ASC`).
			WithArgs("SOLD").
			WillReturnRows(rows)

		sales, err := repo.GetSoldByPrice(context.Background())
		suite.NoError(err)
		suite.Len(sales, 0)
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should return error when db fails", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, created_at, updated_at FROM sales WHERE status = \$1 ORDER BY price ASC`).
			WithArgs("SOLD").
			WillReturnError(errors.New("db error"))

		sales, err := repo.GetSoldByPrice(context.Background())
		suite.Error(err)
		suite.Nil(sales)
		suite.EqualError(err, "db error")
		suite.NoError(mock.ExpectationsWereMet())
	})

	suite.T().Run("should return error when scan fails", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "vehicle_id", "brand", "model", "price", "status", "created_at", "updated_at",
		}).
			AddRow("sale-1", "vehicle-1", "BrandA", "ModelA", "invalid-price", "SOLD", now, now)

		mock.ExpectQuery(`SELECT id, vehicle_id, brand, model, price, status, created_at, updated_at FROM sales WHERE status = \$1 ORDER BY price ASC`).
			WithArgs("SOLD").
			WillReturnRows(rows)

		sales, err := repo.GetSoldByPrice(context.Background())
		suite.Error(err)
		suite.Nil(sales)
		suite.Contains(err.Error(), "Scan error")
		suite.NoError(mock.ExpectationsWereMet())
	})
}
