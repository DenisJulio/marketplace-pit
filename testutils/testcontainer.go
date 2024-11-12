package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbImage = "postgres:15"
)

type DbConfig struct {
	Host string
	Port nat.Port
	Name string
	User string
	Pass string
}

var DefaultDbConfig = DbConfig{
	Host: "localhost",
	Name: "db",
	User: "denis",
	Pass: "password",
}

// StartPGContainer starts a PostgreSQL test container with the specified parameters.
func StartPGContainer(ctx context.Context, dbConf DbConfig, initScriptPath string) (*postgres.PostgresContainer, nat.Port, error) {
	pgC, err := postgres.Run(ctx,
		dbImage,
		postgres.WithInitScripts(initScriptPath),
		postgres.WithDatabase(dbConf.Name),
		postgres.WithUsername(dbConf.User),
		postgres.WithPassword(dbConf.Pass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create container: %w", err)
	}
	dbPort, err := pgC.MappedPort(ctx, "5432")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get mapped port: %w", err)
	}
	return pgC, dbPort, nil
}

// ConnectToDB provides a connection string for the database.
func ConnectToDB(dbConf DbConfig, port int) (*sql.DB, error) {
	psqlStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbConf.Host, port, dbConf.User, dbConf.Pass, dbConf.Name)
	return sql.Open("postgres", psqlStr)
}
