package acc

import (
	"database/sql"
	"fmt"
	"github.com/nullstone-modules/mysql-db-admin/mysql"
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

	rootUsername := "root"
	rootPassword := "mda"
	newDatabase := mysql.Database{
		Name: "test-database",
	}
	firstUser := mysql.User{
		Name:     "first-user",
		Password: "first-password",
	}
	secondUser := mysql.User{
		Name:     "second-user",
		Password: "second-password",
	}

	waitForDb := func(t *testing.T, db *sql.DB, timeout time.Duration) {
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

	connect := func(t *testing.T, database, user, password string) (*sql.DB, *mysql.Store) {
		connUrl := fmt.Sprintf("mysql://%s:%s@localhost:3406/%s?multiStatements=true", user, password, database)
		connConfig, err := mysql.DsnFromUrl(connUrl)
		require.NoError(t, err, "conn config")
		connConfig.MultiStatements = true

		db, err := sql.Open("mysql", connConfig.FormatDSN())
		require.NoError(t, err, "error connecting to mysql")

		waitForDb(t, db, 20*time.Second)

		return db, mysql.NewStore(connUrl)
	}

	ensureFull := func(t *testing.T, store *mysql.Store, database mysql.Database, user mysql.User, testSuffix string) {
		database.UseExisting = true
		_, err := store.Databases.Create(database)
		require.NoError(t, err, "create database")

		user.UseExisting = true
		_, err = store.Users.Create(user)
		require.NoError(t, err, "create user")

		_, err = store.DbPrivileges.Create(mysql.DbPrivilege{Username: user.Name, Database: database.Name})
		require.NoError(t, err, "grant db access")
	}

	t.Run("initial setup", func(t *testing.T) {
		// This connection is used by the admin user from the `mysql` db
		rootDb, store := connect(t, "mysql", rootUsername, rootPassword)
		defer rootDb.Close()

		// Run through creation with first user
		ensureFull(t, store, newDatabase, firstUser, "#1")
	})

	t.Run("connect with first user", func(t *testing.T) {
		db, _ := connect(t, newDatabase.Name, firstUser.Name, firstUser.Password)
		defer db.Close()
		require.NoError(t, db.Ping(), "connect to app db using newly created user")
	})

	t.Run("create second user", func(t *testing.T) {
		// This connection is used by the admin user from the `postgres` db
		rootDb, store := connect(t, "mysql", rootUsername, rootPassword)
		defer rootDb.Close()

		ensureFull(t, store, newDatabase, secondUser, "#2")
	})

	t.Run("create schema using second", func(t *testing.T) {
		db, _ := connect(t, newDatabase.Name, secondUser.Name, secondUser.Password)
		defer db.Close()

		// Attempt to create schema objects
		_, err := db.Exec("CREATE TABLE todos ( id SERIAL NOT NULL, name varchar(255) );")
		require.NoError(t, err, "create table")
	})

	t.Run("insert data using second user", func(t *testing.T) {
		db, _ := connect(t, newDatabase.Name, secondUser.Name, secondUser.Password)
		defer db.Close()

		// Attempt to insert records
		sq := strings.Join([]string{
			`INSERT INTO todos (name) VALUES ('item1');`,
			`INSERT INTO todos (name) VALUES ('item2');`,
			`INSERT INTO todos (name) VALUES ('item3');`,
		}, "")
		_, err := db.Exec(sq)
		require.NoError(t, err, "insert todos")
	})

	t.Run("retrieve data using second user", func(t *testing.T) {
		db, _ := connect(t, newDatabase.Name, secondUser.Name, secondUser.Password)
		defer db.Close()

		// Attempt to retrieve them
		results := make([]string, 0)
		rows, err := db.Query(`SELECT * FROM todos`)
		require.NoError(t, err, "query todos")
		defer rows.Close()
		for rows.Next() {
			var id int
			var name string
			require.NoError(t, rows.Scan(&id, &name), "scan record")
			results = append(results, name)
		}
		assert.Equal(t, []string{"item1", "item2", "item3"}, results)
	})

	t.Run("retrieve data from first user", func(t *testing.T) {
		db, _ := connect(t, newDatabase.Name, firstUser.Name, firstUser.Password)
		defer db.Close()

		// Attempt to retrieve them
		results := make([]string, 0)
		rows, err := db.Query(`SELECT * FROM todos`)
		require.NoError(t, err, "query todos")
		defer rows.Close()
		for rows.Next() {
			var id int
			var name string
			require.NoError(t, rows.Scan(&id, &name), "scan record")
			results = append(results, name)
		}
		assert.Equal(t, []string{"item1", "item2", "item3"}, results)
	})
}
