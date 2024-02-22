.PHONY: all clean

all: ping_windows writecsv filter timecheck ping_reply_map

clean:
	del /Q capture.pcap times.txt filtered_ndtcapture.csv ip_addresses.txt ndtcapture.csv ping_reply.csv

ping_windows:
	go run ping_windows.go util.go

writecsv:
	go run writecsv.go

filter:
	go run filter.go

timecheck:
	go run timecheck.go

ping_reply_map:
	go run ping_reply_map.go
