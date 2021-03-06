/*

environment-exporter
  A webserver for serving a JSON list of processes that are listening on ports.

*/

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	common "github.com/sharingio/environment/pkg/common"
	l "github.com/sharingio/environment/pkg/listening-processes"
)

var (
	AppBuildVersion = "0.0.0"
	AppBuildHash    = "???"
	AppBuildDate    = "yyyy.mm.dd HH:MM"
)

func getListenRoute(w http.ResponseWriter, r *http.Request) {
	listening, err := l.ListListeningProcesses()
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	listeningBytes, err := json.Marshal(listening)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(listeningBytes)
}

func main() {
	log.Printf("launching sharingio/environment:exporter (%v, %v, %v)\n", AppBuildVersion, AppBuildHash, AppBuildDate)
	envFile := common.GetAppEnvFile()
	_ = godotenv.Load(envFile)
	port := common.GetAppPort()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/listening", getListenRoute).Methods(http.MethodGet)
	router.Use(common.Logging)

	srv := &http.Server{
		Handler:      router,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Listening on", port)
	log.Fatal(srv.ListenAndServe())
}
