package routes

import (
	"github.com/DenisJulio/marketplace-pit/handlers"
	"github.com/labstack/echo/v4"
)

type Router struct {
	echo  *echo.Echo
	ancH  handlers.AnunciosHandler
	usuH  handlers.UsuarioHandler
	authH handlers.AuthHandler
	mid   handlers.Middleware
}

func NewRouter(e *echo.Echo, ancH handlers.AnunciosHandler, usuH handlers.UsuarioHandler, authH handlers.AuthHandler, mid handlers.Middleware) *Router {
	return &Router{echo: e, ancH: ancH, authH: authH, usuH: usuH, mid: mid}
}

func (r *Router) RegisterRoutes() {
	r.echo.Static("resources", "public/static")
	r.echo.GET("/login", r.authH.MostraTelaDeLogin)
	r.echo.POST("/login", r.usuH.AutenticaUsuario)
	r.echo.GET("/sair", r.usuH.EncerraSessao)
	r.echo.GET("/entrar-na-conta", r.usuH.MostraBotaoDeEntrarNaConta)
	r.echo.GET("/cadastro", r.authH.MostraTelaDeCadastro)
	r.echo.POST("/cadastro", r.usuH.CadastraNovoUsuario)
	r.echo.POST("/cadastro/nome-de-usuario", r.usuH.ValidaNomeDeUsuarioNaoExistente)
	r.echo.GET("/", r.ancH.MostraTelaDeAnuncios)
	r.echo.GET("/conta", r.usuH.MostraPaginaDeMinhaConta, r.mid.AuthMiddleware)
	r.echo.GET("/conta/formulario-nome", r.usuH.CarregaFormularioNomeDisplay)
	r.echo.GET("/conta/avatar", r.usuH.CarregaAvatar)
	r.echo.POST("/conta/avatar", r.usuH.UploadAvatar)
	r.echo.POST("/conta/nome", r.usuH.AtualizaNomeDisplay)
	r.echo.GET("/conta/anuncios", r.ancH.MostraPaginaDeAnunciosDoUsuario)
	r.echo.GET("/anuncios/:id", r.ancH.MostraDetalhesDoAnuncio)
	r.echo.GET("/anuncios/:id/nova-oferta", r.ancH.MostraTelaDeNovaOferta, r.mid.AuthMiddleware)
	r.echo.POST("/anuncios/:id/nova-oferta", r.ancH.CriaNovaOfertaParaAnuncio)
}
