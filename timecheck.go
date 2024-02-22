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
	filePath := "ip_addresses.txt"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	ipAddresses := strings.Split(string(content), "\n")
	// fmt.Println("IP Addresses:")
	timeArray = [2]string{"", ""}
	serverIP := ipAddresses[6]
	localIPv4 := ipAddresses[4]
	serverIP = strings.TrimSpace(serverIP)
	localIPv4 = strings.TrimSpace(localIPv4)
	// fmt.Println("record[6]= ", ipAddresses[6])
	fmt.Println("serverIP=", serverIP, "h")
	fmt.Println("localipv4=", localIPv4, "h")
	// Iterate over each line in the CSV file

	var prevTime time.Time
	fmt.Println("prev time 1= ", prevTime)


	for {
		// Read the next record
		record, err := reader.Read()
		if err != nil {
			break // End of file
		}
 
		// Extract source and destination IPs from the record
		ip1 := record[1] //
		ip2 := record[2]
		// fmt.Println("ip1= ", ip1, "ip2 = ", ip2)
		
		// Increment counters based on packet direction
		pktSize, _ := strconv.Atoi(record[3])

		

		// fmt.Println("pktsize = ", record[3], pktSize)

		if ip1 == serverIP && ip2 == localIPv4 {
			download+= pktSize
			// fmt.Println("download = ", download)
		} else if ip1 == localIPv4 && ip2 == serverIP {
			upload+= pktSize
			// fmt.Println("upload = ", upload)
		}
		// Increment packet count
		packetCount++

		if packetCount==3{
			fmt.Println("len of ip1= ", len(ip1))
			fmt.Println("len of serverip= ", len(serverIP))
			ip1 = strings.TrimSpace(ip1)

			fmt.Println("ip1=", ip1, "h")
			fmt.Println("ip2=", ip2, "h")
		}

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
				fmt.Println("prevTime  2= ", prevTime)
			}

			diff := currentTime.Sub(prevTime)
			// fmt.Println("prevTime =", prevTime, " current time= ", currentTime, " diff = ", diff)

 
			// Check every 100 packets
			if diff.Milliseconds() >= 100 {
				// fmt.Println("hello")
				// Check if IP1 to IP2 packets are in majority
				// fmt.Println("upload = ", upload, " download = ", download)
				// fmt.Println("upload = ", upload, "  download= ", download)
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
				prevTime= currentTime
				fmt.Println("packetcount= ", packetCount)
				packetCount=0
			}
		}
		

		// Update prevTime to the timestamp from the first record if it's zero value
		

		// Calculate time difference between current and previous timestamps
		
	}
	

	err1 := writeArrayToFile("times.txt")
	if err1 != nil {
		fmt.Println("Error:", err)
	}
}

