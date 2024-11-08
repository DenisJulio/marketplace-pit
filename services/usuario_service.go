package services

import (
	"github.com/DenisJulio/marketplace-pit/model"
	"github.com/DenisJulio/marketplace-pit/store"
)

type UsuarioService struct {
	s store.UsuarioStore
}

func NewUsuarioService(s store.UsuarioStore) *UsuarioService {
	return &UsuarioService{s: s}
}

func (us *UsuarioService) BuscaUsuarioPorId(id int) (model.Usuario, error) {
	return us.s.BuscaUsuarioPorId(id)
}
