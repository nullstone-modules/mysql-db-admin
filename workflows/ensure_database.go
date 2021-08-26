package workflows

import (
	"database/sql"
	"github.com/nullstone-modules/mysql-db-admin/mysql"
	"log"
)

func EnsureDatabase(db *sql.DB, newDatabase mysql.Database) error {
	log.Printf("ensuring database %q\n", newDatabase.Name)

	return newDatabase.Ensure(db)
}
