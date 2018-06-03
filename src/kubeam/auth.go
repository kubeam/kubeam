package main

import (
	"crypto/subtle"
	"net/http"
	"os"

	"github.com/casbin/casbin"
)

// BasicAuth wraps a handler requiring HTTP basic auth for it using the given
// username and password and the specified realm, which shouldn't contain quotes.
//
// Most web browser display a dialog with something like:
//
//    The website says: "<realm>"
//
// Which is really stupid so you may want to set the realm to a message rather than
// an actual realm.
func BasicAuth(handler http.HandlerFunc) http.HandlerFunc {

	realm := "Please enter your username and password for this site"

	return func(w http.ResponseWriter, r *http.Request) {

		if !checkUser(r) {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorized.\n"))
			return
		}

		handler(w, r)
	}
}

// AuthZ allow authz to be handled.
func AuthZ(handler http.HandlerFunc) http.HandlerFunc {

	// handler := BasicAuth(handler)

	e := casbin.NewEnforcer("authz-model.conf", "authz-policy.csv")

	return func(w http.ResponseWriter, r *http.Request) {

		if !checkPermission(e, r) {
			w.WriteHeader(403)
			w.Write([]byte("Fobidden.\n"))
			return
		}

		handler(w, r)
	}
}

func checkUser(r *http.Request) bool {
	username, err := config.GetString("credentials/kubeam/username", "invalid-user")
	if err != nil {
		LogError.Printf("FATAL configuration file: %v\n", err)
		os.Exit(1)
	}
	password, err := config.GetString("credentials/kubeam/password", "invalid-password")
	if err != nil {
		LogError.Printf("FATAL configuration file: %v\n", err)
		os.Exit(1)
	}

	user, pass, ok := r.BasicAuth()

	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
		return false
	}

	return true
}

func checkPermission(e *casbin.Enforcer, r *http.Request) bool {
	user, _, _ := r.BasicAuth()
	method := r.Method
	path := r.URL.Path

	return e.Enforce(user, path, method)
}
