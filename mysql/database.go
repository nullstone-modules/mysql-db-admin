package mysql

import (
	"bytes"
	"fmt"
	"github.com/nullstone-io/go-rest-api"
	"log"
)

type Database struct {
	Name                string `json:"name"`
	DefaultCharacterSet string `json:"defaultCharacterSet"`
	DefaultCollation    string `json:"defaultCollation"`

	// Do not error if trying to create a database that already exists
	// Instead, read the existing and return
	UseExisting bool `json:"useExisting"`
}

var _ rest.DataAccess[string, Database] = &Databases{}

type Databases struct {
	DbOpener DbOpener
}

func (d *Databases) Create(obj Database) (*Database, error) {
	if obj.UseExisting {
		if existing, err := d.Read(obj.Name); err != nil {
			return nil, err
		} else if existing != nil {
			log.Printf("[Create] Database %q already exists, updating...\n", obj.Name)
			return d.Update(obj.Name, obj)
		}
	}

	db, err := d.DbOpener.OpenDatabase("")
	if err != nil {
		return nil, err
	}

	log.Printf("Creating database %q\n", obj.Name)
	if _, err := db.Exec(d.generateCreateSql(obj)); err != nil {
		return nil, fmt.Errorf("error creating database %q: %w", obj.Name, err)
	}
	return &obj, nil
}

func (d *Databases) Read(key string) (*Database, error) {
	db, err := d.DbOpener.OpenDatabase("")
	if err != nil {
		return nil, err
	}

	obj := &Database{Name: key}
	sq := `SELECT schema_name, default_character_set_name, default_collation_name FROM information_schema.schemata WHERE schema_name = ?;`
	row := db.QueryRow(sq, obj.Name)
	var databaseName, charSetName, collationName string
	if err := row.Scan(&databaseName, &charSetName, &collationName); err != nil {
		return nil, err
	}
	obj.DefaultCharacterSet = charSetName
	obj.DefaultCollation = collationName
	return obj, nil
}

func (d *Databases) Update(key string, obj Database) (*Database, error) {
	return d.Read(key)
}

func (d *Databases) Drop(key string) (bool, error) {
	return true, nil
}

func (d *Databases) generateCreateSql(obj Database) string {
	b := bytes.NewBufferString("CREATE DATABASE ")
	fmt.Fprint(b, QuoteIdentifier(obj.Name))

	if obj.DefaultCharacterSet != "" {
		fmt.Fprint(b, " CHARACTER SET ", QuoteIdentifier(obj.DefaultCharacterSet))
	}

	if obj.DefaultCollation != "" {
		fmt.Fprint(b, " COLLATE ", QuoteIdentifier(obj.DefaultCollation))
	}

	return b.String()
}
