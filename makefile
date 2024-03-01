.PHONY: ping nping clean

ping: ping_windows writecsv filter timecheck ping_reply_map ttest

nping: nping_windows writecsv filter timecheck ping_reply_map ttest

clean:
	del /Q capture.pcap times.txt filtered_ndtcapture.csv ip_addresses.txt ndtcapture.csv ping_reply.csv ttest_output.txt

ping_windows:
	go run ping_windows.go util.go

nping_windows:
	go run nping_windows.go util.go

writecsv:
	go run writecsv.go

filter:
	go run filter.go

timecheck:
	go run timecheck.go

ping_reply_map:
	go run ping_reply_map.go

ttest: 
	python ttest.py >> ttest_output.txt
