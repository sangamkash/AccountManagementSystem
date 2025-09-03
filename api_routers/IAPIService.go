package api_routers

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

type HTTPMethod int

// Define enum-like constants
const (
	GET HTTPMethod = iota
	POST
	PUT
	PATCH
	DELETE
	HEAD
	OPTIONS
	TRACE
	CONNECT
)

// String() method to convert enum to string
func (m HTTPMethod) String() string {
	switch m {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case PATCH:
		return "PATCH"
	case DELETE:
		return "DELETE"
	case HEAD:
		return "HEAD"
	case OPTIONS:
		return "OPTIONS"
	default:
		return "UNKNOWN"
	}
}

type APIRoute struct {
	Path      string
	Method    HTTPMethod
	Handler   fiber.Handler
	TimeOutMs time.Duration
}
type IAPIService interface {
	GetFiberRoutes() *[]APIRoute
}

/*
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
			s.App.Delete(route.Path)
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
*/
