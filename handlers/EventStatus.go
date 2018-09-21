package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"fmt"
	"github.com/kubeam/kubeam/services"
)

/*EventStatus parses and saves Jenkins event post data to database*/
func EventStatus(w http.ResponseWriter, r *http.Request) {
	var eventdata map[string]string
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if json.Unmarshal(body, &eventdata) == nil {
			err := services.SaveEventStatus(eventdata)
			if err != nil {
				str := fmt.Sprintf(`{"status": "error", "message": "%v"}`, err.Error())
				w.Write([]byte(str))
				return

			}
		}
	}
	str := fmt.Sprintf(`{"status": "success", "message": "%v"}`, "Event saved")
	w.Write([]byte(str))
	return
}
