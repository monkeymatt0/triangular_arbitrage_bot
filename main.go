package main

import (
	"fmt"
	cs "triangular_arbitrage_bot/crypto_streamer"

	gbub "github.com/monkeymatt0/go-binance-url-builder"
)

func main() {
	channels := []string{gbub.BTCUSDT, gbub.ETHBTC, gbub.ETHUSDT}
	streamers := &cs.CryptoStreamers{}
	streamers.New(channels, true)

	dataCh := make(chan string)

	go streamers.Streams[0].Listen(dataCh)
	go streamers.Streams[1].Listen(dataCh)
	go streamers.Streams[2].Listen(dataCh)

	for data := range dataCh {
		fmt.Println(data)
	}
}
