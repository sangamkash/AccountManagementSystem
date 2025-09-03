package handlers

import (
	"AccountManagementSystem/api_routers"
	"github.com/gofiber/fiber/v2"
)

type AccountHandler struct {
}

func (a AccountHandler) GetFiberRoutes() *[]api_routers.APIRoute {
	return &[]api_routers.APIRoute{
		{"/v1/accounts", api_routers.POST, a.CreateAccount, 3000},
		{"/v1/accounts/id:", api_routers.POST, a.GetAccount, 3000},
	}
}

func (a AccountHandler) CreateAccount(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}

func (a AccountHandler) GetAccount(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}
