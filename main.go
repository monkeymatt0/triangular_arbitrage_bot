package main

import (
	"fmt"
	cs "triangular_arbitrage_bot/crypto_streamer"

	gbub "github.com/monkeymatt0/go-binance-url-builder"
)

func main() {
	channelSymbols := []string{gbub.BTCUSDT, gbub.ETHBTC, gbub.ETHUSDT}
	streamers := &cs.CryptoStreamers{}
	streamers.New(channelSymbols, true)

	dataChs := make([]chan string, 3)

	// Memory allocation for the channels
	for i := 0; i < len(dataChs); i++ {
		dataChs[i] = make(chan string)
	}

	go streamers.Streams[cs.BTCUSDT].Listen(dataChs[cs.BTCUSDT])
	go streamers.Streams[cs.ETHBTC].Listen(dataChs[cs.ETHBTC])
	go streamers.Streams[cs.ETHUSDT].Listen(dataChs[cs.ETHUSDT])

	for {
		select {
		case btcUsdtData := <-dataChs[cs.BTCUSDT]:
			fmt.Println(btcUsdtData)
		case ethBtcData := <-dataChs[cs.ETHBTC]:
			fmt.Println(ethBtcData)
		case ethUsdtData := <-dataChs[cs.ETHUSDT]:
			fmt.Println(ethUsdtData)
		}
	}
}
