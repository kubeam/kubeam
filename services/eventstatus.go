package services

import (
	"github.com/go-sql-driver/mysql"
	"github.com/kubeam/kubeam/common"
)

// SaveEventStatus - stores events to database
func SaveEventStatus(eventdata map[string]string) error {
	db := GetDatabaseConnection()
	defer db.Close()

	sql, err := db.Prepare("INSERT INTO eventstatus VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		common.LogError.Printf("Error preparing SQL statement")
		return err
	}

	_, err = sql.Exec(eventdata["event"],
		eventdata["message"],
		eventdata["timestamp"],
		eventdata["buildtag"],
		eventdata["pipeline"])

	if err != nil {
		if me, ok := err.(*mysql.MySQLError); !ok {
			common.LogError.Printf("Error Executing SQL statement: %s", me.Message)
		}
		return err
	}

	return nil

}
