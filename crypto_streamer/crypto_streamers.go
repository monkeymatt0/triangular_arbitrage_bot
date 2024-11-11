package crypto_streamer

import model "triangular_arbitrage_bot/models"

type CryptoStreamers struct {
	Streams []Streamer
}

func (cs *CryptoStreamers) New(channels []string, side []OrderSide, test bool) (*CryptoStreamers, error) {
	for index, channel := range channels {
		if index != len(channels)-1 {
			cs.Streams = append(cs.Streams, &BuyCryptoStreamer{Testing: test, SymbolChannel: channel, chData: model.ChannelData{Symbol: channel}})
			continue
		}
		cs.Streams = append(cs.Streams, &SellCryptoStreamer{Testing: test, SymbolChannel: channel, chData: model.ChannelData{Symbol: channel}})
	}
	return cs, nil
}
