package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"
)

var workerPool *pool
var httpQueryHost string
var extractedHost *Host
var marshaledMessage []byte

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	workerPool = newPool(runtime.NumCPU() * 2)
	workerPool.start()
	/*	workerPool.addTask("1.1.1.2")
		workerPool.addTask("1.1.1.1")
		workerPool.addTask("8.8.8.8")*/

	go Updater()

	http.HandleFunc("/addHost", addHostHandler)
	http.HandleFunc("/delHost", delHostHandler)
	http.HandleFunc("/getHost", getHostHandler)
	http.HandleFunc("/updateHost", addHostHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func addHostHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if val, ok := request.URL.Query()["ip"]; ok {
		if val[0] != "" {
			httpQueryHost = val[0]
			workerPool.addTask(httpQueryHost)
			responseWriter.Header().Set("Content-Type", "text/html")
			_, _ = io.WriteString(responseWriter, "ok")
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

func Updater() {
	time.Sleep(60 * time.Second)
	for {
		for key := range hostsStorage.hosts {
			workerPool.addTask(key)
		}
		time.Sleep(2 * time.Second)
	}
}
