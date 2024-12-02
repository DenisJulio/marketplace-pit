package model

type Usuario struct {
	ID                                int
	Nome, NomeDeUsuario, Senha, Email string
	Imagem                            *string
}

func NovoUsuario(nome, nomeDeUsuario string) *Usuario {
	return &Usuario{
		Nome:          nome,
		NomeDeUsuario: nomeDeUsuario,
	}
}
