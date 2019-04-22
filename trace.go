package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"syscall"
	"time"
)

type Hop struct {
	IP      [4]byte
	Host    string
	Elapsed time.Duration
}

type traceOptions struct {
	Port       int
	PacketSize int
	MaxHops    int
}

func trace(dest string, options traceOptions) ([]Hop, error) {
	hops := []Hop{}
	ttl := 1
	// lookup
	localAddr, err := getLocalIP()
	if err != nil {
		return hops, err
	}
	log.Println("local adress to use is: ", localAddr)
	remoteAddr, err := getRemoteIP(dest)
	if err != nil {
		return hops, err
	}
	log.Println("remote IP is: ", remoteAddr)
	retry := 0

	for {

		start := time.Now()
		// socket to receive ICMP packets
		recvSocket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
		if err != nil {
			panic(err)
		}

		// socket to send UDP packets out
		sendSocket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
		if err != nil {
			panic(err)
		}
		// set TTL field of IP Headers
		syscall.SetsockoptInt(sendSocket, 0x0, syscall.IP_TTL, ttl)
		// set timeout for socket to a second
		tv := syscall.NsecToTimeval(1000 * 1000 * 1000)
		syscall.SetsockoptTimeval(recvSocket, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv)

		defer syscall.Close(recvSocket)
		defer syscall.Close(sendSocket)

		// listen on port to receive packets
		syscall.Bind(recvSocket, &syscall.SockaddrInet4{Port: options.Port, Addr: getRawIPV4(localAddr)})

		// Send a single null byte UDP packet
		syscall.Sendto(sendSocket, []byte{0x0}, 0, &syscall.SockaddrInet4{Port: options.Port, Addr: getRawIPV4(remoteAddr)})

		var p = make([]byte, options.PacketSize)
		_, from, err := syscall.Recvfrom(recvSocket, p, 0)
		if err == nil {

			currAddr := from.(*syscall.SockaddrInet4).Addr

			hop := Hop{Elapsed: time.Since(start), IP: currAddr}
			//now do reverse DNS to find hostnames associated with this IP
			currHost, err := net.LookupAddr(IPv4ToStr(currAddr))
			if err == nil || len(currHost) > 0 {
				hop.Host = currHost[0]
			}
			hops = append(hops, hop)
			log.Printf("%d    %s    %s    %s\n", ttl, IPv4ToStr(hop.IP), hop.Host, hop.Elapsed)

			ttl += 1
			retry = 0

			if currAddr == getRawIPV4(remoteAddr) {
				// return nil
				fmt.Println("destination reached, SUCCESS!")
				break
			}
			if ttl > options.MaxHops {
				fmt.Println("max hops reached, FAILED!")
				break
			}

		} else {
			retry++
			if retry > 3 {
				return hops, errors.New("can not reach destination after 3 retry")
			}
			log.Println("could not read from ICMP in time, trying agian")
			continue
		}
	}
	return hops, nil
}
