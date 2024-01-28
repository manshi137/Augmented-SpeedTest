package main

import (
	"fmt"
	"net"
	"log"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
    "github.com/google/gopacket/layers"
	// "ndt7module"
	"os/exec"
	"regexp"
	"sync"
	"time"
	"io/ioutil"
	"os"
	// "./utils"
)
// import "github.com/manshi137/COD891/utils"
const (
	port    = 33434
	maxHops = 10
	numThreads = 3
)


func find_server(test_name string, filter_map map[string]string) string {
	localIPv4 := GetLocalIP("v4")
	localIPv6 := GetLocalIP("v6")
	//packet capture params
	var snaplen int32 = 96
	num_pkts := 0
	iface, err := GetDefaultInterface()
	fmt.Println(iface)

	fmt.Println("Interface is ", iface.Name)
	if err != nil {
	  fmt.Println("Failed to get default interface:", err)
	  return ""
	}
  
	capture_filter := filter_map[test_name]
	desiredFriendlyName := "Intel(R) Wi-Fi 6E AX211 160MHz"

	// Find the corresponding device name for the given friendly name
	interfaces, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	var desiredDeviceName string
	for _, dev := range interfaces {
		if dev.Description == desiredFriendlyName {
			desiredDeviceName = dev.Name
			break
		}
	}

	// Check if the desired device name was found
	if desiredDeviceName == "" {
		log.Fatal("Desired network interface not found")
	}

	fmt.Println("Desired device name= ", desiredDeviceName)
	handle, err := pcap.OpenLive(desiredDeviceName, snaplen, false, pcap.BlockForever)
	if err != nil {
	  log.Fatal(err)
	}
  
	defer handle.Close()
  
	// Set the capture filter
	err = handle.SetBPFFilter(capture_filter)
	if err != nil {
	  log.Fatal(err)
	}
  
	ipCountMap := make(map[string]int)
	var localIP string
	// Start capturing packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	var sourceIP string
	var destIP string
	var serverIP string
	for packet := range packetSource.Packets() {
  
	  if IsIPv6Packet(packet) {
		ipPacket, _ := packet.Layer(layers.LayerTypeIPv6).(*layers.IPv6)
		sourceIP = ipPacket.SrcIP.String()
		destIP = ipPacket.DstIP.String()
		localIP = localIPv6
	  } else {
		 ipPacket, _ := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
		 sourceIP = ipPacket.SrcIP.String()
		 destIP = ipPacket.DstIP.String()
		 localIP = localIPv4
	  }
	  //fmt.Println(localIP, srcIP, destIP)
	  if sourceIP == localIP {
		serverIP = destIP
	  } else {
		serverIP = sourceIP
	  }
  
	  if _, ok := ipCountMap[serverIP]; !ok {
		ipCountMap[serverIP] = 0
	  }
	  ipCountMap[serverIP] += 1
	  // Process captured packet
	  num_pkts += 1
	  if num_pkts == 1000 {
		break
	  }
	}
	serverIPMax := GetKeyWithMaxValue(ipCountMap)
	return serverIPMax
}

func getIP(ht string) (string, error) {
	ip, err := net.LookupIP(ht)
	if err != nil || len(ip) == 0 {
		return "", fmt.Errorf("unable to resolve IP address for %s", ht)
	}
	return ip[1].String(), nil
}

func pingWithTTL(ttl int, targetIP string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Ping with TTL %d to target IP: %s\n", ttl, targetIP)
	startTime := time.Now()
	// npingCommand := fmt.Sprintf("sudo ping -c 20 -t %d -i 0.1 %s", ttl, targetIP)
	npingCommand := fmt.Sprintf("nping --tcp -c 20 --ttl %d --delay 0.1 %s", ttl, targetIP)
	npingOutput, err := exec.Command("cmd", "/C", npingCommand).Output()
	if err != nil {
		fmt.Println("Error executing nping:", err)
		return
	}
	endTime := time.Now()
	//print npingOutput
	// fmt.Printf("%s\n", npingOutput)
	fmt.Println("--------------------------------------------------")
	duration := endTime.Sub(startTime)
	fmt.Printf("Execution Time: %v\n", duration)
	fmt.Println("--------------------------------------------------")

	outputFileName := fmt.Sprintf("output%d.txt", ttl)
	err = ioutil.WriteFile(outputFileName, npingOutput, 0644)//write npingOutput to file
	if err != nil {
		fmt.Printf("Error writing to %s: %v\n", outputFileName, err)
	}

	ipMatches := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`).FindAllString(string(npingOutput), -1)
	rttMatches := regexp.MustCompile(`Avg rtt: [-+]?\d*\.\d+`).FindAllString(string(npingOutput), -1)
	
	if len(ipMatches) >= 3 && len(rttMatches) > 0 {
		hopIP := ipMatches[2]
		rtt := rttMatches[0][9:]
		fmt.Printf("Hop %d: ip = %s, rtt = %s ms\n", ttl, hopIP, rtt)
		if hopIP == targetIP {
			fmt.Println("Target IP reached!")
			return
		}
	} else {
		fmt.Printf("Hop %d: *\n", ttl)
	}
}

func runNDT7Speedtest(wg *sync.WaitGroup) {
	defer wg.Done()

	// Replace "ndt7-speedtest" with the actual path or command you want to run
	cmd := exec.Command("ndt7-client")

	// Set the standard output and error to os.Stdout and os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running ndt7-speedtest:", err)
	}
}

func main() {

	var filter_map = map[string]string {
		"mlab": "port 443",
		"ookla": "port 8080 or port 5060",
	  }
	
	var test_name = "mlab"
	targetIP := find_server(test_name, filter_map)
	//print target IP
	fmt.Printf("Target IP = ndt7 server: %s\n", targetIP)
	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup

	wg1.Add(1)
	go runNDT7Speedtest(&wg1) // Run the ndt7-speedtest in a separate goroutine

	// write a function that infer the ndt server; 
	// Check this code: https://github.com/tarunmangla/speedtest-diagnostics/blob/master/tslp/tslp.go#L22
	// findserver function


	// start capturing the packets and store them in a file


	for i := 1; i <= numThreads; i++ { // Run pingWithTTL concurrently in numThreads goroutines
		wg2.Add(1)
		go pingWithTTL(i, targetIP, &wg2)
	}

	// process 
	// Wait for all goroutines to finish
	wg1.Wait() // stop ndt7test

	// wait for 10 more seconds and then stop the pingWithTTL threads 

	// process the pcap file: 1) find out the ping RTTs; 
	// 2) find out the end time for download and the end time of the test;
	// 3) run t-test on the ping data
	
}
