package main

import (
	"errors"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"time"
)

var ResponseTimeout = 50 //ms

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

func (pingerService *PingerService) ping(hostIP string) (error,bool) {

	var pingError error = nil
	_, host := hostsStorage.get(hostIP)

	if host.inWork {
		pingError = errors.New("host is busy")
		return pingError,false
	}

	host.mu.Lock()

	host.inWork = true

	var err error
	res := false

	host.socket, _ = icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	defer host.socket.Close()

	peer := host.peer
	host.mu.Unlock()
	if _, err := host.socket.WriteTo(pingerService.marshalledIcmpMessage, &net.IPAddr{IP: net.ParseIP(hostIP)}); err != nil {
		host.setRTT(-1)
		return pingError,false
	}

	pingerService.startTime = time.Now()
	host.mu.Lock()
	readBuffer := host.readBuffer
	endPosition := 0

	go func() {
		endPosition, peer, err = host.socket.ReadFrom(readBuffer)
		if err != nil {
			res = false
		} else {
			host.channel <- true
		}
	}()

	select {
	case <-host.channel:
		res = true
		//fmt.Println(hostIP, "ok")
	case <-time.After(time.Millisecond * time.Duration(ResponseTimeout)):
		res = false
		//fmt.Println(hostIP, "timeout")
	}
	host.mu.Unlock()
	responseMessage, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), readBuffer[:endPosition])
	if err != nil {
		res = false
		host.setRTT(-1)
		return pingError,res
	}

	switch responseMessage.Type {
	case ipv4.ICMPTypeEchoReply:
		res = true
		host.setRTT(int64(time.Since(pingerService.startTime) / 1e6))
	default:
		res = false
		host.setRTT(-1)
	}
	return pingError,res
}
