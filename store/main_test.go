package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbHost = "localhost"
	dbName = "db"
	dbUser = "denis"
	dbPass = "password"
	dbImage = "postgres:15"
)

var (
	db     *sql.DB
	logger echo.Logger
	pgC    *postgres.PostgresContainer
	dbPort nat.Port
)

func TestMain(m *testing.M) {
	var err error

	logger = echo.New().Logger
	logger.SetLevel(log.INFO)

	pgC, err = startPGContainer(context.Background())
	if err != nil {
		logger.Fatalf("Erro ao iniciar container: %v", err)
	}

	dbPort, err = pgC.MappedPort(context.Background(), "5432")
	if err != nil {
		logger.Fatalf("Erro ao obter porta mapeada: %v", err)
	}
	
	psqlStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort.Int(), dbUser, dbPass, dbName)

	db, err = sql.Open("postgres", psqlStr)
	if err != nil {
		logger.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	defer func() {
		db.Close()
		if err = testcontainers.TerminateContainer(pgC); err != nil {
			logger.Errorf("Erro ao terminar container %v", err)
		}
	}()

	code := m.Run()

	os.Exit(code)
}

func startPGContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	pgC, err := postgres.Run(ctx,
		dbImage,
		postgres.WithInitScripts(filepath.Join("../sql", "schema.sql")), postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		logger.Fatalf("Erro ao criar container: %v", err)
	}
	return pgC, nil
}
