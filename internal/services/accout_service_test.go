package services

import (
	"AccountManagementSystem/internal/repository"
	"AccountManagementSystem/log_color"
	"AccountManagementSystem/test"
	"context"
	"log/slog"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	db := test.GetDB()
	repo := repository.NewAccountRepo(db)
	service := NewAccountService(repo)
	create, err := service.Create(context.Background(), "superman", 2000)
	if err != nil {
		t.Fatal("could not create account", err)
	}
	if create.Username != "superman" || create.Balance != 2000 {
		t.Fatalf("improper create account")
	}
	slog.Info(log_color.BrightGreen("Case1 passed"))
	if acc, getErr := service.Get(context.Background(), create.ID); getErr != nil {
		t.Fatal("failed to get created account", getErr)
	} else if acc.Username != create.Username || acc.Balance != create.Balance || acc.CreatedAt != create.CreatedAt || acc.ID != create.ID {
		t.Fatal("failed to get created account")
	}
	slog.Info(log_color.BrightGreen("Case1 passed"))
}
