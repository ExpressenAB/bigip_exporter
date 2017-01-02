package config

import (
	"flag"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

var (
	bigip_basic_auth      = flag.Bool("bigip.basic_auth", false, "Use HTTP Basic authentication")
	bigip_host            = flag.String("bigip.host", "localhost", "The host on which f5 resides")
	bigip_port            = flag.Int("bigip.port", 443, "The port which f5 listens to")
	bigip_username        = flag.String("bigip.username", "user", "Username")
	bigip_password        = flag.String("bigip.password", "pass", "Password")
	exporter_bind_address = flag.String("exporter.bind_address", "localhost", "Exporter bind address")
	exporter_bind_port    = flag.Int("exporter.bind_port", 9142, "Exporter bind port")
	exporter_partitions   = flag.String("exporter.partitions", "", "A comma separated list of partitions which to export. Default: all")
	exporter_config       = flag.String("exporter.config", "", "bigip_exporter configuration file name.")
	exporter_namespace    = flag.String("exporter.namespace", "bigip", "bigip_exporter namespace.")
	debug                 = flag.Bool("debug", false, "Verbose output, NOTE: prints out the configuration in clear")
)

type bigipConfig struct {
	Bigip_username   string `yaml:"bigip_username"`
	Bigip_password   string `yaml:"bigip_password"`
	Bigip_basic_auth bool   `yaml:"bigip_basic_auth"`
	Bigip_host       string `yaml:"bigip_host"`
	Bigip_port       int    `yaml:"bigip_port"`
}

type exporterConfig struct {
	Exporter_bind_address string `yaml:"exporter_bind_address"`
	Exporter_bind_port    int    `yaml:"exporter_bind_port"`
	Exporter_partitions   string `yaml:"exporter_partitions"`
	Exporter_config       string `yaml:"exporter_config"`
	Exporter_namespace    string `yaml:"exporter_namespace"`
	Exporter_debug        bool   `yaml:"exporter_debug"`
}

type Config struct {
	Bigip_config    bigipConfig    `yaml:"bigip_config"`
	Exporter_config exporterConfig `yaml:"exporter_config"`
}

func defaultConfig() *Config {
	config := &Config{
		Bigip_config: bigipConfig{
			Bigip_username:   "",
			Bigip_password:   "",
			Bigip_basic_auth: false,
			Bigip_host:       "",
			Bigip_port:       443,
		},
		Exporter_config: exporterConfig{
			Exporter_bind_address: "",
			Exporter_bind_port:    9142,
			Exporter_partitions:   "",
			Exporter_config:       "",
			Exporter_namespace:    "",
			Exporter_debug:        false,
		},
	}
	return config
}

func GetConfig() *Config {

	config := defaultConfig()
	flag.Parse()

	if *exporter_config != "" {
		log.Printf("Loading config file %v", *exporter_config)
		yamlFile, err := ioutil.ReadFile(*exporter_config)
		if err != nil {
			log.Fatal(err)
		}

		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			log.Fatal(err)
		}
	}

	if config.Bigip_config.Bigip_username == "" {
		config.Bigip_config.Bigip_username = *bigip_username
	} else {
		log.Printf("Loading bigip_username from configuration file")
	}
	if config.Bigip_config.Bigip_password == "" {
		config.Bigip_config.Bigip_password = *bigip_password
	} else {
		log.Printf("Loading bigip_password from configuration file")
	}
	if config.Bigip_config.Bigip_host == "" {
		config.Bigip_config.Bigip_host = *bigip_host
	} else {
		log.Printf("Loading bigip_host from configuration file")
	}
	if config.Bigip_config.Bigip_port == 443 {
		config.Bigip_config.Bigip_port = *bigip_port
	} else {
		log.Printf("Loading bigip_port from configuration file")
	}
	if !config.Bigip_config.Bigip_basic_auth {
		config.Bigip_config.Bigip_basic_auth = *bigip_basic_auth
	} else {
		log.Printf("Loading bigip_basic_auth from configuration file")
	}

	if config.Exporter_config.Exporter_bind_address == "" {
		config.Exporter_config.Exporter_bind_address = *exporter_bind_address
	} else {
		log.Printf("Loading exporter_bind_address from configuration file")
	}
	if config.Exporter_config.Exporter_bind_port == 9142 {
		config.Exporter_config.Exporter_bind_port = *exporter_bind_port
	} else {
		log.Printf("Loading exporter_bind_port from configuration file")
	}
	if config.Exporter_config.Exporter_partitions == "" {
		config.Exporter_config.Exporter_partitions = *exporter_partitions
	} else {
		log.Printf("Loading exporter_partitions from configuration file")
	}
	if config.Exporter_config.Exporter_namespace == "" {
		config.Exporter_config.Exporter_namespace = *exporter_namespace
	} else {
		log.Printf("Loading exporter_namespace from configuration file")
	}
	if !config.Exporter_config.Exporter_debug {
		config.Exporter_config.Exporter_debug = *debug
	}

	return config
}
