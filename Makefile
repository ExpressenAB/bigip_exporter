#GOROOT := /usr/local/go
GOPATH := $(shell pwd)
GOBIN  := $(GOPATH)
PATH   := $(GOROOT):$(PATH)
DEPS   := github.com/pr8kerl/f5er/f5 github.com/ExpressenAB/bigip_exporter/collector github.com/prometheus/client_golang/prometheus

all: bigip_exporter

deps: $(DEPS)
	GOPATH=$(GOPATH) go get -u $^

bigip_exporter: bigip_exporter.go
		GOPATH=$(GOPATH) go fmt $^
		GOPATH=$(GOPATH) go build -o $@ -v $^
		touch $@

linux64: bigip_exporter.go
		GOPATH=$(GOPATH) go fmt $^
		GOPATH=$(GOPATH) go tool vet $^
		GOOS=linux GOARCH=amd64 GOPATH=$(GOPATH) go build -o bigip_exporter-linux-amd64 -v $^
		touch bigip_exporter-linux-amd64

win64: bigip_exporter.go
		GOPATH=$(GOPATH) go fmt $^
		GOPATH=$(GOPATH) go tool vet $^
		GOOS=windows GOARCH=amd64 GOPATH=$(GOPATH) go build -o bigip_exporter-win-amd64.exe -v $^
		touch bigip_exporter-win-amd64.exe

.PHONY: $(DEPS) clean

clean:
	rm -f bigip_exporter bigip_exporter-win-amd64.exe bigip_exporter-linux-amd64
