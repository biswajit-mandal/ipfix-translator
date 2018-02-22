package kafkaconsumer

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	msghandler "github.com/Juniper/ipfix-translator/msg-handler"
	opts "github.com/Juniper/ipfix-translator/options"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	MHChan = make(chan []byte)
)

func KafkaConsumer() error {
	brokers := opts.KafkaBrokerList
	topic := opts.KafkaTopic

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	k, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               brokers,
		"group.id":                        opts.StrKafkaConGroupID,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"default.topic.config": kafka.ConfigMap{
			"auto.offset.reset": "earliest",
		},
	})
	if err != nil {
		log.Fatalln("Failed to create kafka-consumer ", err)
	}
	k.Subscribe(topic, nil)
	registerMsgHandlers()

	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case sig := <-signalCh:
				log.Println("Interrupt is detected", sig)
				doneCh <- struct{}{}
			case ev := <-k.Events():
				switch e := ev.(type) {
				case kafka.AssignedPartitions:
					fmt.Fprintf(os.Stderr, "%% %v\n", e)
					k.Assign(e.Partitions)
				case kafka.RevokedPartitions:
					fmt.Fprintf(os.Stderr, "%% %v\n", e)
					k.Unassign()
				case *kafka.Message:
					fmt.Printf("%% Message on %s:\n%s\n",
						e.TopicPartition, string(e.Value))
					if opts.Verbose {
						log.Println("Received messages on Kafka-Consumer ", string(e.Value))
					}
					MHChan <- (e.Value)
				case kafka.Error:
					fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				}
			}
		}
	}()

	<-doneCh
	log.Println("Kafka-Consumer Closed")
	return nil
}

func registerMsgHandlers() {
	go func() {
		mh := msghandler.NewMsgHandler(opts.StrDataManager)
		mh.MHChan = MHChan
		if err := mh.Run(); err != nil {
			log.Fatalf("msgHandler run error %v ", err)
		}
	}()
}
