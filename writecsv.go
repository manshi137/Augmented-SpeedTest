package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/layers"
)

func firstCSV(pcapFilePath, csvFilePath string) {
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

		// Format timestamp with milliseconds
		timestampWithMilliseconds := timestamp.Format("15:04:05.999999999")

		// Write data to CSV
		err := csvWriter.Write([]string{
			timestampWithMilliseconds,
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

}




func secondCSV(pcapFilePath, csvFilePath string) {
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
	err = csvWriter.Write([]string{"Timestamp", "SourceIP", "DestinationIP", "Length", "SourcePort", "DestinationPort", "IPProtocol", "ICMPType", "ICMPPayload"})
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

		// Additional headers for specific protocols
		ipProtocol := "Unknown"
		icmpType := ""
		icmpPayload := ""

		// Check if the network layer is IPv4
		if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
			ipProtocol = "IPv4"
			// Check if the transport layer is ICMP
			if icmpLayer := packet.Layer(layers.LayerTypeICMPv4); icmpLayer != nil {
				icmp, _ := icmpLayer.(*layers.ICMPv4)
				icmpType = fmt.Sprintf("%d", icmp.TypeCode.Type())
				icmpPayload = fmt.Sprintf("%v", icmp.Payload)
			}
		}

		// Write data to CSV
		err := csvWriter.Write([]string{
			timestamp.Format(time.RFC3339),
			sourceIP,
			destinationIP,
			fmt.Sprint(length),
			sourcePort,
			destinationPort,
			ipProtocol,
			icmpType,
			icmpPayload,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

}

func main() {
	pcapFilePath := "capture.pcap"


	csvFilePath1 := "ndtcapture.csv"
	// csvFilePath2 := "output2.csv"

	// Call the function to write packets to CSV
	firstCSV(pcapFilePath, csvFilePath1)
	// secondCSV(pcapFilePath, csvFilePath2)
}
