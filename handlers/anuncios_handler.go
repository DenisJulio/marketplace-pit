package handlers

import (
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/DenisJulio/marketplace-pit/components"
	"github.com/DenisJulio/marketplace-pit/model"
	"github.com/DenisJulio/marketplace-pit/services"
	"github.com/DenisJulio/marketplace-pit/store"
	"github.com/DenisJulio/marketplace-pit/utils"
	"github.com/labstack/echo/v4"
)

type AnunciosHandler struct {
	anunSvc services.AnuncioServices
	usuSvc  services.UsuarioService
	oftSvc  services.OfertaService
	imgSvc  services.ImagemService
	ssSvc   services.SessaoService
	logger  utils.Logger
}

func NewAnunciosHandler(anunSvc services.AnuncioServices, usuSvc services.UsuarioService, oftSvc services.OfertaService, imgSvc services.ImagemService, ssSvc services.SessaoService, logger utils.Logger) *AnunciosHandler {
	return &AnunciosHandler{anunSvc: anunSvc, usuSvc: usuSvc, oftSvc: oftSvc, imgSvc: imgSvc, ssSvc: ssSvc, logger: logger}
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
	nome, valor, descricao, imagem, err := h.validarFormularioDeAnuncio(c)

	h.logger.Debugf("Nome: %s, Valor: %f, Descricao: %s, Imagem: %s", nome, valor, descricao, imagem.Filename)

	if err != nil {
		h.logger.Errorf("Erro ao obter o arquivo de imagem: %v", err)
		return c.String(http.StatusBadRequest, "Erro ao obter o arquivo de imagem")
	}

	imagemPath, err := h.imgSvc.SalvalNovaImagem(store.ImagemDeAnuncio, imagem)
	if err != nil {
		h.logger.Errorf("Erro ao salvar a imagem: %v", err)
		return c.String(http.StatusInternalServerError, "Erro ao salvar a imagem")
	}

	nomeDeUsuario, _ := h.ssSvc.BuscaNomeDeUsuarioDaSessao(c)
	usuario, _ := h.usuSvc.BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario)

	// Create the new anuncio
	anuncio := model.NovoAnuncio(nome, usuario.ID, valor, &descricao, &imagemPath)
	if err := h.anunSvc.CriaNovoAnuncio(anuncio); err != nil {
		h.logger.Errorf("Erro ao criar o novo anuncio: %v", err)
		return c.String(http.StatusInternalServerError, "Erro ao criar o novo anuncio")
	}

	// Redirect to the anuncios page
	h.logger.Debugf("Anuncio criado com sucesso: %s", nome)
	return c.NoContent(http.StatusNoContent)
}

func (h *AnunciosHandler) validarFormularioDeAnuncio(c echo.Context) (string, float64, string, *multipart.FileHeader, error) {
	nome := c.FormValue("nome")
	valorStr := c.FormValue("valor")
	valorStr = strings.ReplaceAll(valorStr, ",", ".")
	valor, _ := strconv.ParseFloat(valorStr, 64)
	descricao := c.FormValue("descricao")
	imagem, err := c.FormFile("imagem")

	h.logger.Debugf("Nome: %s, Valor: %f, Descricao: %s, Imagem: %s", nome, valor, descricao, imagem.Filename)

	if nome == "" {
		h.logger.Errorf("O campo nome não pode estar vazio")
		return "", 0, "", nil, echo.NewHTTPError(http.StatusBadRequest, "O campo nome não pode estar vazio")
	}
	if valor == 0 {
		h.logger.Errorf("O valor não pode ser zero")
		return "", 0, "", nil, echo.NewHTTPError(http.StatusBadRequest, "O valor não pode ser zero")
	}
	if descricao == "" {
		h.logger.Errorf("O campo descrição não pode estar vazio")
		return "", 0, "", nil, echo.NewHTTPError(http.StatusBadRequest, "O campo descrição não pode estar vazio")
	}
	if err != nil {
		h.logger.Errorf("O campo imagem não pode estar vazio")
		return "", 0, "", nil, echo.NewHTTPError(http.StatusBadRequest, "O campo imagem não pode estar vazio")
	}
	if imagem.Size > 10*1024*1024 { // 10 MB
		h.logger.Errorf("O tamanho da imagem não pode ser maior que 10 MB")
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
