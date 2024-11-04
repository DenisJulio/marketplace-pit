package main

import (
	"github.com/DenisJulio/marketplace-pit/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	app.HideBanner = true
	logger := app.Logger
	routes.RegisterRoutes(app)
	logger.Fatal(app.Start(":3000"))
}