package service

import (
	"context"
	"selfcart/internal/repository"

	"github.com/google/uuid"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepo
}

func NewTransactionService(transactionRepo *repository.TransactionRepo) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo}
}
func (s *TransactionService) Create(ctx context.Context, cartId int64, totalAmount int) (*repository.Transaction, error) {
	transactionCode := "TRX-" + uuid.New().String()
	return s.transactionRepo.Create(ctx, transactionCode, cartId, totalAmount)
}
