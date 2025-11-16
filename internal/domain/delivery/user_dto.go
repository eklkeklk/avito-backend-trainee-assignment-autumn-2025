package delivery

import "avito-backend-trainee-assignment-autumn-2025/internal/domain/models/pr"

type SetActiveStatusRequest struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type SetActiveStatusResponse struct {
	User struct {
		Id       string `json:"user_id"`
		Name     string `json:"username"`
		TeamName string `json:"team_name"`
		IsActive bool   `json:"is_active"`
	} `json:"user"`
}

type PullRequestShort struct {
	Id     string    `json:"pull_request_id"`
	Name   string    `json:"pull_request_name"`
	Author string    `json:"author_id"`
	Status pr.Status `json:"status"`
}

type UserPullRequest struct {
	UserId       string              `json:"user_id"`
	PullRequests []*PullRequestShort `json:"pull_requests"`
}
