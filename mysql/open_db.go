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

	fmt.Printf("OpenDatabase addr=%s, db=%s, user=%s, length(password)=%d\n", dsn.Addr, dsn.DBName, dsn.User, len(dsn.Passwd))
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
