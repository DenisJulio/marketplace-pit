package db

import (
	"database/sql"
	"fmt"

	"github.com/DenisJulio/marketplace-pit/utils"
	_ "github.com/lib/pq"
)

const (
	dbHost  = "localhost"
	dbName  = "db"
	dbUser  = "denis"
	dbPass  = "password"
	dbImage = "postgres:15"
	dbPort  = 5432
)

func NewDB(logger utils.Logger) *sql.DB {
	psqlStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := sql.Open("postgres", psqlStr)
	if err != nil {
		logger.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	return db
}
