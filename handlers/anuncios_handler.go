package handlers

import (
	"mime/multipart"
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
	imgSvc  services.ImagemService
	logger  utils.Logger
}

func NewAnunciosHandler(anunSvc services.AnuncioServices, usuSvc services.UsuarioService, oftSvc services.OfertaService, logger utils.Logger) *AnunciosHandler {
	return &AnunciosHandler{anunSvc: anunSvc, usuSvc: usuSvc, oftSvc: oftSvc, logger: logger}
}

func (h *AnunciosHandler) MostraTelaDeAnuncios(c echo.Context) error {
	anuncios := h.anunSvc.BuscaTodosAnuncios()
	c.Response().Header().Set("HX-Push-Url", "/")
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

func (h *AnunciosHandler) MostraPaginaDeAnunciosDoUsuario(ctx echo.Context) error {
	// nomeDeUsuario, _ := buscaNomeDeUsuarioDaSessao(ctx, h.logger)
	// busco anuncios para o usuario
	// mostra a pagina de anuncios do usuario
	anun := []model.Anuncio{}
	return render(ctx, http.StatusOK, components.MeusAnuncios(anun))
}

func (h *AnunciosHandler) MostraPaginaDeCriacaoDeNovoAnuncio(ctx echo.Context) error {
	return render(ctx, http.StatusOK, components.NovoAnuncio())
}

func (h *AnunciosHandler) CriaNovaOfertaParaAnuncio(c echo.Context) error {
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

func (h *AnunciosHandler) CriaNovoAnuncio(c echo.Context) error {
	// Validate the form
	nome, valor, descricao, imagem, err := h.validarFormularioDeAnuncio(c)

	if err != nil {
		h.logger.Errorf("Erro ao obter o arquivo de imagem: %v", err)
		return c.String(http.StatusBadRequest, "Erro ao obter o arquivo de imagem")
	}

	// Save the image and get the public path
	imagemPath, err := h.anunSvc.SalvaImagem(imagem)
	if err != nil {
		h.logger.Errorf("Erro ao salvar a imagem: %v", err)
		return c.String(http.StatusInternalServerError, "Erro ao salvar a imagem")
	}

	// Create the new anuncio
	anuncio := model.NewAnuncio(nome, valor, descricao, imagemPath)
	if err := h.anunSvc.CriaNovoAnuncio(anuncio); err != nil {
		h.logger.Errorf("Erro ao criar o novo anuncio: %v", err)
		return c.String(http.StatusInternalServerError, "Erro ao criar o novo anuncio")
	}

	// Redirect to the anuncios page
	return c.Redirect(http.StatusSeeOther, "/anuncios")
}

func (h *AnunciosHandler) validarFormularioDeAnuncio(c echo.Context) (string, float64, string, *multipart.FileHeader, error) {
	nome := c.FormValue("nome")
	valor, _ := strconv.ParseFloat(c.FormValue("valor"), 64)
	descricao := c.FormValue("descricao")
	imagem, err := c.FormFile("imagem")

	if nome == "" {
		return "", 0, "", nil, echo.NewHTTPError(http.StatusBadRequest, "O campo nome não pode estar vazio")
	}
	if valor == 0 {
		return "", 0, "", nil, echo.NewHTTPError(http.StatusBadRequest, "O valor não pode ser zero")
	}
	if descricao == "" {
		return "", 0, "", nil, echo.NewHTTPError(http.StatusBadRequest, "O campo descrição não pode estar vazio")
	}
	if err != nil {
		return "", 0, "", nil, echo.NewHTTPError(http.StatusBadRequest, "O campo imagem não pode estar vazio")
	}
	if imagem.Size > 10*1024*1024 { // 10 MB
		return "", 0, "", nil, echo.NewHTTPError(http.StatusBadRequest, "O tamanho da imagem não pode ser maior que 10 MB")
	}
	return nome, valor, descricao, imagem, nil
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
