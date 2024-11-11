package store

import (
	"database/sql"
	"time"

	"github.com/DenisJulio/marketplace-pit/model"
	"github.com/DenisJulio/marketplace-pit/utils"
)

type OfertaStore interface {
	CriaNovaOfertaParaAnuncio(ofertanteID int, anuncianteID int, anuncioID int, criadoEm time.Time, msg string) (int, error)
	BuscaTodasAsMensagensDaOferta(ofertaID int) ([]model.Mensagem, error)
}

type SqlOfertaStore struct {
	db     *sql.DB
	logger utils.Logger
}

func NewSQLOfertaStore(db *sql.DB, logger utils.Logger) *SqlOfertaStore {
	return &SqlOfertaStore{db: db, logger: logger}
}

func (s *SqlOfertaStore) CriaNovaOfertaParaAnuncio(oferta model.Oferta, msg model.Mensagem) (int, error) {
	// insere a nova oferta no banco de dados
    query := `INSERT INTO ofertas (criado_em, ofertante_id, anunciante_id, anuncio_id) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := s.db.QueryRow(query, oferta.CriadoEm, oferta.OfertanteID, oferta.AnuncianteID, oferta.AnuncioID).Scan(&id)
	if err != nil {
		s.logger.Errorf("Erro ao criar oferta: %v", err)
		return 0, nil
	}
	// insere a mensagem no banco de dados
	query = `INSERT INTO mensagens (criado_em, conteudo, oferta_id, remetente_id) VALUES ($1, $2, $3, $4)`
	_, err = s.db.Exec(query, msg.CriadoEm, msg.Conteudo, id, msg.RemetenteID)
	if err != nil {
		s.logger.Errorf("Erro ao criar oferta: %v", err)
		return 0, nil
	}
	return int(id), nil
}

func (s *SqlOfertaStore) BuscaTodasAsMensagensDaOferta(ofertaID int) ([]model.Mensagem, error) {
	query := `SELECT id, criado_em, conteudo, oferta_id, remetente_id FROM mensagens WHERE oferta_id = $1`
	rows, err := s.db.Query(query, ofertaID)
	if err != nil {
		s.logger.Errorf("Erro ao buscar mensagens da oferta: %v", err)
		return nil, err
	}
	defer rows.Close()
	var msgs []model.Mensagem
	for rows.Next() {
		var m model.Mensagem
		err := rows.Scan(&m.ID, &m.CriadoEm, &m.Conteudo, &m.OfertaID, &m.RemetenteID)
		if err != nil {
			s.logger.Errorf("Erro ao buscar mensagens da oferta: %v", err)
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}
