package main

import (
	"errors"
	"fmt"
	"net"
)

func getLocalIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("can not get local IP addresses: %s", err)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}

	}
	return nil, errors.New("no Local IPV4 avaiable")
}

func getRemoteIP(dest string) (net.IP, error) {
	destAddrs, err := net.LookupHost(dest)
	if err != nil {
		return nil, errors.New("can not resolve destination host")
	}
	//now we have a list of IPs, both v4 and v6, we select the first v4 one
	for _, ip := range destAddrs {
		resolvedIP, err := net.ResolveIPAddr("ip", ip)
		if err == nil && resolvedIP.IP.To4() != nil {
			return resolvedIP.IP, nil
		}
	}
	return nil, errors.New("no remote IPV4 avaiable for host")

}

func getRawIPV4(ip net.IP) [4]byte {
	output := [4]byte{}
	copy(output[:], ip.To4())
	return output
}

func IPv4ToStr(rawIP [4]byte) string {
	return fmt.Sprintf("%v.%v.%v.%v", rawIP[0], rawIP[1], rawIP[2], rawIP[3])
}
