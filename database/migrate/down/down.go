package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	sqlDB, err := sql.Open("postgres", os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = sqlDB.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	migrateWithDatabaseInstance, err := migrate.NewWithDatabaseInstance("file://database/migrate/sql", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	if err = migrateWithDatabaseInstance.Steps(-1); err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}
}
