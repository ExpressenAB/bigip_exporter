package main

import (
	"flag"
	"github.com/ExpressenAB/bigip_exporter/collector"
	"github.com/pr8kerl/f5er/f5"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"strconv"
)

var (
	bigip_host            = flag.String("bigip.host", "localhost", "The host on which f5 resides")
	bigip_port            = flag.Int("bigip.port", 443, "The port which f5 listens to")
	bigip_username        = flag.String("bigip.username", "user", "Username")
	bigip_password        = flag.String("bigip.password", "pass", "Password")
	exporter_bind_address = flag.String("exporter.bind_address", "", "Exporter bind address")
	exporter_bind_port    = flag.Int("exporter.bind_port", 9142, "Exporter bind port")
	exporter_namespace    = "bigip"
)

func main() {
	flag.Parse()
	bigip_endpoint := *bigip_host + ":" + strconv.Itoa(*bigip_port)
	bigip := f5.New(bigip_endpoint, *bigip_username, *bigip_password, f5.TOKEN)
	_, vsCollector := collector.NewVSCollector(bigip, exporter_namespace)
	_, poolCollector := collector.NewPoolCollector(bigip, exporter_namespace)
	prometheus.MustRegister(vsCollector)
	prometheus.MustRegister(poolCollector)
	http.Handle("/metrics", prometheus.Handler())
	exporter_bind := *exporter_bind_address + ":" + strconv.Itoa(*exporter_bind_port)
	log.Fatal(http.ListenAndServe(exporter_bind, nil))
}
