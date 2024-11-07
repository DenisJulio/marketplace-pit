package model

import "time"

type Oferta struct {
	ID int
	CriadoEm time.Time
	AnuncianteID int
	OfertanteID int
	AnuncioID int
}