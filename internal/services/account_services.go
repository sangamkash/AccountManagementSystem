package services

import (
	"AccountManagementSystem/internal/models"
	"AccountManagementSystem/internal/repository"
	"context"
	"time"
)

type AccountService struct{ repo *repository.AccountRepo }

func NewAccountService(r *repository.AccountRepo) *AccountService { return &AccountService{repo: r} }

func (s *AccountService) Create(username string, initial float64) (models.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.Create(ctx, username, initial)
}

func (s *AccountService) Get(id int64) (models.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.GetByID(ctx, id)
}
