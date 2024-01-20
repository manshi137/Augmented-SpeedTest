package main

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"sync"
	"time"
	"io/ioutil"
	"os"
)

const (
	port    = 33434
	maxHops = 10
	numThreads = 3
)
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
	npingCommand := fmt.Sprintf("sudo nping --tcp -c 20 --ttl %d --delay 0.1 %s", ttl, targetIP)
	npingOutput, err := exec.Command("bash", "-c", npingCommand).Output()
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
	var wg sync.WaitGroup

	wg.Add(1)
	go runNDT7Speedtest(&wg) // Run the ndt7-speedtest in a separate goroutine
	for i := 1; i <= numThreads; i++ { // Run pingWithTTL concurrently in numThreads goroutines
		wg.Add(1)
		go pingWithTTL(i, targetIP, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}
