package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/DenisJulio/marketplace-pit/model"
	"github.com/DenisJulio/marketplace-pit/store"
	"github.com/DenisJulio/marketplace-pit/utils"
)

var (
	ErrDadosParaRegistroInvalidos = errors.New("Dados fornecidos sao invalidos.")
	ErrUsuarioExistente           = errors.New("Nome de usuario ja esta cadastrado.")
)

type UsuarioService struct {
	s      store.UsuarioStore
	logger utils.Logger
}

type segredosDeUsuario struct {
	nomeDeUsuario string
	senha         string
}

func NewUsuarioService(s store.UsuarioStore, logger utils.Logger) *UsuarioService {
	return &UsuarioService{s: s, logger: logger}
}

func (us *UsuarioService) BuscaUsuarioPorId(id int) (model.Usuario, error) {
	return us.s.BuscaUsuarioPorId(id)
}

func (us *UsuarioService) BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario string) (model.Usuario, error) {
	return us.s.BuscaUsuarioPorNomeDeUsuario(nomeDeUsuario)
}

func (us *UsuarioService) VerificaUsuarioExistente(nomeDeUsuario string) bool {
	return us.s.VerificaUsuarioExistente(nomeDeUsuario)
}

func (us *UsuarioService) RegistraNovoUsuario(nome, nomeDeUsuario, senha, imagem string) error {
	if err := validaDadosParaRegistro(nome, nomeDeUsuario, senha); err != nil {
		us.logger.Debugf("%v", err)
		return ErrDadosParaRegistroInvalidos
	}
	if us.s.VerificaUsuarioExistente(nomeDeUsuario) {
		us.logger.Debugf("%v", ErrUsuarioExistente)
		return ErrUsuarioExistente
	}
	if err := us.s.InsereNovoUsuario(nomeDeUsuario, nome, senha, imagem); err != nil {
		return err
	}
	return nil
}

func (us *UsuarioService) VerificaSegredosDeUsuario(nomeDeUsuario, senha string) error {
	_, err := us.s.VerificaSegredosDeUsuario(nomeDeUsuario, senha)
	if err != nil {
		return err
	}
	return nil
}

func (us *UsuarioService) SalvalNovaImagemDeAvatar(nomeDeUsuario string, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		us.logger.Errorf("Erro ao abrir o arquivo de imagem: %v", err)
		return "", err
	}
	defer src.Close()

	const uploadDir = "public/static/images/avatars"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		us.logger.Errorf("Erro ao criar o diretorio de upload: %v", err)
		return "", err
	}
	extension := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
	filePath := filepath.Join(uploadDir, filename)
	us.logger.Debugf("Arquivo sera salvo em: %s", filePath)

	// Create the file on the filesystem
	dst, err := os.Create(filePath)
	if err != nil {
		us.logger.Errorf("Erro ao criar o arquivo de imagem: %v", err)
		return "", err
	}
	defer dst.Close()

	// Copy the uploaded file's content to the new file
	if _, err := io.Copy(dst, src); err != nil {
		us.logger.Errorf("Erro ao copiar o arquivo de imagem para o local: %v", err)
		return "", err
	}
	imgPath := fmt.Sprintf("/resources/images/avatars/%s", filename)
	if err := us.s.AtualizaImagemDeUsuario(nomeDeUsuario, imgPath); err != nil {
		return "", err
	}
	return imgPath, nil
}

func (us *UsuarioService) AtualizaNome(nomeDeUsuario, nome string) error {
	return us.s.AtualizaNome(nomeDeUsuario, nome)
}

func validaDadosParaRegistro(nome string, nomeDeUsuario string, senha string) error {
	if nome == "" || nomeDeUsuario == "" || senha == "" {
		return errors.New("Campos obrigatórios não preenchidos")
	}
	if len(senha) < 8 {
		return errors.New("A senha deve ter no mínimo 8 caracteres")
	}
	return nil
}
