package main

import (
	"sync"
)

type pool struct {
	numOfWorkers int
	hostQueue    chan string
	mu           *sync.RWMutex
	closed       bool
}

func newPool(numOfWorkers int) *pool {
	return &pool{
		mu:           &sync.RWMutex{},
		numOfWorkers: numOfWorkers,
		hostQueue:    make(chan string),
	}
}

func (p *pool) Close() {
	p.mu.Lock()
	p.closed = true
	close(p.hostQueue)
	p.mu.Unlock()
}

func (p *pool) addTask(host string) {
	p.hostQueue <- host
}

func (p *pool) start() {
	for i := 0; i < p.numOfWorkers; i++ {
		go p.startWorker()
	}
}

func (p *pool) startWorker() {
	pingerService := newPingerService()
	pingResult := false
	for {
		select {
		case host := <-p.hostQueue:
			hostsStorage.create(host)
			pingResult = pingerService.ping(host)
			hostsStorage.set(host, pingResult)
		}
	}
}
