package root

import (
	"fmt"
	"github.com/spf13/viper"
	"database/sql"
	_ "github.com/lib/pq"

	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/golang-migrate/migrate"
)

func getConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.name"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
	)
}

func upMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	checkErr(err)

	m, err := migrate.NewWithDatabaseInstance(
		viper.GetString("db.migrations"), viper.GetString("db.name"), driver)
	checkErr(err)

	m.Up()
}

func DbConnect() *sql.DB {
	db, err := sql.Open("postgres", getConnectionString())
	checkErr(err)

	err = db.Ping()
	checkErr(err)

	upMigrations(db)
	return db
}
