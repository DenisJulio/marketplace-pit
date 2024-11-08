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
		deveFalhar   bool
	}{
		{usuarioID: 1, nomeEsperado: "Pedro Santos", deveFalhar: false},
		{usuarioID: 2, nomeEsperado: "Maria Antonia", deveFalhar: false},
		{usuarioID: 3, nomeEsperado: "João Marcos", deveFalhar: false},
		{usuarioID: 4, nomeEsperado: "", deveFalhar: true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Buscando por ID de usuario %d", tt.usuarioID), func(t *testing.T) {
			us := &SQLUsuarioStore{db: db, logger: logger}
			u, err := us.BuscaUsuarioPorId(tt.usuarioID)
			assert := assert.New(t)
			if tt.deveFalhar {
				assert.ErrorIs(err, ErrUsuarioNaoEncontrado, "Erro esperado ao buscar usuario nao existente")
			} else {
				assert.NoError(err, "Erro ao conectar ao banco de dados")
			}
			assert.Equal(tt.nomeEsperado, u.Nome, "Resultado diferente do esperado")
		})
	}
}

func TestBuscaUsuarioPorNomeDeUsuario(t *testing.T) {
	tests := []struct {
		nomeDeUsuario string
		nomeEsperado  string
		deveFalhar    bool
	}{
		{nomeDeUsuario: "pedr0", nomeEsperado: "Pedro Santos", deveFalhar: false},
		{nomeDeUsuario: "mari4", nomeEsperado: "Maria Antonia", deveFalhar: false},
		{nomeDeUsuario: "joa0", nomeEsperado: "João Marcos", deveFalhar: false},
		{nomeDeUsuario: "mario", nomeEsperado: "", deveFalhar: true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Buscando por nome de usuario %s", tt.nomeDeUsuario), func(t *testing.T) {
			us := &SQLUsuarioStore{db: db, logger: logger}
			u, err := us.BuscaUsuarioPorNomeDeUsuario(tt.nomeDeUsuario)
			assert := assert.New(t)
			if tt.deveFalhar {
				assert.ErrorIs(err, ErrUsuarioNaoEncontrado, "Erro esperado ao buscar usuario nao existente")
			} else {
				assert.NoError(err, "Erro ao conectar ao banco de dados")
			}
			assert.Equal(tt.nomeEsperado, u.Nome, "Resultado diferente do esperado")
		})
	}
}
