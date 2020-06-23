package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"time"
)

func startHttpServer(config *Config) {

	router := mux.NewRouter()
	router.HandleFunc("/addHost", addHostHandler)
	router.HandleFunc("/delHost", delHostHandler)
	router.HandleFunc("/getHost", getHostHandler)
	router.HandleFunc("/updateHost", addHostHandler)

	s := &http.Server{
		Addr:         ":" + config.httpConfig.port,
		Handler:      router,
		ReadTimeout:  config.httpConfig.readTimeout * time.Second,
		WriteTimeout: config.httpConfig.writeTimeout * time.Second,
		IdleTimeout:  config.httpConfig.idleTimeout * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}

func addHostHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if val, ok := request.URL.Query()["ip"]; ok {
		if val[0] != "" {
			httpQueryHost = val[0]
			responseWriter.Header().Set("Content-Type", "text/html")

			if hostsStorage.hosts[httpQueryHost] != nil && hostsStorage.hosts[httpQueryHost].inWork {
				_, _ = io.WriteString(responseWriter, "host is busy")
			} else if hostsStorage.hosts[httpQueryHost] == nil {
				workerPool.addNewHost(httpQueryHost)
				_, _ = io.WriteString(responseWriter, "ok")
			} else {
				_, _ = io.WriteString(responseWriter, "host exist")
			}
		}
	}
}

func delHostHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if val, ok := request.URL.Query()["ip"]; ok {
		if val[0] != "" {
			httpQueryHost = val[0]
			hostsStorage.delete(httpQueryHost)
			responseWriter.Header().Set("Content-Type", "text/html")
			_, _ = io.WriteString(responseWriter, "ok")
		}
	}
}

func getHostHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if val, ok := request.URL.Query()["ip"]; ok {
		if val[0] != "" {
			httpQueryHost = val[0]
			_, extractedHost = hostsStorage.get(httpQueryHost)
			responseWriter.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(responseWriter, jsonResponse(extractedHost))
			_ = request.Body.Close()
		}
	}
}

func jsonResponse(res *Host) string {
	marshaledMessage, _ = json.Marshal(
		&HttpResponseTemplate{
			Status:     res.status,
			LastUpdate: res.lastUpdate,
			RTT:        res.RTT},
	)
	return string(marshaledMessage)
}

type HttpResponseTemplate struct {
	Status     bool  `json:"status"`
	LastUpdate int64 `json:"lastUpdate"`
	RTT        int64 `json:"rtt"`
}
