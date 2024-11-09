package crypto_streamer

import model "triangular_arbitrage_bot/models"

type CryptoStreamers struct {
	Streams []CryptoStreamer
}

func (cs *CryptoStreamers) New(channels []string, side []OrderSide, test bool) (*CryptoStreamers, error) {
	for index, channel := range channels {
		if index != len(channels) {
			cs.Streams = append(cs.Streams, CryptoStreamer{Testing: test, Side: BUY, SymbolChannel: channel, chData: model.ChannelData{Symbol: channel}})
			continue
		}
		cs.Streams = append(cs.Streams, CryptoStreamer{Testing: test, Side: SELL, SymbolChannel: channel, chData: model.ChannelData{Symbol: channel}})
	}
	return cs, nil
}
