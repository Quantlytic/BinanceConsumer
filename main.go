package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	binanceconsumer "github.com/Quantlytic/BinanceConsumer/internal/BinanceConsumer"
	"github.com/Quantlytic/BinanceConsumer/internal/config"
)

func handler(data []binanceconsumer.TickerData) {
	for _, d := range data {
		out := binanceconsumer.PrettyPrint(d)
		log.Printf("Ticker Event\n%s\n", out)
	}
}

func errHandler(err error) {
	log.Fatal(err)
}

func main() {
	config.Load()

	consumer := binanceconsumer.NewBinanceWSConsumer(handler, errHandler)

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
