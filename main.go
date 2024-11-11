package main

import (
	"fmt"
	"time"
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
	// receivedData := make([]bool, 3)
	// Memory allocation for the channels
	for i := 0; i < len(dataChs); i++ {
		dataChs[i] = make(chan m.ChannelData)
	}

	go streamers.Streams[cs.BTCUSDT].Listen(dataChs[cs.BTCUSDT])
	go streamers.Streams[cs.ETHBTC].Listen(dataChs[cs.ETHBTC])
	go streamers.Streams[cs.ETHUSDT].Listen(dataChs[cs.ETHUSDT])

	btcUsdtPrevPrice := 0.0
	ethBtcPrevPrice := 0.0
	ethUsdtPrevPrice := 0.0

	opportunityChecker.FirstCoinPrice = btcUsdtPrevPrice
	opportunityChecker.SecondCoinPrice = ethBtcPrevPrice
	opportunityChecker.ThirdCoinPrice = ethUsdtPrevPrice

	timeout := 100 * time.Millisecond // 0.2 seconds
	for {
		fmt.Println(time.Now())
		if opportunityChecker.FirstCoinPrice != 0.0 && opportunityChecker.SecondCoinPrice != 0.0 && opportunityChecker.ThirdCoinPrice != 0.0 {
			fmt.Println("******************* RECEIVED ALL *******************")
			isPorfitable, err := opportunityChecker.IsProfitable(100)
			if err != nil {
				fmt.Println(err)
				return
			}
			if isPorfitable {
				fmt.Println("#################### Profitable ####################")
				continue
			}
			fmt.Println("#################### Not Good ####################")
		}
		select {
		case ethBtcData := <-dataChs[cs.ETHBTC]:
			if ethBtcPrevPrice == 0.0 || ethBtcPrevPrice != opportunityChecker.SecondCoinPrice {
				ethBtcPrevPrice = ethBtcData.Price
			} else {
				fmt.Println("ethBtc same")
			}
			opportunityChecker.SecondCoinPrice = ethBtcData.Price
			opportunityChecker.SecondCoinQty = ethBtcData.Quantity
			fmt.Println(ethBtcData)
		case btcUsdtData := <-dataChs[cs.BTCUSDT]:
			if btcUsdtPrevPrice == 0.0 || btcUsdtPrevPrice != opportunityChecker.FirstCoinPrice {
				btcUsdtPrevPrice = btcUsdtData.Price
			} else {
				fmt.Println("btcUsdt same")
			}
			opportunityChecker.FirstCoinPrice = btcUsdtData.Price
			opportunityChecker.FirstCoinQty = btcUsdtData.Quantity
			fmt.Println(btcUsdtData)
		case ethUsdtData := <-dataChs[cs.ETHUSDT]:
			if ethUsdtPrevPrice == 0.0 || ethUsdtPrevPrice != opportunityChecker.ThirdCoinPrice {
				ethUsdtPrevPrice = ethUsdtData.Price
			} else {
				fmt.Println("ethUsdt same")
			}
			opportunityChecker.ThirdCoinPrice = ethUsdtData.Price
			opportunityChecker.ThirdCoinQty = ethUsdtData.Quantity
			fmt.Println(ethUsdtData)
		case <-time.After(timeout):
			fmt.Println("Timeout, skipping this cycle")
		}

	}
}
