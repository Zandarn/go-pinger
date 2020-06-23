package main

import (
	"sync"
	"time"
)

type Updater struct {
	hosts []string
	mu    *sync.RWMutex
}

func (updater *Updater) start() {
	time.Sleep(60 * time.Second)
	for {
		updater.mu.RLock()
		for _, v := range updater.hosts {
			workerPool.addTask(v)
		}
		updater.mu.RUnlock()
		time.Sleep(2 * time.Second)
	}
}
