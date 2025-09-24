package repository

import (
	"context"
	"database/sql"

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
