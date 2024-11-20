package store

import (
	"database/sql"
	"errors"

	"github.com/DenisJulio/marketplace-pit/model"
	"github.com/DenisJulio/marketplace-pit/utils"
)

type UsuarioStore interface {
	BuscaUsuarioPorId(ID int) (model.Usuario, error)
	BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario string) (model.Usuario, error)
	VerificaUsuarioExistente(nomeDeUsuario string) bool
	InsereNovoUsuario(nomeDeUsuario, nome, senha string) error
	VerificaSegredosDeUsuario(nomeDeUsuario, senha string) (model.Usuario, error)
	AtualizaImagemDeUsuario(nomeDeUsuario, imagem string) error
	AtualizaNome(nomeDeUsuario, nome string) error
}

var (
	ErrUsuarioNaoEncontrado = errors.New("usuario nao encontrado")
	ErrSenhaInvalida        = errors.New("senha invalida")
)

type SQLUsuarioStore struct {
	db     *sql.DB
	logger utils.Logger
}

func NewSQLUsuarioStore(db *sql.DB, logger utils.Logger) *SQLUsuarioStore {
	return &SQLUsuarioStore{db: db, logger: logger}
}

func (s *SQLUsuarioStore) BuscaUsuarioPorId(ID int) (model.Usuario, error) {
	row := s.db.QueryRow("SELECT id, nome_de_usuario, nome, imagem FROM usuarios WHERE id = $1", ID)
	var u model.Usuario
	err := row.Scan(&u.ID, &u.NomeDeUsuario, &u.Nome, &u.Imagem)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Debugf("Usuario com ID=%d nao encontrado. %v", ID, err)
			return model.Usuario{}, ErrUsuarioNaoEncontrado
		}
		s.logger.Errorf("Erro ao buscar usuario por id=%d. %v", ID, err)
		return model.Usuario{}, err
	}
	return u, nil
}

func (s *SQLUsuarioStore) BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario string) (model.Usuario, error) {
	row := s.db.QueryRow("SELECT id, nome_de_usuario, nome, imagem FROM usuarios WHERE nome_de_usuario = $1", nomeDeUsuario)
	var u model.Usuario
	err := row.Scan(&u.ID, &u.NomeDeUsuario, &u.Nome, &u.Imagem)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Debugf("Usuario com nome de usuario=%s nao encontrado. %v", nomeDeUsuario, err)
			return model.Usuario{}, ErrUsuarioNaoEncontrado
		}
		s.logger.Errorf("Erro ao buscar usuario por nome de usuario:%s. %v", nomeDeUsuario, err)
		return model.Usuario{}, err
	}
	return u, nil
}

func (s *SQLUsuarioStore) VerificaUsuarioExistente(nomeDeUsuario string) bool {
	q := `SELECT 1 FROM usuarios WHERE nome_de_usuario = $1`
	var exists int
	err := s.db.QueryRow(q, nomeDeUsuario).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		s.logger.Errorf("Erro ao verificar se o usuário %s já existe", err)
		return false
	}
	return true
}

func (s *SQLUsuarioStore) InsereNovoUsuario(nomeDeUsuario, nome, senha string) error {
	q := `INSERT INTO usuarios (nome_de_usuario, nome, senha) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(q, nomeDeUsuario, nome, senha)
	if err != nil {
		s.logger.Errorf("Erro ao inserir os dados para novo segredos de usuario. %v", err)
		return err
	}
	return nil
}

func (s *SQLUsuarioStore) VerificaSegredosDeUsuario(nomeDeUsuario, senha string) (model.Usuario, error) {
	q := `SELECT id, nome_de_usuario, nome, senha, imagem FROM usuarios WHERE nome_de_usuario = $1`
	row := s.db.QueryRow(q, nomeDeUsuario)
	var u model.Usuario
	err := row.Scan(&u.ID, &u.NomeDeUsuario, &u.Nome, &u.Senha, &u.Imagem)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Usuario{}, ErrUsuarioNaoEncontrado
		}
		s.logger.Errorf("Erro ao buscar usuario por nome de usuario:%s. %v", nomeDeUsuario, err)
		return model.Usuario{}, err
	}
	if u.Senha != senha {
		return model.Usuario{}, ErrSenhaInvalida
	}
	return u, nil
}

func (s SQLUsuarioStore) AtualizaImagemDeUsuario(nomeDeUsuario, imagem string) error {
	s.logger.Debugf("Atualizando imagem de usuario:%s", nomeDeUsuario)
	q := `UPDATE usuarios SET imagem = $1 WHERE nome_de_usuario = $2`
	_, err := s.db.Exec(q, imagem, nomeDeUsuario)
	if err != nil {
		s.logger.Errorf("Erro ao atualizar imagem de usuario:%s. %v", nomeDeUsuario, err)
		return err
	}
	return nil
}

func (s *SQLUsuarioStore) AtualizaNome(nomeDeUsuario, nome string) error {
	q := `UPDATE usuarios SET nome = $1 WHERE nome_de_usuario = $2`
	_, err := s.db.Exec(q, nome, nomeDeUsuario)
	if err != nil {
		s.logger.Errorf("Erro ao atualizar nome de usuario:%s com o nome:%s. %v", nomeDeUsuario, nome, err)
		return err
	}
	return nil
}
