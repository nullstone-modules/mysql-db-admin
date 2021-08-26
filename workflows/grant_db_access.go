package workflows

import (
	"database/sql"
	"github.com/nullstone-modules/mysql-db-admin/mysql"
	"log"
)

func GrantDbAccess(db *sql.DB, appDb *sql.DB, user mysql.User, database mysql.Database) error {
	log.Printf("Granting user %q db access to %q\n", user.Name, database.Name)
	panic("not implemented")
}
