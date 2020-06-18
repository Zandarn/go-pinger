package main

import (
	"fmt"
	"sync"
)

type pool struct {
	numOfWorkers int
	taskQueue    chan string
	newHostQueue chan string
	mu           *sync.RWMutex
	closed       bool
}

func newPool(numOfWorkers int) *pool {
	return &pool{
		mu:           &sync.RWMutex{},
		numOfWorkers: numOfWorkers,
		taskQueue:    make(chan string),
		newHostQueue: make(chan string),
	}
}

func (p *pool) Close() {
	p.mu.Lock()
	p.closed = true
	close(p.taskQueue)
	close(p.newHostQueue)
	p.mu.Unlock()
}

func (p *pool) addTask(host string) {
	p.taskQueue <- host
}

func (p *pool) addNewHost(host string) {
	p.newHostQueue <- host
}

func (p *pool) start() {
	for i := 0; i < p.numOfWorkers; i++ {
		go p.startWorker()
	}

	for i := 0; i < p.numOfWorkers; i++ {
		go p.newHostWorker()
	}
}

func (p *pool) startWorker() {
	pingerService := newPingerService()
	pingResult := false
	for {
		select {
		case host := <-p.taskQueue:
			fmt.Println("host from queue:",host)
			pingResult = pingerService.ping(host)
			fmt.Println("host result from queue:",pingResult)
			hostsStorage.set(host, pingResult)
			fmt.Println("host ok")
		}
	}
}

func (p *pool) newHostWorker() {
	pingerService := newPingerService()
	pingResult := false
	for {
		select {
		case host := <-p.newHostQueue:
			hostsStorage.create(host)
			fmt.Println("new host from queue:",host)
			pingResult = pingerService.ping(host)
			fmt.Println("new host result from queue:",pingResult)
			hostsStorage.set(host, pingResult)
		}
	}
}
