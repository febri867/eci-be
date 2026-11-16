package service

import (
	"context"

	historyModel "eci-backend/internal/history/model"
	historyRepository "eci-backend/internal/history/repository"
)

type HistoryService struct {
	repository *historyRepository.ElasticsearchRepository
}

func NewHistoryService(
	repo *historyRepository.ElasticsearchRepository,
) *HistoryService {
	return &HistoryService{
		repository: repo,
	}
}

func (s *HistoryService) FindAll(
	ctx context.Context,
	page int,
	limit int,
	filter historyRepository.HistoryFilter,
) ([]historyModel.OutboxEvent, int64, historyModel.HistorySummary, error) {
	return s.repository.FindAll(ctx, page, limit, filter)
}
