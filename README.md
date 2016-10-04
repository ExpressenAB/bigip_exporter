# BIG-IP exporter
Prometheus exporter for BIG-IP statistics. Uses iControl REST API.

## Get it
The latest version is 0.2.2. All releases can be found under [Releases](https://github.com/ExpressenAB/bigip_exporter/releases) and docker images are available at [Docker Hub](https://hub.docker.com/r/expressenab/bigip_exporter/tags/).

## Usage
The bigip_exporter is easy to use. Example: 
```
./bigip_exporter -bigip.host <bigip-host> -bigip.port 443 -bigip.username admin -bigip.password admin
```
Or, if you prefer, you can run it in a docker container
```
docker run -p 9142:9142 expressenab/bigip_exporter -bigip.host <bigip-host> -bigip.port 443 -bigip.username admin -bigip.password admin
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
### Building locally
This project uses [govendor](https://github.com/kardianos/govendor). If you do not already have that installed, take a detour and install that beforehand.
```
# This assumes that you already have go and govendor installed and $GOPATH configured
go get github.com/ExpressenAB/bigip_exporter
cd $GOPATH/src/github.com/ExpressenAB/bigip_exporter
govendor build +p
```
### Cross compilation
Go offers possibility to cross compile the application for different use on a different OS and architecture. This is achieved by setting the environment valiables `GOOS` and `GOARCH`. If you for example want to build for linux on an amd64 architecture the `go build` step can be replaced with the following:
```
GOOS=linux GOARCH=amd64 govendor build +p
```
A list of available options for `GOOS` and `GOARCH` is available in the [documentation](https://golang.org/doc/install/source#environment)

## Possible improvements
### Gather data in the background
Currently the data is gathered when the `/metrics` endpoint is called. This causes the request to take about 4-6 seconds before completing. This could be fixed by having a go thread that gathers data at regular intervals and that is returned upon a call to the `/metrics` endpoint. This would however go against the [guidelines](https://prometheus.io/docs/instrumenting/writing_exporters/#scheduling).
