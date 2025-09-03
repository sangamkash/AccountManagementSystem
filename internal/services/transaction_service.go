package services

import (
	"AccountManagementSystem/internal/models"
	"AccountManagementSystem/internal/queue"
	"AccountManagementSystem/internal/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrInsufficient = errors.New("insufficient funds")

type TransactionService struct {
	accRepo *repository.AccountRepo
	txRepo  *repository.TransactionRepo
	rb      *queue.Rabbit
}

func NewTransactionService(a *repository.AccountRepo, t *repository.TransactionRepo, r *queue.Rabbit) *TransactionService {
	return &TransactionService{accRepo: a, txRepo: t, rb: r}
}

func (s *TransactionService) EnqueueDeposit(accountID int64, amount float64, key string) error {
	msg := queue.TransactionMessage{AccountID: accountID, Type: "deposit", Amount: amount, IdempotencyKey: key}
	return s.rb.Publish(context.Background(), msg)
}

func (s *TransactionService) EnqueueWithdraw(accountID int64, amount float64, key string) error {
	msg := queue.TransactionMessage{AccountID: accountID, Type: "withdraw", Amount: amount, IdempotencyKey: key}
	return s.rb.Publish(context.Background(), msg)
}

func (s *TransactionService) List(accountID int64, limit int) ([]models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.txRepo.ListByAccount(ctx, accountID, limit)
}

// ProcessMessage is used by the processor to apply the transaction in a DB transaction
func (s *TransactionService) ProcessMessage(ctx context.Context, msg queue.TransactionMessage) error {
	tx, err := s.accRepo.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	// ensure commit/rollback
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	acc, err := s.accRepo.GetForUpdate(ctx, tx, msg.AccountID)
	if err != nil {
		return err
	}

	switch msg.Type {
	case "deposit":
		acc.Balance += msg.Amount
	case "withdraw":
		if acc.Balance < msg.Amount {
			return ErrInsufficient
		}
		acc.Balance -= msg.Amount
	default:
		return errors.New("unknown transaction type")
	}

	if err := s.accRepo.UpdateBalance(ctx, tx, acc.ID, acc.Balance); err != nil {
		return err
	}

	// insert transaction record
	_, err = s.txRepo.Insert(ctx, tx, models.Transaction{
		AccountID: acc.ID,
		Type:      msg.Type,
		Amount:    msg.Amount,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	committed = true
	return nil
}
