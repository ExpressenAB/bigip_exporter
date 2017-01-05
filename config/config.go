package config

import (
	"os"
	"strings"

	"github.com/juju/loggo"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type bigipConfig struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	BasicAuth bool   `yaml:"basic_auth"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
}

type exporterConfig struct {
	BindAddress string `yaml:"bind_address"`
	BindPort    int    `yaml:"bind_port"`
	Partitions  string `yaml:"partitions"`
	Config      string `yaml:"config"`
	Namespace   string `yaml:"namespace"`
	LogLevel    string `yaml:"log_level"`
}

// Config is a container for settings modifiable by the user
type Config struct {
	Bigip    bigipConfig    `yaml:"bigip"`
	Exporter exporterConfig `yaml:"exporter"`
}

var (
	logger = loggo.GetLogger("")
)

func init() {
	loggo.ConfigureLoggers("<root>=INFO")
	registerFlags()
	bindFlags()
	bindEnvs()
	flag.Parse()

	if viper.GetString("exporter.config") != "" {
		readConfigFile(viper.GetString("exporter.config"))
	}

	logLevel := viper.GetString("exporter.log_level")

	if _, validLevel := loggo.ParseLevel(logLevel); validLevel {
		loggo.ConfigureLoggers("<root>=" + strings.ToUpper(logLevel))
		return
	}

	logger.Warningf("Invalid log level - Using info")
}

func registerFlags() {
	flag.Bool("bigip.basic_auth", false, "Use HTTP Basic authentication")
	flag.String("bigip.host", "localhost", "The host on which f5 resides")
	flag.Int("bigip.port", 443, "The port which f5 listens to")
	flag.String("bigip.username", "user", "Username")
	flag.String("bigip.password", "pass", "Password")
	flag.String("exporter.bind_address", "localhost", "Exporter bind address")
	flag.Int("exporter.bind_port", 9142, "Exporter bind port")
	flag.String("exporter.partitions", "", "A comma separated list of partitions which to export. (default: all)")
	flag.String("exporter.config", "", "bigip_exporter configuration file name.")
	flag.String("exporter.namespace", "bigip", "bigip_exporter namespace.")
	flag.String("exporter.log_level", "info", "Available options are trace, debug, info, warning, error and critical")
}

func bindFlags() {
	flag.VisitAll(func(f *flag.Flag) {
		err := viper.BindPFlag(f.Name, f)
		if err != nil {
			logger.Warningf("Failed to bind flag (%s)", err)
		}
	})
}

func bindEnvs() {
	viper.SetEnvPrefix("be")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	flag.VisitAll(func(f *flag.Flag) {
		err := viper.BindEnv(f.Name)
		if err != nil {
			logger.Warningf("Failed to bind environment variable BE_%s (%s)", strings.ToUpper(strings.Replace(f.Name, ".", "_", -1)), err)
		}
	})
}

func readConfigFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		logger.Warningf("Failed to open configuration file (%s)", err)
		return
	}
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(file)
	if err != nil {
		logger.Warningf("Failed to read configuration file (%s)", err)
	}
}

// GetConfig returns an instance of Config containing the resulting parameters
// to the program
func GetConfig() *Config {
	return &Config{
		Bigip: bigipConfig{
			Username:  viper.GetString("bigip.username"),
			Password:  viper.GetString("bigip.password"),
			BasicAuth: viper.GetBool("bigip.basic_auth"),
			Host:      viper.GetString("bigip.host"),
			Port:      viper.GetInt("bigip.port"),
		},
		Exporter: exporterConfig{
			BindAddress: viper.GetString("exporter.bind_address"),
			BindPort:    viper.GetInt("exporter.bind_port"),
			Partitions:  viper.GetString("exporter.partitions"),
			Config:      viper.GetString("exporter.config"),
			Namespace:   viper.GetString("exporter.namespace"),
			LogLevel:    viper.GetString("exporter.log_level"),
		},
	}
}
