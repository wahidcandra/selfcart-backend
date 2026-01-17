package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartRepo struct {
	db *pgxpool.Pool
}

func NewCartRepo(db *pgxpool.Pool) *CartRepo {
	return &CartRepo{db: db}
}

type Cart struct {
	ID         int64       `json:"id"`
	CustomerID int64       `json:"customer_id"`
	CartCode   string      `json:"cart_code"`
	Status     string      `json:"status"`
	Discount   int         `json:"discount"`
	Total      int         `json:"total"`
	Items      *[]CartItem `json:"items"`
}

type CartItem struct {
	ID        int64   `json:"id"`
	ProductID int64   `json:"product_id"`
	Qty       int     `json:"qty"`
	Price     int     `json:"price"`
	Discount  int     `json:"discount"`
	SubTotal  int     `json:"sub_total"`
	Name      *string `json:"name"`
	Barcode   *string `json:"barcode"`
}

func (r *CartRepo) GetCart(ctx context.Context, cartID int64) (*Cart, error) {
	var cart Cart

	err := r.db.QueryRow(ctx, `
		SELECT a.id, a.customer_id, a.cart_code, a.status,sum(b.discount) as discount,sum(b.sub_total) as total
		FROM carts a
		LEFT JOIN cart_items b ON a.id = b.cart_id
		WHERE a.id = $1
		GROUP BY a.id, a.customer_id, a.cart_code, a.status`, cartID).Scan(&cart.ID, &cart.CustomerID, &cart.CartCode, &cart.Status, &cart.Discount, &cart.Total)
	return &cart, err
}
func (r *CartRepo) UpdateStatus(ctx context.Context, cartID int64, status string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE carts
		SET status = $1
		WHERE id = $2`, status, cartID)
	return err
}

func (r *CartRepo) Create(ctx context.Context, customerID int64, cartCode string) (*Cart, error) {
	var cart Cart

	err := r.db.QueryRow(ctx, `
		INSERT INTO carts (customer_id, cart_code, status)
		VALUES ($1, $2, 'ACTIVE')
		RETURNING id,customer_id,cart_code,status`, customerID, cartCode).Scan(&cart.ID, &cart.CustomerID, &cart.CartCode, &cart.Status)
	return &cart, err
}

func (r *CartRepo) AddItem(ctx context.Context, cartId int64, item CartItem) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO cart_items (cart_id, product_id, qty,price,discount,sub_total)
		VALUES ($1, $2, $3,$4,$5,$6)`,
		cartId, item.ProductID, item.Qty, item.Price, item.Discount, item.SubTotal)

	return err
}

func (r *CartRepo) GetItem(ctx context.Context, cartID, productID int64) (*CartItem, error) {
	var item CartItem

	err := r.db.QueryRow(ctx, `
		SELECT id, product_id, qty,price,discount,sub_total
		FROM cart_items
		WHERE cart_id = $1 AND product_id = $2`, cartID, productID).Scan(&item.ID, &item.ProductID, &item.Qty, &item.Price, &item.Discount, &item.SubTotal)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Item not found - return empty item, not an error
			return &CartItem{}, nil
		}
		return nil, err
	}
	return &item, nil
}
func (r *CartRepo) UpdateItem(ctx context.Context, id int64, item CartItem) error {
	_, err := r.db.Exec(ctx, `
		UPDATE cart_items
		SET qty = $1, price = $2, discount = $3, sub_total = $4
		WHERE id = $5`, item.Qty, item.Price, item.Discount, item.SubTotal, id)
	return err
}
func (r *CartRepo) DeleteItem(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM cart_items
		WHERE id = $1`, id)
	return err
}

func (r *CartRepo) GetItems(ctx context.Context, cartID int64) ([]CartItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT a.id, a.product_id, a.qty,a.price,a.discount,a.sub_total,b.name,b.upc
		FROM cart_items a
		LEFT JOIN products b ON a.product_id = b.id
		WHERE a.cart_id = $1`, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem
	for rows.Next() {
		var i CartItem
		if err := rows.Scan(&i.ID, &i.ProductID, &i.Qty, &i.Price, &i.Discount, &i.SubTotal, &i.Name, &i.Barcode); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}
