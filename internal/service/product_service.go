package service

import (
	"context"
	"selfcart/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepo
}

func NewProductService(r *repository.ProductRepo) *ProductService {
	return &ProductService{repo: r}
}

func (s *ProductService) GetAll(ctx context.Context) ([]repository.Product, error) {
	return s.repo.GetAll(ctx)
}
func (s *ProductService) GetByBarcode(ctx context.Context, barcode string) (*repository.Product, error) {
	return s.repo.GetByBarcode(ctx, barcode)
}
func (s *ProductService) GetByCategory(ctx context.Context, categoryID int64) ([]repository.Product, error) {
	return s.repo.GetByCategory(ctx, categoryID)
}
