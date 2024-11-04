package handlers

import (
	"net/http"

	"github.com/DenisJulio/marketplace-pit/components"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)
	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return nil
	}
	return ctx.HTML(statusCode, buf.String())
}

func HomeHandler(c echo.Context) error {
	return render(c, http.StatusOK, components.AnunciosPage())
}
