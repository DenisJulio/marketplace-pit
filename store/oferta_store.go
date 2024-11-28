package store

import (
	"database/sql"

	"github.com/DenisJulio/marketplace-pit/model"
	"github.com/DenisJulio/marketplace-pit/utils"
)

type OfertaStore interface {
	CriaNovaOfertaParaAnuncio(oferta model.Oferta, msg model.Mensagem) (int, error)
	BuscaTodasAsMensagensDaOferta(ofertaID int) ([]model.Mensagem, error)
	BuscaTodasAsOfertasExpandidasDoUsuario(usuarioId int) ([]model.OfertaExpandida, error)
}

type SqlOfertaStore struct {
	db     *sql.DB
	logger utils.Logger
}

func NewSQLOfertaStore(db *sql.DB, logger utils.Logger) *SqlOfertaStore {
	return &SqlOfertaStore{db: db, logger: logger}
}

// CriaNovaOfertaParaAnuncio cria uma nova oferta para um anúncio específico.
//
// Parâmetros:
//   - oferta: objeto do tipo model.Oferta contendo os detalhes da oferta.
//   - msg: objeto do tipo model.Mensagem contendo a mensagem associada à oferta.
//
// Retorna:
//   - int: o ID da nova oferta criada.
//   - error: um erro, caso ocorra algum problema durante a criação da oferta.
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

func (s *SqlOfertaStore) BuscaTodasAsOfertasExpandidasDoUsuario(usuarioId int) ([]model.OfertaExpandida, error) {
	query := `
	SELECT 
		o.id, o.criado_em,
		a.id, a.nome_de_usuario, a.nome, a.senha, a.imagem,
		of.id, of.nome_de_usuario, of.nome, of.senha, of.imagem,
		an.id, an.nome, an.criado_em, an.anunciante_id, an.valor, an.descricao, an.imagem,
		CASE WHEN o.ofertante_id = $1 THEN true ELSE false END AS e_ofertante
	FROM ofertas o
	INNER JOIN usuarios a ON o.anunciante_id = a.id
	INNER JOIN usuarios of ON o.ofertante_id = of.id
	INNER JOIN anuncios an ON o.anuncio_id = an.id
	WHERE o.anunciante_id = $1 OR o.ofertante_id = $1;
	`

	rows, err := s.db.Query(query, usuarioId)
	if err != nil {
		s.logger.Errorf("Erro ao buscar ofertas do usuário: %v", err)
		return nil, err
	}
	defer rows.Close()

	var ofertas []model.OfertaExpandida

	for rows.Next() {
		var oferta model.OfertaExpandida
		var anuncianteImagem, ofertanteImagem, anuncioImagem, anuncioDescricao *string
		var eOfertante bool

		err := rows.Scan(
			&oferta.ID, &oferta.CriadoEm,
			&oferta.Anunciante.ID, &oferta.Anunciante.NomeDeUsuario, &oferta.Anunciante.Nome, &oferta.Anunciante.Senha, &anuncianteImagem,
			&oferta.Ofertante.ID, &oferta.Ofertante.NomeDeUsuario, &oferta.Ofertante.Nome, &oferta.Ofertante.Senha, &ofertanteImagem,
			&oferta.Anuncio.ID, &oferta.Anuncio.Nome, &oferta.Anuncio.CriadoEm, &oferta.Anuncio.AnuncianteId, &oferta.Anuncio.Valor, &anuncioDescricao, &anuncioImagem,
			&eOfertante,
		)
		if err != nil {
			s.logger.Errorf("Erro ao escanear as ofertas do usuário: %v", err)
			return nil, err
		}

		// Assign nullable fields if they are not nil
		oferta.Anunciante.Imagem = anuncianteImagem
		oferta.Ofertante.Imagem = ofertanteImagem
		oferta.Anuncio.Imagem = anuncioImagem
		oferta.Anuncio.Descricao = anuncioDescricao

		// Add the e_ofertante field
		oferta.EOfertante = eOfertante

		// Append the result to the list
		ofertas = append(ofertas, oferta)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		s.logger.Errorf("Erro ao iterar sobre as ofertas: %v", err)
		return nil, err
	}

	return ofertas, nil
}
