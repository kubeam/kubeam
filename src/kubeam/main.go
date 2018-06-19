package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/creamdog/gonfig"
	"github.com/gorilla/mux"
)

var config gonfig.Gonfig

var (
	//LogTrace ...
	LogTrace *log.Logger
	//LogDebug ...
	LogDebug *log.Logger
	//LogInfo ...
	LogInfo *log.Logger
	//LogWarning ...
	LogWarning *log.Logger
	//LogError ...
	LogError *log.Logger
)

/*InitLogger initializes different logging handlers*/
func InitLogger(
	traceHandle io.Writer,
	debugHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	LogTrace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	LogDebug = log.New(debugHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	LogInfo = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	LogWarning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	LogError = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

/*LowerCaseURI returns a lower cases converted URL path*/
func LowerCaseURI(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func main() {

	// Init Loggers:
	// File descriptors in order: Trace, Debug, Info, Warning, Error
	// set to ** ioutil.Discard ** to stop recording those logs
	InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	//
	// Read application config from file
	f, err := os.Open("config-sample.yaml")
	if err != nil {
		LogError.Printf("Error configuration file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	config, err = gonfig.FromYml(f)
	if err != nil {
		LogError.Printf("Error configuration file: %v\n", err)
		os.Exit(1)
	}
	//fmt.Println(reflect.TypeOf(config))
	redisClient = NewDBClient()

	//StartScheduler()

	router := mux.NewRouter().StrictSlash(true)
	setRoutes(router)

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		},
	}
	//server_port, err := config.GetString("https/port", 8443)
	//server_addr := fmt.Sprintf("%v:%v", "", server_port)

	serverMuxHC := http.NewServeMux()
	serverMuxHC.HandleFunc("/health-check", HealthCheck)

	// Health Check plan http listener
	go func() {
		http.ListenAndServe(":8081", serverMuxHC)
	}()

	srv := &http.Server{
		Addr:         ":443",
		Handler:      LowerCaseURI(router),
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	sslcert, err := config.GetString("https/certificate", "server.crt")
	sslkey, err := config.GetString("https/key", "server.key")

	log.Fatal(srv.ListenAndServeTLS(sslcert, sslkey))
	//log.Fatal(http.ListenAndServe(":8080", router))

}
