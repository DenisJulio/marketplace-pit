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
