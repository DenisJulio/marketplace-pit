package store

import (
	"database/sql"
	"time"

	"github.com/DenisJulio/marketplace-pit/utils"
)

type SessaoStore struct {
	db     *sql.DB
	logger utils.Logger
}

func NovaSessaoStore(db *sql.DB, logger utils.Logger) *SessaoStore {
	return &SessaoStore{db: db, logger: logger}
}

func (s *SessaoStore) SalvaSessao(sessaoID, nomeDeUsuario string, expiraEm time.Time) error {
	_, err := s.db.Exec(
		`INSERT INTO sessoes (sessao_id, nome_de_usuario, expira_em) 
		 VALUES ($1, $2, $3) 
		 ON CONFLICT (sessao_id) DO UPDATE 
		 SET nome_de_usuario = $2, expira_em = $3`,
		sessaoID, nomeDeUsuario, expiraEm,
	)
	return err
}

func (s *SessaoStore) BuscaSessao(sessaoID string) (string, error) {
	var nomeDeUsuario string
	var expiraEm time.Time
	err := s.db.QueryRow(
		`SELECT nome_de_usuario, expira_em FROM sessoes WHERE sessao_id = $1`,
		sessaoID,
	).Scan(&nomeDeUsuario, &expiraEm)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Debugf("Sessao nao encontrada para")
			return "", nil // No sessao found
		}
		s.logger.Debugf("Sessao nao encontrada para")
		return "", err
	}

	if time.Now().After(expiraEm) {
		s.logger.Debugf("Sessao expirada")
		return "", nil // sessao expired
	}

	return nomeDeUsuario, nil
}

func (s *SessaoStore) RemoveSessao(sessaoID string) error {
	_, err := s.db.Exec(`DELETE FROM sessoes WHERE sessao_id = $1`, sessaoID)
	return err
}
