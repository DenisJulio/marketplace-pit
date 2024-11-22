package services

import (
	"net/http"
	"time"

	"github.com/DenisJulio/marketplace-pit/store"
	"github.com/DenisJulio/marketplace-pit/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SessaoService struct {
	s      store.SessaoStore
	logger utils.Logger
}

func NovaSessaoService(s store.SessaoStore, logger utils.Logger) *SessaoService {
	return &SessaoService{s: s, logger: logger}
}

func (ss *SessaoService) IniciarSessao(ctx echo.Context, nomeDeUsuario string) error {
	ss.logger.Debugf("Iniciando sessao para: %s", ctx.Request().URL.Path)
	sessaoID := uuid.NewString()
	expiraEm := time.Now().Add(24 * time.Hour)
	// Salva sessao no banco de dados
	if err := ss.s.SalvaSessao(sessaoID, nomeDeUsuario, expiraEm); err != nil {
		ss.logger.Errorf("Erro ao salvar sessao: %v", err)
		return err
	}
	// Cria cookie de sessao
	cookie := &http.Cookie{
		Name:     "sessaoID",
		Value:    sessaoID,
		Path:     "/",
		Expires:  expiraEm,
		HttpOnly: true,
		Secure:   true,
	}
	ctx.SetCookie(cookie)
	ss.logger.Debugf("Sessao iniciada para: usuario %s acessar %s", nomeDeUsuario, ctx.Request().URL.Path)
	return nil
}

func (ss *SessaoService) BuscaNomeDeUsuarioDaSessao(ctx echo.Context) (string, error) {
	ss.logger.Debugf("Buscando sessao para: %s", ctx.Request().URL.Path)
	cookie, err := ctx.Cookie("sessaoID")
	if err != nil {
		return "", err
	}
	// Busca pela sessao no banco de dados
	nomeDeUsuario, err := ss.s.BuscaSessao(cookie.Value)
	if err != nil {
		ss.logger.Errorf("Erro ao buscar sessao: %v", err)
		return "", err
	}
	if nomeDeUsuario == "" {
		ss.logger.Debugf("Sessao nao encontrada para: %s", ctx.Request().URL.Path)
		return "", echo.ErrUnauthorized
	}
	ss.logger.Debugf("Sessao encontrada para: usuario %s acessar %s", nomeDeUsuario, ctx.Request().URL.Path)
	return nomeDeUsuario, nil
}

func (ss *SessaoService) EncerraSessao(ctx echo.Context) error {
	ss.logger.Debugf("Iniciando encerramento de sessao")
	cookie, err := ctx.Cookie("sessaoID")
	if err != nil {
		if err == http.ErrNoCookie {
			ss.logger.Debugf("Nenhum cookie de sessao encontrado para logout")
			return echo.ErrUnauthorized
		}
		ss.logger.Errorf("Erro ao recuperar cookie de sessao: %v", err)
		return err
	}
	// Encerra sessao no banco de dados
	if err := ss.s.RemoveSessao(cookie.Value); err != nil {
		ss.logger.Errorf("Erro ao remover sessao no banco de dados: %v", err)
		return err
	}
	// Remove cookie no browser
	expiredCookie := &http.Cookie{
		Name:     "sessaoID",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0), // Expira imediatamente
		HttpOnly: true,
		Secure:   true,
	}
	ctx.SetCookie(expiredCookie)
	ss.logger.Debugf("Sessao encerrada com sucesso")
	return nil
}
