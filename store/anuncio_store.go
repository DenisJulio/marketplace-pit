package store

import (
	"database/sql"

	"github.com/DenisJulio/marketplace-pit/model"
	"github.com/DenisJulio/marketplace-pit/utils"
)

type AnuncioStore interface {
	BuscaTodosAnuncios() ([]model.Anuncio, error)
}

type SQLAnuncioStore struct {
	db     *sql.DB
	logger utils.Logger
}

func (s *SQLAnuncioStore) BuscaTodosAnuncios() ([]model.Anuncio, error) {
	q := "SELECT id, nome, criado_em, anunciante_id, valor, descricao, imagem FROM anuncios"
	rows, err := s.db.Query(q)
	if err != nil {
		s.logger.Errorf("Erro ao buscar todos os anuncios. %v", err)
		return nil, err
	}
	defer rows.Close()

	var anuncios []model.Anuncio
	for rows.Next() {
		var a model.Anuncio
		err := rows.Scan(&a.ID, &a.Nome, &a.CriadoEm, &a.AnuncianteId, &a.Valor, &a.Descricao, &a.Imagem)
		if err != nil {
			s.logger.Errorf("Erro ao buscar anuncio. %v", err)
			return nil, err
		}
		anuncios = append(anuncios, a)
	}
	return anuncios, nil
}
