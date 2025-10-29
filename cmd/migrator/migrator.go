package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	storagePath, migrationsPath, migrationsTable string
	down                                         bool
)

func init() {
	flag.StringVar(&storagePath, "storage-path", "", "Path to a directory containing the migration files")
	flag.StringVar(&migrationsPath, "migrations-path", "", "Path to a directory containing the migration files")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "Name of the migrations table in database")
	flag.BoolVar(&down, "down", false, "Downgrade the database to the latest version")
}

func main() {
	flag.Parse()

	if migrationsPath == "" {
		panic("migrations path is required")
	}

	if storagePath == "" {
		panic("storage path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		"sqlite3://"+storagePath+"?x-migrations-table="+migrationsTable,
	)
	if err != nil {
		panic(err)
	}
	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		if err != nil {
			panic(err)
		}
	}(m)

	if down {
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no changes to apply")
				return
			}
			panic(err)
		}
	} else {
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no changes to apply")
				return
			}
			panic(err)
		}
	}

	fmt.Println("applied all migrations")
}
