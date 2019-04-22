## structure and architecture

the application consiss of a main trace funcion that collects all the hops from local to destination IP, and returns it. 
it uses a few utility functions in utils.go file.

in utils we have facilities to get local IP addresses, and also to resolve remote host to it's IP address.

the implementation is based on UDP based traceroute. it is still the default behaviour on POSIX systems, we could use ICMPEcho based solution though.

## trace function

the act of tracing is to start from local Ip and travel through all routers to get to destination IP.
to achive that goal we use two mechanisms in IP networking.

1. we set the TTL fileld of IP header to specify how many hops we intend the packet to travel. we start from 1 and go upward. each router decrease it by one, and when it reaches to zero, the router discards the packet. so when for example we set th ttl to 5, it could ravel only through 5 routers.
2. when a router discards a packet because of TTL reaching zero, it emits an ICMP "time exceeded" message to origina IP.

now the trick is simple. we set TTL to 1 and then wait for ICMP message of first router, from the ICMP message we aquire the IP address of the router and we can do a reverse DNS to get its hostname. then we se it to 2 to get information of second router, and so on, until the router IP is same as destination IP address.


## limtations

I did not implement some features because of lack of time, so now we have following limitaions. but implementig them is easy!

1. no IPv6 for now.
2. we use the first locl IP of system that is not a loopback. we can implement a mechanism to use other IP addresses if there is no route form that IP, and also if other IPs have better route to destination.
3. we use the first remote IP (A record) of destination host, agian it is possible that other IP addreses have a better path to destination
4. for reverse DNS we use first result, it is just to show the result, there could be more than one DNS associated with that IP.
5. we probe the routers sequentially for simplicity, but we can run multiple probes in parallel using multiple ports for faster response times.

## how to use it

the application is written in Golang and compiled with Go 1.12 under linux. there is no dependancy. just compile the application and invoke the app with following parameters:

*  --host (mandatory) the host name to find a route to
*  --port (optional) the port to which app should bind
*  --max_hops (optional) maximum number of hops
*  --packet_size (optioanal) size of UDP packet to send to destination

