package options

import (
	"io/ioutil"
	"log"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type ConfigOptions struct {
	Verbose          bool   `yaml:verbose" env:IPFIX_TRANSLATOR_LOG_ENABLE`
	KafkaBrokerList  string `yaml:"kafka-broker-list" env:"KAFKA_BROKER_LIST"`
	KafkaPartition   int    `yaml:kafka-partition" env:"KAFKA_PARTITION"`
	DataMgrIpAddress string `yaml:"data-manager-ip" env:"DATA_MANAGER_IP_ADDRESS"`
	DataMgrPort      string `yaml:"data-manager-port" env:"DATA_MANAGER_PORT"`
}

var (
	Verbose          = false
	KafkaBrokerList  = "127.0.0.1:9092"
	KafkaTopic       = "vflow.ipfix"
	KafkaPartition   = 0
	DataMgrIpAddress = "127.0.0.1"
	DataMgrPort      = "9000"

	StrDataManager  = "data-manager"
    StrKafkaConGroupID  = "ipfixConsGrpID"
	MHConfigFileStr = "config-file"
	MHConfigFile    = "/etc/ipfix-translator-config.yml"
	IPFixCollection = "ipfix_flow_collection"
)

func ParseArgs(c *cli.Context) error {
	MHConfigFile = c.String(MHConfigFileStr)
	b, err := ioutil.ReadFile(MHConfigFile)
	if err != nil {
		log.Fatalln("config file read error ", err)
	}
	config := ConfigOptions{
		Verbose:          Verbose,
		KafkaBrokerList:  KafkaBrokerList,//strings.Join(KafkaBrokerList, ","),
		KafkaPartition:   KafkaPartition,
		DataMgrIpAddress: DataMgrIpAddress,
		DataMgrPort:      DataMgrPort,
	}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		log.Fatalf("Config file %v parse error: %v", MHConfigFile, err)
	}
	KafkaBrokerList = config.KafkaBrokerList//strings.Split(config.KafkaBrokerList, ",")
	KafkaPartition = config.KafkaPartition
	DataMgrIpAddress = config.DataMgrIpAddress
	DataMgrPort = config.DataMgrPort
	Verbose = config.Verbose
	return nil
}
