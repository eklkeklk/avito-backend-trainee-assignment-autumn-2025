package database

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models/team"
	"context"
)

type TeamRepository interface {
	Create(ctx context.Context, team *team.Team) (*team.Team, error)
	GetByName(ctx context.Context, name string) (*team.Team, error)
}
