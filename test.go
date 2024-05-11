package main

import (
	"fmt"
	// "net"
	// "log"
	// "github.com/google/gopacket"
	// "github.com/google/gopacket/pcap"
    // "github.com/google/gopacket/layers"
	// "github.com/google/gopacket/pcapgo"
	"os/exec"
	// "sync"
	"time"
	// "regexp"
	// "io/ioutil"
	// "os"
	// "./utils"
)
// import "github.com/manshi137/COD891/utils"



func pingWithTTL(ttl int, targetIP string) {
	numPacket := 1
	
	
	npingCommand := fmt.Sprintf("ping -n %d -i %d  %s", numPacket, ttl, targetIP)
	pingOutput1, err := exec.Command("cmd", "/C", npingCommand).Output()
	if err != nil {
		fmt.Println("Error executing ping:", err, ttl)
		return
	}
	fmt.Printf("ping with TTL %d, TARGETIP= %s \n", ttl, targetIP)
	fmt.Println("ping reply = ")
	fmt.Println(string(pingOutput1))
	// return 
}





func main() {
	targetIP := "2403:0:400:56::229"
	// targetIP2 := "115.113.240.203"
	ttl:=1
	pingWithTTL(ttl, targetIP)
	time.Sleep(1 * time.Second)
	// pingWithTTL(ttl, targetIP2)
	
}
// ping with TTL 1, TARGETIP= 115.113.240.203 
// ping reply =

// Pinging 115.113.240.203 with 32 bytes of data:
// Reply from 10.194.32.13: TTL expired in transit.

// Ping statistics for 115.113.240.203:
//     Packets: Sent = 1, Received = 1, Lost = 0 (0% loss),

// ---------------------------
// ping with TTL 1, numPacket= 1 
// ping reply = 

// Pinging 2403:0:400:56::216 with 32 bytes of data:
// Reply from 2403:0:400:56::216: TTL expired in transit.

// Ping statistics for 2403:0:400:56::216:
//     Packets: Sent = 1, Received = 1, Lost = 0 (0% loss),

