package repository

import (
	"AccountManagementSystem/log_color"
	"AccountManagementSystem/test"
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	db := test.GetDB()
	test.ClearAccount(db)
	accountRepo := NewAccountRepo(db)
	//case 1
	create, err := accountRepo.Create(context.Background(), "abc", 100)
	if err != nil {
		t.Fatal("could not create account", err)
	}
	if create.Username != "abc" || create.Balance != 100 {
		t.Fatalf("improper create account")
	}
	slog.Info(log_color.BrightGreen("Case1 passed"))

	//case 2
	if acc, getErr := accountRepo.GetByID(context.Background(), create.ID); getErr != nil {
		t.Fatal("failed to get created account", getErr)
	} else if acc.Username != create.Username || acc.Balance != create.Balance || acc.CreatedAt != create.CreatedAt || acc.ID != create.ID {
		t.Fatal("failed to get created account")
	}
	slog.Info(log_color.BrightGreen("Case1 passed"))
}
