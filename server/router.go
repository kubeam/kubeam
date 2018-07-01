package server

import (
	"github.com/gorilla/mux"
)

func setRoutes(router *mux.Router) {
	router.HandleFunc("/", BasicAuth(AuthZ(Index)))
	router.HandleFunc("/health-check", BasicAuth(AuthZ(HealthCheck)))
	router.HandleFunc("/v1/create/{application}/{environment}/{cluster}/{tag}",
		BasicAuth(AuthZ(ApplicationCreate))).Methods("PUSH")
	router.HandleFunc("/v1/create/{application}/{environment}/{cluster}/{tag}",
		BasicAuth(AuthZ(ApplicationCreate))).Methods("POST")

	router.HandleFunc("/v1/provision/self/{environment}/{tag}",
		BasicAuth(AuthZ(SelfProvision))).Methods("PUSH")
	router.HandleFunc("/v1/provision/self/{environment}/{tag}",
		BasicAuth(AuthZ(SelfProvision))).Methods("POST")

	// Provision methods are used by ephemeral environments. This is they have a timer
	router.HandleFunc("/v1/provision/{application}/{environment}/{tag}",
		BasicAuth(AuthZ(ApplicationProvision))).Methods("PUSH")
	router.HandleFunc("/v1/provision/{application}/{environment}/{tag}",
		BasicAuth(AuthZ(ApplicationProvision))).Methods("POST")
	router.HandleFunc("/v1/provision/{application}/{environment}/{cluster}",
		BasicAuth(AuthZ(ApplicationDelete))).Methods("DELETE")
	router.HandleFunc("/v1/provision/{application}/{environment}/{cluster}",
		BasicAuth(AuthZ(ApplicationStatus))).Methods("GET")

	// Methors use by non ephemeral environments. They are non destructive. Managed resources are long lived (no timer)
	router.HandleFunc("/v1/deploy/{application}/{environment}/{cluster}/{tag}",
		BasicAuth(AuthZ(ApplicationDeploy))).Methods("POST")

	// Gets a list of all cluster reservations
	router.HandleFunc("/v1/listclusters/{application}/{environment}",
		BasicAuth(AuthZ(ApplicationListClusters))).Methods("GET")

	// Gets detailed information about a cluster. (PODS, Docker tags, etc)
	router.HandleFunc("/v1/getclusterdetail/{application}/{environment}/{cluster}",
		BasicAuth(AuthZ(ApplicationGetClusterDetail))).Methods("GET")

	router.HandleFunc("/v1/getallclustersdetail/{application}/{environment}",
		BasicAuth(AuthZ(ApplicationGetAllClustersDetail))).Methods("GET")

	router.HandleFunc("/v1/waitforready/{application}/{environment}/{cluster}",
		BasicAuth(AuthZ(ApplicationWaitForReady))).Methods("GET")

	router.HandleFunc("/v1/provision/{application}/{environment}/{cluster}/{tag}/{ttl}",
		BasicAuth(AuthZ(ApplicationDelete))).Methods("PATCH")

	router.HandleFunc("/v1/event/{application}/{environment}/{cluster}/{tag}",
		BasicAuth(AuthZ(EventStatus))).Methods("POST")

	// Manage Kubernetes Jobs
	router.HandleFunc("/v1/kubejob/{application}/{environment}/{cluster}/{jobname}",
		BasicAuth(AuthZ(RunJob))).Methods("POST")
	router.HandleFunc("/v1/kubejob/{application}/{environment}/{cluster}/{jobname}",
		BasicAuth(AuthZ(GetJobStatus))).Methods("GET")
	router.HandleFunc("/v1/kubejob/{application}/{environment}/{cluster}/{jobname}",
		BasicAuth(AuthZ(DeleteJob))).Methods("DELETE")

	// Github hook to checkout git repos
	router.HandleFunc("/v1/githubhook", GetGithubRepos).Methods("POST")
	router.HandleFunc("/v1/githubhook", LoadGitRepos).Methods("PUSH")

	// Get Clusters for Feature branches
	router.HandleFunc("/v1/featurecluster/{application}/{environment}/{branch}",
		BasicAuth(AuthZ(ReserveFeatureCluster))).Methods("POST")
	router.HandleFunc("/v1/featurecluster/{application}/{environment}/{branch}",
		BasicAuth(AuthZ(GetFeatureCluster))).Methods("GET")
	router.HandleFunc("/v1/featurecluster/{application}/{environment}/{branch}",
		BasicAuth(AuthZ(FreeFeatureCluster))).Methods("DELETE")

}
