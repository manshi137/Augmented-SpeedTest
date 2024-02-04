package main
 
import (
	"encoding/csv"
	"fmt"
	"os"
	// "strconv"
)
 
func main() {
	// Open the CSV file
	file, err := os.Open("filteredndtcap.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
 
	// Create a CSV reader
	reader := csv.NewReader(file)
 
	// Initialize counters
	download := 0
	upload := 0
	packetCount := 0
 
	// Iterate over each line in the CSV file
	for {
		// Read the next record
		record, err := reader.Read()
		if err != nil {
			break // End of file
		}
 
		// Extract source and destination IPs from the record
		ip1 := record[1] //
		ip2 := record[2]
		serverIP := "34.131.204.19"
		localIP := "10.184.52.48"
		// Increment counters based on packet direction
		if ip1 == serverIP && ip2 == localIP {
			download++
		} else if ip1 == localIP && ip2 == serverIP {
			upload++
		}
 
		// Increment packet count
		packetCount++
 
		// Check every 100 packets
		if packetCount%500 == 0 {
			// Check if IP1 to IP2 packets are in majority
			if download < upload {
				fmt.Println("upload packets are in majority at this timestamp = ",record[0], " upload packets= ", upload, " download packets= ", download )
				// Record timestamp or perform other actions
			}
			// Reset counters
			download = 0
			upload = 0
		}
	}
}

