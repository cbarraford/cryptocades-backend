package context

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

// Migrates the database to the latest version
func MigrateDB(url string, migrationDir string) error {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationDir, "postgres", driver)
	if err != nil {
		return err
	}

	version, dirty, err := m.Version()
	if version > 0 && err != nil {
		return err
	}
	if !dirty {
		err = m.Up()
		if err != nil {
			return err
		}
	}
	return nil
}
