package services

import (
	"errors"

	"github.com/DenisJulio/marketplace-pit/model"
	"github.com/DenisJulio/marketplace-pit/store"
	"github.com/DenisJulio/marketplace-pit/utils"
)

var (
	ErrDadosParaRegistroInvalidos = errors.New("Dados fornecidos sao invalidos.")
	ErrUsuarioExistente           = errors.New("Nome de usuario ja esta cadastrado.")
)

type UsuarioService struct {
	s      store.UsuarioStore
	logger utils.Logger
}

type segredosDeUsuario struct {
	nomeDeUsuario string
	senha         string
}

func NewUsuarioService(s store.UsuarioStore, logger utils.Logger) *UsuarioService {
	return &UsuarioService{s: s, logger: logger}
}

func (us *UsuarioService) BuscaUsuarioPorId(id int) (model.Usuario, error) {
	return us.s.BuscaUsuarioPorId(id)
}

func (us *UsuarioService) BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario string) (model.Usuario, error) {
	return us.s.BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario)
}

func (us *UsuarioService) VerificaUsuarioExistente(nomeDeUsuario string) bool {
	return us.s.VerificaUsuarioExistente(nomeDeUsuario)
}

func (us *UsuarioService) RegistraNovoUsuario(nome, email, nomeDeUsuario, senha, imagem string) error {
	if err := validaDadosParaRegistro(nome, nomeDeUsuario, senha); err != nil {
		us.logger.Debugf("%v", err)
		return ErrDadosParaRegistroInvalidos
	}
	if us.s.VerificaUsuarioExistente(nomeDeUsuario) {
		us.logger.Debugf("%v", ErrUsuarioExistente)
		return ErrUsuarioExistente
	}
	if err := us.s.InsereNovoUsuario(nomeDeUsuario, email, nome, senha, imagem); err != nil {
		return err
	}
	return nil
}

func (us *UsuarioService) VerificaSegredosDeUsuario(nomeDeUsuario, senha string) error {
	_, err := us.s.VerificaSegredosDeUsuario(nomeDeUsuario, senha)
	if err != nil {
		return err
	}
	return nil
}

func (us *UsuarioService) AtualizaNome(nomeDeUsuario, nome string) error {
	return us.s.AtualizaNome(nomeDeUsuario, nome)
}

func (us *UsuarioService) AtualizaImagemDeUsuario(nomeDeUsuario, imagem string) (string, error) {
	return us.s.AtualizaImagemDeUsuario(nomeDeUsuario, imagem)
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
