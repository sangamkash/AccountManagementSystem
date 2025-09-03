package server

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"

	"AccountManagementSystem/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "AccountManagementSystem",
			AppName:      "AccountManagementSystem",
		}),

		db: database.New(),
	}

	return server
}

func (s *FiberServer) DB() *sql.DB {
	return s.db.DB()
}
