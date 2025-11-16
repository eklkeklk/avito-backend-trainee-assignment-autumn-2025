package team

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models/user"
)

type Team struct {
	Name    string             `json:"team_name"`
	Members []*user.TeamMember `json:"members"`
}
