package crypto_streamer

import (
	"encoding/json"
	"fmt"
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

/*
-------------- CONSIDERATION
You should have a a way to send (via channel) the last price for each operation buy/sell, so you should tell
during creation phase which type of streamer this should be if a buy streamer or a sell streamer
*/
type Streamer interface {
	Listen(dataCh chan m.ChannelData)
}

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

type BuyCryptoStreamer struct {
	Testing       bool
	SymbolChannel string
	depthUpdate   m.StreamDepthModel
	chData        m.ChannelData
}

func (bcs *BuyCryptoStreamer) Listen(ch chan m.ChannelData) {
	u := url.URL{Scheme: gbub.SECURE_WEB_SOCKET}
	if bcs.Testing {
		u.Host = gbub.TEST_WSS_HOST
	} else {
		u.Host = gbub.PRODUCTION_WSS_HOST
	}
	u.Path = strings.Join([]string{"", gbub.WSS_API, bcs.SymbolChannel}, "/")

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
		if err := json.Unmarshal(message, &bcs.depthUpdate); err != nil {
			log.Println(err)
		}
		fmt.Println("Bid: ", bcs.depthUpdate.Asks)
		if len(bcs.depthUpdate.Asks) > 0 {
			if bcs.chData.Price, err = strconv.ParseFloat(bcs.depthUpdate.Asks[0][0], 64); err != nil {
				fmt.Println("BUY[0][0]: ", bcs.depthUpdate)
				log.Println(err)
			}
			if bcs.chData.Quantity, err = strconv.ParseFloat(bcs.depthUpdate.Asks[0][1], 64); err != nil {
				fmt.Println("BUY[0][1]: ", bcs.depthUpdate)
				log.Println(err)
			}
		} else {
			fmt.Println("Zero data Asks")
		}
		ch <- bcs.chData
	}
}

type SellCryptoStreamer struct {
	Testing       bool
	SymbolChannel string
	depthUpdate   m.StreamDepthModel
	chData        m.ChannelData
}

func (scs *SellCryptoStreamer) Listen(ch chan m.ChannelData) {
	u := url.URL{Scheme: gbub.SECURE_WEB_SOCKET}
	if scs.Testing {
		u.Host = gbub.TEST_WSS_HOST
	} else {
		u.Host = gbub.PRODUCTION_WSS_HOST
	}
	u.Path = strings.Join([]string{"", gbub.WSS_API, scs.SymbolChannel}, "/")

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
		if err := json.Unmarshal(message, &scs.depthUpdate); err != nil {
			log.Println(err)
		}
		fmt.Println("Sell: ", scs.depthUpdate.Bids)
		if len(scs.depthUpdate.Bids) > 0 {
			if scs.chData.Price, err = strconv.ParseFloat(scs.depthUpdate.Bids[0][0], 64); err != nil {
				fmt.Println("SELL[0][0]: ", scs.depthUpdate)
				log.Println(err)
			}
			if scs.chData.Quantity, err = strconv.ParseFloat(scs.depthUpdate.Bids[0][1], 64); err != nil {
				fmt.Println("SELL[0][1]: ", scs.depthUpdate)
				log.Println(err)
			}
		} else {
			fmt.Println("Zero data Bids")
		}
		ch <- scs.chData
	}
}
