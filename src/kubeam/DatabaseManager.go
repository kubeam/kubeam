package main

import (
	"database/sql"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

/*CheckDatabaseConnection opens and pings a new database connection for test*/
func CheckDatabaseConnection() *sql.DB {
	databaseIP, err := config.GetString("database/host", "127.0.0.1")
	databasePort, err := config.GetString("database/port", "3306")
	databasePassword, err := config.GetString("database/password", "")
	databaseType, err := config.GetString("database/type", "mysql")
	databaseName, err := config.GetString("database/name", "kubeam")
	databaseConn := "root:" + databasePassword +
		"@tcp(" + databaseIP + ":" + databasePort + ")/" + databaseName

	db, err := sql.Open(databaseType, databaseConn)
	if err != nil {
		LogError.Println(err.Error())
	}
	err = db.Ping()
	if err != nil {
		LogError.Println(err.Error())
	}
	return db
}

/*InsertIntoDatabase accepts a dictionary of POST parameters and
Inserts into MySQL database*/
func InsertIntoDatabase(
	eventName,
	msg,
	timestamp,
	dockerTag,
	pipelineName string,
	db *sql.DB) (string, int) {
	insertStmt, err := db.Prepare("INSERT INTO deployStatus VALUES ( ?, ?, ?, ?, ? )")
	response, statusCode := HandleDatabaseError(err)
	if statusCode == 200 {
		_, err = insertStmt.Exec(
			eventName,
			msg,
			timestamp,
			dockerTag,
			pipelineName)
		response, statusCode = HandleDatabaseError(err)
		defer insertStmt.Close()
	}
	LogInfo.Println(response, statusCode)
	return response, statusCode
}

/*HandleDatabaseError handles the mysql returned error codes and returns
http response text and http status code*/
func HandleDatabaseError(err error) (string, int) {
	var statusCode int
	errCodeRegex, _ := regexp.Compile("\\s[0-9]{4}:")
	response := "Query Succeeded"
	statusCode = 200
	errorResponse := map[string]string{
		"1064": "Syntax Error",
		"1146": "Table Does not Exist",
		"1292": "Incorrect Timestamp Format",
	}
	if err != nil {
		LogInfo.Println(err.Error())
		errCode := errCodeRegex.FindString(err.Error())[1:5]
		response = errorResponse[errCode]
		statusCode = 500
	}

	return response, statusCode
}
