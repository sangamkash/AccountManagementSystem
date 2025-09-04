package handlers

import (
	"AccountManagementSystem/api_routers"
	"AccountManagementSystem/internal/requestData"
	"AccountManagementSystem/internal/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: transactionService}
}

func (t TransactionHandler) GetFiberRoutes() *[]api_routers.APIRoute {

	return &[]api_routers.APIRoute{
		{"/v1/transactions/deposit", api_routers.POST, t.Deposit, 3000},
		{"/v1/transactions/withdraw", api_routers.POST, t.Withdraw, 3000},
		{"/v1/transactions/:accountId", api_routers.GET, t.ListTransactions, 3000},
	}
}

// Deposit godoc
// @Summary Deposit into account
// @Description Enqueues a deposit transaction for processing
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body requestData.TransactionReq true "Deposit request"
// @Success 202 {object} map[string]interface{} "accepted"
// @Failure 400 {object} map[string]string "bad request"
// @Failure 500 {object} map[string]string "server error"
// @Router /v1/transactions/deposit [post]
func (t TransactionHandler) Deposit(ctx *fiber.Ctx) error {
	var req requestData.TransactionReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return t.service.EnqueueDeposit(req.AccountId, req.Amount, req.IdempotencyKey)
}

// Withdraw godoc
// @Summary Withdraw from account
// @Description Enqueues a withdrawal transaction for processing
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body requestData.TransactionReq true "Withdraw request"
// @Success 202 {object} map[string]interface{} "accepted"
// @Failure 400 {object} map[string]string "bad request"
// @Failure 500 {object} map[string]string "server error"
// @Router /v1/transactions/withdraw [post]
func (t TransactionHandler) Withdraw(ctx *fiber.Ctx) error {
	var req requestData.TransactionReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return t.service.EnqueueDeposit(req.AccountId, req.Amount, req.IdempotencyKey)
}

// ListTransactions godoc
// @Summary List transactions of an account
// @Description Returns recent transactions for an account
// @Tags transactions
// @Accept json
// @Produce json
// @Param accountId path int true "Account ID"
// @Param limit query int false "Number of transactions (default 10, max 20)"
// @Success 200 {object} map[string]interface{} "list of transactions"
// @Failure 400 {object} map[string]string "bad request"
// @Failure 500 {object} map[string]string "server error"
// @Router /v1/transactions/{accountId} [get]
func (t TransactionHandler) ListTransactions(ctx *fiber.Ctx) error {
	accountIdStr := ctx.Params("accountId")
	if accountIdStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid account id"})
	}
	accountId, err := strconv.ParseInt(accountIdStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	limit := 10
	limitStr := ctx.Params("limit")
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
		if limit <= 0 || limit > 20 {
			limit = 10
		}
	}
	list, err := t.service.List(accountId, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    list,
	})
}
