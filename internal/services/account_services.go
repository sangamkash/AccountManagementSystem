package services

import (
	"AccountManagementSystem/internal/models"
	"AccountManagementSystem/internal/repository"
	"context"
	"errors"
	"time"
)

type AccountService struct{ repo *repository.AccountRepo }

func NewAccountService(repo *repository.AccountRepo) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) Create(ctx context.Context, username string, initial float64) (*models.Account, error) {
	if initial < 0 {
		return nil, errors.New("balance should be greater than zero")
	}
	timedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.repo.Create(timedCtx, username, initial)
}

func (s *AccountService) Get(ctx context.Context, id int64) (*models.Account, error) {
	timedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.repo.GetByID(timedCtx, id)
}
