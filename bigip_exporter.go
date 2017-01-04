package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ExpressenAB/bigip_exporter/collector"
	"github.com/ExpressenAB/bigip_exporter/config"
	"github.com/juju/loggo"
	"github.com/pr8kerl/f5er/f5"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	logger = loggo.GetLogger("")
)

func listen(exporter_bind_address string, exporter_bind_port int) {
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
	exporter_bind := exporter_bind_address + ":" + strconv.Itoa(exporter_bind_port)
	logger.Criticalf("Process failed: %s", http.ListenAndServe(exporter_bind, nil))
}

func main() {
	config := config.GetConfig()
	logger.Debugf("Config: %v", config)

	bigip_endpoint := config.Bigip.Host + ":" + strconv.Itoa(config.Bigip.Port)
	var exporter_partitions_list []string
	if config.Exporter.Partitions != "" {
		exporter_partitions_list = strings.Split(config.Exporter.Partitions, ",")
	} else {
		exporter_partitions_list = nil
	}
	auth_method := f5.TOKEN
	if config.Bigip.BasicAuth {
		auth_method = f5.BASIC_AUTH
	}

	bigip := f5.New(bigip_endpoint, config.Bigip.Username, config.Bigip.Password, auth_method)

	_, bigipCollector := collector.NewBigIpCollector(bigip, config.Exporter.Namespace, exporter_partitions_list)

	prometheus.MustRegister(bigipCollector)
	listen(config.Exporter.BindAddress, config.Exporter.BindPort)
}
