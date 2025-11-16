package service

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/stats"
	"avito-backend-trainee-assignment-autumn-2025/internal/infrastructure/persistance/database"
	"context"
)

type StatsServicer interface {
	GetUserStats(ctx context.Context) ([]stats.UserStats, error)
	GetPrStats(ctx context.Context) ([]stats.PrStats, error)
	GetOverallStats(ctx context.Context) (*stats.StatsResponse, error)
}

type StatsRepositoryPostgres struct {
	repository database.StatsRepository
}

func NewStatsService(repository database.StatsRepository) *StatsRepositoryPostgres {
	return &StatsRepositoryPostgres{
		repository: repository,
	}
}

func (s *StatsRepositoryPostgres) GetUserStats(ctx context.Context) ([]stats.UserStats, error) {
	return s.repository.GetUserAssignmentStats(ctx)
}

func (s *StatsRepositoryPostgres) GetPrStats(ctx context.Context) ([]stats.PrStats, error) {
	return s.repository.GetPrReviewersStats(ctx)
}

func (s *StatsRepositoryPostgres) GetOverallStats(ctx context.Context) (*stats.StatsResponse, error) {
	return s.repository.GetOverallStats(ctx)
}
