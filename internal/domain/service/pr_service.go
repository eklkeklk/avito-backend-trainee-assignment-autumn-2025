package service

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models/pr"
	"avito-backend-trainee-assignment-autumn-2025/internal/infrastructure/persistance/database"
	"context"
	"fmt"
	"slices"
	"time"
)

type PullRequestServicer interface {
	CreatePullRequest(ctx context.Context, id string, name string, authorId string) (*pr.PullRequest, error)
	MergePullRequest(ctx context.Context, id string) (*pr.PullRequest, error)
	ReassignPullRequest(ctx context.Context, id string, reviewerId string) (*pr.PullRequest, string, error)
}

type PullRequestService struct {
	repository     database.PullRequestRepository
	userRepository database.UserRepository
}

func NewPullRequestService(repository database.PullRequestRepository, userPrRepository database.UserRepository) *PullRequestService {
	return &PullRequestService{
		repository:     repository,
		userRepository: userPrRepository,
	}
}

func (p *PullRequestService) CreatePullRequest(ctx context.Context, id string, name string, authorId string) (*pr.PullRequest, error) {
	if id == "" {
		return nil, fmt.Errorf("invalid request")
	}
	if name == "" {
		return nil, fmt.Errorf("invalid request")
	}
	if authorId == "" {
		return nil, fmt.Errorf("invalid request")
	}
	pullRequest := &pr.PullRequest{
		Id:                id,
		Name:              name,
		Author:            authorId,
		Status:            pr.Open,
		Reviewers:         []string{},
		NeedMoreReviewers: false,
		CreatedAt:         time.Now().UTC(),
	}
	team, err := p.userRepository.FindUserTeamById(ctx, authorId)
	if err != nil {
		return nil, err
	}
	availableReviewers, err := p.userRepository.FindReviewers(ctx, team, authorId)
	if err != nil {
		return nil, err
	}
	if len(availableReviewers) >= 2 {
		pullRequest.Reviewers = append(pullRequest.Reviewers, availableReviewers[0])
		pullRequest.Reviewers = append(pullRequest.Reviewers, availableReviewers[1])
	} else {
		pullRequest.Reviewers = availableReviewers
		pullRequest.NeedMoreReviewers = true
	}
	createdPr, err := p.repository.Create(ctx, pullRequest)
	if err != nil {
		return nil, err
	}
	return createdPr, nil
}

func (p *PullRequestService) MergePullRequest(ctx context.Context, id string) (*pr.PullRequest, error) {
	if id == "" {
		return nil, fmt.Errorf("invalid request")
	}
	pullRequest := &pr.PullRequest{
		Id:       id,
		MergedAt: time.Now().UTC(),
	}
	updatedPr, err := p.repository.Merge(ctx, pullRequest)
	if err != nil {
		return nil, err
	}
	return updatedPr, nil
}

func (p *PullRequestService) ReassignPullRequest(ctx context.Context, id string, reviewerId string) (*pr.PullRequest, string, error) {
	if id == "" {
		return nil, "", fmt.Errorf("invalid request")
	}
	if reviewerId == "" {
		return nil, "", fmt.Errorf("invalid request")
	}
	pullRequest := &pr.PullRequest{
		Id: id,
	}
	isOpen, err := p.repository.IsOpen(ctx, id)
	if err != nil {
		return nil, "", err
	}
	if !isOpen {
		return nil, "", fmt.Errorf("cannot reassign on merged PR")
	}
	currentReviewers, err := p.repository.FindReviewers(ctx, id)
	if err != nil {
		return nil, "", err
	}
	if !slices.Contains(currentReviewers, reviewerId) {
		return nil, "", fmt.Errorf("reviewer is not assigned to this PR")
	}
	team, err := p.userRepository.FindUserTeamById(ctx, reviewerId)
	if err != nil {
		return nil, "", err
	}
	authorId, err := p.repository.FindAuthor(ctx, pullRequest.Id)
	if err != nil {
		return nil, "", err
	}
	reviewers, err := p.userRepository.FindNewReviewers(ctx, team, authorId, reviewerId)
	if err != nil {
		return nil, "", err
	}
	var newReviewerId string
	for _, reviewer := range reviewers {
		if !slices.Contains(currentReviewers, reviewer) {
			newReviewerId = reviewer
			break
		}
	}
	if newReviewerId == "" {
		return nil, "", fmt.Errorf("no active replacement candidate in team")
	}
	updatedPr, reviewerId, err := p.repository.Reassign(ctx, pullRequest, reviewerId, newReviewerId)
	if err != nil {
		return nil, "", err
	}

	return updatedPr, reviewerId, nil
}
