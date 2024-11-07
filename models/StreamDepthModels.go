package model

// Model representing the object sent during the stream
type StreamDepthModel struct {
	LastUpdateID uint64     `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}
