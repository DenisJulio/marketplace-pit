package store

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/DenisJulio/marketplace-pit/utils"
)

type TipoDeImagem string

const (
	ImagemDeAvatar  TipoDeImagem = "avatars"
	ImagemDeAnuncio TipoDeImagem = "anuncios"
	AvatarPadrao    string       = "/resources/images/avatars/avatar-padrao.png"
)

type ImagemStore interface {
	SalvaImagem(tipoDeImagem TipoDeImagem, arquivo *multipart.FileHeader) (string, error)
	RemoveImagem(tipoDeImagem TipoDeImagem, nomeDoArquivo string) error
}

type FileSystemImagemStore struct {
	BaseUploadDir  string
	BasePublicPath string
	logger         utils.Logger
}

func NewFileSystemImageStore(baseUploadDir, basePublicPath string, logger utils.Logger) *FileSystemImagemStore {
	return &FileSystemImagemStore{
		BaseUploadDir:  baseUploadDir,
		BasePublicPath: basePublicPath,
		logger:         logger,
	}
}

func (fs *FileSystemImagemStore) SalvaImagem(tipoDeImagem TipoDeImagem, arquivo *multipart.FileHeader) (string, error) {
	uploadDir := filepath.Join(fs.BaseUploadDir, string(tipoDeImagem))
	publicPath := filepath.Join(fs.BasePublicPath, string(tipoDeImagem))

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		fs.logger.Errorf("Erro ao criar o diretorio de upload: %v", err)
		return "", err
	}

	// Cria um nome unico para o arquivo
	extension := filepath.Ext(arquivo.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
	filePath := filepath.Join(uploadDir, filename)

	src, err := arquivo.Open()
	if err != nil {
		fs.logger.Errorf("Erro ao abrir o arquivo de imagem: %v", err)
		return "", err
	}
	defer src.Close()

	// Cria o arquivo no sistema de arquivos
	dst, err := os.Create(filePath)
	if err != nil {
		fs.logger.Errorf("Erro ao criar o arquivo de imagem: %v", err)
		return "", err
	}
	defer dst.Close()

	// Copia o conteudo do arquivo para o novo arquivo
	if _, err := io.Copy(dst, src); err != nil {
		fs.logger.Errorf("Erro ao copiar o arquivo de imagem para o local: %v", err)
		return "", err
	}

	return filepath.Join(publicPath, filename), nil
}

func (fs *FileSystemImagemStore) RemoveImagem(tipoDeImagem TipoDeImagem, nomeDoArquivo string) error {
	relativePath := strings.TrimPrefix(nomeDoArquivo, fs.BasePublicPath)
	fileName := filepath.Base(relativePath)
	filePath := filepath.Join(fs.BaseUploadDir, string(tipoDeImagem), fileName)

	if err := os.Remove(filePath); err != nil {
		fs.logger.Errorf("Erro ao remover o arquivo de imagem: %v", err)
		return err
	}
	return nil
}
