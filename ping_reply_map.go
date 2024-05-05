package main

import (
	"fmt"
	"log"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/layers"
	"encoding/binary"
	"strings"
	"os"
	"io/ioutil"
	"encoding/csv"
	"time"
	"bufio"
)
var times []time.Time

func readTimesFromFile(filename string) ([]time.Time, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var times []time.Time

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Read each line as a string
		line := scanner.Text()

		// Parse the string into a time.Time value
		t, err := time.Parse("15:04:05.999999", line)
		if err != nil {
			return nil, fmt.Errorf("error parsing time: %w", err)
		}

		// Append the parsed time to the times slice
		times = append(times, t)
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return times, nil
}
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
		return fmt.Sprintf("%d", icmpPacket.Id)
	}
	return ""
}
func writeMatchingPacketsToCSV(echoRequests, echoReply map[string]gopacket.Packet, ipAddresses []string) error {
	// Open the CSV file for writing
	file, err := os.Create("ping_reply.csv")
	if err != nil {
		return fmt.Errorf("error creating CSV file: %w", err)
	}
	defer file.Close()
	upload_start:= times[0]
	idle_start:= times[1]
	fmt.Println("Size of the req:", len(echoRequests))
	fmt.Println("Size of the reply:", len(echoReply))
	fmt.Println("uploadstart=", upload_start)
	fmt.Println("idle start= ", idle_start)
	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header to CSV file
	header := []string{"SequenceNumber","RequestTime", "ReplyTime", "RequestSourceIP", "RequestDestIP", "ReplySourceIP", "ReplyDestIP", "ReplyTTLexpiredIP", "TTL", "Download/Upload/Idle"}
	err = writer.Write(header)
	if err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}
	count := 0
	// Iterate over echoRequests and echoReply maps
	for key, request := range echoRequests {
		reply, ok := echoReply[key]
		count +=1
		if ok {
			count+=1
			// Extract required fields from request and reply packets
			sequenceNumber := getRequestSequenceNumber(request)
			requestTime := request.Metadata().Timestamp.Format("15:04:05.999999999")
			replyTime := reply.Metadata().Timestamp.Format("15:04:05.999999999")
			requestSourceIP, requestDestIP := getRequestIPs(request)
			replySourceIP, replyDestIP := getRequestIPs(reply)
			replyTTLExpiredIP := getRequestTTLExpiredIP(reply)
			ttl:= ""
			dui:=""

			for ind, ip := range ipAddresses {
				if ip == replySourceIP {
					ttl = fmt.Sprintf("%d", ind)
				}
			}
			requestTimeParsed, _ := time.Parse("15:04:05.999999999", requestTime)
			
			if requestTimeParsed.Before(upload_start) {
				dui = "download"
				fmt.Println("download")
			} else if (requestTimeParsed.After(upload_start) && requestTimeParsed.Before(idle_start)) ||  requestTimeParsed.Equal(upload_start) || requestTimeParsed.Equal(idle_start)  {
				dui = "upload"
				fmt.Println("upload")
			} else {
				dui = "idle"
				fmt.Println("idle")
			}
			// Write fields to CSV file
			record := []string{sequenceNumber, requestTime, replyTime, requestSourceIP, requestDestIP, replySourceIP, replyDestIP, replyTTLExpiredIP, ttl, dui}
			err := writer.Write(record)
			if err != nil {
				return fmt.Errorf("error writing CSV record: %w", err)
			}
		}
	}
	fmt.Println(count)
	return nil
}


func main() {
	// Open the pcap file
	handle, err := pcap.OpenOffline("capture.pcap")
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// ------------------------------------------------------------------
	filePath := "ip_addresses.txt"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	ipAddresses := strings.Split(string(content), "\n")

	// -------------------------------------------------------------------


	// Create a packet source from the handle
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Map to store ICMP Echo Request packets by their identifier and sequence number
	echoRequests := make(map[string]gopacket.Packet)
	echoReply := make(map[string]gopacket.Packet)
	count:= 0
	// Loop through each packet in the pcap file
	for packet := range packetSource.Packets() {
		// Check if the packet is an IPv4 packet containing ICMP
		if ipv4Layer := packet.Layer(layers.LayerTypeIPv4); ipv4Layer != nil {
			icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
			// payloadLength := len(icmpLayer.LayerPayload())
			// fmt.Println("length =====", payloadLength)
			if icmpLayer != nil {
				icmp, _ := icmpLayer.(*layers.ICMPv4)
				
				// Check if the ICMP packet is an Echo Request (Type 8)
				if icmp.TypeCode.Type() == layers.ICMPv4TypeEchoRequest {
					// Store Echo Request packet by identifier and sequence number
					keyreq := fmt.Sprintf("%d", icmp.Id)
					// fmt.Println("length req= ", len((icmp.Payload)))
					// lastsixtyBytes := icmp.Payload[len(icmp.Payload)-4 : len(icmp.Payload)-2]
					// sequenceNumberreq := binary.LittleEndian.Uint16(lastsixtyBytes)
					// keyreq := fmt.Sprintf("%d", sequenceNumberreq)
					
					// fmt.Println("req", keyreq)
					echoRequests[keyreq] = packet
					fmt.Println("key", keyreq)

					// count +=1
				
				}

				// Check if the ICMP packet is a Time Exceeded (Type 11)
				if icmp.TypeCode.Type() == layers.ICMPv4TypeTimeExceeded {
					// Ensure the payload is at least 2 bytes long
					if len(icmp.Payload) < 2 {
						fmt.Println("Error: Payload is too short")
						return
					}
				
					// Extract the last 2 bytes from the payload to match with sequence number
					lastTwoBytes := icmp.Payload[len(icmp.Payload)-4 : len(icmp.Payload)-2]
					sequenceNumber := binary.BigEndian.Uint16(lastTwoBytes)
					key := fmt.Sprintf("%d", sequenceNumber)
    				echoReply[key] = packet
					// fmt.Println("reply ", key)
					count+= 1
				}
			}
		}
	}
	fmt.Print(count)
	times, _ = readTimesFromFile("times.txt")
	fmt.Println("Size of the req:", len(echoRequests))
	fmt.Println("Size of the reply:", len(echoReply))

	err1 := writeMatchingPacketsToCSV(echoRequests, echoReply, ipAddresses)
	if err1 != nil {
		fmt.Println("Error writing to CSV:", err1)
	}

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	


}