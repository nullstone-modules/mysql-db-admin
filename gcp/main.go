package gcp

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/nullstone-modules/mysql-db-admin/api"
	"github.com/nullstone-modules/mysql-db-admin/mysql"
	"os"
)

var (
	dbConnUrlEnvVar = "DB_CONN_URL"
)

func init() {
	store := mysql.NewStore(os.Getenv(dbConnUrlEnvVar))
	router := api.CreateRouter(store)
	functions.HTTP("mysql-db-admin", router.ServeHTTP)
}
