package main
 
import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	// "time"
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

	fmt.Println("Upload / Download time written to times.txt successfully!")
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
		serverIP := "115.113.240.203"
		localIP := "10.184.59.62"
		// Increment counters based on packet direction
		pktSize, _ := strconv.Atoi(record[3])
		// fmt.Println("pktsize = ", record[3], pktSize)

		if ip1 == serverIP && ip2 == localIP {
			download+= pktSize
			// fmt.Println("download = ", download)
		} else if ip1 == localIP && ip2 == serverIP {
			upload+= pktSize
			// fmt.Println("upload = ", upload)
		}
		// Increment packet count
		packetCount++
 
		// Check every 100 packets
		if packetCount%500 == 0 {
			// Check if IP1 to IP2 packets are in majority
			// fmt.Println("upload = ", upload, " download = ", download)
			if download < upload {
				if(timeArray[0] == "" ){
					timeArray[0] = record[0]
				}
				timeArray[1] = record[0]
				fmt.Println("upload packets are in majority at this timestamp = ", record[0], " upload packets = ", upload, " download packets = ", download )
				// Record timestamp or perform other actions
			}
			// Reset counters
			download = 0
			upload = 0
		}
	}
	timeArray = [2]string{"", ""}

	err1 := writeArrayToFile("times.txt")
	if err1 != nil {
		fmt.Println("Error:", err)
	}
}

