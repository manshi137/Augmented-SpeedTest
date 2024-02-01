package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func main() {
	// Replace "your_input.pcap" with the path to your pcap file
	pcapFilePath := "C:/Users/richi/OneDrive/Documents/COD891/COD891/capture.pcap"

	// Replace "output.csv" with the desired output CSV file path
	csvFilePath := "output.csv"

	// Open the pcap file
	handle, err := pcap.OpenOffline(pcapFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Open the CSV file for writing
	csvFile, err := os.Create(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	// Create a CSV writer
	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// Write CSV header
	err = csvWriter.Write([]string{"Timestamp", "SourceIP", "DestinationIP", "Length", "SourcePort", "DestinationPort"})
	if err != nil {
		log.Fatal(err)
	}

	// Packet source
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Iterate through packets
	for packet := range packetSource.Packets() {
		// Extract relevant information from the packet
		timestamp := packet.Metadata().Timestamp
		networkLayer := packet.NetworkLayer()
		transportLayer := packet.TransportLayer()

		if networkLayer == nil || transportLayer == nil {
			continue // Skip packets without both network and transport layers
		}

		sourceIP := networkLayer.NetworkFlow().Src().String()
		destinationIP := networkLayer.NetworkFlow().Dst().String()
		length := len(packet.Data())
		sourcePort := transportLayer.TransportFlow().Src().String()
		destinationPort := transportLayer.TransportFlow().Dst().String()

		// Write data to CSV
		err := csvWriter.Write([]string{
			timestamp.Format(time.RFC3339),
			sourceIP,
			destinationIP,
			fmt.Sprint(length),
			sourcePort,
			destinationPort,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("CSV file created successfully.")
}
