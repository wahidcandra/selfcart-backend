package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Transaction struct {
	ID              int64       `json:"id"`
	TransactionCode string      `json:"transaction_code"`
	CartID          int64       `json:"cart_id"`
	TotalAmount     int         `json:"total_amount"`
	PaymentStatus   int         `json:"payment_status"`
	PaymentMethod   int         `json:"payment_method"`
	PaymentBank     int         `json:"payment_amount"`
	CashierID       int         `json:"cashier_id"`
	CreatedAt       time.Time   `json:"created_at"`
	PaidAt          time.Time   `json:"paid_at"`
	Items           *[]CartItem `json:"items"`
}

type TransactionItem struct {
	ID            int64     `json:"id"`
	TransactionID int64     `json:"transaction_id"`
	ProductID     int64     `json:"product_id"`
	Qty           int       `json:"qty"`
	Price         int       `json:"price"`
	Discount      int       `json:"discount"`
	SubTotal      int       `json:"sub_total"`
	CreatedAt     time.Time `json:"created_at"`
}

type TransactionRepo struct {
	db *pgxpool.Pool
}

func NewTransactionRepo(db *pgxpool.Pool) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) Create(ctx context.Context, transactionCode string, cartId int64, totalAmount int) (*Transaction, error) {
	var transaction Transaction

	err := r.db.QueryRow(ctx, `
		INSERT INTO transactions (transaction_code, cart_id, total_amount,payment_status,cashier_id)
		VALUES ($1, $2, $3,'UNPAID',1)
		RETURNING id,transaction_code,cart_id,total_amount`, transactionCode, cartId, totalAmount).Scan(&transaction.ID, &transaction.TransactionCode, &transaction.CartID, &transaction.TotalAmount)
	return &transaction, err
}
