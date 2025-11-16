package service

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models/team"
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models/user"
	"avito-backend-trainee-assignment-autumn-2025/internal/infrastructure/persistance/database"
	"context"
	"fmt"
)

type TeamServicer interface {
	GetTeam(ctx context.Context, name string) (*team.Team, error)
	CreateTeam(ctx context.Context, name string, members []*user.TeamMember) (*team.Team, error)
}

type TeamService struct {
	repository database.TeamRepository
}

func NewTeamService(repository database.TeamRepository) *TeamService {
	return &TeamService{
		repository: repository,
	}
}

func (t *TeamService) CreateTeam(ctx context.Context, name string, members []*user.TeamMember) (*team.Team, error) {
	if name == "" {
		return nil, fmt.Errorf("invalid request")
	}
	if len(members) == 0 {
		return nil, fmt.Errorf("invalid request")
	}
	newTeam := &team.Team{Name: name, Members: members}

	return t.repository.Create(ctx, newTeam)
}

func (t *TeamService) GetTeam(ctx context.Context, name string) (*team.Team, error) {
	if name == "" {
		return nil, fmt.Errorf("invalid request")
	}
	return t.repository.GetByName(ctx, name)
}
