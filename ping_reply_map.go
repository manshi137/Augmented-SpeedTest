package main

import (
	"fmt"
	"log"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/layers"
	"encoding/binary"
	// "strings"
	"os"
	"encoding/csv"
	// "time"
)

// Helper function to extract source and destination IPs from packet
func getRequestIPs(packet gopacket.Packet) (string, string) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ipPacket, _ := ipLayer.(*layers.IPv4)
		return ipPacket.SrcIP.String(), ipPacket.DstIP.String()
	}
	return "", ""
}

// Helper function to extract TTL expired IP from packet
func getRequestTTLExpiredIP(packet gopacket.Packet) string {
	icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
	if icmpLayer != nil {
		icmpPacket, _ := icmpLayer.(*layers.ICMPv4)
		if icmpPacket.TypeCode.Type() == layers.ICMPv4TypeTimeExceeded {
			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			if ipLayer != nil {
				ipPacket, _ := ipLayer.(*layers.IPv4)
				return ipPacket.SrcIP.String()
			}
		}
	}
	return ""
}

// Helper function to extract sequence number from packet
func getRequestSequenceNumber(packet gopacket.Packet) string {
	icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
	if icmpLayer != nil {
		icmpPacket, _ := icmpLayer.(*layers.ICMPv4)
		return fmt.Sprintf("%d", icmpPacket.Seq)
	}
	return ""
}
func writeMatchingPacketsToCSV(echoRequests, echoReply map[string]gopacket.Packet) error {
	// Open the CSV file for writing
	file, err := os.Create("ping_reply.csv")
	if err != nil {
		return fmt.Errorf("error creating CSV file: %w", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header to CSV file
	header := []string{"SequenceNumber","RequestTime", "ReplyTime", "RequestSourceIP", "RequestDestIP", "ReplySourceIP", "ReplyDestIP", "ReplyTTLexpiredIP"}
	err = writer.Write(header)
	if err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	// Iterate over echoRequests and echoReply maps
	for key, request := range echoRequests {
		reply, ok := echoReply[key]
		if ok {
			// Extract required fields from request and reply packets
			sequenceNumber := getRequestSequenceNumber(request)
			requestTime := request.Metadata().Timestamp.Format("15:04:05.999999999")
			replyTime := reply.Metadata().Timestamp.Format("15:04:05.999999999")
			requestSourceIP, requestDestIP := getRequestIPs(request)
			replySourceIP, replyDestIP := getRequestIPs(reply)
			replyTTLExpiredIP := getRequestTTLExpiredIP(reply)

			// Write fields to CSV file
			record := []string{sequenceNumber, requestTime, replyTime, requestSourceIP, requestDestIP, replySourceIP, replyDestIP, replyTTLExpiredIP}
			err := writer.Write(record)
			if err != nil {
				return fmt.Errorf("error writing CSV record: %w", err)
			}
		}
	}
	fmt.Println("CSV file written successfully.")
	return nil
}

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
	err1 := writeMatchingPacketsToCSV(echoRequests, echoReply)
	if err1 != nil {
		fmt.Println("Error writing to CSV:", err1)
	}

}
