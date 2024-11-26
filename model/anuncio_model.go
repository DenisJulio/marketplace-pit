package model

import "time"

type Anuncio struct {
	ID           int
	Nome         string
	CriadoEm     time.Time
	AnuncianteId int
	Valor        float64
	Descricao    *string
	Imagem       *string
}

func NovoAnuncio(nome string, anuncianteId int, valor float64, descricao, imagem *string) Anuncio {
	return Anuncio{
		Nome:         nome,
		CriadoEm:     time.Now(),
		AnuncianteId: anuncianteId,
		Valor:        valor,
		Descricao:    descricao,
		Imagem:       imagem,
	}
}
