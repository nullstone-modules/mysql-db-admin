package mysql

import (
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	Name     string
	Password string
}

func (u User) Ensure(db *sql.DB) error {
	if exists, err := u.Exists(db); exists {
		log.Printf("User %q already exists\n", u.Name)
		return nil
	} else if err != nil {
		return fmt.Errorf("error checking for user %q: %w", u.Name, err)
	}
	if err := u.Create(db); err != nil {
		return fmt.Errorf("error creating user %q: %w", u.Name, err)
	}
	return nil
}

func (u User) Create(db *sql.DB) error {
	fmt.Printf("Creating user %q\n", u.Name)
	sq := u.generateCreateSql()
	if _, err := db.Exec(sq); err != nil {
		return fmt.Errorf("error creating user %q: %w", u.Name, err)
	}
	return nil
}

func (u User) generateCreateSql() string {
	panic("not implemented")
}

func (u User) Exists(db *sql.DB) (bool, error) {
	check := User{Name: u.Name}
	if err := check.Read(db); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (u User) Read(db *sql.DB) error {
	panic("not implemented")
}
