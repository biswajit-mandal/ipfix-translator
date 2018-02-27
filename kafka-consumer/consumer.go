/*
 * Copyright (c) 2018 Juniper Networks, Inc. All rights reserved.
 *
 * file:    main.go
 * details: Kafka Consumer based on confluent-kafka-go library
 *
 */
package kafkaconsumer

import (
	"os"
	"os/signal"

	msghandler "github.com/Juniper/ipfix-translator/msg-handler"
	opts "github.com/Juniper/ipfix-translator/options"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	MHChan = make(chan []byte)
	inCh   = make(chan []byte)
	outCh  = make(chan []byte)
)

//KafkaConsumer constructs Kafka-Consumer based on confluent-kafka-go library
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
		opts.Logger.Fatalln("Failed to create kafka-consumer ", err)
	}
	k.Subscribe(topic, nil)
	inCh, outCh := manageChannels()
	registerMsgHandlers(outCh)

	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case sig := <-signalCh:
				opts.Logger.Println("Interrupt is detected", sig)
				doneCh <- struct{}{}
			case ev := <-k.Events():
				switch e := ev.(type) {
				case kafka.AssignedPartitions:
					opts.Logger.Println(e)
					k.Assign(e.Partitions)
				case kafka.RevokedPartitions:
					opts.Logger.Println(e)
					k.Unassign()
				case *kafka.Message:
					if opts.Verbose {
						opts.Logger.Printf("Received on [%s] messages %s\n", e.TopicPartition, string(e.Value))
					}
					inCh <- (e.Value)
				case kafka.Error:
					opts.Logger.Println(e)
				}
			}
		}
	}()

	<-doneCh
	opts.Logger.Println("Kafka-Consumer Closed")
	return nil
}

func manageChannels() (chan []byte, chan []byte) {
	go func() {
		var inQueue [][]byte
		outChannel := func() chan []byte {
			if len(inQueue) == 0 {
				return nil
			}
			return outCh
		}
		enQueue := func(val []byte) {
			inQueue = append(inQueue, val)
			if opts.Verbose {
				opts.Logger.Println("Data Queued:", string(val))
			}
		}
		deQueue := func() []byte {
			if len(inQueue) == 0 {
				return nil
			}
			curVal := inQueue[0]
			inQueue = inQueue[1:]
			if opts.Verbose {
				opts.Logger.Println("Data served from Queue:", string(curVal))
			}
			return curVal
		}
		for len(inQueue) > 0 || inCh != nil {
			select {
			case v, ok := <-inCh:
				if !ok {
					inCh = nil
				} else {
					enQueue(v)
				}
			case outChannel() <- deQueue():
			}
		}
		close(outCh)
	}()
	return inCh, outCh
}

func registerMsgHandlers(outCh chan []byte) {
	go func(outCh chan []byte) {
		mh := msghandler.NewMsgHandler(opts.StrDataManager)
		mh.MHChan = outCh
		if err := mh.Run(); err != nil {
			opts.Logger.Fatalf("msgHandler run error %v ", err)
		}
	}(outCh)
}
