package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (app *application) migrate() {
	m, err := migrate.New("file://db/migrations", app.config.dbUrl)

	if err != nil {
		app.errorLog.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		app.errorLog.Fatal(err)
	}
}
