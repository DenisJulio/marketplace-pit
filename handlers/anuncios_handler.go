package handlers

import (
	"net/http"

	"github.com/DenisJulio/marketplace-pit/components"
	"github.com/DenisJulio/marketplace-pit/services"
	"github.com/labstack/echo/v4"
)

type AnunciosHandler struct {
	service services.AnuncioServices
}

func NewAnunciosHandler(service services.AnuncioServices) *AnunciosHandler {
	return &AnunciosHandler{service: service}
}

func (h *AnunciosHandler) HomeHandler(c echo.Context) error {
	anuncios := h.service.BuscaTodosAnuncios()
	return render(c, http.StatusOK, components.AnunciosPage(anuncios))
}
