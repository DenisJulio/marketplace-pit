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

func (h *UsuarioHandler) RegistraNovoUsuario(ctx echo.Context) error {
	nome := ctx.FormValue("nome")
	nomeDeUsuario := ctx.FormValue("nomeDeUsuario")
	senha := ctx.FormValue("senha")

	if err := h.usuSvc.RegistraNovoUsuario(nome, nomeDeUsuario, senha); err != nil {
		if errors.Is(err, services.ErrDadosParaRegistroInvalidos) {
			return ctx.NoContent(http.StatusBadRequest)
		}
		if errors.Is(err, services.ErrDadosParaRegistroInvalidos) {
			return ctx.NoContent(http.StatusConflict)
		}
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return render(ctx, http.StatusCreated, components.PaginaDeLogin())
}
