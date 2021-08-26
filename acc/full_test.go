package acc

import (
	"database/sql"
	dbmysql "github.com/go-sql-driver/mysql"
	"github.com/nullstone-modules/mysql-db-admin/mysql"
	"github.com/nullstone-modules/mysql-db-admin/workflows"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestFull(t *testing.T) {
	if os.Getenv("ACC") != "1" {
		t.Skip("Set ACC=1 to run e2e tests")
	}

	connUrl := "mysql://mda:mda@localhost:3406/admin"
	adminConnConfig, err := dbmysql.ParseDSN(connUrl)
	require.NoError(t, err, "admin conn config")
	db, err := sql.Open("mysql", adminConnConfig.FormatDSN())
	require.NoError(t, err, "error connecting to mysql")
	defer db.Close()

	newDatabase := mysql.Database{
		Name: "test-database",
	}
	newUser := mysql.User{
		Name:     "test-user",
		Password: "test-password",
	}

	appConnConfig := *adminConnConfig
	appConnConfig.User = newUser.Name
	appConnConfig.Passwd = newUser.Password
	appConnConfig.DBName = newDatabase.Name
	appDb, err := sql.Open("mysql", appConnConfig.FormatDSN())
	defer appDb.Close()

	require.NoError(t, workflows.EnsureDatabase(db, newDatabase), "ensure database")
	require.NoError(t, workflows.EnsureUser(db, newUser), "ensure user")
	require.NoError(t, workflows.GrantDbAccess(db, appDb, newUser, newDatabase), "grant db access")

	// Attempt to create schema

	// Attempt to insert rows into collection

	// Attempt to retrieve them

}
