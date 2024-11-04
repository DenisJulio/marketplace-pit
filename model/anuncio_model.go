package model

import "time"

type Anuncio struct {
	ID int
	Nome string
    CriadoEm time.Time
    AnuncianteId int
    Valor float64
    Descricao *string
    Imagem *string
}