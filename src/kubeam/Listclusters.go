package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

/*ApplicationGetClusterDetail ...*/
func ApplicationGetClusterDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	application := vars["application"]
	appEnv := vars["environment"]
	appCluster := vars["cluster"]

	clusterList, err := DBClientGetSingleClusterDetail(redisClient, application, appEnv, appCluster)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		str := `{"status": "error", "description": "Unable to select cluster"}`
		w.Write([]byte(str))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(clusterList))
}

/*ApplicationGetAllClustersDetail ...*/
func ApplicationGetAllClustersDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	application := vars["application"]
	appEnv := vars["environment"]

	clusterList, err := DBClientGetAllClustersDetail(redisClient, application, appEnv)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		str := `{"status": "error", "description": "Unable get detail on environment"}`
		w.Write([]byte(str))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(clusterList))
}

/*ApplicationListClusters ...*/
func ApplicationListClusters(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	application := vars["application"]
	appEnv := vars["environment"]

	clusterList, err := DBClientGetAllClusters(redisClient, application, appEnv)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		str := `{"status": "error", "description": "Unable to list clusters for environment"}`
		w.Write([]byte(str))
		//w.Write([]byte( "ERROR: Unable to select cluster for specified environment, No free slots?\n" ))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(clusterList))
}
