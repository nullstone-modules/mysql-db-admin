package mysql

import (
	"fmt"
	"github.com/nullstone-io/go-rest-api"
	"log"
)

// DbPrivilege grants access to User on Database
//
//	ALL PRIVILEGES
type DbPrivilege struct {
	Username string `json:"username"`
	Database string `json:"database"`
}

func (p DbPrivilege) Key() DbPrivilegeKey {
	return DbPrivilegeKey{
		Username: p.Username,
		Database: p.Database,
	}
}

type DbPrivilegeKey struct {
	Username string
	Database string
}

var _ rest.DataAccess[DbPrivilegeKey, DbPrivilege] = &DbPrivileges{}

type DbPrivileges struct {
	DbOpener DbOpener
}

func (r *DbPrivileges) Create(obj DbPrivilege) (*DbPrivilege, error) {
	return r.Update(obj.Key(), obj)
}

func (r *DbPrivileges) Read(key DbPrivilegeKey) (*DbPrivilege, error) {
	// TODO: Introspect
	obj := DbPrivilege{
		Username: key.Username,
		Database: key.Database,
	}
	return &obj, nil
}

func (r *DbPrivileges) Update(key DbPrivilegeKey, obj DbPrivilege) (*DbPrivilege, error) {
	db, err := r.DbOpener.OpenDatabase(obj.Database)
	if err != nil {
		return nil, err
	}

	log.Printf("Granting user %q db access to %q\n", obj.Username, obj.Database)
	sq := fmt.Sprintf(`GRANT ALL PRIVILEGES ON %s.* TO %s@'%%';`, QuoteIdentifier(obj.Database), QuoteLiteral(obj.Username))
	if _, err := db.Exec(sq); err != nil {
		return nil, fmt.Errorf("error granting privileges: %w", err)
	}
	return &obj, err
}

func (r *DbPrivileges) Drop(key DbPrivilegeKey) (bool, error) {
	return true, nil
}
