package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

var (
	dbOpenConnTimeout = 3 * time.Second
)

func OpenDatabase(connUrl string, databaseName string) (*sql.DB, error) {
	dsn, err := DsnFromUrl(connUrl)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	dsn.DBName = databaseName

	db, err := sql.Open("mysql", dsn.FormatDSN())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbOpenConnTimeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error establishing connection to mysql: %w", err)
	}
	log.Println("Mysql connection established")
	return db, nil
}
