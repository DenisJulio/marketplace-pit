package handlers

import (
	"errors"
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
	nomeDeUsuario, _ := h.ssSvc.BuscaNomeDeUsuarioDaSessao(c)
	usuario, _ := h.usuSvc.BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario)
	id, err := h.extrairId(c.Param("id"))
	if err != nil {
		// TODO: render not found page
	}
	anuncio, err := h.anunSvc.BuscaAnuncioPorID(id)
	if err != nil {
		h.logger.Errorf("Anuncio com id=%d nao encontrado", id)
		// TODO: render not found page
	}
	eAnunciante := usuario.ID == anuncio.AnuncianteId
	return render(c, http.StatusOK, components.DetalhesDoAnuncio(anuncio, eAnunciante))
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
	nomeDeUsuario, _ := h.ssSvc.BuscaNomeDeUsuarioDaSessao(ctx)
	anun, _ := h.anunSvc.BuscaAnunciosPorNomeDeUsuario(nomeDeUsuario)
	return render(ctx, http.StatusOK, components.MeusAnuncios(anun))
}

func (h *AnunciosHandler) MostraPaginaDeCriacaoDeNovoAnuncio(ctx echo.Context) error {
	return render(ctx, http.StatusOK, components.NovoAnuncio())
}

func (h *AnunciosHandler) CriaNovaOfertaParaAnuncio(c echo.Context) error {
	nUsOfertante, _ := h.ssSvc.BuscaNomeDeUsuarioDaSessao(c)
	ofertante, _ := h.usuSvc.BuscaUsuarioPorNomeDeUsuario(nUsOfertante)

	// extrair e validar o formulario
	msgForm, err := h.validarFormulario(c, "mensagem")
	if err != nil {
		c = enviaNotificacaoToast(c, toastErro, "Erro no envio", "Voce deve enviar uma mensagem para o anunciante.")
		c.NoContent(http.StatusBadRequest)
	}

	id, _ := h.extrairId(c.Param("id"))
	anuncio, _ := h.anunSvc.BuscaAnuncioPorID(id)
	anunciante, _ := h.usuSvc.BuscaUsuarioPorId(anuncio.AnuncianteId)

	// create the oferta
	oferta := model.NewOferta(anunciante.ID, ofertante.ID, anuncio.ID)
	mensagen := model.NewMensagem(ofertante.ID, msgForm)
	_, errOft := h.oftSvc.CriaNovaOfertaParaAnuncio(*oferta, *mensagen)
	if errOft != nil {
		h.logger.Errorf("Erro ao criar a nova oferta: %v", errOft)
		return c.NoContent(http.StatusInternalServerError)
	}

	// redirect to Home
	c = enviaNotificacaoToast(c, toastSucesso, "Oferta enviada", "Sua oferta foi enviada com sucesso.")
	return c.NoContent(http.StatusNoContent)
}

func (h *AnunciosHandler) CriaNovoAnuncio(c echo.Context) error {
	nome, valor, descricao, imagem, err := h.validarFormularioDeAnuncio(c)

	if err != nil {
		c = enviaNotificacaoToast(c, toastErro, "Erro no envio", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	imagemPath, err := h.imgSvc.SalvalNovaImagem(store.ImagemDeAnuncio, imagem)
	if err != nil {
		h.logger.Errorf("Erro ao salvar a imagem: %v", err)
		c = enviaNotificacaoToast(c, toastErro, "Erro interno", "Erro ao salvar a imagem")
		return c.NoContent(http.StatusInternalServerError)
	}

	nomeDeUsuario, _ := h.ssSvc.BuscaNomeDeUsuarioDaSessao(c)
	usuario, _ := h.usuSvc.BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario)

	anuncio := model.NovoAnuncio(nome, usuario.ID, valor, &descricao, &imagemPath)
	if err := h.anunSvc.CriaNovoAnuncio(anuncio); err != nil {
		h.logger.Errorf("Erro ao criar o novo anuncio: %v", err)
		c = enviaNotificacaoToast(c, toastErro, "Erro interno", "Erro ao criar o anuncio")
		return c.NoContent(http.StatusInternalServerError)
	}

	h.logger.Debugf("Anuncio criado com sucesso: %s", nome)
	return c.NoContent(http.StatusNoContent)
}

func (h *AnunciosHandler) RemoveAnuncio(c echo.Context) error {
	id, err := h.extrairId(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	nomeDeUsuario, _ := h.ssSvc.BuscaNomeDeUsuarioDaSessao(c)
	usuario, _ := h.usuSvc.BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario)

	anuncio, err := h.anunSvc.BuscaAnuncioPorID(id)
	if err != nil {
		c = enviaNotificacaoToast(c, toastErro, "Erro interno", "Erro ao buscar o anuncio")
		return c.NoContent(http.StatusInternalServerError)
	}

	if anuncio.AnuncianteId != usuario.ID {
		c = enviaNotificacaoToast(c, toastErro, "Erro interno", "Voce nao tem permissao para remover este anuncio")
		return c.NoContent(http.StatusForbidden)
	}

	if err := h.anunSvc.RemoveAnuncio(anuncio.ID); err != nil {
		c = enviaNotificacaoToast(c, toastErro, "Erro interno", "Erro ao remover o anuncio")
		return c.NoContent(http.StatusInternalServerError)
	}
	h.imgSvc.RemoveImagem(store.ImagemDeAnuncio, *anuncio.Imagem)

	h.logger.Debugf("Anuncio removido com sucesso: %d", id)
	return c.NoContent(http.StatusOK)
}

func (h *AnunciosHandler) validarFormularioDeAnuncio(c echo.Context) (string, float64, string, *multipart.FileHeader, error) {
	nome := c.FormValue("nome")
	valorStr := c.FormValue("valor")
	valorStr = strings.ReplaceAll(valorStr, ",", ".")
	valor, _ := strconv.ParseFloat(valorStr, 64)
	descricao := c.FormValue("descricao")
	imagem, err := c.FormFile("imagem")

	if strings.TrimSpace(nome) == "" {
		h.logger.Errorf("O campo nome não pode estar vazio")
		return "", 0, "", nil, errors.New("Um nome para o anuncio deve ser fornecido")
	}
	if valor == 0 {
		h.logger.Errorf("O valor não pode ser zero")
		return "", 0, "", nil, errors.New("O valor do anuncio deve ser maior que zero")
	}
	if strings.TrimSpace(descricao) == "" {
		h.logger.Errorf("O campo descrição não pode estar vazio")
		return "", 0, "", nil, errors.New("O anuncio deve ter uma descricao")
	}
	if err != nil {
		h.logger.Errorf("Erro ao obter o arquivo de imagem: %v", err)
		return "", 0, "", nil, errors.New("Voce deve enviar uma imagem para o produto")
	}
	if imagem.Size > 10*1024*1024 { // 10 MB
		h.logger.Errorf("Imagem excede o tamanho maximo permitido de 10 MB")
		return "", 0, "", nil, errors.New("O tamanho da imagem nao pode ser maior que 10 MB")
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
	if strings.TrimSpace(val) == "" {
		h.logger.Errorf("Campo do formalario: %s nao pode ser vazio", campoForm)
		return "", echo.ErrBadRequest
	}
	return val, nil
}
