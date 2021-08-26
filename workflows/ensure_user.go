package workflows

import (
	"database/sql"
	"github.com/nullstone-modules/mysql-db-admin/mysql"
	"log"
)

func EnsureUser(db *sql.DB, newUser mysql.User) error {
	log.Printf("ensuring user %q\n", newUser.Name)

	return newUser.Ensure(db)
}
