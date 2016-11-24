package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/ExpressenAB/bigip_exporter/collector"
	"github.com/pr8kerl/f5er/f5"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	bigip_basic_auth      = flag.Bool("bigip.basic_auth", false, "Use HTTP Basic authentication")
	bigip_host            = flag.String("bigip.host", "localhost", "The host on which f5 resides")
	bigip_port            = flag.Int("bigip.port", 443, "The port which f5 listens to")
	bigip_username        = flag.String("bigip.username", "", "Username")
	bigip_password        = flag.String("bigip.password", "", "Password")
	exporter_bind_address = flag.String("exporter.bind_address", "", "Exporter bind address")
	exporter_bind_port    = flag.Int("exporter.bind_port", 9142, "Exporter bind port")
	exporter_partitions   = flag.String("exporter.partitions", "", "A comma separated list of partitions which to export. Default: all")
	exporter_config       = flag.String("exporter.config", "", "bigip_exporter configuration file name.")
	exporter_namespace    = "bigip"
	username              *string
	password              *string
)

type Config struct {
	Bigip_username string
	Bigip_password string
}

func DefaultConfig() *Config {
	config := &Config{
		Bigip_username: "",
		Bigip_password: "",
	}
	return config
}

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
	log.Fatal(http.ListenAndServe(exporter_bind, nil))
}

func main() {
	flag.Parse()

	config := DefaultConfig()

	if *exporter_config != "" {
		yamlFile, err := ioutil.ReadFile(*exporter_config)
		if err != nil {
			log.Fatal(err)
		}

		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *bigip_username != "" {
		username = bigip_username
	} else if config.Bigip_username != "" {
		username = &config.Bigip_username
	} else {
		log.Fatal("ERROR: Missing argument username")
	}

	if *bigip_password != "" {
		password = bigip_password
	} else if config.Bigip_password != "" {
		password = &config.Bigip_password
	} else {
		log.Fatal("ERROR: Missing argument password")
	}

	bigip_endpoint := *bigip_host + ":" + strconv.Itoa(*bigip_port)
	var exporter_partitions_list []string
	if *exporter_partitions != "" {
		exporter_partitions_list = strings.Split(*exporter_partitions, ",")
	} else {
		exporter_partitions_list = nil
	}
	auth_method := f5.TOKEN
	if *bigip_basic_auth {
		auth_method = f5.BASIC_AUTH
	}

	bigip := f5.New(bigip_endpoint, *username, *password, auth_method)

	_, bigipCollector := collector.NewBigIpCollector(bigip, exporter_namespace, exporter_partitions_list)

	prometheus.MustRegister(bigipCollector)
	listen(*exporter_bind_address, *exporter_bind_port)
}
