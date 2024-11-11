package model

import "time"

type Mensagem struct {
	ID          int
	CriadoEm    time.Time
	Conteudo    string
	OfertaID    int
	RemetenteID int
}

func NewMensagem(remetenteID int, conteudo string) *Mensagem {
	return &Mensagem{
		CriadoEm:    time.Now(),
		Conteudo:    conteudo,
		RemetenteID: remetenteID,
	}
}
