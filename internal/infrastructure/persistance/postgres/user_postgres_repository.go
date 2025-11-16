package postgres

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/delivery"
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models/pr"
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models/user"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type UserRepositoryPostgres struct {
	db *sql.DB
}

func NewUserRepositoryPostgres(db *sql.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{
		db: db,
	}
}

func (u UserRepositoryPostgres) Update(ctx context.Context, curUser *user.User) (*user.User, error) {
	query := `
	UPDATE users 
	SET is_active=$1
	WHERE user_id=$2
	RETURNING user_id, username, team_name, is_active
	`
	var updatedUser user.User
	err := u.db.QueryRowContext(
		ctx,
		query,
		curUser.IsActive,
		curUser.Id,
	).Scan(
		&updatedUser.Id,
		&updatedUser.Name,
		&updatedUser.TeamName,
		&updatedUser.IsActive,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("resource not found")
		}
		return nil, fmt.Errorf("internal error")
	}

	return &updatedUser, nil
}

func (u UserRepositoryPostgres) GetReviewsByID(ctx context.Context, id string) (*delivery.UserPullRequest, error) {
	var userExists bool
	checkUserQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)`
	err := u.db.QueryRowContext(ctx, checkUserQuery, id).Scan(&userExists)
	if err != nil {
		return nil, fmt.Errorf("internal error")
	}
	if !userExists {
		return nil, fmt.Errorf("resource not found")
	}
	pullRequests, err := u.getPullRequestById(ctx, id)
	if err != nil {
		return nil, err
	}
	result := &delivery.UserPullRequest{
		UserId:       id,
		PullRequests: pullRequests,
	}

	return result, nil
}

func (u UserRepositoryPostgres) FindReviewers(ctx context.Context, team string, author string) ([]string, error) {
	query := `
        SELECT user_id 
        FROM users 
        WHERE team_name = $1 AND is_active = true AND user_id != $2
    `
	rows, err := u.db.QueryContext(ctx, query, team, author)
	if err != nil {
		return nil, fmt.Errorf("internal error")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)
	var userIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("internal error")
		}
		userIDs = append(userIDs, userID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("internal error")
	}

	return userIDs, nil
}

func (u UserRepositoryPostgres) FindNewReviewers(ctx context.Context, team string, author string, replaced string) ([]string, error) {
	query := `
        SELECT user_id 
        FROM users 
        WHERE team_name = $1 AND is_active = true AND user_id != $2 AND user_id != $3
    `
	rows, err := u.db.QueryContext(ctx, query, team, author, replaced)
	if err != nil {
		return nil, fmt.Errorf("internal error")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)
	var userIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("internal error")
		}
		userIDs = append(userIDs, userID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("internal error")
	}
	if len(userIDs) == 0 {
		return nil, fmt.Errorf("no active replacement candidate in team")
	}

	return userIDs, nil
}

func (u UserRepositoryPostgres) FindUserTeamById(ctx context.Context, id string) (string, error) {
	query := `
        SELECT team_name 
        FROM users 
        WHERE user_id = $1
    `
	var teamName string
	err := u.db.QueryRowContext(ctx, query, id).Scan(&teamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("resource not found")
		}
		return "", fmt.Errorf("internal error")
	}

	return teamName, nil
}

func (u UserRepositoryPostgres) getPullRequestById(ctx context.Context, id string) ([]*delivery.PullRequestShort, error) {
	query := `
        SELECT 
            pr.pull_request_id,
            pr.pull_request_name,
            pr.author_id,
            pr.status
        FROM pull_requests pr
        INNER JOIN pull_request_reviewers r ON pr.pull_request_id = r.pull_request_id
        WHERE r.reviewer_id = $1
        ORDER BY pr.created_at DESC
    `
	rows, err := u.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("internal error")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)
	var pullRequests []*delivery.PullRequestShort
	for rows.Next() {
		var prShort delivery.PullRequestShort
		var statusStr string
		err := rows.Scan(&prShort.Id, &prShort.Name, &prShort.Author, &statusStr)
		if err != nil {
			return nil, fmt.Errorf("internal error")
		}
		prShort.Status = pr.Status(statusStr)
		pullRequests = append(pullRequests, &prShort)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("internal error")
	}

	return pullRequests, nil
}
