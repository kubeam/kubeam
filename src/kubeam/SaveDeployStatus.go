package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

/*SaveDeployStatus handles POST request and return status code*/
func SaveDeployStatus(w http.ResponseWriter, r *http.Request) {
	eventName, msg, timestamp, dockerTag, pipelineName := ExtractPostData(r)
	db := CheckDatabaseConnection()
	defer db.Close()
	response, statusCode := InsertIntoDatabase(
		eventName,
		msg,
		timestamp,
		dockerTag,
		pipelineName,
		db)
	LogInfo.Println(response, statusCode)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(statusCode)
	w.Write([]byte(response))
}

/*ExtractPostData accepts HTTP request and extract POST data into a dictionary*/
func ExtractPostData(r *http.Request) (string, string, string, string, string) {
	var postData map[string]string
	var eventName, msg, timestamp, dockerTag, pipelineName string

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		LogError.Printf("HTTP Body read error: %s", err.Error())
	}

	// verifies the structure of json received; that is parameter count and keys
	err = json.Unmarshal(body, &postData)
	if err != nil {
		LogError.Printf("POST Data Extract Error: %s", err.Error())
	}
	LogInfo.Println(postData)

	// extract post data
	eventName = postData["Event"]
	msg = postData["Message"]
	timestamp = postData["Timestamp"]
	dockerTag = postData["DockerTag"]
	pipelineName = postData["Pipeline"]

	return eventName, msg, timestamp, dockerTag, pipelineName
}
