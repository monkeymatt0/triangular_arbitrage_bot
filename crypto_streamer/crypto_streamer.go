package crypto_streamer

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"strings"

	m "triangular_arbitrage_bot/models"

	gbub "github.com/monkeymatt0/go-binance-url-builder"

	"github.com/gorilla/websocket"
)

type OrderSide uint

const (
	BUY OrderSide = iota
	SELL
)

// @remind : Optimization: Create an interface with Listen function and implement specifically for buy streamer and sell streamer
// In this way you can be faster on the unmarshal phase

/*
-------------- CONSIDERATION
You should have a a way to send (via channel) the last price for each operation buy/sell, so you should tell
during creation phase which type of streamer this should be if a buy streamer or a sell streamer
*/
// This crypto streamr will own the connection to binance stream
type CryptoStreamer struct {
	Testing       bool
	Side          OrderSide
	SymbolChannel string
	depthUpdate   m.StreamDepthModel
	// @remind : all these data can be replaced from a ChannelData
	chData m.ChannelData
}

// Listen function just listen to the channel using a gorilla web socket
//
// @param detaCh is the channel where the data will be sent once the arrive
func (cs *CryptoStreamer) Listen(dataCh chan m.ChannelData) {
	u := url.URL{Scheme: gbub.SECURE_WEB_SOCKET}
	if cs.Testing {
		u.Host = gbub.TEST_WSS_HOST
	} else {
		u.Host = gbub.PRODUCTION_WSS_HOST
	}
	u.Path = strings.Join([]string{"", gbub.WSS_API, cs.SymbolChannel}, "/")

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error during connection with the socket: ", err)
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("Error while reading: ", err)
			return
		}
		if err := json.Unmarshal(message, &cs.depthUpdate); err != nil {
			log.Println(err)
		}
		switch cs.Side {
		case BUY: // This will set the Buy price and quantity
			if cs.chData.Price, err = strconv.ParseFloat(cs.depthUpdate.Asks[0][0], 64); err != nil {
				log.Println(err)
			}
			if cs.chData.Quantity, err = strconv.ParseFloat(cs.depthUpdate.Asks[0][1], 64); err != nil {
				log.Println(err)
			}
			break
		case SELL: // This will set the Sell price and quantity
			if cs.chData.Price, err = strconv.ParseFloat(cs.depthUpdate.Bids[0][0], 64); err != nil {
				log.Println(err)
			}
			if cs.chData.Quantity, err = strconv.ParseFloat(cs.depthUpdate.Bids[0][1], 64); err != nil {
				log.Println(err)
			}
			break
		}
		dataCh <- cs.chData
	}
}
