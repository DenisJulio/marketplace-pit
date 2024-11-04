package routes

import (
	"github.com/DenisJulio/marketplace-pit/handlers"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.Static("resources", "public/static")
	e.GET("/", handlers.HomeHandler)
}
