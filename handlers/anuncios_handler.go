package handlers

import (
	"net/http"
	"strconv"

	"github.com/DenisJulio/marketplace-pit/components"
	"github.com/DenisJulio/marketplace-pit/model"
	"github.com/DenisJulio/marketplace-pit/services"
	"github.com/DenisJulio/marketplace-pit/utils"
	"github.com/labstack/echo/v4"
)

type AnunciosHandler struct {
	anunSvc services.AnuncioServices
	usuSvc  services.UsuarioService
	oftSvc  services.OfertaService
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
	id, err := h.extrairId(c.Param("id"))
	if err != nil {
		// TODO: render not found page
	}
	anuncio, err := h.anunSvc.BuscaAnuncioPorID(id)
	if err != nil {
		h.logger.Errorf("Anuncio com id=%d nao encontrado", id)
		// TODO: render not found page
	}
	return render(c, http.StatusOK, components.DetalhesDoAnuncio(anuncio))
}

func (h *AnunciosHandler) MostraTelaDeNovaOferta(c echo.Context) error {
	id, err := h.extrairId(c.Param("id"))
	if err != nil {
		// TODO: render not found page
	}
	anuncio, err := h.anunSvc.BuscaAnuncioPorID(id)
	if err != nil {
		h.logger.Errorf("Anuncio com id=%d nao encontrado", id)
		// TODO: render not found page
	}
	anunciante, err := h.usuSvc.BuscaUsuarioPorId(anuncio.AnuncianteId)
	if err != nil {
		h.logger.Errorf("Anunciante com id=%d nao encontrado", anuncio.AnuncianteId)
	}
	return render(c, http.StatusOK, components.NovaOferta(anuncio, anunciante))
}

func (h *AnunciosHandler) CriaOferta(c echo.Context) error {
	// TODO: extrair o id do anunciante da sessao
	var ofertanteId int

	// extrair e validar o formulario
	msgForm, err := h.validarFormulario(c, "mensagem")
	if err != nil {
		// TODO: retornar um snackbar com htmx
	}

	id, _ := h.extrairId(c.Param("id"))
	anuncio, _ := h.anunSvc.BuscaAnuncioPorID(id)
	anunciante, _ := h.usuSvc.BuscaUsuarioPorId(anuncio.AnuncianteId)

	// create the oferta
	oferta := model.NewOferta(anunciante.ID, ofertanteId, anuncio.ID)
	mensagen := model.NewMensagem(ofertanteId, msgForm)
	h.oftSvc.CriaNovaOfertaParaAnuncio(*oferta, *mensagen)

	// redirect to Home
	return c.Redirect(http.StatusSeeOther, "/") // TODO: idealmete, redirecionar pelo cliente
}

func (h *AnunciosHandler) extrairId(id string) (int, error) {
	convId, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Errorf("Erro ao converter id recebido: %s para inteiro. %v", id, err)
		return 0, err
	}
	return convId, nil
}

func (h *AnunciosHandler) validarFormulario(c echo.Context, campoForm string) (string, error) {
	val := c.FormValue(campoForm)
	if val == "" {
		h.logger.Errorf("Campo do formalario: %s nao pode ser vazio", campoForm)
		return "", echo.ErrBadRequest
	}
	return val, nil
}
