package main

import (
	"github.com/gorilla/mux"
)

func setRoutes(router *mux.Router) {
	router.HandleFunc("/", AuthZ(Index))
	router.HandleFunc("/health-check", AuthZ(HealthCheck))
	router.HandleFunc("/v1/create/{application}/{environment}/{cluster}/{tag}",
		AuthZ(ApplicationCreate)).Methods("PUSH")
	router.HandleFunc("/v1/create/{application}/{environment}/{cluster}/{tag}",
		AuthZ(ApplicationCreate)).Methods("POST")

	router.HandleFunc("/v1/provision/self/{environment}/{tag}",
		AuthZ(SelfProvision)).Methods("PUSH")
	router.HandleFunc("/v1/provision/self/{environment}/{tag}",
		AuthZ(SelfProvision)).Methods("POST")

	// Provision methods are used by ephemeral environments. This is they have a timer
	router.HandleFunc("/v1/provision/{application}/{environment}/{tag}",
		AuthZ(ApplicationProvision)).Methods("PUSH")
	router.HandleFunc("/v1/provision/{application}/{environment}/{tag}",
		AuthZ(ApplicationProvision)).Methods("POST")
	router.HandleFunc("/v1/provision/{application}/{environment}/{cluster}",
		AuthZ(ApplicationDelete)).Methods("DELETE")
	router.HandleFunc("/v1/provision/{application}/{environment}/{cluster}",
		AuthZ(ApplicationStatus)).Methods("GET")

	// Methors use by non ephemeral environments. They are non destructive. Managed resources are long lived (no timer)
	router.HandleFunc("/v1/deploy/{application}/{environment}/{cluster}/{tag}",
		AuthZ(ApplicationDeploy)).Methods("POST")

	// Gets a list of all cluster reservations
	router.HandleFunc("/v1/listclusters/{application}/{environment}",
		AuthZ(ApplicationListClusters)).Methods("GET")

	// Gets detailed information about a cluster. (PODS, Docker tags, etc)
	router.HandleFunc("/v1/getclusterdetail/{application}/{environment}/{cluster}",
		AuthZ(ApplicationGetClusterDetail)).Methods("GET")

	router.HandleFunc("/v1/getallclustersdetail/{application}/{environment}",
		AuthZ(ApplicationGetAllClustersDetail)).Methods("GET")

	router.HandleFunc("/v1/waitforready/{application}/{environment}/{cluster}",
		AuthZ(ApplicationWaitForReady)).Methods("GET")

	router.HandleFunc("/v1/provision/{application}/{environment}/{cluster}/{tag}/{ttl}",
		AuthZ(ApplicationDelete)).Methods("PATCH")

	router.HandleFunc("/v1/event/{application}/{environment}/{cluster}/{tag}",
		AuthZ(EventStatus)).Methods("POST")

	// Manage kubernetes Job objects
	router.HandleFunc("/v1/kubejob/{application}/{environment}/{cluster}/{tag}",
		AuthZ(RunJob)).Methods("POST")
}
