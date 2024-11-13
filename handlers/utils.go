package handlers

import (
	"net/http"
	"time"

	"github.com/DenisJulio/marketplace-pit/utils"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var sessaoStore = map[string]string{}

func render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)
	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return nil
	}
	return ctx.HTML(statusCode, buf.String())
}

func iniciarSessao(ctx echo.Context, nomeDeUsuario string) error {
	sessaoID := uuid.NewString()
	sessaoStore[sessaoID] = nomeDeUsuario
	cookie := &http.Cookie{
		Name:     "sessaoID",
		Value:    sessaoID,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}
	ctx.SetCookie(cookie)
	return nil
}

func buscaNomeDeUsuarioDaSessao(ctx echo.Context, logger utils.Logger) (string, error) {
	logger.Debugf("Buscando sessao para: %s", ctx.Request().URL.Path)
	cookie, err := ctx.Cookie("sessaoID")
	if err != nil {
		return "", err
	}
	un, ok := sessaoStore[cookie.Value]
	if !ok {
		return "", echo.ErrUnauthorized
	}
	return un, nil
}
