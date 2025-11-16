package stats

type UserStats struct {
	UserId           string `json:"user_id"`
	Username         string `json:"username"`
	IsActive         bool   `json:"is_active"`
	AssignmentsCount int    `json:"assignments_count"`
}

type PrStats struct {
	PrId           string `json:"pr_id"`
	PrName         string `json:"pr_name"`
	ReviewersCount int    `json:"reviewers_count"`
}

type StatsResponse struct {
	UserAssignments  []UserStats `json:"user_assignments"`
	PrReviewers      []PrStats   `json:"pr_reviewers"`
	TotalAssignments int         `json:"total_assignments"`
	TotalPrs         int         `json:"total_prs"`
}
