package routes

import (
	"github.com/DenisJulio/marketplace-pit/handlers"
	"github.com/labstack/echo/v4"
)

type Router struct {
	echo *echo.Echo
	ancH handlers.AnunciosHandler
}

func NewRouter(e *echo.Echo, ancH handlers.AnunciosHandler) *Router {
	return &Router{echo: e, ancH: ancH}
}

func (r *Router) RegisterRoutes() {
	r.echo.Static("resources", "public/static")
	r.echo.GET("/", r.ancH.HomeHandler)
}
