package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"os/exec"
)

// Function to modify the timestamp in the output
func modifyTimestamp(npingOutputStr string) string {
	// Extract timestamp from the output using regular expressions
	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2})`)
	matches := re.FindStringSubmatch(npingOutputStr)
	if len(matches) != 2 {
		fmt.Println("Timestamp not found in Nping output")
		return npingOutputStr
	}

	// Parse the timestamp string into a time.Time object
	timestampStr := matches[1]
	timestamp, err := time.Parse("2006-01-02 15:04", timestampStr)
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return npingOutputStr
	}

	// Format the timestamp with seconds
	formattedTimestamp := timestamp.Format("2006-01-02 15:04:05")

	// Replace the original timestamp in the output string with the modified one
	modifiedOutput := strings.Replace(npingOutputStr, timestampStr, formattedTimestamp, 1)
	return modifiedOutput
}
func main() {
	// Example output of Nping
	numPacket := 2.0
	delay := 0.1
	timeout := numPacket * delay 
	ttl:= 1
	fmt.Printf("ping with TTL %d and timeout %f\n", ttl, timeout)
	targetIP := "10.255.107.3"
	// startTime := time.Now()
	// npingCommand := fmt.Sprintf("sudo ping -c 20 -t %d -i 0.1 %s", ttl, targetIP)
	npingCommand := fmt.Sprintf("sudo nping --tcp -c %d --ttl %d --delay %f %s", int(numPacket), ttl, delay, targetIP)
	// npingCommand := fmt.Sprintf("sudo nping --tcp -c 20 --ttl %d --delay 0.1 %s", ttl, targetIP)
	npingOutput, err := exec.Command("bash", "-c", npingCommand).Output()
	if err != nil {
		fmt.Println("Error executing nping:", err)
		return
	}
	// Convert npingOutput to a string
	npingOutputStr := string(npingOutput)

	// Modify timestamp in the output
	modifiedOutput := modifyTimestamp(npingOutputStr)

	// Print the modified output
	fmt.Println("Modified output:")
	fmt.Println(modifiedOutput)
}
