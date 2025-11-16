package pr

import "time"

type Status string

const (
	Open   Status = "OPEN"
	Merged Status = "MERGED"
)

type PullRequest struct {
	Id                string    `json:"pull_request_id"`
	Name              string    `json:"pull_request_name"`
	Author            string    `json:"author_id"`
	Status            Status    `json:"status"`
	Reviewers         []string  `json:"assigned_reviewers"`
	NeedMoreReviewers bool      `json:"need_more_reviewers"`
	CreatedAt         time.Time `json:"created_at"`
	MergedAt          time.Time `json:"merged_at"`
}
