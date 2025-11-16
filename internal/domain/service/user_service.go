package service

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/delivery"
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models/user"
	"avito-backend-trainee-assignment-autumn-2025/internal/infrastructure/persistance/database"
	"context"
	"fmt"
)

type UserServicer interface {
	SetActiveStatus(ctx context.Context, userId string, isActive bool) (*user.User, error)
	GetReviewList(ctx context.Context, userId string) (*delivery.UserPullRequest, error)
}

type UserService struct {
	repository database.UserRepository
}

func NewUserService(repository database.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (u *UserService) GetReviewList(ctx context.Context, userId string) (*delivery.UserPullRequest, error) {
	if userId == "" {
		return nil, fmt.Errorf("invalid request")
	}
	return u.repository.GetReviewsByID(ctx, userId)
}

func (u *UserService) SetActiveStatus(ctx context.Context, userId string, isActive bool) (*user.User, error) {
	if userId == "" {
		return nil, fmt.Errorf("invalid request")
	}
	updatedUser := &user.User{
		Id:       userId,
		Name:     "",
		TeamName: "",
		IsActive: isActive,
	}
	return u.repository.Update(ctx, updatedUser)
}
