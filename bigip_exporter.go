package main

import (
	"net/http"
	"strconv"
	//"strings"

	//"github.com/ExpressenAB/bigip_exporter/collector"
	"github.com/ExpressenAB/bigip_exporter/config"
	"github.com/juju/loggo"
	//"github.com/pr8kerl/f5er/f5"
	//"github.com/prometheus/client_golang/prometheus"
   	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strings"
	"github.com/pr8kerl/f5er/f5"
	"github.com/ExpressenAB/bigip_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	logger = loggo.GetLogger("")
	configuration = config.GetConfig()
)

func listen(exporterBindAddress string, exporterBindPort int) {
	http.Handle("/metrics", getTarget(promhttp.Handler()))
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

//Wrapper around the handler.
//Get the key from the handler
//Return/Call the handler that is passed in as Parameter

func getTarget(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Query().Get("target")) == 0 {
			logger.Errorf("Missing target")
			http.Error(w, "Missing target", http.StatusUnprocessableEntity)
			return // don't call original handler
		}else {
			if val, ok := configuration.Lookup[r.URL.Query().Get("target")]; ok {
				bigipEndpoint := val.Host + ":" + strconv.Itoa(val.Port)
				var exporterPartitionsList []string
				if configuration.Exporter.Partitions != "" {
					exporterPartitionsList = strings.Split(configuration.Exporter.Partitions, ",")
				} else {
					exporterPartitionsList = nil
				}
				authMethod := f5.TOKEN
				if val.BasicAuth {
					authMethod = f5.BASIC_AUTH
				}

				bigip := f5.New(bigipEndpoint,val.Username,val.Password,authMethod)
				bigipCollector, _ := collector.NewBigipCollector(bigip, configuration.Exporter.Namespace, exporterPartitionsList)
				prometheus.MustRegister(bigipCollector)
			} else {
				//Target not found
				logger.Errorf("Exporter does not have the configuration for target [%v]", r.URL.Query().Get("target"))
				http.Error(w, "Target not supported", http.StatusUnprocessableEntity)
			}
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	logger.Debugf("Config: [%v]", configuration)
	listen(configuration.Exporter.BindAddress, configuration.Exporter.BindPort)
}
