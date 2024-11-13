package model

type SegredosUsuario struct {
	NomeDeUsuario string
	Senha         string
}

func NovoSegredosUsuario() *SegredosUsuario {
	return &SegredosUsuario{}
}
