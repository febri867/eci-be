package service

import (
	"context"
	"errors"

	"eci-backend/internal/region/model"
	"eci-backend/internal/region/repository"
)

type RegionService struct {
	repo *repository.RegionRepository
}

func NewRegionService(repo *repository.RegionRepository) *RegionService {
	return &RegionService{
		repo: repo,
	}
}

func (s *RegionService) GetProvinces(
	ctx context.Context,
) ([]model.Region, error) {
	return s.repo.FindByLevelAndParent(ctx, 1, "")
}

func (s *RegionService) GetCities(
	ctx context.Context,
	provinceID string,
) ([]model.Region, error) {
	if provinceID == "" {
		return nil, errors.New("province ID is required")
	}

	return s.repo.FindByLevelAndParent(ctx, 2, provinceID)
}

func (s *RegionService) GetDistricts(
	ctx context.Context,
	cityID string,
) ([]model.Region, error) {
	if cityID == "" {
		return nil, errors.New("city ID is required")
	}

	return s.repo.FindByLevelAndParent(ctx, 3, cityID)
}
