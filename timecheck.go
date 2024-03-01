package main
 
import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"io/ioutil"
	"strings"
	"time"
)
var timeArray [2]string
func writeArrayToFile(filename string) error {
	// Create a new file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Write each element of the array to the file
	for _, item := range timeArray {
		_, err := file.WriteString(item + "\n")
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	return nil
}


func main() {
	// Open the CSV file
	file, err := os.Open("filtered_ndtcapture.csv")
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
	filePath := "ip_addresses.txt"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	ipAddresses := strings.Split(string(content), "\n")
	
	timeArray = [2]string{"", ""}
	serverIP := ipAddresses[6]
	localIPv4 := ipAddresses[4]
	serverIP = strings.TrimSpace(serverIP)
	localIPv4 = strings.TrimSpace(localIPv4)
	var prevTime time.Time



	for {
		// Read the next record
		record, err := reader.Read()
		if err != nil {
			break // End of file
		}
 
		// Extract source and destination IPs from the record
		ip1 := record[1] //
		ip2 := record[2]
		
		// Increment counters based on packet direction
		pktSize, _ := strconv.Atoi(record[3])

		if ip1 == serverIP && ip2 == localIPv4 {
			download+= pktSize
		} else if ip1 == localIPv4 && ip2 == serverIP {
			upload+= pktSize
		}
		// Increment packet count
		packetCount++



		timestamp := record[0]

		// Convert timestamp string to time.Time object
		if timestamp != "Timestamp"{
			currentTime, err := time.Parse("15:04:05.999999", timestamp)
			if err != nil {
				fmt.Println("Error parsing timestamp:", err)
				continue
			}
			if prevTime.IsZero(){
				prevTime = currentTime
			}

			diff := currentTime.Sub(prevTime)

 
			// Check every 100 packets
			if diff.Milliseconds() >= 100 {
				if download < upload {
					if(timeArray[0] == "" ){
						timeArray[0] = record[0]
					}
					timeArray[1] = record[0]
				}
				// Reset counters
				download = 0
				upload = 0
				prevTime= currentTime
				packetCount=0
			}
		}
		
		
	}
	

	err1 := writeArrayToFile("times.txt")
	if err1 != nil {
		fmt.Println("Error:", err)
	}
}

