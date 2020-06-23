package main

import (
	"errors"
	"golang.org/x/net/icmp"
	"net"
	"sync"
	"time"
)

var hostsStorage = HostsStorage{mu: &sync.RWMutex{}, hosts: make(map[string]*Host)}

type HostsStorage struct {
	mu    *sync.RWMutex
	hosts map[string]*Host
}

type Host struct {
	mu         *sync.RWMutex
	inWork     bool
	RTT        int64
	lastUpdate int64
	status     bool
	socket     *icmp.PacketConn
	channel    chan bool
	peer       net.Addr
	readBuffer []byte
}

func (hostsStorage *HostsStorage) create(host string) {
	hostsStorage.mu.Lock()
	if hostsStorage.hosts[host] == nil {
		hostsStorage.hosts[host] = &Host{
			mu:         &sync.RWMutex{},
			inWork:     false,
			status:     false,
			RTT:        0,
			lastUpdate: 0,
			channel:    make(chan bool),
			socket:     &icmp.PacketConn{},
			readBuffer: make([]byte, 100),
		}
	}
	hostsStorage.mu.Unlock()

	updater.mu.Lock()
	updater.hosts = append(updater.hosts, host)
	updater.mu.Unlock()
}

func (hostsStorage *HostsStorage) set(host string, status bool) {
	hostsStorage.mu.Lock()
	if hostsStorage.hosts[host] != nil {
		hostsStorage.hosts[host].status = status
		hostsStorage.hosts[host].lastUpdate = time.Now().Unix()
	}
	hostsStorage.hosts[host].inWork = false
	hostsStorage.mu.Unlock()
}

func (hostsStorage *HostsStorage) get(host string) (error, *Host) {
	hostsStorage.mu.RLock()
	defer hostsStorage.mu.RUnlock()
	if hostsStorage.hosts[host] != nil {
		return nil, hostsStorage.hosts[host]
	} else {
		return errors.New("host doesn't exists"), &Host{}
	}
}

func (hostsStorage *HostsStorage) delete(host string) {
	hostsStorage.mu.Lock()
	delete(hostsStorage.hosts, host)
	hostsStorage.mu.Unlock()

	updater.mu.Lock()
	for i := range updater.hosts {
		if host == updater.hosts[i] {
			updater.hosts[i] = updater.hosts[len(updater.hosts)-1]
			updater.hosts[len(updater.hosts)-1] = ""
			updater.hosts = updater.hosts[:len(updater.hosts)-1]
		}
	}
	updater.mu.Unlock()
}

func (host *Host) setRTT(rtt int64) {
	host.mu.Lock()
	host.RTT = rtt
	host.mu.Unlock()
}
