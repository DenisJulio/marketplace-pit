package store

import (
	"database/sql"

	"github.com/DenisJulio/marketplace-pit/model"
	"github.com/DenisJulio/marketplace-pit/utils"
)

type AnuncioStore interface {
	BuscaTodosAnuncios() ([]model.Anuncio, error)
	BuscaAnuncioPorID(id int) (model.Anuncio, error)
	SalvaNovoAnuncio(anuncio model.Anuncio) error
	BuscaAnunciosPorNomeDeUsuario(nomeDeUsuario string) ([]model.Anuncio, error)
	RemoveAnuncio(id int) error
}

type SQLAnuncioStore struct {
	DB     *sql.DB
	Logger utils.Logger
}

func NewSQLAnuncioStore(db *sql.DB, logger utils.Logger) *SQLAnuncioStore {
	return &SQLAnuncioStore{db, logger}
}

func (s *SQLAnuncioStore) BuscaTodosAnuncios() ([]model.Anuncio, error) {
	q := "SELECT id, nome, criado_em, anunciante_id, valor, descricao, imagem FROM anuncios"
	rows, err := s.DB.Query(q)
	if err != nil {
		s.Logger.Errorf("Erro ao buscar todos os anuncios. %v", err)
		return []model.Anuncio{}, err
	}
	defer rows.Close()

	var anuncios []model.Anuncio
	for rows.Next() {
		var a model.Anuncio
		err := rows.Scan(&a.ID, &a.Nome, &a.CriadoEm, &a.AnuncianteId, &a.Valor, &a.Descricao, &a.Imagem)
		if err != nil {
			s.Logger.Errorf("Erro ao buscar anuncio. %v", err)
			continue
		}
		anuncios = append(anuncios, a)
	}
	return anuncios, nil
}

func (s *SQLAnuncioStore) BuscaAnuncioPorID(id int) (model.Anuncio, error) {
	s.Logger.Debugf("Buscando anuncio por id=%d", id)
	q := "SELECT id, nome, criado_em, anunciante_id, valor, descricao, imagem FROM anuncios WHERE id = $1"
	row := s.DB.QueryRow(q, id)

	var a model.Anuncio
	err := row.Scan(&a.ID, &a.Nome, &a.CriadoEm, &a.AnuncianteId, &a.Valor, &a.Descricao, &a.Imagem)
	if err != nil {
		s.Logger.Errorf("Erro ao buscar anuncio por id=%d. %v", id, err)
		return model.Anuncio{}, err
	}
	return a, nil
}

func (s *SQLAnuncioStore) SalvaNovoAnuncio(anuncio model.Anuncio) error {
	q := `
	INSERT INTO anuncios 
		(nome, criado_em, anunciante_id, valor, descricao, imagem)
	VALUES 
		($1, $2, $3, $4, $5, $6)`
	_, err := s.DB.Exec(q, anuncio.Nome, anuncio.CriadoEm, anuncio.AnuncianteId, anuncio.Valor, anuncio.Descricao, anuncio.Imagem)
	if err != nil {
		s.Logger.Errorf("Erro ao salvar novo anuncio. %v", err)
		return err
	}
	return nil
}

func (s *SQLAnuncioStore) BuscaAnunciosPorNomeDeUsuario(nomeDeUsuario string) ([]model.Anuncio, error) {
	s.Logger.Debugf("Buscando anuncios por nome de usuario=%s", nomeDeUsuario)
	q := `
		SELECT a.id, a.nome, a.criado_em, a.anunciante_id, a.valor, a.descricao, a.imagem 
		FROM anuncios a
		JOIN usuarios u ON a.anunciante_id = u.id
		WHERE u.nome_de_usuario = $1`
	rows, err := s.DB.Query(q, nomeDeUsuario)
	if err != nil {
		s.Logger.Errorf("Erro ao buscar anuncios por nome de usuario=%s. %v", nomeDeUsuario, err)
		return []model.Anuncio{}, err
	}
	defer rows.Close()

	var anuncios []model.Anuncio
	for rows.Next() {
		var a model.Anuncio
		err := rows.Scan(&a.ID, &a.Nome, &a.CriadoEm, &a.AnuncianteId, &a.Valor, &a.Descricao, &a.Imagem)
		if err != nil {
			s.Logger.Errorf("Erro ao buscar anuncio. %v", err)
			continue
		}
		anuncios = append(anuncios, a)
	}
	return anuncios, nil
}

func (s *SQLAnuncioStore) RemoveAnuncio(id int) error {
	q := "DELETE FROM anuncios WHERE id = $1"
	_, err := s.DB.Exec(q, id)
	if err != nil {
		s.Logger.Errorf("Erro ao remover anuncio id=%d. %v", id, err)
	}
	return err
}
