package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func filterCSV(inputCSV, outputCSV, sourceIP, destinationIP string) error {
	// Open the input CSV file
	inputFile, err := os.Open(inputCSV)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Open the output CSV file
	outputFile, err := os.Create(outputCSV)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Create a CSV reader for the input file
	reader := csv.NewReader(inputFile)
	// Create a CSV writer for the output file
	writer := csv.NewWriter(outputFile)

	// Read and write the header
	header, err := reader.Read()
	if err != nil {
		return err
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Column indices for source IP and destination IP
	sourceIPIndex := 1
	destinationIPIndex := 2

	// Find the column indices for source IP and destination IP
	for i, col := range header {
		if col == "SourceIP" {
			sourceIPIndex = i
		} else if col == "DestinationIP" {
			destinationIPIndex = i
		}
	}

	if sourceIPIndex == -1 || destinationIPIndex == -1 {
		return fmt.Errorf("required columns not found in the CSV file")
	}

	// Iterate through rows, filter based on conditions, and write to output CSV
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		if	((row[sourceIPIndex] == sourceIP && row[destinationIPIndex] == destinationIP) ||
				(row[destinationIPIndex] == sourceIP && row[sourceIPIndex] == destinationIP)) {
			if err := writer.Write(row); err != nil {
				return err
			}
		}
	}

	// Flush the writer to ensure all data is written to the file
	writer.Flush()

	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}

func main() {
	inputCSV := "ndtcapture.csv"
	outputCSV := "filtered_ndtcapture.csv"
	serverIP := "115.113.240.203"
	localIP := "10.184.59.62"

	err := filterCSV(inputCSV, outputCSV, serverIP, localIP)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Filtered CSV file created successfully.")
	}
}
