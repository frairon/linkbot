package schema

import (
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var schemas embed.FS

// ApplySchema applies the schema to passed database and drops all values if required
func ApplySchema(dropAll bool, driver database.Driver) error {
	schemaSource, err := iofs.New(schemas, ".")
	if err != nil {
		return fmt.Errorf("error loading schema files: %v", err)
	}

	log.Printf("creating migrator")
	dbMigrate, err := migrate.NewWithInstance("test", schemaSource, "", driver)
	if err != nil {
		return err
	}

	// migrating everything
	if dropAll {
		log.Printf("Dropping all values in database as requested by flag")
		err := dbMigrate.Drop()
		if err != nil {
			return fmt.Errorf("error dropping: %v", err)
		}
	}
	log.Printf("applying migration")

	err = dbMigrate.Up()
	if err != migrate.ErrNoChange {
		return err
	} else {
		log.Printf("db is up to date")
	}
	return nil
}
