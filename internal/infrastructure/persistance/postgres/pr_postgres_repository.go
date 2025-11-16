package postgres

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models/pr"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type PullRequestRepositoryPostgres struct {
	db *sql.DB
}

func NewPullRequestRepositoryPostgres(db *sql.DB) *PullRequestRepositoryPostgres {
	return &PullRequestRepositoryPostgres{
		db: db,
	}
}

func (p PullRequestRepositoryPostgres) Create(ctx context.Context, request *pr.PullRequest) (*pr.PullRequest, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("internal error")
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			return
		}
	}(tx)
	checkQuery := `SELECT 1 FROM pull_requests WHERE pull_request_id = $1`
	var exists int
	err = tx.QueryRowContext(ctx, checkQuery, request.Id).Scan(&exists)
	if err == nil {
		return nil, fmt.Errorf("PR id already exists")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("internal error")
	}
	prQuery := `
        INSERT INTO pull_requests 
            (pull_request_id, pull_request_name, author_id, status, need_more_reviewers, created_at, merged_at)
        VALUES 
            ($1, $2, $3, $4, $5, $6, $7)
        RETURNING 
            pull_request_id, pull_request_name, author_id, status, need_more_reviewers, created_at, merged_at
    `
	var createdPR pr.PullRequest
	err = tx.QueryRowContext(
		ctx,
		prQuery,
		request.Id,
		request.Name,
		request.Author,
		request.Status,
		request.NeedMoreReviewers,
		request.CreatedAt,
		request.MergedAt,
	).Scan(
		&createdPR.Id,
		&createdPR.Name,
		&createdPR.Author,
		&createdPR.Status,
		&createdPR.NeedMoreReviewers,
		&createdPR.CreatedAt,
		&createdPR.MergedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("internal error")
	}
	err = p.setReviewers(ctx, tx, request.Id, request.Reviewers)
	if err != nil {
		return nil, err
	}
	createdPR.Reviewers = request.Reviewers

	return &createdPR, nil
}

func (p PullRequestRepositoryPostgres) Merge(ctx context.Context, request *pr.PullRequest) (*pr.PullRequest, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("internal error")
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			return
		}
	}(tx)
	checkQuery := `
        SELECT status, merged_at 
        FROM pull_requests 
        WHERE pull_request_id = $1
    `
	var currentStatus string
	var currentMergedAt sql.NullTime
	err = tx.QueryRowContext(ctx, checkQuery, request.Id).Scan(&currentStatus, &currentMergedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("resource not found")
		}
		return nil, fmt.Errorf("internal error")
	}
	if currentStatus == "MERGED" {
		request.MergedAt = currentMergedAt.Time
	}
	query := `
        UPDATE pull_requests 
        SET 
            status = 'MERGED',
            merged_at = $1
        WHERE 
            pull_request_id = $2
        RETURNING 
            pull_request_id, pull_request_name, author_id, status, need_more_reviewers, created_at, merged_at
    `
	var mergedPR pr.PullRequest
	err = tx.QueryRowContext(ctx, query, request.MergedAt, request.Id).Scan(
		&mergedPR.Id,
		&mergedPR.Name,
		&mergedPR.Author,
		&mergedPR.Status,
		&mergedPR.NeedMoreReviewers,
		&mergedPR.CreatedAt,
		&mergedPR.MergedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("resource not found")
		}
		return nil, fmt.Errorf("internal error")
	}
	reviewersIds, err := p.getReviewers(ctx, tx, request.Id)
	if err != nil {
		return nil, err
	}
	mergedPR.Reviewers = reviewersIds
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("internal error")
	}

	return &mergedPR, nil
}

func (p PullRequestRepositoryPostgres) FindAuthor(ctx context.Context, id string) (string, error) {
	query := `
        SELECT author_id 
        FROM pull_requests 
        WHERE pull_request_id = $1
    `
	var authorID string
	err := p.db.QueryRowContext(ctx, query, id).Scan(&authorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("resource not found")
		}
		return "", fmt.Errorf("internal error")
	}

	return authorID, nil
}

func (p PullRequestRepositoryPostgres) FindReviewers(ctx context.Context, id string) ([]string, error) {
	query := `
        SELECT reviewer_id 
        FROM pull_request_reviewers 
        WHERE pull_request_id = $1
    `
	rows, err := p.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("internal error")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var reviewerIds []string
	for rows.Next() {
		var reviewerId string
		if err := rows.Scan(&reviewerId); err != nil {
			return nil, fmt.Errorf("internal error")
		}
		reviewerIds = append(reviewerIds, reviewerId)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("internal error")
	}

	return reviewerIds, nil
}

func (p PullRequestRepositoryPostgres) Reassign(ctx context.Context, request *pr.PullRequest, reviewerId string, newReviewerId string) (*pr.PullRequest, string, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, "", fmt.Errorf("internal error")
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			return
		}
	}(tx)
	updateQuery := `
        UPDATE pull_request_reviewers 
        SET reviewer_id = $3
        WHERE pull_request_id = $1 AND reviewer_id = $2
        RETURNING reviewer_id
    `
	var oldReviewerId string
	err = tx.QueryRowContext(ctx, updateQuery, request.Id, reviewerId, newReviewerId).Scan(&oldReviewerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", fmt.Errorf("resource not found")
		}
		return nil, "", fmt.Errorf("internal error")
	}
	reviewerIds, err := p.getReviewers(ctx, tx, request.Id)
	if err != nil {
		return nil, "", err
	}
	updatedPr, err := p.getPullRequestWithoutReviewers(ctx, request.Id)
	if err != nil {
		return nil, "", fmt.Errorf("internal error")
	}
	updatedPr.Reviewers = reviewerIds
	if err := tx.Commit(); err != nil {
		return nil, "", fmt.Errorf("internal error")
	}

	return updatedPr, oldReviewerId, nil
}

func (p PullRequestRepositoryPostgres) IsOpen(ctx context.Context, id string) (bool, error) {
	query := `
        SELECT status 
        FROM pull_requests 
        WHERE pull_request_id = $1
    `
	var status string
	err := p.db.QueryRowContext(ctx, query, id).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("resource not found")
		}
		return false, fmt.Errorf("internal error")
	}

	return status == "OPEN", nil
}

func (p PullRequestRepositoryPostgres) getPullRequestWithoutReviewers(ctx context.Context, id string) (*pr.PullRequest, error) {
	prQuery := `
        SELECT 
            pull_request_id, 
            pull_request_name, 
            author_id, 
            status, 
            need_more_reviewers, 
            created_at, 
            merged_at
        FROM pull_requests 
        WHERE pull_request_id = $1
    `
	var pullRequest pr.PullRequest
	err := p.db.QueryRowContext(ctx, prQuery, id).Scan(
		&pullRequest.Id,
		&pullRequest.Name,
		&pullRequest.Author,
		&pullRequest.Status,
		&pullRequest.NeedMoreReviewers,
		&pullRequest.CreatedAt,
		&pullRequest.MergedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("resource not found")
		}
		return nil, fmt.Errorf("internal error")
	}

	return &pullRequest, nil
}

func (p PullRequestRepositoryPostgres) setReviewers(ctx context.Context, tx *sql.Tx, id string, reviewers []string) error {
	if len(reviewers) > 0 {
		reviewersQuery := `
            INSERT INTO pull_request_reviewers 
                (pull_request_id, reviewer_id)
            VALUES 
                ($1, $2)
        `

		stmt, err := tx.PrepareContext(ctx, reviewersQuery)
		if err != nil {
			return fmt.Errorf("internal error")
		}
		defer func(stmt *sql.Stmt) {
			err := stmt.Close()
			if err != nil {
				return
			}
		}(stmt)

		for _, reviewerID := range reviewers {
			_, err := stmt.ExecContext(ctx, id, reviewerID)
			if err != nil {
				return fmt.Errorf("internal error")
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("internal error")
	}

	return nil
}

func (p PullRequestRepositoryPostgres) getReviewers(ctx context.Context, tx *sql.Tx, id string) ([]string, error) {
	reviewersQuery := `
        SELECT reviewer_id 
        FROM pull_request_reviewers 
        WHERE pull_request_id = $1
    `
	rows, err := tx.QueryContext(ctx, reviewersQuery, id)
	if err != nil {
		return nil, fmt.Errorf("internal error")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)
	var reviewerIds []string
	for rows.Next() {
		var reviewerId string
		err = rows.Scan(&reviewerId)
		if err != nil {
			return nil, fmt.Errorf("internal error")
		}
		reviewerIds = append(reviewerIds, reviewerId)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("internal error")
	}

	return reviewerIds, nil
}
