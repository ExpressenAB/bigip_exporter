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
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	logger = loggo.GetLogger("")
)

func listen(exporterBindAddress string, exporterBindPort int) {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>BIG-IP Exporter</title></head>
			<body>
			<h1>BIG-IP Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})
	exporterBind := exporterBindAddress + ":" + strconv.Itoa(exporterBindPort)
	logger.Criticalf("Process failed: %s", http.ListenAndServe(exporterBind, nil))
}

func main() {
	config := config.GetConfig()
	logger.Debugf("Config: %v", config)

	bigipEndpoint := config.Bigip.Host + ":" + strconv.Itoa(config.Bigip.Port)
	var exporterPartitionsList []string
	if config.Exporter.Partitions != "" {
		exporterPartitionsList = strings.Split(config.Exporter.Partitions, ",")
	} else {
		exporterPartitionsList = nil
	}
	authMethod := f5.TOKEN
	if config.Bigip.BasicAuth {
		authMethod = f5.BASIC_AUTH
	}

	bigip := f5.New(bigipEndpoint, config.Bigip.Username, config.Bigip.Password, authMethod)

	bigipCollector, _ := collector.NewBigipCollector(bigip, config.Exporter.Namespace, exporterPartitionsList)

	prometheus.MustRegister(bigipCollector)
	listen(config.Exporter.BindAddress, config.Exporter.BindPort)
}
