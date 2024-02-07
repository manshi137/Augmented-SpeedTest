package main

import (
	"fmt"
	"log"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/layers"
	"encoding/binary"
)

func main() {
	// Open the pcap file
	handle, err := pcap.OpenOffline("capture.pcap")
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Create a packet source from the handle
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Map to store ICMP Echo Request packets by their identifier and sequence number
	echoRequests := make(map[string]gopacket.Packet)
	echoReply := make(map[string]gopacket.Packet)

	// Loop through each packet in the pcap file
	for packet := range packetSource.Packets() {
		// Check if the packet is an IPv4 packet containing ICMP
		if ipv4Layer := packet.Layer(layers.LayerTypeIPv4); ipv4Layer != nil {
			icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
			if icmpLayer != nil {
				icmp, _ := icmpLayer.(*layers.ICMPv4)

				// Check if the ICMP packet is an Echo Request (Type 8)
				if icmp.TypeCode.Type() == layers.ICMPv4TypeEchoRequest {
					// Store Echo Request packet by identifier and sequence number
					key := fmt.Sprintf("%d", icmp.Seq)
					echoRequests[key] = packet
				
					// Print identifier and sequence number of this echo request
					fmt.Printf("Echo Request Sequence Number: %d\n",icmp.Seq)
				}

				// Check if the ICMP packet is a Time Exceeded (Type 11)
				if icmp.TypeCode.Type() == layers.ICMPv4TypeTimeExceeded {
					// Ensure the payload is at least 2 bytes long
					if len(icmp.Payload) < 2 {
						fmt.Println("Payload is too short")
						return
					}
				
					// Extract the last 2 bytes from the payload
					lastTwoBytes := icmp.Payload[len(icmp.Payload)-2:]
				
					// Convert the bytes to a 16-bit unsigned integer
					sequenceNumber := binary.BigEndian.Uint16(lastTwoBytes)
				
					// Print the extracted sequence number
					fmt.Printf("Reply Sequence Number: %d\n", sequenceNumber)
					key := fmt.Sprintf("%d", sequenceNumber)
    				echoReply[key] = packet
				}
			}
		}
	}
}
