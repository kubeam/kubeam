package main

import (
	"fmt"
	"net/http"
)

/*Index ...*/
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "These aren't the droids you're looking for.")
}

/*HealthCheck is a health check for kubeAM*/
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
