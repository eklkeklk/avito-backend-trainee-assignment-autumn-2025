package database

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/stats"
	"context"
)

type StatsRepository interface {
	GetUserAssignmentStats(ctx context.Context) ([]stats.UserStats, error)
	GetPrReviewersStats(ctx context.Context) ([]stats.PrStats, error)
	GetOverallStats(ctx context.Context) (*stats.StatsResponse, error)
}
