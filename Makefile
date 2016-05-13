GOPATH := $(shell pwd)
GOBIN  := $(GOPATH)
PATH   := $(GOROOT):$(PATH)
DEPS   := github.com/pr8kerl/f5er/f5 github.com/ExpressenAB/bigip_exporter/collector github.com/prometheus/client_golang/prometheus
VER    := 0.1.0

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
		GOOS=linux GOARCH=amd64 GOPATH=$(GOPATH) go build -o bigip_exporter -v $^
		mkdir bigip_exporter-$(VER).linux-amd64
		cp bigip_exporter bigip_exporter-$(VER).linux-amd64/
		cp LICENSE bigip_exporter-$(VER).linux-amd64/
		tar -zcvf bigip_exporter-$(VER).linux-amd64.tar.gz bigip_exporter-$(VER).linux-amd64

win64: bigip_exporter.go
		GOPATH=$(GOPATH) go fmt $^
		GOPATH=$(GOPATH) go tool vet $^
		GOOS=windows GOARCH=amd64 GOPATH=$(GOPATH) go build -o bigip_exporter.exe -v $^
		mkdir bigip_exporter-$(VER).win-amd64
		cp bigip_exporter.exe bigip_exporter-$(VER).win-amd64/
		cp LICENSE bigip_exporter-$(VER).win-amd64/
		tar -zcvf bigip_exporter-$(VER).win-amd64.tar.gz bigip_exporter-$(VER).win-amd64

.PHONY: $(DEPS) clean

clean:
	rm -rf bigip_exporter bigip_exporter.exe bigip_exporter*-amd64
