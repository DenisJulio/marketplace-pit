package services

import "github.com/DenisJulio/marketplace-pit/store"

type OfertaService struct {
	s store.OfertaStore
}

func NewOfertaService(s store.OfertaStore) *OfertaService {
	return &OfertaService{s: s}
}

func (s *OfertaService) CriaNovaOfertaParaAnuncio(msg string) {
	// TODO
}
