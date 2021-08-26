package mysql

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
)

type Database struct {
	Name                string
	DefaultCharacterSet string
	DefaultCollation    string
}

func (d Database) Create(db *sql.DB) error {
	sq := d.generateCreateSql()
	log.Printf("Creating database %q", d.Name)
	if _, err := db.Exec(sq); err != nil {
		return fmt.Errorf("error creating database %q: %w", d.Name, err)
	}
	return nil
}

func (d Database) generateCreateSql() string {
	b := bytes.NewBufferString("CREATE DATABASE ")
	fmt.Fprint(b, QuoteIdentifier(d.Name))

	if d.DefaultCharacterSet != "" {
		fmt.Fprint(b, " CHARACTER SET ", QuoteIdentifier(d.DefaultCharacterSet))
	}

	if d.DefaultCollation != "" {
		fmt.Fprint(b, " COLLATE ", QuoteIdentifier(d.DefaultCollation))
	}

	return b.String()
}

func (d Database) Ensure(db *sql.DB) error {
	if exists, err := d.Exists(db); exists {
		log.Printf("database %q already exists\n", d.Name)
		return nil
	} else if err != nil {
		return fmt.Errorf("error checking for database %q: %w", d.Name, err)
	}
	if err := d.Create(db); err != nil {
		return fmt.Errorf("error creating database %q: %w", d.Name, err)
	}
	return nil
}

func (d Database) Exists(db *sql.DB) (bool, error) {
	check := Database{Name: d.Name}
	if err := check.Read(db); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (d *Database) Read(db *sql.DB) error {
	row := db.QueryRow(fmt.Sprintf(`SHOW CREATE DATABASE %s`, QuoteIdentifier(d.Name)))
	var databaseName, createSql string
	if err := row.Scan(&databaseName, &createSql); err != nil {
		return err
	}
	return nil
}

func (d Database) Update(db *sql.DB) error {
	return nil
}

func (d Database) Drop(db *sql.DB) error {
	return nil
}
