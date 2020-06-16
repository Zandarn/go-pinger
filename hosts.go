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
			status:     false,
			RTT:        0,
			lastUpdate: 0,
			channel:    make(chan bool),
			socket:     &icmp.PacketConn{},
			readBuffer: make([]byte, 100),
		}
	}
	hostsStorage.mu.Unlock()
}

func (hostsStorage *HostsStorage) set(host string, status bool) {
	hostsStorage.mu.Lock()
	if hostsStorage.hosts[host] != nil {
		hostsStorage.hosts[host].status = status
		hostsStorage.hosts[host].lastUpdate = time.Now().Unix()
	}

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
}
