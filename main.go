package main

import (
	"github.com/DenisJulio/marketplace-pit/db"
	"github.com/DenisJulio/marketplace-pit/handlers"
	"github.com/DenisJulio/marketplace-pit/routes"
	"github.com/DenisJulio/marketplace-pit/services"
	"github.com/DenisJulio/marketplace-pit/store"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	app := echo.New()
	app.HideBanner = true
	logger := app.Logger

	db := db.NewDB(logger)
	anuncioStore := &store.SQLAnuncioStore{DB: db, Logger: logger}
	anuncioService := services.NewAnuncioService(anuncioStore)
	anuncioHandler := handlers.NewAnunciosHandler(*anuncioService)
	router := routes.NewRouter(app, *anuncioHandler)
	router.RegisterRoutes()
	logger.Fatal(app.Start(":3000"))
}
