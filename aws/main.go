package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nullstone-modules/mysql-db-admin/aws/secrets"
	"github.com/nullstone-modules/mysql-db-admin/mysql"
	"log"
	"os"
	"time"
)

const (
	dbConnUrlSecretIdEnvVar = "DB_CONN_URL_SECRET_ID"

	eventTypeCreateDatabase = "create-database"
	eventTypeCreateUser     = "create-user"
	eventTypeCreateDbAccess = "create-db-access"
)

type AdminEvent struct {
	Type     string            `json:"type"`
	Metadata map[string]string `json:"metadata"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	dbConnurl, err := secrets.GetString(ctx, os.Getenv(dbConnUrlSecretIdEnvVar))
	if err != nil {
		log.Println(err.Error())
	}

	store := mysql.NewStore(dbConnurl)
	defer store.Close()

	lambda.Start(HandleRequest(store))
}

func HandleRequest(store *mysql.Store) func(ctx context.Context, event AdminEvent) error {
	return func(ctx context.Context, event AdminEvent) error {
		switch event.Type {
		case eventTypeCreateDatabase:
			databaseName, _ := event.Metadata["databaseName"]
			if databaseName == "" {
				return fmt.Errorf("cannot create database: databaseName is required")
			}
			_, err := store.Databases.Create(mysql.Database{Name: databaseName, UseExisting: true})
			return err
		case eventTypeCreateUser:
			username, _ := event.Metadata["username"]
			password, _ := event.Metadata["password"]
			if username == "" {
				return fmt.Errorf("cannot create user: username is required")
			}
			if password == "" {
				return fmt.Errorf("cannot create user: password is required")
			}
			_, err := store.Users.Create(mysql.User{Name: username, Password: password, UseExisting: true})
			return err
		case eventTypeCreateDbAccess:
			username, _ := event.Metadata["username"]
			if username == "" {
				return fmt.Errorf("cannot grant user access to db: username is required")
			}
			database, _ := event.Metadata["databaseName"]
			if database == "" {
				return fmt.Errorf("cannot grant user access to db: database name is required")
			}
			_, err := store.DbPrivileges.Create(mysql.DbPrivilege{Username: username, Database: database})
			return err
		default:
			return fmt.Errorf("unknown event %q", event.Type)
		}
	}
}
