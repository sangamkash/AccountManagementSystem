package services

import (
	"AccountManagementSystem/internal/repository"
	"AccountManagementSystem/test"
	"context"
	"github.com/gofiber/fiber/v2/utils"
	"testing"
)

func TestTransactionDeposit(t *testing.T) {
	db := test.GetDB()
	test.ClearAccount(db)
	test.ClearTransaction(db)
	kafka := test.GetKafka()
	tansRepo := repository.NewTransactionRepo(db)
	accRepo := repository.NewAccountRepo(db)
	accService := NewAccountService(accRepo)
	transService := NewTransactionService(accRepo, tansRepo, kafka)
	test.KafkaConsumer(transService)
	acc1, createErr := accService.Create(context.Background(), "batman", 2000)
	if createErr != nil {
		t.Fatal("could not acc1 account", createErr)
	}
	if err := transService.EnqueueDeposit(acc1.ID, 9, utils.UUID()); err != nil {
		t.Fatal("could not enqueue deposit", err)
	}
	if getAcc, err := accService.Get(context.Background(), acc1.ID); err != nil {
		t.Fatal("could not get account to verify transaction", err)
	} else if getAcc.Balance != 2009 {
		t.Fatal("improper account balance")
	}
	if err := transService.EnqueueDeposit(acc1.ID, -9, utils.UUID()); err == nil {
		t.Fatal("negative values are not allowed", err)
	}
}

func TestTransactionWithdraw(t *testing.T) {
	db := test.GetDB()
	test.ClearAccount(db)
	test.ClearTransaction(db)
	kafka := test.GetKafka()
	tansRepo := repository.NewTransactionRepo(db)
	accRepo := repository.NewAccountRepo(db)
	accService := NewAccountService(accRepo)
	transService := NewTransactionService(accRepo, tansRepo, kafka)
	test.KafkaConsumer(transService)
	acc1, createErr := accService.Create(context.Background(), "batman", 2000)
	if createErr != nil {
		t.Fatal("could not acc1 account", createErr)
	}
	if err := transService.EnqueueWithdraw(acc1.ID, 9, utils.UUID()); err != nil {
		t.Fatal("could not enqueue deposit", err)
	}
	if getAcc, err := accService.Get(context.Background(), acc1.ID); err != nil {
		t.Fatal("could not get account to verify transaction", err)
	} else if getAcc.Balance != 1991 {
		t.Fatal("improper account balance")
	}
	if err := transService.EnqueueWithdraw(acc1.ID, -9, utils.UUID()); err == nil {
		t.Fatal("negative values are not allowed", err)
	}
}
