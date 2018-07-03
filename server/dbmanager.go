package server

import (
	"database/sql"
	"fmt"
)

/*GetDatabaseConnection opens, tests and returns a new database connection*/
func GetDatabaseConnection() *sql.DB {
	var err error
	var db *sql.DB

	databaseIP, err := config.GetString("database/host", "localhost")
	databasePort, err := config.GetString("database/port", "3306")
	databasePassword, err := config.GetString("database/password", "")
	databaseType, err := config.GetString("database/type", "mysql")
	databaseName, err := config.GetString("database/name", "kubeam")
	databaseUser, err := config.GetString("database/user", "sample-user")

	if err == nil {
		databaseConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", databaseUser, databasePassword, databaseIP, databasePort, databaseName)
		if db, err = sql.Open(databaseType, databaseConn); err == nil {
			if err = db.Ping(); err == nil {
				return db
			}
		}
	}

	LogError.Println(err.Error())
	return nil
}
