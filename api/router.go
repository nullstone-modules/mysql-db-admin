package api

import (
	"github.com/gorilla/mux"
	"github.com/nullstone-io/go-rest-api"
	"github.com/nullstone-modules/mysql-db-admin/mysql"
	"log"
	"net/http"
)

func CreateRouter(store *mysql.Store) *mux.Router {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%d %s %s\n", http.StatusNotFound, r.Method, r.RequestURI)
		http.NotFound(w, r)
	})
	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%d %s %s\n", http.StatusMethodNotAllowed, r.Method, r.RequestURI)
		http.Error(w, "", http.StatusMethodNotAllowed)
	})

	r.Methods(http.MethodDelete).Path("/skip").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	databases := &rest.Resource[string, mysql.Database]{
		DataAccess: store.Databases,
		KeyParser:  rest.PathParameterKeyParser("name"),
	}
	r.Methods(http.MethodPost).Path("/databases").HandlerFunc(databases.Create)
	r.Methods(http.MethodGet).Path("/databases/{name}").HandlerFunc(databases.Get)
	r.Methods(http.MethodPut).Path("/databases/{name}").HandlerFunc(databases.Update)
	r.Methods(http.MethodDelete).Path("/databases/{name}").HandlerFunc(databases.Delete)

	users := rest.Resource[string, mysql.User]{
		DataAccess: store.Users,
		KeyParser:  rest.PathParameterKeyParser("name"),
	}
	r.Methods(http.MethodPost).Path("/users").HandlerFunc(users.Create)
	r.Methods(http.MethodGet).Path("/users/{name}").HandlerFunc(users.Get)
	r.Methods(http.MethodPut).Path("/users/{name}").HandlerFunc(users.Update)
	r.Methods(http.MethodDelete).Path("/users/{name}").HandlerFunc(users.Delete)

	dbPrivileges := rest.Resource[mysql.DbPrivilegeKey, mysql.DbPrivilege]{
		DataAccess: store.DbPrivileges,
		KeyParser: func(r *http.Request) (mysql.DbPrivilegeKey, error) {
			vars := mux.Vars(r)
			return mysql.DbPrivilegeKey{
				Database: vars["database"],
				Username: vars["username"],
			}, nil
		},
	}
	r.Methods(http.MethodPost).Path("/databases/{database}/db_privileges").HandlerFunc(dbPrivileges.Create)
	r.Methods(http.MethodGet).Path("/databases/{database}/db_privileges/{username}").HandlerFunc(dbPrivileges.Get)
	r.Methods(http.MethodPut).Path("/databases/{database}/db_privileges/{username}").HandlerFunc(dbPrivileges.Update)
	r.Methods(http.MethodDelete).Path("/databases/{database}/db_privileges/{username}").HandlerFunc(dbPrivileges.Delete)

	return r
}
