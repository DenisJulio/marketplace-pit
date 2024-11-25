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

	basePathImagens := "public/static/images"
	basePathResourceImagens := "/resources/images"

	db := db.NewDB(logger)
	usuarioStore := store.NewSQLUsuarioStore(db, logger)
	anuncioStore := &store.SQLAnuncioStore{DB: db, Logger: logger}
	ofertaStore := store.NewSQLOfertaStore(db, logger)
	sessaoStore := store.NovaSessaoStore(db, logger)
	fsStore := store.NewFileSystemImageStore(basePathImagens, basePathResourceImagens, logger)
	usuarioService := services.NewUsuarioService(usuarioStore, logger)
	anuncioService := services.NewAnuncioService(anuncioStore)
	ofertaService := services.NewOfertaService(*&ofertaStore)
	sessaoService := services.NovaSessaoService(*sessaoStore, logger)
	imagemService := services.NovoImagemService(fsStore, logger)
	anuncioHandler := handlers.NewAnunciosHandler(*anuncioService, *usuarioService, *ofertaService, logger)
	authHandler := handlers.NovoAuthHandler(logger)
	usuarioHandler := handlers.NovoUsuarioHandler(*usuarioService, *sessaoService, *imagemService, logger)
	mid := handlers.NovoMiddleware(*sessaoService, logger)
	router := routes.NewRouter(app, *anuncioHandler, *usuarioHandler, *authHandler, *mid)
	router.RegisterRoutes()
	logger.Fatal(app.Start(":3000"))
}
