package store

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuscaUsuarioPorID(t *testing.T) {
	tests := []struct {
		usuarioID    int
		nomeEsperado string
	}{
		{usuarioID: 1, nomeEsperado: "Pedro Santos"},
		{usuarioID: 2, nomeEsperado: "Maria Antonia"},
		{usuarioID: 3, nomeEsperado: "João Marcos"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Buscando por ID de usuario %d", tt.usuarioID), func(t *testing.T) {
			us := &SQLUsuarioStore{DB: db, Logger: logger}
			u, err := us.BuscaUsuarioPorID(tt.usuarioID)
			assert := assert.New(t)
			assert.NoError(err, "Erro ao conectar ao banco de dados")
			assert.Equal(tt.nomeEsperado, u.Nome, "Nome do usuario incorreto")
		})
	}
}

func TestBuscaUsuarioPorNomeDeUsuario(t *testing.T) {
	tests := []struct {
		nomeDeUsuario string
		nomeEsperado  string
	}{
		{nomeDeUsuario: "pedr0", nomeEsperado: "Pedro Santos"},
		{nomeDeUsuario: "mari4", nomeEsperado: "Maria Antonia"},
		{nomeDeUsuario: "joa0", nomeEsperado: "João Marcos"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Buscando por nome de usuario %s", tt.nomeDeUsuario), func(t *testing.T) {
			us := &SQLUsuarioStore{DB: db, Logger: logger}
			u, err := us.BuscaUsuarioPorNomeDeUsuario(tt.nomeDeUsuario)
			assert := assert.New(t)
			assert.NoError(err, "Erro ao conectar ao banco de dados")
			assert.Equal(tt.nomeEsperado, u.Nome, "Nome do usuario incorreto")
		})
	}
}