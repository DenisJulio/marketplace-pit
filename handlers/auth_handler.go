package handlers

import (
	"net/http"

	"github.com/DenisJulio/marketplace-pit/components"
	"github.com/DenisJulio/marketplace-pit/utils"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	logger utils.Logger	
}

func NovoAuthHandler(logger utils.Logger) *AuthHandler {
	return &AuthHandler{logger: logger}
}

func (a *AuthHandler) MostraTelaDeLogin(c echo.Context) error {
	redirectTo := c.QueryParam("redirect_to")
	return render(c, http.StatusOK, components.PaginaDeLogin(redirectTo))
}

func (a *AuthHandler) MostraTelaDeCadastro(c echo.Context) error {
	a.logger.Debugf("Exibindo formulário de cadastro para novo usuario")
	return render(c, http.StatusOK, components.PaginaDeCadastro())
}

