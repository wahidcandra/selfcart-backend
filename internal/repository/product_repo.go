package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{db: db}
}

type Product struct {
	ID         int64     `json:"id"`
	Sku        string    `json:"sku"`
	Upc        string    `json:"upc"`
	Name       string    `json:"name"`
	Image      *string   `json:"image"`
	Price      float64   `json:"price"`
	CategoryID int64     `json:"category_id"`
	RackID     int64     `json:"rack_id"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

func (r *ProductRepo) GetAll(ctx context.Context) ([]Product, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, sku, upc, name, image, price, category_id, rack_id, is_active, created_at
		FROM products
		WHERE is_active = true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Sku, &p.Upc, &p.Name, &p.Image, &p.Price, &p.CategoryID, &p.RackID, &p.IsActive, &p.CreatedAt); err != nil {
			return nil, errors.New("product not found")
		}
		data = append(data, p)
	}
	return data, nil
}
func (r *ProductRepo) GetByBarcode(ctx context.Context, barcode string) (*Product, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, sku, upc, name, image, price, category_id, rack_id, is_active, created_at
		FROM products
		WHERE upc = $1 AND is_active = true`, barcode)

	var p Product
	if err := row.Scan(&p.ID, &p.Sku, &p.Upc, &p.Name, &p.Image, &p.Price, &p.CategoryID, &p.RackID, &p.IsActive, &p.CreatedAt); err != nil {
		return nil, errors.New("product not found")
	}
	return &p, nil
}
func (r *ProductRepo) GetByCategory(ctx context.Context, categoryID int64) ([]Product, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, sku, upc, name, image, price, category_id, rack_id, is_active, created_at
		FROM products
		WHERE category_id = $1 AND is_active = true`, categoryID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Sku, &p.Upc, &p.Name, &p.Image, &p.Price, &p.CategoryID, &p.RackID, &p.IsActive, &p.CreatedAt); err != nil {
			return nil, errors.New("product not found")
		}
		data = append(data, p)
	}
	return data, nil
}
