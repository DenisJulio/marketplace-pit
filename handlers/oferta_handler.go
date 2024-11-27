package handlers

import (
	"net/http"

	"github.com/DenisJulio/marketplace-pit/components"
	"github.com/DenisJulio/marketplace-pit/model"
	"github.com/DenisJulio/marketplace-pit/services"
	"github.com/DenisJulio/marketplace-pit/utils"
	"github.com/labstack/echo/v4"
)

type OfertaHandler struct {
	oftSvc services.OfertaService
	logger utils.Logger
}

func NewOfertaHandler(oftSvc services.OfertaService, logger utils.Logger) *OfertaHandler {
	return &OfertaHandler{oftSvc: oftSvc, logger: logger}
}

func (oh *OfertaHandler) MostraListaDeOfertasDoUsuario(c echo.Context) error {
	var ofertas []model.Oferta
	return render(c, http.StatusOK, components.ListaDeOfertasDeUsuario(ofertas))
}
