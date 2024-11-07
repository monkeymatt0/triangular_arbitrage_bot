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

// @remind : This one should implement only one comunication stream, always in this package you
// have to create a composition of crypto streamer and then use them in a go routine in parallel
// @todo : Utilizzare go-binance-url-builder con isWebSocketStream = true
// @todo : Creare mecanismo di comunicazione con binance tramite stream

// @todo : Step 1: Get proper URL (ALMOST DONE)
// @todo : Step 2: Fetch the stream from binance (DONE)
// @todo : Step 3: Create a channel that send a signal each time it receive a stream from binance (TODO)
// @todo : Step 4: Create data structure do make marshal and unmarshal the json object(NEXT TO DO)

/*
-------------- CONSIDERATION
You should have a a way to send (via channel) the last price for each operation buy/sell, so you should tell
during creation phase which type of streamer this should be if a buy streamer or a sell streamer
*/

type CryptoStreamer struct {
	Testing       bool
	SymbolChannel string
	depthUpdate   m.StreamDepthModel
	buyPrice      float64
	buyQty        float64
	sellPrice     float64
	sellQty       float64
}

func (cs *CryptoStreamer) Listen(dataCh chan string) {
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
		if cs.buyPrice, err = strconv.ParseFloat(cs.depthUpdate.Asks[0][0], 64); err != nil {
			log.Println(err)
		}
		if cs.buyQty, err = strconv.ParseFloat(cs.depthUpdate.Asks[0][1], 64); err != nil {
			log.Println(err)
		}
		if cs.sellPrice, err = strconv.ParseFloat(cs.depthUpdate.Bids[0][0], 64); err != nil {
			log.Println(err)
		}
		if cs.sellQty, err = strconv.ParseFloat(cs.depthUpdate.Bids[0][1], 64); err != nil {
			log.Println(err)
		}
		dataCh <- fmt.Sprintf("++++++++++++++++++++++++++++++\n%s\n------------------------------\nbuyPrice: %f - Quantity: %f\nsellPrice: %f - Quantity: %f\n++++++++++++++++++++++++++++++\n", cs.SymbolChannel, cs.buyPrice, cs.buyQty, cs.sellPrice, cs.sellQty)
	}
}
