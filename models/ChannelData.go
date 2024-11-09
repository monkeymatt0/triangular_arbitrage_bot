package model

// This data structure will hold the value of the price
// and the quantity of crypto in that price
// This means that this data structure will hold information about
// the liquidity on a specific level
type ChannelData struct {
	Symbol   string
	Price    float64
	Quantity float64
}
