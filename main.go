package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	kafkaproducer "github.com/Quantlytic/AlpacaConsumer/pkg/kafkaproducer"
	binanceconsumer "github.com/Quantlytic/BinanceConsumer/internal/BinanceConsumer"
	"github.com/Quantlytic/BinanceConsumer/internal/config"
)

func errHandler(err error) {
	log.Fatal(err)
}

func main() {
	log.Printf("KAFKA_BROKERS env: %s", os.Getenv("KAFKA_BROKERS"))

	cfg := config.Load()

	p, err := kafkaproducer.NewKafkaProducer(kafkaproducer.KafkaProducerConfig{
		Servers:  cfg.KafkaBrokers,
		ClientId: cfg.KafkaClientId,
		Acks:     "all",
	})
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}

	// Helper function to handle JSON marshaling and Kafka producing
	produceToKafka := func(topic, symbol string, data interface{}) error {
		dataJSON, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error marshaling data to JSON: %v", err)
			return err
		}

		err = p.Produce(topic, []byte(symbol), dataJSON)
		if err != nil {
			log.Printf("Error producing message to Kafka: %v", err)
			return err
		}
		return nil
	}

	handler := func(data []binanceconsumer.TickerData) {
		for _, d := range data {
			// Pretty print for logging
			out := binanceconsumer.PrettyPrint(d)
			log.Printf("Ticker Event\n%s\n", out)

			// Send to Kafka
			produceToKafka(cfg.KafkaTopic, d.Symbol, d)
		}
	}

	consumer := binanceconsumer.NewBinanceWSConsumer(handler, errHandler)

	// Gracefully shutdown on SIGINT or SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	consumer.SubscribeAll()
	log.Println("Binance Consumer started. Press Ctrl+C to gracefully shutdown...")

	<-sigChan

	if consumer.IsSubscribed() {
		consumer.Unsubscribe()
		log.Println("Successfully unsubscribed from Binance WebSocket")
	}
}
