package binanceconsumer

import (
	"encoding/json"

	binance_connector "github.com/binance/binance-connector-go"
)

type TickerData = binance_connector.WsMarketTickerStatEvent

type Handler func(data []TickerData)
type ErrHandler func(err error)

func PrettyPrint(data interface{}) string {
	s, _ := json.MarshalIndent(data, "", "\t")
	return string(s)
}

type BinanceWSConsumer struct {
	client         *binance_connector.WebsocketStreamClient
	userHandler    Handler
	userErrHandler ErrHandler

	subscribed bool
	doneCh     chan struct{}
	stopCh     chan struct{}
}

// Create Binance ws Consumer with handlers
func NewBinanceWSConsumer(h Handler, eh ErrHandler) *BinanceWSConsumer {
	client := binance_connector.NewWebsocketStreamClient(false, "wss://stream.binance.us:9443")
	return &BinanceWSConsumer{
		client:         client,
		userHandler:    h,
		userErrHandler: eh,
	}
}

// wrapHandler converts the user's handler to the Binance connector expected format
func (b *BinanceWSConsumer) wrapHandler() binance_connector.WsAllMarketTickersStatHandler {
	return func(event binance_connector.WsAllMarketTickersStatEvent) {
		var tickerData []TickerData
		for _, ticker := range event {
			tickerData = append(tickerData, *ticker)
		}
		b.userHandler(tickerData)
	}
}

func (b *BinanceWSConsumer) wrapErrHandler() binance_connector.ErrHandler {
	return func(err error) {
		b.userErrHandler(err)
	}
}

func (b *BinanceWSConsumer) SubscribeAll() {
	doneCh, stopCh, err := b.client.WsAllMarketTickersStatServe(b.wrapHandler(), b.wrapErrHandler())
	if err != nil {
		b.userErrHandler(err)
		return
	}

	b.doneCh = doneCh
	b.stopCh = stopCh
	b.subscribed = true
}

func (b *BinanceWSConsumer) IsSubscribed() bool {
	return b.subscribed
}

func (b *BinanceWSConsumer) Unsubscribe() {
	if b.subscribed {
		b.stopCh <- struct{}{}
		<-b.doneCh
		b.subscribed = false
	}
}
