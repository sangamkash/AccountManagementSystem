package services

import (
	"AccountManagementSystem/internal/models"
	"AccountManagementSystem/internal/repository"
	"context"
	"database/sql"
	"time"
)

type AccountService struct{ repo *repository.AccountRepo }

func NewAccountService(db *sql.DB) *AccountService {
	return &AccountService{repo: repository.NewAccountRepo(db)}
}

func (s *AccountService) Create(ctx context.Context, username string, initial float64) (models.Account, error) {
	timedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.repo.Create(timedCtx, username, initial)
}

func (s *AccountService) Get(ctx context.Context, id int64) (models.Account, error) {
	timedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.repo.GetByID(timedCtx, id)
}
