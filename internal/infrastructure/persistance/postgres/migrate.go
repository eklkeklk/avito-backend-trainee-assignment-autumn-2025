package postgres

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"os"
)

func Migrate(db *sql.DB, migDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if _, err := os.Stat(migDir); os.IsNotExist(err) {
		return fmt.Errorf("directory doesnt exist: %s", migDir)
	}
	if err := goose.Up(db, migDir); err != nil {
		return err
	}

	return nil
}
