package store

import (
	"database/sql"

	"github.com/DenisJulio/marketplace-pit/model"
	"github.com/DenisJulio/marketplace-pit/utils"
)

type UsuarioStore interface {
	BuscaUsuarioPorId(ID int) (model.Usuario, error)
}

type SQLUsuarioStore struct {
	DB     *sql.DB
	Logger utils.Logger
}

func (s *SQLUsuarioStore) BuscaUsuarioPorID(ID int) (model.Usuario, error) {
	row := s.DB.QueryRow("SELECT id, nome_de_usuario, nome, imagem FROM usuarios WHERE id = $1", ID)
	var u model.Usuario
	err := row.Scan(&u.ID, &u.NomeDeUsuario, &u.Nome, &u.Imagem)
	if err != nil {
		s.Logger.Errorf("Erro ao buscar usuario por id: %v", err)
		return model.Usuario{}, err
	}
	return u, nil
}

func (s *SQLUsuarioStore) BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario string) (model.Usuario, error) {
	row := s.DB.QueryRow("SELECT id, nome_de_usuario, nome, imagem FROM usuarios WHERE nome_de_usuario = $1", nomeDeUsuario)
	var u model.Usuario
	err := row.Scan(&u.ID, &u.NomeDeUsuario, &u.Nome, &u.Imagem)
	if err != nil {
		s.Logger.Errorf("Erro ao buscar usuario por nome de usuario: %v", err)
		return model.Usuario{}, err
	}
	return u, nil
}