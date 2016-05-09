package main

import (
	"flag"
	"github.com/ExpressenAB/bigip_exporter/collector"
	"github.com/pr8kerl/f5er/f5"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	bigip_host            = flag.String("bigip.host", "localhost", "The host on which f5 resides")
	bigip_port            = flag.Int("bigip.port", 443, "The port which f5 listens to")
	bigip_username        = flag.String("bigip.username", "user", "Username")
	bigip_password        = flag.String("bigip.password", "pass", "Password")
	exporter_bind_address = flag.String("exporter.bind_address", "", "Exporter bind address")
	exporter_bind_port    = flag.Int("exporter.bind_port", 9142, "Exporter bind port")
	exporter_partitions   = flag.String("exporter.partitions", "", "A comma separated list of partitions which to export. Default: all")
	exporter_namespace    = "bigip"
)

func main() {
	flag.Parse()
	bigip_endpoint := *bigip_host + ":" + strconv.Itoa(*bigip_port)
	var exporter_partitions_list []string
	if *exporter_partitions != "" {
		exporter_partitions_list = strings.Split(*exporter_partitions, ",")
	} else {
		exporter_partitions_list = nil
	}
	bigip := f5.New(bigip_endpoint, *bigip_username, *bigip_password, f5.TOKEN)
	_, vsCollector := collector.NewVSCollector(bigip, exporter_namespace, exporter_partitions_list)
	_, poolCollector := collector.NewPoolCollector(bigip, exporter_namespace, exporter_partitions_list)
	_, nodeCollector := collector.NewNodeCollector(bigip, exporter_namespace, exporter_partitions_list)
	_, ruleCollector := collector.NewRuleCollector(bigip, exporter_namespace, exporter_partitions_list)
	prometheus.MustRegister(vsCollector)
	prometheus.MustRegister(poolCollector)
	prometheus.MustRegister(nodeCollector)
	prometheus.MustRegister(ruleCollector)
	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>BIG-IP Exporter</title></head>
			<body>
			<h1>BIG-IP Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})
	exporter_bind := *exporter_bind_address + ":" + strconv.Itoa(*exporter_bind_port)
	log.Fatal(http.ListenAndServe(exporter_bind, nil))
}
