package model

import "time"

type Oferta struct {
	ID           int
	CriadoEm     time.Time
	AnuncianteID int
	OfertanteID  int
	AnuncioID    int
}

func NewOferta(anuncianteID int, ofertanteID int, anuncioID int) *Oferta {
	return &Oferta{
		CriadoEm:     time.Now(),
		AnuncianteID: anuncianteID,
		OfertanteID:  ofertanteID,
		AnuncioID:    anuncioID,
	}
}

type OfertaExpandida struct {
	ID         int
	CriadoEm   time.Time
	EOfertante bool
	Anunciante Usuario
	Ofertante  Usuario
	Anuncio    Anuncio
}
