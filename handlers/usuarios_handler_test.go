package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DenisJulio/marketplace-pit/services"
	"github.com/DenisJulio/marketplace-pit/store"
	"github.com/DenisJulio/marketplace-pit/testutils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegistraNovoUsuario(t *testing.T) {
	// setup
	ctx := context.Background()
	pg, port, _ := testutils.StartPGContainer(ctx, testutils.DefaultDbConfig, filepath.Join("../sql", "schema.sql"))
	db, _ := testutils.ConnectToDB(testutils.DefaultDbConfig, port.Int())
	defer pg.Terminate(ctx)
	defer db.Close()
	e := echo.New()
	usSto := store.NewSQLUsuarioStore(db, e.Logger)
	usSvc := services.NewUsuarioService(usSto)
	h := NovoUsuarioHandler(*usSvc, e.Logger)

	// given
	formData := url.Values{}
	formData.Add("nome", "Paulo Santos")
	formData.Add("senha", "paulo123")
	formData.Add("nomeDeUsuario", "paulo_santos")

	// when
	req := httptest.NewRequest(http.MethodPost, "/registrar", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := h.RegistraNovoUsuario(c)
	
	// then
	assert := assert.New(t)
	assert.NoError(err, "Erro ao registrar novo usuario")
	assert.True(h.usuSvc.VerificaUsuarioExistente("paulo_santos"), "Usuario deve existir na base de dados")
	assert.Equal(http.StatusCreated, rec.Code)
}
