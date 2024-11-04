package store

import (
	"testing"

	"github.com/DenisJulio/marketplace-pit/model"
	"github.com/stretchr/testify/assert"
)

func TestBuscaTodosAnuncios(t *testing.T) {
	as := &SQLAnuncioStore{db: db, logger: logger}
	anuncios, err := as.BuscaTodosAnuncios()
	assert := assert.New(t)
	assert.NoError(err, "Erro ao buscar todos os anuncios")
	assert.NotEmpty(anuncios, "Nenhum anuncio encontrado")
	assert.True(len(anuncios) == 5, "Quantidade de anuncios diferente do esperado")
	var a model.Anuncio	
	for _, anuncio := range anuncios {
		if anuncio.Nome == "Carro" {
			a = anuncio
		}
	}
	assert.Equal(float64(15750), a.Valor, "Valor do anuncio diferente do esperado")
}
