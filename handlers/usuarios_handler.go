package handlers

import (
	"errors"
	"net/http"

	"github.com/DenisJulio/marketplace-pit/components"
	"github.com/DenisJulio/marketplace-pit/services"
	"github.com/DenisJulio/marketplace-pit/utils"
	"github.com/labstack/echo/v4"
)

type UsuarioHandler struct {
	logger utils.Logger
	usuSvc services.UsuarioService
}

func NovoUsuarioHandler(usuSvc services.UsuarioService, logger utils.Logger) *UsuarioHandler {
	return &UsuarioHandler{usuSvc: usuSvc, logger: logger}
}

func (h *UsuarioHandler) CadastraNovoUsuario(ctx echo.Context) error {
	h.logger.Debugf("Iniciando cadastro de novo usuario")

	nome := ctx.FormValue("nome")
	nomeDeUsuario := ctx.FormValue("nomeDeUsuario")
	senha := ctx.FormValue("senha")

	h.logger.Debugf("Recebendo dados para cadastro: %s, %s, %s", nome, nomeDeUsuario, senha)

	if err := h.usuSvc.RegistraNovoUsuario(nome, nomeDeUsuario, senha); err != nil {
		if errors.Is(err, services.ErrDadosParaRegistroInvalidos) {
			return ctx.NoContent(http.StatusBadRequest)
		}
		if errors.Is(err, services.ErrUsuarioExistente) {
			return ctx.NoContent(http.StatusConflict)
		}
		return ctx.NoContent(http.StatusInternalServerError)
	}
	ctx.Response().Header().Set("HX-Redirect", "/login")
	return ctx.NoContent(http.StatusCreated)
}

func (h *UsuarioHandler) AutenticaUsuario(c echo.Context) error {
	h.logger.Debugf("Autenticando usuario")

	nomeDeUsuario := c.FormValue("nomeDeUsuario")
	senha := c.FormValue("senha")
	redirectTo := c.FormValue("redirect_to")

	h.logger.Debugf("Recebendo dados para login: %s, %s", nomeDeUsuario, senha)

	var err error
	if err = validaDadosParaLogin(nomeDeUsuario, senha); err != nil {
		h.logger.Errorf("%v", err)
		render(c, http.StatusOK, components.AlertaErroAutenticacao())
	}
	err = h.usuSvc.VerificaSegredosDeUsuario(nomeDeUsuario, senha)
	if err != nil {
		h.logger.Errorf("Erro ao autenticar o usuario: %s. %v", nomeDeUsuario, err)
		render(c, http.StatusOK, components.AlertaErroAutenticacao())
	}

	if err = iniciarSessao(c, nomeDeUsuario); err != nil {
		h.logger.Errorf("Erro ao iniciar a sessao para o usuario: %s. %v", nomeDeUsuario, err)
	}

	if redirectTo == "" {
		redirectTo = "/"
	}
	c.Response().Header().Set("HX-Redirect", redirectTo)
	return c.NoContent(http.StatusOK)
}

func (h *UsuarioHandler) MostraBotaoDeEntrarNaConta(c echo.Context) error {
	nomeDeUsuario, err := buscaNomeDeUsuarioDaSessao(c, h.logger)
	if err != nil || nomeDeUsuario == "" {
		return render(c, http.StatusOK, components.EntrarNaConta(false, ""))
	}
	usuario, err := h.usuSvc.BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario)
	if err != nil {
		return render(c, http.StatusOK, components.EntrarNaConta(false, ""))
	}
	return render(c, http.StatusOK, components.EntrarNaConta(true, *usuario.Imagem))
}

func validaDadosParaLogin(nomeDeUsuario, senha string) error {
	if nomeDeUsuario == "" || senha == "" {
		return errors.New("Dados para login nao podem ser vazios")
	}
	return nil
}
