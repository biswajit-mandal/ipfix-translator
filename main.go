package main

import (
	"log"
	"os"

	kc "github.com/Juniper/ipfix-translator/kafka-consumer"
	opts "github.com/Juniper/ipfix-translator/options"
	"github.com/urfave/cli"
)

func handleKafkaConsumer(c *cli.Context) error {
	opts.ParseArgs(c)
	kc.KafkaConsumer()
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "Kafka Consumer CLI"
	app.Commands = []cli.Command{
		{
			Name:  "kafka-consumer",
			Usage: "Kafka Consumer",
			Flags: []cli.Flag{
				cli.StringFlag{Name: opts.MHConfigFileStr, Value: opts.MHConfigFile,
					Usage: "The config file"},
			},
			Action: handleKafkaConsumer,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Application error: %s", err)
	}
}
