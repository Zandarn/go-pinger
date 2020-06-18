package main

import (
	"github.com/gorilla/mux"
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
