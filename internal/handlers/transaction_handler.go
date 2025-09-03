package handlers

import (
	"AccountManagementSystem/api_routers"
	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
}

func (t TransactionHandler) GetFiberRoutes() *[]api_routers.APIRoute {

	return &[]api_routers.APIRoute{
		{"/v1/transactions/deposit", api_routers.POST, t.Deposit, 3000},
		{"/v1/transactions/withdraw", api_routers.POST, t.Withdraw, 3000},
		{"/v1/transactions/:accountId", api_routers.POST, t.ListTransactions, 3000},
	}
}

func (t TransactionHandler) Deposit(ctx *fiber.Ctx) error {

	return nil
}

func (t TransactionHandler) Withdraw(ctx *fiber.Ctx) error {
	return nil
}

func (t TransactionHandler) ListTransactions(ctx *fiber.Ctx) error {
	return nil
}
