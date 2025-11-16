package postgres

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/stats"
	"context"
	"database/sql"
)

type StatsPostgresRepository struct {
	db *sql.DB
}

func NewStatsPostgresRepository(db *sql.DB) *StatsPostgresRepository {
	return &StatsPostgresRepository{
		db: db,
	}
}

func (s StatsPostgresRepository) GetUserAssignmentStats(ctx context.Context) ([]stats.UserStats, error) {
	query := `
		SELECT 
			u.user_id,
			u.username,
			u.is_active,
			COUNT(pr.reviewer_id) as assignments_count
		FROM users u
		LEFT JOIN pull_request_reviewers pr ON u.user_id = pr.reviewer_id
		GROUP BY u.user_id, u.username
		ORDER BY assignments_count DESC
	`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)
	var userStats []stats.UserStats
	for rows.Next() {
		var stat stats.UserStats
		err := rows.Scan(&stat.UserId, &stat.Username, &stat.IsActive, &stat.AssignmentsCount)
		if err != nil {
			return nil, err
		}
		userStats = append(userStats, stat)
	}

	return userStats, nil
}

func (s StatsPostgresRepository) GetPrReviewersStats(ctx context.Context) ([]stats.PrStats, error) {
	query := `
		SELECT 
			p.pull_request_id,
			p.pull_request_name,
			COUNT(pr.reviewer_id) as reviewers_count
		FROM pull_requests p
		LEFT JOIN pull_request_reviewers pr ON p.pull_request_id = pr.pull_request_id
		WHERE p.status = 'OPEN'
		GROUP BY p.pull_request_id, p.pull_request_name
		ORDER BY reviewers_count DESC
	`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)
	var prStats []stats.PrStats
	for rows.Next() {
		var stat stats.PrStats
		err := rows.Scan(&stat.PrId, &stat.PrName, &stat.ReviewersCount)
		if err != nil {
			return nil, err
		}
		prStats = append(prStats, stat)
	}

	return prStats, nil
}

func (s StatsPostgresRepository) GetOverallStats(ctx context.Context) (*stats.StatsResponse, error) {
	userStats, err := s.GetUserAssignmentStats(ctx)
	if err != nil {
		return nil, err
	}
	prStats, err := s.GetPrReviewersStats(ctx)
	if err != nil {
		return nil, err
	}
	totalAssignments := 0
	for _, user := range userStats {
		totalAssignments += user.AssignmentsCount
	}
	return &stats.StatsResponse{
		UserAssignments:  userStats,
		PrReviewers:      prStats,
		TotalAssignments: totalAssignments,
		TotalPrs:         len(prStats),
	}, nil
}
