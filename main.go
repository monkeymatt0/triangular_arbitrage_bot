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

	go streamers.Streams[cs.BTCUSDT].Listen(dataCh)
	go streamers.Streams[cs.ETHBTC].Listen(dataCh)
	go streamers.Streams[cs.ETHUSDT].Listen(dataCh)

	// @todo : Parallelize the channels so that I don't have bottleneck on a single channel
	for data := range dataCh {
		fmt.Println(data)
	}
}
