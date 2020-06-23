package main

import (
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
		go p.worker()
	}

	for i := 0; i < p.numOfWorkers; i++ {
		go p.newHostWorker()
	}
}

func (p *pool) worker() {
	pingerService := newPingerService()
	pingResult := false
	var err error
	for {
		for host := range p.taskQueue {
			err, pingResult = pingerService.ping(host)
			if err == nil {
				hostsStorage.set(host, pingResult)
			}
		}
	}
}

func (p *pool) newHostWorker() {
	pingerService := newPingerService()
	pingResult := false
	var err error
	for {
		for host := range p.newHostQueue {
			hostsStorage.create(host)
			err, pingResult = pingerService.ping(host)
			if err == nil {
				hostsStorage.set(host, pingResult)
			}
		}
	}
}
