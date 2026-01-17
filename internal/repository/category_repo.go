package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{db: db}
}

type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Image     *string   `json:"image"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *CategoryRepo) GetAll(ctx context.Context) ([]Category, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, image, created_at
		FROM categories
		ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Image, &c.CreatedAt); err != nil {
			return nil, err
		}
		data = append(data, c)
	}
	return data, nil
}
