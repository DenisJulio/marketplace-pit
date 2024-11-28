package services

import (
	"github.com/DenisJulio/marketplace-pit/model"
	"github.com/DenisJulio/marketplace-pit/store"
)

type OfertaService struct {
	s store.OfertaStore
}

func NewOfertaService(s store.OfertaStore) *OfertaService {
	return &OfertaService{s: s}
}

func (os *OfertaService) CriaNovaOfertaParaAnuncio(oferta model.Oferta, msg model.Mensagem) (int, error) {
	return os.s.CriaNovaOfertaParaAnuncio(oferta, msg)
}

func (os *OfertaService) BuscaTodasAsOfertasExpandidasDoUsuario(usuarioId int) ([]model.OfertaExpandida, error) {
	return os.s.BuscaTodasAsOfertasExpandidasDoUsuario(usuarioId)
}
