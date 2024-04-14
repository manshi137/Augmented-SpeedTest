.PHONY: ping nping clean

ping: ping_linux writecsv filter timecheck ping_reply_map ttest

nping: nping_linux writecsv filter timecheck ping_reply_map ttest

clean:
	del /Q capture.pcap times.txt filtered_ndtcapture.csv ip_addresses.txt ndtcapture.csv ping_reply.csv ttest_output.txt

ping_linux:
	sudo env "PATH=$(PATH)" "GOPATH=$(GOPATH)" go run ping_linux.go util.go

nping_linux:
	sudo env "PATH=$(PATH)" "GOPATH=$(GOPATH)"  go run nping_linux.go util.go

writecsv:
	sudo env "PATH=$(PATH)" "GOPATH=$(GOPATH)" go run writecsv.go

filter:
	sudo env "PATH=$(PATH)" "GOPATH=$(GOPATH)" go run filter.go

timecheck:
	sudo env "PATH=$(PATH)" "GOPATH=$(GOPATH)" go run timecheck.go

ping_reply_map:
	sudo env "PATH=$(PATH)" "GOPATH=$(GOPATH)" go run ping_reply_map.go

ttest: 
	python3 ttest.py >> ttest_output.txt
