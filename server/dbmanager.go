package server

import (
	"database/sql"
	"fmt"
	"github.com/kubeam/kubeam/common"
)

/*GetDatabaseConnection opens, tests and returns a new database connection*/
func GetDatabaseConnection() *sql.DB {
	var err error
	var db *sql.DB

	databaseIP, err := common.Config.GetString("database/host", "localhost")
	databasePort, err := common.Config.GetString("database/port", "3306")
	databasePassword, err := common.Config.GetString("database/password", "")
	databaseType, err := common.Config.GetString("database/type", "mysql")
	databaseName, err := common.Config.GetString("database/name", "kubeam")
	databaseUser, err := common.Config.GetString("database/user", "sample-user")

	if err == nil {
		databaseConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", databaseUser, databasePassword, databaseIP, databasePort, databaseName)
		if db, err = sql.Open(databaseType, databaseConn); err == nil {
			if err = db.Ping(); err == nil {
				return db
			}
		}
	}

	common.LogError.Println(err.Error())
	return nil
}
