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
	ssSvs  services.SessaoService
}

func NovoUsuarioHandler(usuSvc services.UsuarioService, ssSvc services.SessaoService, logger utils.Logger) *UsuarioHandler {
	return &UsuarioHandler{usuSvc: usuSvc, logger: logger, ssSvs: ssSvc}
}

func (h *UsuarioHandler) CadastraNovoUsuario(ctx echo.Context) error {
	h.logger.Debugf("Iniciando cadastro de novo usuario")

	nome := ctx.FormValue("nome")
	nomeDeUsuario := ctx.FormValue("nomeDeUsuario")
	senha := ctx.FormValue("senha")

	h.logger.Debugf("Recebendo dados para cadastro: %s, %s, %s", nome, nomeDeUsuario, senha)

	imagemPadrao := "/resources/images/avatars/avatar-padrao.png"

	if err := h.usuSvc.RegistraNovoUsuario(nome, nomeDeUsuario, senha, imagemPadrao); err != nil {
		if errors.Is(err, services.ErrDadosParaRegistroInvalidos) {
			return ctx.NoContent(http.StatusBadRequest)
		}
		if errors.Is(err, services.ErrUsuarioExistente) {
			return ctx.NoContent(http.StatusConflict)
		}
		return ctx.NoContent(http.StatusInternalServerError)
	}
	ctx = enviaNotificacaoToast(ctx, toastSucesso, "Cadastro concluido", "Cadastro concluido com sucesso")
	return ctx.NoContent(http.StatusNoContent)
}

func (h *UsuarioHandler) ValidaNomeDeUsuarioNaoExistente(ctx echo.Context) error {
	nomeDeUsuario := ctx.FormValue("nomeDeUsuario")
	usuarioExiste := h.usuSvc.VerificaUsuarioExistente(nomeDeUsuario)
	return render(ctx, http.StatusOK, components.AlertaValidacaoNomeDeUsuario(usuarioExiste))
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
		ctx := enviaNotificacaoToast(c, toastErro, "Erro de Login", "Usuario ou senha invalidos")
		return ctx.NoContent(http.StatusBadRequest)
	}

	if err = h.ssSvs.IniciarSessao(c, nomeDeUsuario); err != nil {
		h.logger.Errorf("Erro ao iniciar a sessao para o usuario: %s. %v", nomeDeUsuario, err)
	}

	if redirectTo == "" {
		redirectTo = "/"
	}
	c.Response().Header().Set("HX-Redirect", redirectTo)
	return c.NoContent(http.StatusOK)
}

func (h *UsuarioHandler) EncerraSessao(ctx echo.Context) error {
	if err := h.ssSvs.EncerraSessao(ctx); err != nil {
		h.logger.Errorf("Erro ao encerrar a sessao: %v", err)		
		ctx = enviaNotificacaoToast(ctx, toastErro, "Erro", "Erro ao encerrar sessão")
		return ctx.NoContent(http.StatusInternalServerError)
	}	
	ctx = enviaNotificacaoToast(ctx, toastSucesso, "Ate mais", "Sessao encerrada com sucesso")
	return ctx.NoContent(http.StatusNoContent)
}

func (h *UsuarioHandler) MostraPaginaDeMinhaConta(ctx echo.Context) error {
	nomeDeUsuario, err := h.ssSvs.BuscaNomeDeUsuarioDaSessao(ctx)
	if err != nil {
		return render(ctx, http.StatusOK, components.EntrarNaConta(false, ""))
	}
	usuario, _ := h.usuSvc.BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario)
	return render(ctx, http.StatusOK, components.MinhaConta(usuario))
}

func (h *UsuarioHandler) MostraBotaoDeEntrarNaConta(ctx echo.Context) error {
	nomeDeUsuario, err := h.ssSvs.BuscaNomeDeUsuarioDaSessao(ctx)
	if err != nil {
		return render(ctx, http.StatusOK, components.EntrarNaConta(false, ""))
	}
	usuario, _ := h.usuSvc.BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario)
	return render(ctx, http.StatusOK, components.EntrarNaConta(true, *usuario.Imagem))
}

func (h *UsuarioHandler) CarregaFormularioNomeDisplay(c echo.Context) error {
	return render(c, http.StatusOK, components.NomeLabelForm())
}

func (h *UsuarioHandler) AtualizaNomeDisplay(ctx echo.Context) error {
	nomeDeUsuario, err := h.ssSvs.BuscaNomeDeUsuarioDaSessao(ctx)
	if err != nil {
		h.logger.Errorf("Erro %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	usuario, _ := h.usuSvc.BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario)

	nomeDisplay := ctx.FormValue("nome")
	if err := h.usuSvc.AtualizaNome(usuario.NomeDeUsuario, nomeDisplay); err != nil {
		h.logger.Errorf("Erro ao atualizar o nome: %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return render(ctx, http.StatusOK, components.NomeLabel(nomeDisplay))
}

func (h *UsuarioHandler) UploadAvatar(ctx echo.Context) error {
	nomeDeUsuario, err := h.ssSvs.BuscaNomeDeUsuarioDaSessao(ctx)
	if err != nil {
		h.logger.Errorf("Erro %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	usuario, _ := h.usuSvc.BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario)

	h.logger.Debugf("Iniciando upload de imagem para avatar")
	const maxUploadSize = 5 * 1024 * 1024
	file, err := ctx.FormFile("avatar-image")
	if err != nil {
		h.logger.Errorf("Erro ao obter o arquivo de imagem: %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	imgPath, err := h.usuSvc.SalvalNovaImagemDeAvatar(usuario.NomeDeUsuario, file)

	return render(ctx, http.StatusOK, components.ImagemAvatar(imgPath))
}

func (h *UsuarioHandler) CarregaAvatar(ctx echo.Context) error {
	nomeDeUsuario, err := h.ssSvs.BuscaNomeDeUsuarioDaSessao(ctx)
	if err != nil {
		h.logger.Errorf("Erro %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	usuario, _ := h.usuSvc.BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario)
	return render(ctx, http.StatusOK, components.ImagemAvatarNav(*usuario.Imagem))
}

func validaDadosParaLogin(nomeDeUsuario, senha string) error {
	if nomeDeUsuario == "" || senha == "" {
		return errors.New("Dados para login nao podem ser vazios")
	}
	return nil
}
