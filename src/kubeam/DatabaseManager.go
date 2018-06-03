package main

import (
	"database/sql"
)

/*GetDatabaseConnection opens, test and returns a new database connection*/
func GetDatabaseConnection() *sql.DB {
	var err error
	var db *sql.DB

	databaseIP, err := config.GetString("database/host", "127.0.0.1")
	databasePort, err := config.GetString("database/port", "3306")
	databasePassword, err := config.GetString("database/password", "")
	databaseType, err := config.GetString("database/type", "mysql")
	databaseName, err := config.GetString("database/name", "kubeam")
	databaseConn := "root:" + databasePassword +
		"@tcp(" + databaseIP + ":" + databasePort + ")/" + databaseName

	if db, err = sql.Open(databaseType, databaseConn); err == nil {
		if err = db.Ping(); err == nil {
			return db
		}
	}
	LogError.Println(err.Error())
	return nil
}
