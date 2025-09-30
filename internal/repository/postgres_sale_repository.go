package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/NicolasNSC/showcase-service-fiap/internal/domain"
)

type postgresSaleRepository struct {
	db *sql.DB
}

func NewPostgresSaleRepository(db *sql.DB) SaleRepository {
	return &postgresSaleRepository{
		db: db,
	}
}

func (r *postgresSaleRepository) Save(ctx context.Context, sale *domain.Sale) error {
	query := `INSERT INTO sales (id, vehicle_id, brand, model, price, status, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.ExecContext(ctx, query,
		sale.ID,
		sale.VehicleID,
		sale.Brand,
		sale.Model,
		sale.Price,
		sale.Status,
		sale.CreatedAt,
		sale.UpdatedAt,
	)

	return err
}

func (r *postgresSaleRepository) Update(ctx context.Context, sale *domain.Sale) error {
	query := `UPDATE sales 
	          SET vehicle_id = $1, brand = $2, model = $3, price = $4, status = $5, 
	              payment_id = $6, buyer_cpf = $7, sale_date = $8, updated_at = $9
	          WHERE id = $10`

	var paymentID, buyerCPF sql.NullString
	var saleDate sql.NullTime

	if sale.PaymentID != "" {
		paymentID = sql.NullString{String: sale.PaymentID, Valid: true}
	}
	if sale.BuyerCPF != nil {
		buyerCPF = sql.NullString{String: *sale.BuyerCPF, Valid: true}
	}
	if sale.SaleDate != nil {
		saleDate = sql.NullTime{Time: *sale.SaleDate, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query,
		sale.VehicleID,
		sale.Brand,
		sale.Model,
		sale.Price,
		sale.Status,
		paymentID,
		buyerCPF,
		saleDate,
		sale.UpdatedAt,
		sale.ID,
	)

	return err
}

func (r *postgresSaleRepository) GetByID(ctx context.Context, id string) (*domain.Sale, error) {
	query := `SELECT id, vehicle_id, brand, model, price, status, payment_id, buyer_cpf, sale_date, created_at, updated_at 
	          FROM sales 
	          WHERE id = $1`

	var s domain.Sale
	var paymentID, buyerCPF sql.NullString
	var saleDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.VehicleID, &s.Brand, &s.Model, &s.Price, &s.Status,
		&paymentID, &buyerCPF, &saleDate,
		&s.CreatedAt, &s.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("sale not found")
		}

		return nil, err
	}

	if paymentID.Valid {
		s.PaymentID = paymentID.String
	}
	if buyerCPF.Valid {
		s.BuyerCPF = &buyerCPF.String
	}
	if saleDate.Valid {
		s.SaleDate = &saleDate.Time
	}

	return &s, nil
}

func (r *postgresSaleRepository) GetByVehicleID(ctx context.Context, vehicleID string) (*domain.Sale, error) {
	query := `SELECT id, vehicle_id, brand, model, price, status, payment_id, buyer_cpf, sale_date, created_at, updated_at 
	          FROM sales 
	          WHERE vehicle_id = $1`

	var s domain.Sale
	var paymentID, buyerCPF sql.NullString
	var saleDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, vehicleID).Scan(
		&s.ID, &s.VehicleID, &s.Brand, &s.Model, &s.Price, &s.Status,
		&paymentID, &buyerCPF, &saleDate,
		&s.CreatedAt, &s.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("sale listing for the given vehicle_id not found")
		}
		return nil, err
	}

	if paymentID.Valid {
		s.PaymentID = paymentID.String
	}
	if buyerCPF.Valid {
		s.BuyerCPF = &buyerCPF.String
	}
	if saleDate.Valid {
		s.SaleDate = &saleDate.Time
	}

	return &s, nil
}

func (r *postgresSaleRepository) GetByPaymentID(ctx context.Context, paymentID string) (*domain.Sale, error) {
	query := `SELECT id, vehicle_id, brand, model, price, status, payment_id, buyer_cpf, sale_date, created_at, updated_at 
	          FROM sales 
	          WHERE payment_id = $1`

	var s domain.Sale
	var pID, buyerCPF sql.NullString
	var saleDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, paymentID).Scan(
		&s.ID, &s.VehicleID, &s.Brand, &s.Model, &s.Price, &s.Status,
		&pID, &buyerCPF, &saleDate,
		&s.CreatedAt, &s.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("sale not found for the given payment_id")
		}
		return nil, err
	}

	if pID.Valid {
		s.PaymentID = pID.String
	}
	if buyerCPF.Valid {
		s.BuyerCPF = &buyerCPF.String
	}
	if saleDate.Valid {
		s.SaleDate = &saleDate.Time
	}

	return &s, nil
}

func (r *postgresSaleRepository) GetAvailableByPrice(ctx context.Context) ([]*domain.Sale, error) {
	query := `SELECT id, vehicle_id, brand, model, price, status, created_at, updated_at 
	          FROM sales 
	          WHERE status = $1 
	          ORDER BY price ASC`

	rows, err := r.db.QueryContext(ctx, query, domain.StatusAvailable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []*domain.Sale
	for rows.Next() {
		var s domain.Sale
		if err := rows.Scan(&s.ID, &s.VehicleID, &s.Brand, &s.Model, &s.Price, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sales = append(sales, &s)
	}

	return sales, nil
}

func (r *postgresSaleRepository) GetSoldByPrice(ctx context.Context) ([]*domain.Sale, error) {
	query := `SELECT id, vehicle_id, brand, model, price, status, created_at, updated_at 
	          FROM sales 
	          WHERE status = $1 
	          ORDER BY price ASC`

	rows, err := r.db.QueryContext(ctx, query, domain.StatusSold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []*domain.Sale
	for rows.Next() {
		var s domain.Sale
		if err := rows.Scan(&s.ID, &s.VehicleID, &s.Brand, &s.Model, &s.Price, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sales = append(sales, &s)
	}

	return sales, nil
}
