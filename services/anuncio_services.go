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