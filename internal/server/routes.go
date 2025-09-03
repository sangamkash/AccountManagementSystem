package server

import (
	"AccountManagementSystem/api_routers"
	"AccountManagementSystem/env_helper"
	"AccountManagementSystem/log_color"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/gofiber/swagger"
	"log/slog"
	"time"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	s.App.Get("/swagger/*", swagger.New(swagger.Config{
		Title: "AccountManagementSystem",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	}))

	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)
	slog.Info(log_color.Orange("==Check Swagger UI for API documentation=="))
	port := env_helper.ReadString("PORT", "8080")
	slog.Info(log_color.Green("http://127.0.0.1:" + port + "/swagger/index.html"))
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) RegisterAPIRoutes(apiService api_routers.IAPIService) {
	apiRoute := apiService.GetFiberRoutes()

	for _, route := range *apiRoute {
		handler := timeout.NewWithContext(route.Handler, route.TimeOutMs*time.Millisecond)
		switch route.Method {
		case api_routers.GET:
			s.App.Get(route.Path, handler)
			slog.Info(log_color.BrightGreen(route.Method.String()) + ":" + log_color.Pink(route.Path))
			break
		case api_routers.POST:
			s.App.Post(route.Path, handler)
			slog.Info(log_color.BrightCyan(route.Method.String()) + ":" + log_color.Pink(route.Path))
			break
		case api_routers.PUT:
			s.App.Put(route.Path, handler)
			slog.Info(log_color.BrightBlue(route.Method.String()) + ":" + log_color.Pink(route.Path))
			break
		case api_routers.PATCH:
			s.App.Patch(route.Path, handler)
			slog.Info(log_color.BrightMagenta(route.Method.String()) + ":" + log_color.Pink(route.Path))
			break
		case api_routers.DELETE:
			s.App.Delete(route.Path, handler)
			slog.Info(log_color.BrightRed(route.Method.String()) + ":" + log_color.Pink(route.Path))
			break
		case api_routers.HEAD:
			s.App.Head(route.Path, handler)
			slog.Info(log_color.Orange(route.Method.String()) + ":" + log_color.Pink(route.Path))
			break
		case api_routers.OPTIONS:
			s.App.Options(route.Path, handler)
			slog.Info(log_color.BrightYellow(route.Method.String()) + ":" + log_color.Pink(route.Path))
			break
		case api_routers.TRACE:
			s.App.Trace(route.Path, handler)
			slog.Info(log_color.Peach(route.Method.String()) + ":" + log_color.Pink(route.Path))
			break
		case api_routers.CONNECT:
			s.App.Connect(route.Path, handler)
			slog.Info(log_color.Violet(route.Method.String()) + ":" + log_color.Pink(route.Path))
			break
		}
	}
}

// healthHandler godoc
// @Summary      Health check
// @Description  Returns the health status of the database/server
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /health [get]
func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
