package server

import (
	"github.com/gorilla/mux"
	"github.com/kubeam/kubeam/handlers"
)

// SetRoutes - configures all https routes
func SetRoutes(router *mux.Router) {
	router.HandleFunc("/", BasicAuth(AuthZ(handlers.Index)))
	router.HandleFunc("/health-check", BasicAuth(AuthZ(handlers.HealthCheck)))
	router.HandleFunc("/v1/create/{application}/{environment}/{cluster}/{tag}",
		BasicAuth(AuthZ(handlers.ApplicationCreate))).Methods("PUSH")
	router.HandleFunc("/v1/create/{application}/{environment}/{cluster}/{tag}",
		BasicAuth(AuthZ(handlers.ApplicationCreate))).Methods("POST")

	router.HandleFunc("/v1/provision/self/{environment}/{tag}",
		BasicAuth(AuthZ(handlers.SelfProvision))).Methods("PUSH")
	router.HandleFunc("/v1/provision/self/{environment}/{tag}",
		BasicAuth(AuthZ(handlers.SelfProvision))).Methods("POST")

	// Provision methods are used by ephemeral environments. This is they have a timer
	router.HandleFunc("/v1/provision/{application}/{environment}/{tag}",
		BasicAuth(AuthZ(handlers.ApplicationProvision))).Methods("PUSH")
	router.HandleFunc("/v1/provision/{application}/{environment}/{tag}",
		BasicAuth(AuthZ(handlers.ApplicationProvision))).Methods("POST")
	router.HandleFunc("/v1/provision/{application}/{environment}/{cluster}",
		BasicAuth(AuthZ(handlers.ApplicationDelete))).Methods("DELETE")
	router.HandleFunc("/v1/provision/{application}/{environment}/{cluster}",
		BasicAuth(AuthZ(handlers.ApplicationStatus))).Methods("GET")

	// Methors use by non ephemeral environments. They are non destructive. Managed resources are long lived (no timer)
	router.HandleFunc("/v1/deploy/{application}/{environment}/{cluster}/{tag}",
		BasicAuth(AuthZ(handlers.ApplicationDeploy))).Methods("POST")

	// Gets a list of all cluster reservations
	router.HandleFunc("/v1/listclusters/{application}/{environment}",
		BasicAuth(AuthZ(handlers.ApplicationListClusters))).Methods("GET")

	// Gets detailed information about a cluster. (PODS, Docker tags, etc)
	router.HandleFunc("/v1/getclusterdetail/{application}/{environment}/{cluster}",
		BasicAuth(AuthZ(handlers.ApplicationGetClusterDetail))).Methods("GET")

	router.HandleFunc("/v1/getallclustersdetail/{application}/{environment}",
		BasicAuth(AuthZ(handlers.ApplicationGetAllClustersDetail))).Methods("GET")

	router.HandleFunc("/v1/waitforready/{application}/{environment}/{cluster}",
		BasicAuth(AuthZ(handlers.ApplicationWaitForReady))).Methods("GET")

	router.HandleFunc("/v1/provision/{application}/{environment}/{cluster}/{tag}/{ttl}",
		BasicAuth(AuthZ(handlers.ApplicationDelete))).Methods("PATCH")

	router.HandleFunc("/v1/event/{application}/{environment}/{cluster}/{tag}",
		BasicAuth(AuthZ(handlers.EventStatus))).Methods("POST")

	// Manage Kubernetes Jobs
	router.HandleFunc("/v1/kubejob/{application}/{environment}/{cluster}/{jobname}",
		BasicAuth(AuthZ(handlers.RunJob))).Methods("POST")
	router.HandleFunc("/v1/kubejob/{application}/{environment}/{cluster}/{jobname}",
		BasicAuth(AuthZ(handlers.GetJobStatus))).Methods("GET")
	router.HandleFunc("/v1/kubejob/{application}/{environment}/{cluster}/{jobname}",
		BasicAuth(AuthZ(handlers.DeleteJob))).Methods("DELETE")

	// Github hook to checkout git repos
	router.HandleFunc("/v1/githubhook", handlers.GetGithubRepos).Methods("POST")
	router.HandleFunc("/v1/githubhook", handlers.LoadGitRepos).Methods("PUSH")

	// Get Clusters for Feature branches
	router.HandleFunc("/v1/featurecluster/{application}/{environment}/{branch}",
		BasicAuth(AuthZ(handlers.ReserveFeatureCluster))).Methods("POST")
	router.HandleFunc("/v1/featurecluster/{application}/{environment}/{branch}",
		BasicAuth(AuthZ(handlers.GetFeatureCluster))).Methods("GET")
	router.HandleFunc("/v1/featurecluster/{application}/{environment}/{branch}",
		BasicAuth(AuthZ(handlers.FreeFeatureCluster))).Methods("DELETE")

}
