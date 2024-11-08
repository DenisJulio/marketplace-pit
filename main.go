package main

import (
	"github.com/DenisJulio/marketplace-pit/db"
	"github.com/DenisJulio/marketplace-pit/handlers"
	"github.com/DenisJulio/marketplace-pit/routes"
	"github.com/DenisJulio/marketplace-pit/services"
	"github.com/DenisJulio/marketplace-pit/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
)

func main() {
	app := echo.New()
	app.HideBanner = true
	logger := app.Logger
	logger.SetLevel(log.DEBUG)

	db := db.NewDB(logger)
	usuarioStore := store.NewSQLUsuarioStore(db, logger)
	anuncioStore := &store.SQLAnuncioStore{DB: db, Logger: logger}
	usuarioService := services.NewUsuarioService(usuarioStore)
	anuncioService := services.NewAnuncioService(anuncioStore)
	anuncioHandler := handlers.NewAnunciosHandler(*anuncioService, *usuarioService, logger)
	router := routes.NewRouter(app, *anuncioHandler)
	router.RegisterRoutes()
	logger.Fatal(app.Start(":3000"))
}
