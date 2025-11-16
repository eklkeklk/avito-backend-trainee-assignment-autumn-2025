package database

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models/pr"
	"context"
)

type PullRequestRepository interface {
	Create(ctx context.Context, request *pr.PullRequest) (*pr.PullRequest, error)
	Merge(ctx context.Context, request *pr.PullRequest) (*pr.PullRequest, error)
	Reassign(ctx context.Context, request *pr.PullRequest, reviewerId string, newReviewerId string) (*pr.PullRequest, string, error)
	FindAuthor(ctx context.Context, id string) (string, error)
	FindReviewers(ctx context.Context, id string) ([]string, error)
	IsOpen(ctx context.Context, id string) (bool, error)
}
