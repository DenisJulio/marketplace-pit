package handlers

import (
	"net/http"
	"strconv"

	"github.com/DenisJulio/marketplace-pit/components"
	"github.com/DenisJulio/marketplace-pit/services"
	"github.com/DenisJulio/marketplace-pit/utils"
	"github.com/labstack/echo/v4"
)

type AnunciosHandler struct {
	anunSvc services.AnuncioServices
	usuSvc  services.UsuarioService
	oftSvc	services.OfertaService
	logger  utils.Logger
}

func NewAnunciosHandler(anunSvc services.AnuncioServices, usuSvc services.UsuarioService, oftSvc services.OfertaService, logger utils.Logger) *AnunciosHandler {
	return &AnunciosHandler{anunSvc: anunSvc, usuSvc: usuSvc, oftSvc: oftSvc, logger: logger}
}

func (h *AnunciosHandler) MostraTelaDeAnuncios(c echo.Context) error {
	anuncios := h.anunSvc.BuscaTodosAnuncios()
	return render(c, http.StatusOK, components.AnunciosPage(anuncios))
}

func (h *AnunciosHandler) MostraDetalhesDoAnuncio(c echo.Context) error {
	id := c.Param("id")
	h.logger.Debugf("Path id recebido: %s", id)
	convId, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Errorf("Erro ao converter id recebido: %s para inteiro. %v", id, err)
		// TODO: reder not founc page
	}
	anuncio, err := h.anunSvc.BuscaAnuncioPorID(convId)
	if err != nil {
		h.logger.Errorf("Anuncio com id=%d nao encontrado", convId)
		// TODO: render not found page
	}
	return render(c, http.StatusOK, components.DetalhesDoAnuncio(anuncio))
}

func (h *AnunciosHandler) MostraTelaDeNovaOferta(c echo.Context) error {
	id := c.Param("id")
	h.logger.Debugf("Path id recebido: %s", id)
	convId, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Errorf("Erro ao converter id recebido: %s para inteiro. %v", id, err)
	}
	anuncio, err := h.anunSvc.BuscaAnuncioPorID(convId)
	if err != nil {
		h.logger.Errorf("Anuncio com id=%d nao encontrado", convId)
	}
	anunciante, err := h.usuSvc.BuscaUsuarioPorId(anuncio.AnuncianteId)
	if err != nil {
		h.logger.Errorf("Anunciante com id=%d nao encontrado", anuncio.AnuncianteId)
		// TODO: Tratar esse erro melhor
	}
	return render(c, http.StatusOK, components.NovaOferta(anuncio, anunciante))
}

func (h *AnunciosHandler) CriaOferta(c echo.Context) error {
	// parse the form request
	c.Request().ParseForm()
	msg := c.FormValue("mensagem")
	h.logger.Debugf("POST request recebido com mensagem=%s", msg)

	// validate form
	if msg == "" {
		h.logger.Errorf("Mensagem nao pode ser vazia")
	}

	// create the oferta
	h.oftSvc.CriaNovaOfertaParaAnuncio(msg)

	// redirect to Home
	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}
