package main

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
)

const (
	port    = 33434
	maxHops = 30
)
func getIP(ht string) (string, error) {
	ip, err := net.LookupIP(ht)
	//print output of net.loolupip
	// fmt.Printf("net.LookupIP(%s) = %s, %v\n", ht, ip, err)

	if err != nil || len(ip) == 0 {
		return "", fmt.Errorf("unable to resolve IP address for %s", ht)
	}
	return ip[1].String(), nil
}


func main() {
  	destination := "www.google.com"
	var targetIP string
	if destination[0] == 'w' {
		ip, err := getIP(destination)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		targetIP = ip
	} else {
		targetIP = destination
	}
	//print target IP
	fmt.Printf("Target IP: %s\n", targetIP)
	timeToLive := 1
	for timeToLive < maxHops {
		ttl := timeToLive
		npingCommand := fmt.Sprintf("sudo nping --icmp -c 1 --ttl %d %s", ttl, destination)

		npingOutput, err := exec.Command("bash", "-c", npingCommand).Output()
		
		if err != nil {
			fmt.Println("Error executing nping:", err)
			break
		}

		ipMatches := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`).FindAllString(string(npingOutput), -1)
		rttMatches := regexp.MustCompile(`Avg rtt: [-+]?\d*\.\d+`).FindAllString(string(npingOutput), -1)

		if len(ipMatches) >= 3 && len(rttMatches) > 0 {
			hopIP := ipMatches[2]
			rtt := rttMatches[0][9:]
			fmt.Printf("Hop %d: ip = %s, rtt = %s ms\n", ttl, hopIP, rtt)
			if hopIP == targetIP {
				break
			}
		} else {
			fmt.Printf("Hop %d: *\n", ttl)
		}
		timeToLive++
	}
}
