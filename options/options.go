package options

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

// ConfigOptions configuration options
type ConfigOptions struct {
	Verbose          bool   `yaml:"verbose" env:"IPFIX_TRANSLATOR_LOG_ENABLE"`
	KafkaBrokerList  string `yaml:"kafka-broker-list" env:"KAFKA_BROKER_LIST"`
	KafkaPartition   int    `yaml:"kafka-partition" env:"KAFKA_PARTITION"`
	DataMgrIpAddress string `yaml:"data-manager-ip" env:"DATA_MANAGER_IP_ADDRESS"`
	DataMgrPort      string `yaml:"data-manager-port" env:"DATA_MANAGER_PORT"`
	LogFile          string `yaml:"log-file" env:"IPFIX_LOG_FILE"`
}

var (
	Verbose          = false
	KafkaBrokerList  = "127.0.0.1:9092"
	KafkaTopic       = "vflow.ipfix"
	KafkaPartition   = 0
	DataMgrIpAddress = "127.0.0.1"
	DataMgrPort      = "9000"
	LogFile          = "/var/log/ipfix-translator.log"

	StrDataManager     = "data-manager"
	StrKafkaConGroupID = "ipfixConsGrpID"
	MHConfigFileStr    = "config-file"
	MHConfigFile       = "/etc/ipfix-translator-config.yml"
	IPFixCollection    = "ipfix_flow_collection"
)

var Logger *log.Logger

// ParseArgs parses the arguments as passed from Command Line Interface (CLI)
func ParseArgs(c *cli.Context) error {
	MHConfigFile = c.String(MHConfigFileStr)
	b, err := ioutil.ReadFile(MHConfigFile)
	if err != nil {
		log.Fatalln("config file read error ", err)
	}
	config := ConfigOptions{
		Verbose:          Verbose,
		KafkaBrokerList:  KafkaBrokerList,
		KafkaPartition:   KafkaPartition,
		DataMgrIpAddress: DataMgrIpAddress,
		DataMgrPort:      DataMgrPort,
		LogFile:          LogFile,
	}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		log.Fatalf("Config file %v parse error: %v", MHConfigFile, err)
	}
	KafkaBrokerList = config.KafkaBrokerList //strings.Split(config.KafkaBrokerList, ",")
	KafkaPartition = config.KafkaPartition
	DataMgrIpAddress = config.DataMgrIpAddress
	DataMgrPort = config.DataMgrPort
	Verbose = config.Verbose
	LogFile = config.LogFile
	if LogFile != "" {
		Logger = log.New(os.Stderr, "[ipfix] ", log.Ldate|log.Ltime)
		f, err := os.OpenFile(LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			Logger.Println(err)
		} else {
			Logger.SetOutput(f)
		}
	}
	return nil
}
