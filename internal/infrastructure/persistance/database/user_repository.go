package database

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/delivery"
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models/user"
	"context"
)

type UserRepository interface {
	Update(ctx context.Context, curUser *user.User) (*user.User, error)
	GetReviewsByID(ctx context.Context, id string) (*delivery.UserPullRequest, error)
	FindReviewers(ctx context.Context, team string, author string) ([]string, error)
	FindNewReviewers(ctx context.Context, team string, author string, replaced string) ([]string, error)
	FindUserTeamById(ctx context.Context, id string) (string, error)
}
