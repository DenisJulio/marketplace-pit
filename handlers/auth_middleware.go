package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/DenisJulio/marketplace-pit/services"
	"github.com/DenisJulio/marketplace-pit/utils"
	"github.com/labstack/echo/v4"
)

type Middleware struct {
	ssVs  services.SessaoService
	logger utils.Logger
}

func NovoMiddleware(ssVs services.SessaoService,logger utils.Logger) *Middleware {
	return &Middleware{ssVs: ssVs,logger: logger}
}

func (m *Middleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		m.logger.Debugf("Iniciando autenticacao para: %s", c.Request().URL.Path)

		// nomeDeUsuario, err := buscaNomeDeUsuarioDaSessao(c, m.logger)
		nomeDeUsuario, err := m.ssVs.BuscaNomeDeUsuarioDaSessao(c)
		if err != nil || nomeDeUsuario == "" {
			m.logger.Debugf("Request nao autenticado para: %s", c.Request().URL.Path)
			reqUrl := c.Request().URL.Path
			loginURL := fmt.Sprintf("/login?redirect_to=%s", url.QueryEscape(reqUrl))
			if c.Request().Header.Get("HX-Request") == "true" {
				c.Response().Header().Set("HX-Redirect", loginURL)
				return c.NoContent(http.StatusUnauthorized)
			}
			return c.Redirect(http.StatusSeeOther, loginURL)
		}
		return next(c)
	}
}
