package store

import (
	"testing"

	"github.com/DenisJulio/marketplace-pit/model"
	"github.com/stretchr/testify/assert"
)

func TestCriaNovaOfertaParaAnuncio(t *testing.T) {
	os := NewSQLOfertaStore(db, logger)
	as := NewSQLAnuncioStore(db, logger)
	of := model.NewOferta(2, 1, 3)
	msg := model.NewMensagem(1, "Olá, gostaria de fazer uma oferta neste anuncio.")
	id, err := os.CriaNovaOfertaParaAnuncio(*of, *msg)
	if err != nil {
		t.Errorf("Erro ao criar oferta: %v", err)
	}
	assert := assert.New(t)
	assert.Greater(id, 0, "ID da oferta deve ser maior que zero")
	anuncios, _ := as.BuscaTodosAnuncios()
	assert.Equal(len(anuncios), 5, "Quantidade de anuncios diferente do esperado")
	msgs, _ := os.BuscaTodasAsMensagensDaOferta(id)
	assert.Equal(len(msgs), 1, "Quantidade de mensagens diferente do esperado")
}
