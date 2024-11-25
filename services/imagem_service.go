package services

import (
	"mime/multipart"

	"github.com/DenisJulio/marketplace-pit/store"
	"github.com/DenisJulio/marketplace-pit/utils"
)

type ImagemService struct {
	s      store.ImagemStore
	logger utils.Logger
}

func NovoImagemService(s store.ImagemStore, logger utils.Logger) *ImagemService {
	return &ImagemService{s: s, logger: logger}
}

func (is *ImagemService) SalvalNovaImagem(tipoDeImagem store.TipoDeImagem, file *multipart.FileHeader) (string, error) {
	imgPath, err := is.s.SalvaImagem(tipoDeImagem, file)
	if err != nil {
		is.logger.Errorf("Falha ao salvar o arquivo de imagem: %v", tipoDeImagem, err)
		return "", err
	}
	is.logger.Debugf("Arquivo de imagem salvo com sucesso e disponivel em: %s", imgPath)
	return imgPath, nil
}
