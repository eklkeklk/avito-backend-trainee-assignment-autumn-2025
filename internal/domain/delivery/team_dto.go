package delivery

import "avito-backend-trainee-assignment-autumn-2025/internal/domain/models/user"

type CreateTeamResponse struct {
	Team struct {
		Name    string             `json:"team_name"`
		Members []*user.TeamMember `json:"members"`
	} `json:"team"`
}
