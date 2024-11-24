package handlers

import (
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

type tipoDeToast int

const (
	toastSucesso tipoDeToast = iota
	toastErro
)

func (t tipoDeToast) String() string {
	return [...]string{"toastSucesso", "toastErro"}[t]
}

func enviaNotificacaoToast(ctx echo.Context, t tipoDeToast, titulo, msg string) echo.Context {
	ctx.Response().Header().Set("X-Toast-Titulo", titulo)
	switch t {
	case toastSucesso:
		ctx.Response().Header().Set("X-Toast-Sucesso", msg)
	case toastErro:
		ctx.Response().Header().Set("X-Toast-Erro", msg)
	}
	return ctx
}
