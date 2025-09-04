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
	mq      *queue.KafkaQueue
}

func NewTransactionService(a *repository.AccountRepo, t *repository.TransactionRepo, kq *queue.KafkaQueue) *TransactionService {
	return &TransactionService{accRepo: a, txRepo: t, mq: kq}
}

func (s *TransactionService) EnqueueDeposit(accountID int64, amount float64, key string) error {
	msg := models.TransactionMessage{
		AccountID:      accountID,
		Type:           "deposit",
		Amount:         amount,
		IdempotencyKey: key,
	}
	return s.mq.PublishMessage(msg, key) // Kafka publish
}

func (s *TransactionService) EnqueueWithdraw(accountID int64, amount float64, key string) error {
	msg := models.TransactionMessage{
		AccountID:      accountID,
		Type:           "withdraw",
		Amount:         amount,
		IdempotencyKey: key,
	}
	return s.mq.PublishMessage(msg, key) // Kafka publish
}

func (s *TransactionService) List(accountID int64, limit int) ([]models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.txRepo.ListByAccount(ctx, accountID, limit)
}

// ProcessMessage is used by the processor to apply the transaction in a db transaction
func (s *TransactionService) ProcessMessage(ctx context.Context, msg models.TransactionMessage) error {
	tx, beginTransErr := s.accRepo.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if beginTransErr != nil {
		return beginTransErr
	}
	// ensure commit/rollback
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	acc, getAccErr := s.accRepo.GetForUpdate(ctx, tx, msg.AccountID)
	if getAccErr != nil {
		return getAccErr
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

	if updatedBalErr := s.accRepo.UpdateBalance(ctx, tx, acc.ID, acc.Balance); updatedBalErr != nil {
		return updatedBalErr
	}

	// insert transaction record
	_, insertTranErr := s.txRepo.Insert(ctx, tx, models.Transaction{
		AccountID: acc.ID,
		Type:      msg.Type,
		Amount:    msg.Amount,
		CreatedAt: time.Now(),
	})
	if insertTranErr != nil {
		return insertTranErr
	}

	if commitError := tx.Commit(); commitError != nil {
		return commitError
	}
	committed = true
	return nil
}
