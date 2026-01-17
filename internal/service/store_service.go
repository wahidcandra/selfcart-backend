package service

import (
	"context"
	"selfcart/internal/repository"
)

type StoreService struct {
	repo *repository.StoreRepo
}

func NewStoreService(r *repository.StoreRepo) *StoreService {
	return &StoreService{repo: r}
}

func (s *StoreService) GetStore(ctx context.Context) (*repository.Store, error) {
	return s.repo.GetStore(ctx)
}
func (s *StoreService) GetRack(ctx context.Context, storeID int64) ([]repository.Rack, error) {
	return s.repo.GetRack(ctx, storeID)
}
func (s *StoreService) GetRackZone(ctx context.Context, storeID int64, zone string) ([]repository.Rack, error) {
	return s.repo.GetRackZone(ctx, storeID, zone)
}
