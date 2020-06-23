package main

import (
	"runtime"
	"sync"
)

var workerPool *pool
var httpQueryHost string
var extractedHost *Host
var marshaledMessage []byte
var updater = Updater{mu: &sync.RWMutex{}}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	config := createConfig()
	config.parse()

	workerPool = newPool(config.pingerConfig.numbersOfWorker)
	workerPool.start()

	go updater.start()
	startHttpServer(config)
}
