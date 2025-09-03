package handlers

import (
	"AccountManagementSystem/api_routers"
	"AccountManagementSystem/internal/requestData"
	"AccountManagementSystem/internal/services"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type AccountHandler struct {
	accountService *services.AccountService
}

func NewAccountHandler(db *sql.DB) *AccountHandler {
	return &AccountHandler{accountService: services.NewAccountService(db)}
}

// CreateAccount godoc
// @Summary      Create a new account
// @Description  Creates a new account with username and initial amount
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        account  body      requestData.CreateAccountReq  true  "Account request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /v1/accounts [post]
func (a AccountHandler) GetFiberRoutes() *[]api_routers.APIRoute {
	return &[]api_routers.APIRoute{
		{"/v1/accounts", api_routers.POST, a.CreateAccount, 3000},
		{"/v1/accounts/:id", api_routers.GET, a.GetAccount, 3000},
	}
}

// GetAccount godoc
// @Summary      Get account by ID
// @Description  Retrieves account details using ID
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /v1/accounts/{id} [get]
func (a AccountHandler) CreateAccount(ctx *fiber.Ctx) error {
	var req requestData.CreateAccountReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "improper request body"})
	}
	account, err := a.accountService.Create(ctx.Context(), req.Username, req.InitialAmount)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(fiber.Map{"message": "Account created successfully", "data": account})
}

// GetAccount godoc
// @Summary      Get account by ID
// @Description  Retrieves account details using ID
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /v1/accounts/{id} [get]
func (a AccountHandler) GetAccount(ctx *fiber.Ctx) error {
	idstr := ctx.Params("id")
	if idstr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	data, err := a.accountService.Get(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(fiber.Map{"msg": "successfully retrieved account", "data": data})
}
