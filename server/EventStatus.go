package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

/*EventStatus parses and saves Jenkins event post data to database*/
func EventStatus(w http.ResponseWriter, r *http.Request) {
	var eventdata map[string]string
	defer r.Body.Close()

	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if json.Unmarshal(body, &eventdata) == nil {
			saveEventStatus(eventdata)
		}
	}
}

func saveEventStatus(eventdata map[string]string) {
	db := GetDatabaseConnection()
	defer db.Close()

	sql, err := db.Prepare("INSERT INTO eventstatus VALUES (?, ?, ?, ?, ?)")
	errorz(err)

	_, err = sql.Exec(eventdata["event"],
		eventdata["message"],
		eventdata["timestamp"],
		eventdata["buildtag"],
		eventdata["pipeline"])

	errorz(err)
}

func errorz(err error) {
	if err != nil {
		if me, ok := err.(*mysql.MySQLError); !ok {
			LogError.Fatalf("mysql error for event status: %s", me.Message)
		}
		LogError.Println(err.Error())
	}
}
