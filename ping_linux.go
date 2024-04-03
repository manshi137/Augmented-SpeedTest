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
	
)

const (
	numThreads = 3
)
var stopPingFlag bool
var stopPingMutex sync.Mutex
var ipAddressArray [numThreads+1+3]string
var locks [numThreads+1+3]sync.Mutex

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
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Default network interface:")
	fmt.Printf("Name: %s\n", iface.Name)


	if err != nil {
	  fmt.Println("Failed to get default interface:", err)
	  return ""
	}
  
	capture_filter := filter_map[test_name]
	// desiredFriendlyName := "Intel(R) Wi-Fi 6E AX211 160MHz"

	// // Find the corresponding device name for the given friendly name
	// interfaces, err := pcap.FindAllDevs()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("List of network interfaces:")
	// for i, dev := range interfaces {
	// 	fmt.Printf("[%d] Name: %s\n", i+1, dev.Name)
	// 	fmt.Printf("    Description: %s\n", dev.Description)
	// 	fmt.Printf("    Flags: %v\n", dev.Flags)
	// 	fmt.Printf("    Addresses:\n")
	// 	for _, address := range dev.Addresses {
	// 		fmt.Printf("        %s\n", address)
	// 	}
	// 	fmt.Println()
	// }
	// var desiredDeviceName string
	// for _, dev := range interfaces {
	// 	if dev.Description == desiredFriendlyName {
	// 		desiredDeviceName = dev.Name
	// 		break
	// 	}
	// }

	// // Check if the desired device name was found
	// if desiredDeviceName == "" {
	// 	log.Fatal("Desired network interface not found")
	// }
	

	handle, err := pcap.OpenLive(iface.Name, snaplen, false, pcap.BlockForever)
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

func runping(ch chan<- string, npingCommand string, ttl int) {
    time.Sleep(1 * time.Millisecond)

    // cmd := exec.Command("cmd", "/C", npingCommand)
	cmd := exec.Command("bash", "-c", npingCommand)

    pingOutput, err := cmd.CombinedOutput()

    if err != nil {
        fmt.Println("runping error:", err)
        ch <- fmt.Sprintf("Error: %v", err)
        return
    }
	ipv4Regex := regexp.MustCompile(`(?:[0-9]{1,3}\.){3}[0-9]{1,3}`)
	ipv4Matches := ipv4Regex.FindAllString(string(pingOutput), -1)
	if len(ipv4Matches) > 0 {
		locks[ttl].Lock()
		ipAddressArray[ttl]=ipv4Matches[1];
		locks[ttl].Unlock()
	}
	ipv6Regex := regexp.MustCompile(`(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}`)
	ipv6Matches := ipv6Regex.FindAllString(string(pingOutput), -1)
	if len(ipv6Matches) > 0 {
		locks[ttl].Lock()
		ipAddressArray[ttl]=ipv6Matches[1];
		locks[ttl].Unlock()
	}

    ch <- string(pingOutput)
}

func pingWithTTL(ttl int, targetIP string, wg *sync.WaitGroup) {
	defer wg.Done()
	numPacket := 1
	
	startTime := time.Now()
	npingCommand := fmt.Sprintf("ping -n %d -i %d  %s", numPacket, ttl, targetIP)
	interval := 100*time.Millisecond
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	// var pingOutput []byte
	for {
		select {
		case <-ticker.C:
			if getStopPingFlag() {
				endTime := time.Now()
				fmt.Println("--------------------------------------------------")
				duration := endTime.Sub(startTime)
				fmt.Printf("Execution Time of ping: %v , ttl= %d \n", duration, ttl)
				fmt.Println("--------------------------------------------------")
				
				return
			}
			ch := make(chan string)
			go runping(ch, npingCommand, ttl)
		}
	}
}



func runNDT7Speedtest(wg *sync.WaitGroup) {
	defer wg.Done()
	startTime := time.Now()
	fmt.Println("speedtest start at: ",startTime)
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
	endTime := time.Now()
	fmt.Println("speedtest end at: ",endTime)
	fmt.Println("Done running ndt7-speedtest.")
}

func capturePacket(test_name string, filter_map map[string]string, time_sec int, wg *sync.WaitGroup) {
	defer wg.Done()

	// packet capture params
	var snaplen int32 = 1600
	iface, err := GetDefaultInterface()
	if err != nil {
		fmt.Println("Failed to get default interface:", err)
		return
	}

	capture_filter := filter_map[test_name]

	// desiredFriendlyName := "Intel(R) Wi-Fi 6E AX211 160MHz"

	// // Find the corresponding device name for the given friendly name
	// interfaces, err := pcap.FindAllDevs()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// var desiredDeviceName string
	// for _, dev := range interfaces {
	// 	if dev.Description == desiredFriendlyName {
	// 		desiredDeviceName = dev.Name
	// 		break
	// 	}
	// }

	// // Check if the desired device name was found
	// if desiredDeviceName == "" {
	// 	log.Fatal("Desired network interface not found")
	// }
	// handle, err := pcap.OpenLive(desiredDeviceName, snaplen, false, pcap.BlockForever)

	handle, err := pcap.OpenLive(iface.Name, snaplen, false, pcap.BlockForever)
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

	for packet := range packetSource.Packets() {
		if getStopPingFlag() {
			return
		}

		// Write the packet to the pcap file
		err := pcapWriter.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		if err != nil {
			fmt.Println("Error writing packet to pcap file:", err)
		}
	}
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
	// wait for 10 more seconds and then stop the pingWithTTL threads
	time.Sleep(10 * time.Second)
	setStopPingFlag(true)
	//stop nping now
	

	wg2.Wait()
	endTime:= time.Now()
	fmt.Println("Pings ending at: ", endTime)
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

	
}