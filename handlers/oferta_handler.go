package handlers

import (
	"net/http"

	"github.com/DenisJulio/marketplace-pit/components"
	"github.com/DenisJulio/marketplace-pit/services"
	"github.com/DenisJulio/marketplace-pit/utils"
	"github.com/labstack/echo/v4"
)

type OfertaHandler struct {
	oftSvc services.OfertaService
	ssSvc  services.SessaoService
	usSvc  services.UsuarioService
	logger utils.Logger
}

func NewOfertaHandler(oftSvc services.OfertaService, ssSvc services.SessaoService, usSvc services.UsuarioService, logger utils.Logger) *OfertaHandler {
	return &OfertaHandler{oftSvc: oftSvc, ssSvc: ssSvc, usSvc: usSvc, logger: logger}
}

func (oh *OfertaHandler) MostraListaDeOfertasDoUsuario(c echo.Context) error {
	nomeDeUsuario, _ := oh.ssSvc.BuscaNomeDeUsuarioDaSessao(c)
	usuario, _ := oh.usSvc.BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario)
	ofertas, err := oh.oftSvc.BuscaTodasAsOfertasExpandidasDoUsuario(usuario.ID)
	if err != nil {
		oh.logger.Errorf("Erro ao buscar ofertas do usuario %s: %v", nomeDeUsuario, err)
		c.NoContent(http.StatusInternalServerError)
	}
	return render(c, http.StatusOK, components.ListaDeOfertasDeUsuario(ofertas))
}
