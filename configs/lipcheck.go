package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	ipv4           bool
	ipv6           bool
	allowPrivate   bool
	allowLinkLocal bool
)

func param() {
	flag.BoolVar(&ipv4, "4", false, "ipv4")
	flag.BoolVar(&ipv6, "6", false, "ipv6")
	flag.BoolVar(&allowPrivate, "allowprivate", false, "allow private ip address")
	flag.BoolVar(&allowLinkLocal, "allowlinklocal", false, "allow link local ip address")
	flag.Parse()
}

func main() {
	param()
	netInterfaceAddresses, _ := net.InterfaceAddrs()
	for _, netInterfaceAddress := range netInterfaceAddresses {
		ipAddr, ok := netInterfaceAddress.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if ipAddr.IP.To4() != nil && Is4(ipAddr.IP) && ipv4 {
			if IsPrivate(ipAddr.IP) && !allowPrivate {
				continue
			}
			fmt.Println("global ipv4 " + ipAddr.IP.String())
		} else if ipAddr.IP.To16() != nil && Is6(ipAddr.IP) && ipv6 {
			if IsPrivate(ipAddr.IP) && !allowPrivate {
				continue
			}
			if IsLinkLocal(ipAddr.IP) && !allowLinkLocal {
				continue
			}
			fmt.Println("global ipv6 " + ipAddr.IP.String())
		}
	}
}

func Is4(ip net.IP) bool {
	//fmt.Printf("internal 0:%d 10:%d 11:%d\n", ip[0], ip[10], ip[11])
	return ip[0] == 0 &&
		ip[10] == 255 &&
		ip[11] == 255
}

func Is6(ip net.IP) bool {
	return !Is4(ip) && !Is0(ip)
}

func Is0(ip net.IP) bool {
	return ip[0] == 0 &&
		ip[1] == 0 &&
		ip[2] == 0 &&
		ip[3] == 0 &&
		ip[4] == 0 &&
		ip[5] == 0 &&
		ip[6] == 0 &&
		ip[7] == 0 &&
		ip[8] == 0 &&
		ip[9] == 0 &&
		ip[10] == 0 &&
		ip[11] == 0 &&
		ip[12] == 0 &&
		ip[13] == 0 &&
		ip[14] == 0 &&
		ip[15] == 0
}

func IsPrivate(ip net.IP) bool {
	if Is4(ip) {
		// Match the stdlib's IsPrivate logic.
		// RFC 1918 allocates 10.0.0.0/8, 172.16.0.0/12, and 192.168.0.0/16 as
		// private IPv4 address subnets.
		return ip[12] == 10 ||
			(ip[12] == 172 && ip[13]&0xf0 == 16) ||
			(ip[12] == 192 && ip[13] == 168)
	}
	if Is6(ip) {
		// RFC 4193 allocates fc00::/7 as the unique local unicast IPv6 address
		// subnet.
		// fmt.Printf("internal %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d\n", ip[0], ip[1], ip[2], ip[3], ip[4], ip[5], ip[6], ip[7], ip[8], ip[9], ip[10], ip[11], ip[12], ip[13], ip[14], ip[15])
		return ip[0]&0xfe == 0xfc
	}

	return false
}

func IsLinkLocal(ip net.IP) bool {
	if Is6(ip) {
		// RFC 4193 allocates fc00::/7 as the unique local unicast IPv6 address
		// subnet.
		// fmt.Printf("internal %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d\n", ip[0], ip[1], ip[2], ip[3], ip[4], ip[5], ip[6], ip[7], ip[8], ip[9], ip[10], ip[11], ip[12], ip[13], ip[14], ip[15])
		return ip[0] == 0xfe && ip[1] == 0x80
	}

	return false
}
