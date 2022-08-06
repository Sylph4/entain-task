package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sylph4/entain-task/storage/postgres"
	"gopkg.in/gorp.v1"
)

func Connect() (*gorp.DbMap, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Warning: %s environment variable not set.\n", k)
		}
		return v
	}

	var (
		dbUser = mustGetenv("DB_USER")
		dbPwd  = mustGetenv("DB_PASS")
		dbHost = mustGetenv("DB_HOST")
		dbPort = mustGetenv("DB_PORT")
		dbName = mustGetenv("DB_NAME")
	)

	port, _ := strconv.Atoi(dbPort)
	dbURI := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, port, dbUser, dbPwd, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	configureConnectionPool(dbPool)

	db := &gorp.DbMap{Db: dbPool, Dialect: gorp.PostgresDialect{}}
	db.AddTableWithName(postgres.Transaction{}, "transaction").SetKeys(false, "id")
	db.AddTableWithName(postgres.User{}, "users").SetKeys(false, "id")

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	n, err := migrate.Exec(dbPool, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Println("Could not apply migrations", err)
	} else {
		log.Printf("Applied %d migrations!\n", n)
	}

	return db, nil
}

func configureConnectionPool(db *sql.DB) {
	db.SetMaxIdleConns(5)

	db.SetMaxOpenConns(7)

	db.SetConnMaxLifetime(1800 * time.Second)
}
