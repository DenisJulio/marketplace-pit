package services

import (
	"errors"

	"github.com/DenisJulio/marketplace-pit/model"
	"github.com/DenisJulio/marketplace-pit/store"
)

var (
	ErrDadosParaRegistroInvalidos = errors.New("Dados fornecidos sao invalidos.")
	ErrUsuarioExistente           = errors.New("Nome de usuario ja esta cadastrado.")
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

func (us *UsuarioService) VerificaUsuarioExistente(nomeDeUsuario string) bool {
	return us.s.VerificaUsuarioExistente(nomeDeUsuario)
}

func (us *UsuarioService) RegistraNovoUsuario(nome, nomeDeUsuario, senha string) error {
	if err := validaDadosParaRegistro(nome, nomeDeUsuario, senha); err != nil {
		return ErrDadosParaRegistroInvalidos
	}
	if us.s.VerificaUsuarioExistente(nomeDeUsuario) {
		return ErrUsuarioExistente
	}
	if err := us.s.InsereNovoUsuario(nomeDeUsuario, nome, senha); err != nil {
		return err
	}
	return nil
}

func validaDadosParaRegistro(nome string, nomeDeUsuario string, senha string) error {
	if nome == "" || nomeDeUsuario == "" || senha == "" {
		return errors.New("Campos obrigatórios não preenchidos")
	}
	if len(senha) < 8 {
		return errors.New("A senha deve ter no mínimo 8 caracteres")
	}
	return nil
}
