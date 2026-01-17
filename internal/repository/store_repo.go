package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StoreRepo struct {
	db *pgxpool.Pool
}

func NewStoreRepo(db *pgxpool.Pool) *StoreRepo {
	return &StoreRepo{db: db}
}

type Store struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Logo      *string   `json:"logo"`
	Address   *string   `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}
type Rack struct {
	ID          int64     `json:"id"`
	StoreID     int64     `json:"store_id"`
	RackCode    string    `json:"rack_code"`
	Zone        *string   `json:"zone"`
	XCoordinate *string   `json:"x_coordinate"`
	YCoordinate *string   `json:"y_coordinate"`
	CreatedAt   time.Time `json:"created_at"`
}

func (r *StoreRepo) GetStore(ctx context.Context) (*Store, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, logo, address, created_at
		FROM stores
		ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data *Store
	for rows.Next() {
		var s Store
		if err := rows.Scan(&s.ID, &s.Name, &s.Logo, &s.Address, &s.CreatedAt); err != nil {
			return nil, err
		}
		data = &s
	}
	return data, nil
}

func (r *StoreRepo) GetRack(ctx context.Context, storeID int64) ([]Rack, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, store_id, rack_code, zone, x_coordinate, y_coordinate, created_at
		FROM store_racks
		WHERE store_id = $1
		ORDER BY id ASC`, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []Rack
	for rows.Next() {
		var s Rack
		if err := rows.Scan(&s.ID, &s.StoreID, &s.RackCode, &s.Zone, &s.XCoordinate, &s.YCoordinate, &s.CreatedAt); err != nil {
			return nil, err
		}
		data = append(data, s)
	}
	return data, nil
}

func (r *StoreRepo) GetRackZone(ctx context.Context, storeID int64, zone string) ([]Rack, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, store_id, rack_code, zone, x_coordinate, y_coordinate, created_at
		FROM store_racks
		WHERE store_id = $1 AND zone = $2
		ORDER BY id ASC`, storeID, zone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []Rack
	for rows.Next() {
		var s Rack
		if err := rows.Scan(&s.ID, &s.StoreID, &s.RackCode, &s.Zone, &s.XCoordinate, &s.YCoordinate, &s.CreatedAt); err != nil {
			return nil, err
		}
		data = append(data, s)
	}
	return data, nil
}
