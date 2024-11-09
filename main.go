package main

import (
	"fmt"
	cs "triangular_arbitrage_bot/crypto_streamer"
	m "triangular_arbitrage_bot/models"
	o "triangular_arbitrage_bot/opportunity"

	gbub "github.com/monkeymatt0/go-binance-url-builder"
)

// For now the quantity of USDT will be fixed
// For later implementation this could become a parameter

func main() {
	channelSymbols := []string{gbub.BTCUSDT, gbub.ETHBTC, gbub.ETHUSDT}
	streamers := &cs.CryptoStreamers{}
	sides := []cs.OrderSide{cs.BUY, cs.BUY, cs.SELL}
	streamers.New(channelSymbols, sides, true)
	opportunityChecker := &o.Opportunity{}
	opportunityChecker.New() // This will set the fee

	dataChs := make([]chan m.ChannelData, 3)
	receivedData := make([]bool, 3)
	fmt.Println(receivedData)
	// Memory allocation for the channels
	for i := 0; i < len(dataChs); i++ {
		dataChs[i] = make(chan m.ChannelData)
	}

	go streamers.Streams[cs.BTCUSDT].Listen(dataChs[cs.BTCUSDT])
	go streamers.Streams[cs.ETHBTC].Listen(dataChs[cs.ETHBTC])
	go streamers.Streams[cs.ETHUSDT].Listen(dataChs[cs.ETHUSDT])

	for {
		// @todo : whenever I receive a signal I should
		select {
		case ethBtcData := <-dataChs[cs.ETHBTC]:
			fmt.Println(ethBtcData)
		case ethUsdtData := <-dataChs[cs.ETHUSDT]:
			fmt.Println(ethUsdtData)
		case btcUsdtData := <-dataChs[cs.BTCUSDT]:
			receivedData[cs.BTCUSDT] = true
			fmt.Println(btcUsdtData)
		}

		if receivedData[cs.BTCUSDT] && receivedData[cs.ETHBTC] && receivedData[cs.ETHUSDT] {
			isPorfitable, err := opportunityChecker.IsProfitable(1000) // @remind : This is mocked and need to be replaced with correct values
			if err != nil {
				fmt.Println(err)
				return
			}
			if isPorfitable {
				// @todo : Place the orders in order to execute them in a triangular way
			}
		}
	}
}
