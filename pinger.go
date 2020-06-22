package main

import (
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"time"
)

const ResponseTimeout = 200 //ms

type PingerService struct {
	startTime             time.Time
	icmpMessage           icmp.Message
	marshalledIcmpMessage []byte
}

func newPingerService() *PingerService {
	a := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("1"),
		},
	}
	b, _ := a.Marshal(nil)
	return &PingerService{
		icmpMessage:           a,
		marshalledIcmpMessage: b,
	}
}

func (pingerService *PingerService) ping(host string) bool {
	if hostsStorage.hosts[host].inWork {
		return false
	}
	hostsStorage.hosts[host].inWork = true

	var err error
	res := false

	hostsStorage.hosts[host].mu.Lock()
	defer hostsStorage.hosts[host].mu.Unlock()

	hostsStorage.hosts[host].socket, _ = icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	defer hostsStorage.hosts[host].socket.Close()

	peer := hostsStorage.hosts[host].peer

	if _, err := hostsStorage.hosts[host].socket.WriteTo(pingerService.marshalledIcmpMessage, &net.IPAddr{IP: net.ParseIP(host)}); err != nil {
		hostsStorage.hosts[host].RTT = -1
		return false
	}

	pingerService.startTime = time.Now()

	readBuffer := hostsStorage.hosts[host].readBuffer
	endPosition := 0

	go func() {
		endPosition, peer, err = hostsStorage.hosts[host].socket.ReadFrom(readBuffer)
		if err != nil {
			res = false
		} else {
			hostsStorage.hosts[host].channel <- true
		}
	}()

	select {
	case <-hostsStorage.hosts[host].channel:
		res = true
		//fmt.Println(host, "ok")
	case <-time.After(time.Millisecond * ResponseTimeout):
		res = false
		//fmt.Println(host, "timeout")
	}

	responseMessage, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), readBuffer[:endPosition])
	if err != nil {
		res = false
		hostsStorage.hosts[host].RTT = -1
		return res
	}

	switch responseMessage.Type {
	case ipv4.ICMPTypeEchoReply:
		res = true
		hostsStorage.hosts[host].RTT = int64(time.Since(pingerService.startTime)/1e6)
	default:
		res = false
		hostsStorage.hosts[host].RTT = -1
	}
	return res
}
