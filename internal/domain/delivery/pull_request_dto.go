package delivery

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models/pr"
	"time"
)

type CreatePullRequestRequest struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
}

type CreatePullRequestResponse struct {
	Pr struct {
		PullRequestId     string    `json:"pull_request_id"`
		PullRequestName   string    `json:"pull_request_name"`
		AuthorId          string    `json:"author_id"`
		Status            pr.Status `json:"status"`
		AssignedReviewers []string  `json:"assigned_reviewers"`
	} `json:"pr"`
}

type MergePullRequestRequest struct {
	PullRequestId string `json:"pull_request_id"`
}

type MergePullRequestResponse struct {
	Pr struct {
		PullRequestId     string    `json:"pull_request_id"`
		PullRequestName   string    `json:"pull_request_name"`
		AuthorId          string    `json:"author_id"`
		Status            pr.Status `json:"status"`
		AssignedReviewers []string  `json:"assigned_reviewers"`
		MergedAt          time.Time `json:"mergedAt"`
	} `json:"pr"`
}

type ReassignPullRequestRequest struct {
	PullRequestId string `json:"pull_request_id"`
	OldReviewerId string `json:"old_reviewer_id"`
}

type ReassignPullRequestResponse struct {
	Pr struct {
		PullRequestId     string    `json:"pull_request_id"`
		PullRequestName   string    `json:"pull_request_name"`
		AuthorId          string    `json:"author_id"`
		Status            pr.Status `json:"status"`
		AssignedReviewers []string  `json:"assigned_reviewers"`
	} `json:"pr"`
	ReplacedBy string `json:"replaced_by"`
}
