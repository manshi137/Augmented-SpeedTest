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
	// Open the pcap file
	handle, err := pcap.OpenOffline("capture.pcap")
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Set up packet decoding
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Create a CSV file for writing
	file, err := os.Create("icmp_packets.csv")
	if err != nil {
		log.Fatal("Error creating CSV file:", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Timestamp", "Type", "Code", "Identifier", "SequenceNum"})

	// Iterate through each packet in the pcap file
	for packet := range packetSource.Packets() {
		// Extract packet timestamp
		timestamp := packet.Metadata().Timestamp
		timestampFormatted := timestamp.Format("15:04:05.999999999")

		// Check if the packet is an ICMP packet
		networkLayer := packet.NetworkLayer()
		if networkLayer == nil {
			continue // Skip packet if no network layer
		}
		icmpLayer := packet.Layer(gopacket.LayerTypeICMPv4)
		if icmpLayer == nil {
			icmpLayer = packet.Layer(gopacket.LayerTypeICMPv6)
		}
		if icmpLayer == nil {
			continue // Skip packet if not ICMP
		}

		// Extract ICMP packet fields
		icmp := icmpLayer.(*gopacket.ICMPv4)
		if icmp == nil {
			icmp = icmpLayer.(*gopacket.ICMPv6)
		}
		if icmp != nil {
			// Create a slice to hold the CSV data
			data := []string{
				timestampFormatted,
				fmt.Sprintf("%d", int(icmp.TypeCode.Type())),
				fmt.Sprintf("%d", int(icmp.TypeCode.Code())),
				fmt.Sprintf("%d", int(icmp.Id)),
				fmt.Sprintf("%d", int(icmp.Seq)),
			}
			// Write data to CSV file
			if err := writer.Write(data); err != nil {
				log.Fatal("Error writing CSV data:", err)
			}
		}
	}
}
