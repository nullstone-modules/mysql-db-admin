package mysql

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/nullstone-io/go-rest-api"
	"log"
)

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	// Do not error if trying to create a role that already exists
	// Instead, read the existing, set the password, and return
	UseExisting bool `json:"useExisting"`
}

var _ rest.DataAccess[string, User] = &Users{}

type Users struct {
	DbOpener DbOpener
}

func (u *Users) Create(obj User) (*User, error) {
	if obj.UseExisting {
		if existing, err := u.Read(obj.Name); err != nil {
			return nil, err
		} else if existing != nil {
			log.Printf("[Create] User %q already exists, updating...\n", obj.Name)
			return u.Update(obj.Name, obj)
		}
	}

	db, err := u.DbOpener.OpenDatabase("")
	if err != nil {
		return nil, err
	}

	log.Printf("Creating user %q\n", obj.Name)
	if _, err := db.Exec(u.generateCreateSql(obj)); err != nil {
		return nil, fmt.Errorf("error creating user %q: %w", obj.Name, err)
	}
	return &obj, nil
}

func (u *Users) Read(key string) (*User, error) {
	db, err := u.DbOpener.OpenDatabase("")
	if err != nil {
		return nil, err
	}

	sq := `select User, Host from mysql.user where User = ?;`
	row := db.QueryRow(sq, key)
	var user, host string
	if err := row.Scan(&user, &host); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &User{Name: key}, nil
}

func (u *Users) Update(key string, obj User) (*User, error) {
	return u.Read(key)
}

func (u *Users) Drop(key string) (bool, error) {
	return true, nil
}

func (u *Users) generateCreateSql(obj User) string {
	b := bytes.NewBufferString("CREATE USER ")
	fmt.Fprint(b, QuoteLiteral(obj.Name), "@'%'")
	fmt.Fprintf(b, " IDENTIFIED BY '%s'", obj.Password)
	return b.String()
}
