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
}

var ErrUsuarioNaoEncontrado = errors.New("usuario nao encontrado")

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
