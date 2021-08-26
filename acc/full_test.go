package acc

import (
	"database/sql"
	"github.com/nullstone-modules/mysql-db-admin/mysql"
	"github.com/nullstone-modules/mysql-db-admin/workflows"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestFull(t *testing.T) {
	if os.Getenv("ACC") != "1" {
		t.Skip("Set ACC=1 to run e2e tests")
	}

	connUrl := "mysql://root:mda@localhost:3406/"
	adminConnConfig, err := mysql.DsnFromUrl(connUrl)
	require.NoError(t, err, "admin conn config")
	adminConnConfig.MultiStatements = true
	db, err := sql.Open("mysql", adminConnConfig.FormatDSN())
	require.NoError(t, err, "error connecting to mysql")
	defer db.Close()
	// Wait for mysql to launch in docker
	waitForDb(t, db,10*time.Second)

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
	require.NoError(t, workflows.GrantDbAccess(db, newUser, newDatabase), "grant db access")

	// Attempt to create schema
	_, err = appDb.Exec("CREATE TABLE todos ( id SERIAL NOT NULL, name varchar(255) );")
	require.NoError(t, err, "create table")

	// Attempt to insert rows into todos table
	sq := strings.Join([]string{
		`INSERT INTO todos (name) VALUES ('item1');`,
		`INSERT INTO todos (name) VALUES ('item2');`,
		`INSERT INTO todos (name) VALUES ('item3');`,
	}, "")
	_, err = appDb.Exec(sq)
	require.NoError(t, err, "insert todos")

	// Attempt to retrieve them
	results := make([]string, 0)
	rows, err := appDb.Query(`SELECT * FROM todos`)
	require.NoError(t, err, "query todos")
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		require.NoError(t, rows.Scan(&id, &name), "scan record")
		results = append(results, name)
	}
	assert.Equal(t, []string{"item1", "item2", "item3"}, results)
}

func waitForDb(t *testing.T, db *sql.DB, timeout time.Duration) {
	healthy := make(chan bool)
	go func() {
		defer close(healthy)
		for {
			if err := db.Ping(); err == nil {
				healthy <- true
			}
			log.Println("waiting for db...")
			time.Sleep(500 * time.Millisecond)
		}
	}()
	select {
	case <-time.After(timeout):
		t.Fatalf("timed out waiting for database to launch")
	case <-healthy:
	}
}
