package main

import (
	"fmt"
	// "net"
	"log"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
    "github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"os/exec"
	"sync"
	"time"
	"regexp"
	// "io/ioutil"
	"os"
	// "./utils"
)
// import "github.com/manshi137/COD891/utils"
const (
	numThreads = 3
)
var stopPingFlag bool
var stopPingMutex sync.Mutex
var ipAddressArray [numThreads+1+3]string

func setStopPingFlag(value bool) {
	stopPingMutex.Lock()
	stopPingFlag = value
	stopPingMutex.Unlock()
}

func getStopPingFlag() bool {
	stopPingMutex.Lock()
	defer stopPingMutex.Unlock()
	return stopPingFlag
}

func find_server(test_name string, filter_map map[string]string, wg *sync.WaitGroup) string {
	defer wg.Done()
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

func pingWithTTL(ttl int, targetIP string, wg *sync.WaitGroup) {
	defer wg.Done()
	numPacket := 1
	
	fmt.Printf("ping with TTL %d, numPackat= %d \n", ttl, int(numPacket))

	startTime := time.Now()
	npingCommand := fmt.Sprintf("ping -n %d -i %d  %s", numPacket, ttl, targetIP)

	interval := 1000*time.Millisecond
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	var pingOutput []byte
	for {
		select {
		case <-ticker.C:
			if getStopPingFlag() {
				fmt.Println("Stopping continuous ping due to StopFlag...")
				endTime := time.Now()
				fmt.Println("--------------------------------------------------")
				duration := endTime.Sub(startTime)
				fmt.Printf("Execution Time of ping: %v , ttl= %d \n", duration, ttl)
				fmt.Println("--------------------------------------------------")
				ipRegex := regexp.MustCompile(`(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|(?:[0-9]{1,3}\.){3}[0-9]{1,3}`)

				//get ip of hop 
				ipMatches := ipRegex.FindAllString(string(pingOutput), -1)
				fmt.Println("IP Address:", ipMatches[1], " ttl= ", ttl)
				ipAddressArray[ttl]=ipMatches[1];
				return
			}
			pingOutput1, err := exec.Command("cmd", "/C", npingCommand).Output()
			if err != nil {
				fmt.Println("Error executing ping:", err, ttl)
				return
			}
			pingOutput = pingOutput1
		}
	}
}



func runNDT7Speedtest(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Running ndt7-speedtest...")
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
	fmt.Println("Done running ndt7-speedtest.")
}

func capturePacket(test_name string, filter_map map[string]string, time_sec int, wg *sync.WaitGroup) {
	defer wg.Done()

	// packet capture params
	fmt.Println("Starting capturepackets...")
	var snaplen int32 = 1600
	iface, err := GetDefaultInterface()
	fmt.Println("Interface is ", iface.Name)
	if err != nil {
		fmt.Println("Failed to get default interface:", err)
		return
	}

	capture_filter := filter_map[test_name]

	fmt.Println("Capture filter is ", capture_filter)
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
	handle, err := pcap.OpenLive(desiredDeviceName, snaplen, false, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}

	defer handle.Close()

	// Set the capture filter
	err = handle.SetBPFFilter(capture_filter + " or icmp") // only capture packets from/to "port 443"
	if err != nil {
		log.Fatal(err)
	}

	outputFileName := fmt.Sprintf("capture.pcap")
	file, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Error creating pcap file:", err)
		return
	}
	defer file.Close()

	// Create a pcapgo writer
	pcapWriter := pcapgo.NewWriter(file)
	pcapWriter.WriteFileHeader(1600, layers.LinkTypeEthernet)

	// Start capturing packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType()) //packetSource.Packets() is a channel

	fmt.Println("Start capturing packets...")
	for packet := range packetSource.Packets() {
		if getStopPingFlag() {
			fmt.Println("Stopping capturePacket due to StopFlag...")
			return
		}

		// Write the packet to the pcap file
		err := pcapWriter.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		if err != nil {
			fmt.Println("Error writing packet to pcap file:", err)
		}
	}
	fmt.Println("Done capturepackets w/o StopFlag.")
}



func main() {
	var filter_map = map[string]string {
		"mlab": "port 443",
		"ookla": "port 8080 or port 5060",
	}
	
	//print target IP
	var wg1 sync.WaitGroup
	var wg3 sync.WaitGroup
	wg1.Add(1)
	go runNDT7Speedtest(&wg1) // Run the ndt7-speedtest in a separate goroutine
	var test_name = "mlab"
	wg3.Add(1)
	targetIP := find_server(test_name, filter_map, &wg3)
	fmt.Printf("Target IP: %s\n", targetIP)

	var wg2 sync.WaitGroup


	// write a function that infer the ndt server; 
	// Check this code: https://github.com/tarunmangla/speedtest-diagnostics/blob/master/tslp/tslp.go#L22
	// findserver function


	// start capturing the packets and store them in a file
	wg2.Add(1)
	go capturePacket(test_name, filter_map, 10, &wg2)

	wg3.Wait()
	for i := 1; i <= numThreads; i++ { // Run pingWithTTL concurrently in numThreads goroutines
		wg2.Add(1)
		go pingWithTTL(i, targetIP, &wg2)
	}

	wg1.Wait()
	fmt.Println("Done ndt7test....")
	// wait for 10 more seconds and then stop the pingWithTTL threads
	fmt.Println("Wait for 10 seconds...") 
	time.Sleep(10 * time.Second)
	setStopPingFlag(true)
	//stop nping now
	

	wg2.Wait()
	fmt.Println("Done pingWithTTL and capturePacket....")
	localIPv4 := GetLocalIP("v4")
	localIPv6 := GetLocalIP("v6")
	fmt.Println("Local IPv4: ", localIPv4)
	fmt.Println("Local IPv6: ", localIPv6)
	fmt.Println("Target IP: ", targetIP)
	ipAddressArray[numThreads+1] = localIPv4
	ipAddressArray[numThreads+2] = localIPv6
	ipAddressArray[numThreads+3] = targetIP
	// process the pcap file: 
	// 1) find out the ping RTTs; 
	// 2) find out the end time for download and the end time of the test;
	// 3) run t-test on the ping data
	// Remove the file if it exists
	filePath := "ip_addresses.txt"
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			fmt.Println("Error deleting file:", err)
			return
		}
	}

	// Open a file for appending (create if it doesn't exist)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Append each IP address to the file
	for _, ip := range ipAddressArray {
		// Append IP address followed by a newline to the file
		if _, err := file.WriteString(ip + "\n"); err != nil {
			fmt.Println("Error appending to file:", err)
			return
		}
	}

	fmt.Println("IP addresses appended to ip_addresses.txt successfully.")
	
}