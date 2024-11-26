package services

import (
	"github.com/DenisJulio/marketplace-pit/model"
	"github.com/DenisJulio/marketplace-pit/store"
)

type AnuncioServices struct {
	store store.AnuncioStore
}

func NewAnuncioService(s store.AnuncioStore) *AnuncioServices {
	return &AnuncioServices{store: s}
}

func (s *AnuncioServices) BuscaTodosAnuncios() []model.Anuncio {
	a, _ := s.store.BuscaTodosAnuncios()
	return a
}

func (s *AnuncioServices) BuscaAnuncioPorID(id int) (model.Anuncio, error) {
	return s.store.BuscaAnuncioPorID(id)
}

func (s *AnuncioServices) CriaNovoAnuncio(anuncio model.Anuncio) error {
	return s.store.SalvaNovoAnuncio(anuncio)
}
