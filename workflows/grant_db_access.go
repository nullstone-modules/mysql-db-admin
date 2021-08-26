package workflows

import (
	"database/sql"
	"fmt"
	"github.com/nullstone-modules/mysql-db-admin/mysql"
	"log"
)

func GrantDbAccess(db *sql.DB, user mysql.User, database mysql.Database) error {
	log.Printf("Granting user %q db access to %q\n", user.Name, database.Name)
	return grantAllPrivileges(db, user, database)
}

func grantAllPrivileges(db *sql.DB, user mysql.User, database mysql.Database) error {
	sq := fmt.Sprintf(`GRANT ALL PRIVILEGES ON %s.* TO %s@'%%';`, mysql.QuoteIdentifier(database.Name), mysql.QuoteLiteral(user.Name))
	if _, err := db.Exec(sq); err != nil {
		return fmt.Errorf("error granting privileges: %w", err)
	}
	return nil
}
