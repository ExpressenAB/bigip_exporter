# BIG-IP exporter
Prometheus exporter for BIG-IP statistics. Uses iControl REST API.

## Get it
The latest version is 0.2.1 and all releases can be found under [Releases](https://github.com/ExpressenAB/bigip_exporter/releases).

## Usage
The bigip_exporter is easy to use. Example: 
```
./bigip_exporter -bigip.host 127.0.0.1 -bigip.port 443 -bigip.username admin -bigip.password admin
```

### Flags
Flag | Description | Default
-----|-------------|---------
-bigip.host | BIG-IP host | localhost
-bigip.port | BIG-IP port | 443
-bigip.username | BIG-IP username | user
-bigip.password | BIG-IP password | pass
-exporter.bind_address | The address the exporter should bind to | All interfaces
-exporter.bind_port | Which port the exporter should listen on | 9142
-exporter.partitions | A comma separated list containing the partitions that should be exported | All partitions
-exporter.namespace | The namespace used in prometheus labels | bigip

## Implemented metrics
* Virtual Server
* Rule
* Pool
* Node

## Prerequisites
* User with read access to iControl REST API

## Tested versions of iControl REST API
Currently only version 12.0.0 is tested. If you experience any problems with other versions, create an issue explaining the problem and I'll look at it as soon as possible or if you'd like to contribute with a pull request that would be greatly appreciated.

## Building
### Building for your local machine
```
# This presumes that you already have go installed and $GOPATH configured
go get github.com/ExpressenAB/bigip_exporter
cd $GOPATH/src/github.com/ExpressenAB/bigip_exporter
go build
```
### Cross compilation
Go offers possibility to cross compile the application for different use on a different OS and architecture. This is achieved by setting the environment valiables `GOOS` and `GOARCH`. If you for example want to build for linux on an amd64 architecture the `go build` step can be replaced with the following:
```
GOOS=linux GOARCH=amd64 go build
```
A list of available options for `GOOS` and `GOARCH` is available in the [documentation](https://golang.org/doc/install/source#environment)

## Possible improvements
### Gather data in the background
Currently the data is gathered when the `/metrics` endpoint is called. This causes the request to take about 4-6 seconds before completing. This could be fixed by having a go thread that gathers data at regular intervals and that is returned upon a call to the `/metrics` endpoint. This would however go against the [guidelines](https://prometheus.io/docs/instrumenting/writing_exporters/#scheduling).
