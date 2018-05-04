package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/creamdog/gonfig"
	"github.com/gorilla/mux"
)

var router *mux.Router

func TestAuthZ(t *testing.T) {
	type args struct {
		handler http.HandlerFunc
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AuthZ(tt.args.handler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthZ() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setUp(t *testing.T) {
	router = mux.NewRouter().StrictSlash(true)
	setRoutes(router)

	// Init Loggers:
	// File descriptors in order: Trace, Debug, Info, Warning, Error
	// set to ** ioutil.Discard ** to stop recording those logs
	InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	// Read application config from file
	f, err := os.Open("config-test.yaml")
	if err != nil {
		t.Error(err)
		os.Exit(1)
	}
	defer f.Close()
	config, err = gonfig.FromYml(f)
	if err != nil {
		t.Errorf("Error: %v", err)
		os.Exit(2)
	}
}

func testAuthzRequest(t *testing.T, router *mux.Router, user string, path string, method string, expectedCode int) {
	r, _ := http.NewRequest(method, path, nil)
	r.SetBasicAuth(user, "123456")
	r.URL.Path = strings.ToLower(r.URL.Path)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != expectedCode {
		t.Errorf("%s, %s, %s: %d, supposed to be %d", user, path, method, w.Code, expectedCode)
	}
}

func TestBasic(t *testing.T) {
	setUp(t)

	// Unauthorized request.
	testAuthzRequest(t, router, "nosuchuser", "/v1/provision/self/QA/pr-build-12345", "PUSH", 401)

	// Authorized request.
	// testAuthzRequest(t, router, "admin", "/v1/provision/self/QA/pr-build-12345", "PUSH", 200)
}

func TestRBAC(t *testing.T) {
	setUp(t)

	testAuthzRequest(t, router, "admin", "/v1/provision/self/QA/pr-build-12345", "PUSH", 200)
}
