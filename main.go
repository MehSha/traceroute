package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("starting trace")
	hostPtr := flag.String("host", "google.com", "the host name")
	portPtr := flag.Int("port", 33434, "listening port to listen to ICMP packets")
	maxHopsPtr := flag.Int("max_hops", 255, "maximum hops to check between local to remote IPs")
	packetSizePtr := flag.Int("packet_size", 52, "listening port to listen to ICMP packets")
	flag.Parse()

	hops, err := trace(*hostPtr, traceOptions{
		Port:       *portPtr,
		PacketSize: *packetSizePtr,
		MaxHops:    *maxHopsPtr,
	})
	if err != nil {
		fmt.Println("can not find the route to destination.", err)
		os.Exit(1)
	}

	baseHop, nextHop, duration := getMaxHopDistance(hops)
	fmt.Printf("the maximum hop distance is between %s(%s) and %s(%s) with %s duration\n", IPv4ToStr(baseHop.IP), baseHop.Host,
		IPv4ToStr(nextHop.IP), nextHop.Host, duration)
}

func getMaxHopDistance(hops []Hop) (Hop, Hop, time.Duration) {
	baseHop := hops[0]
	nextHop := hops[1]
	duration := nextHop.Elapsed - baseHop.Elapsed

	for i := 2; i < len(hops); i++ {
		diff := hops[i].Elapsed - hops[i-1].Elapsed
		if diff > duration {
			//found a bigger difference
			duration = diff
			baseHop = hops[i-1]
			nextHop = hops[i]
		}
	}

	return baseHop, nextHop, duration
}
