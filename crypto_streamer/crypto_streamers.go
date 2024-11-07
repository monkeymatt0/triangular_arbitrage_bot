package crypto_streamer

type CryptoStreamers struct {
	Streams []CryptoStreamer
}

func (cs *CryptoStreamers) New(channels []string, test bool) (*CryptoStreamers, error) {
	for _, channel := range channels {
		cs.Streams = append(cs.Streams, CryptoStreamer{Testing: test, SymbolChannel: channel})
	}
	return cs, nil
}
