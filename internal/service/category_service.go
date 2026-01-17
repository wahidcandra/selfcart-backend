package service

import (
	"context"
	"selfcart/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepo
}

func NewCategoryService(r *repository.CategoryRepo) *CategoryService {
	return &CategoryService{repo: r}
}

func (s *CategoryService) GetAll(ctx context.Context) ([]repository.Category, error) {
	return s.repo.GetAll(ctx)
}
