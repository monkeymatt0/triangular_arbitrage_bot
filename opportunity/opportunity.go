package opportunity

import (
	"fmt"
	ce "triangular_arbitrage_bot/custom_errors"
)

// This pacakge is the responsable for the calculation of the oportunity
// It took a triangular coins and caluculate the possible opportunity to profit
// from a market discrepancy

// COIN PAIRS EXAMPLE:
// BTC/USDT
// ETH/BTC
// ETH/USDT

// We will use the above coins as reference
type Opportunity struct {
	FirstCoinPrice  float64 // This will be the price of BTC in USDT
	FirstCoinQty    float64 // This will be the quantity of BTC in USDT
	SecondCoinPrice float64 // This will be the price of ETH in BTC
	SecondCoinQty   float64 // This will be the quantity of ETH in BTC
	ThirdCoinPrice  float64 // This will be the price of ETH in USDT
	ThirdCoinQty    float64 // This will be the quantity of ETH in USDT

	// Since we are interacting with an exchange we need also to consider
	// the fees that the exange takes when executing an order.
	// In this case this will value will be 0.1000% since the exchange for now is binance
	Fee float64
}

func (o *Opportunity) New() {
	o.Fee = 0.1000 // Binance Fees @remind : For now it's fine like this
}

// This function will check if actually the opportunity is profitable
// usdt will be the quantity we want to use for that opportunity
// @param buyCoinQty is the coin I will use to start check the opportunity, following the example -> USDT
func (o *Opportunity) IsProfitable(buyCoinQty float64) (bool, error) {
	fmt.Println("################## Check profitablity ##################")
	// 1) Buy the first coin, we will use this formula to have the coins bought and check for it's liquidity on that level
	if !o.haveLiquidity(buyCoinQty, o.FirstCoinPrice, o.FirstCoinQty) {
		fmt.Println("################## No qty for first ##################")
		return false, nil
	}
	firstCoinQty, err1 := o.simulateExchangeOrder("BUY", buyCoinQty, o.FirstCoinPrice) // This will be the quantity of BTC(Following the example)
	if err1 != nil {
		return false, err1
	}

	// 2) Buy the second coin with the coin bought in the previous step, and check for the liquidity on that level
	if !o.haveLiquidity(firstCoinQty, o.SecondCoinPrice, o.SecondCoinQty) {
		fmt.Println("################## No qty for second ##################")
		return false, nil
	}
	secondCoinQty, err2 := o.simulateExchangeOrder("BUY", firstCoinQty, o.SecondCoinPrice) // This will be the quantity of ETH(Following the example)
	if err2 != nil {
		return false, err2
	}

	// 3) Sell the coin and check for the liquidity on that level
	if !o.haveLiquidity(secondCoinQty, o.ThirdCoinPrice, o.ThirdCoinQty) {
		fmt.Println("################## No qty for third ##################")
		return false, nil
	}
	stableCoinFinalQty, err3 := o.simulateExchangeOrder("SELL", secondCoinQty, o.ThirdCoinPrice) // This will be the final amount of stable coin after the exchanges
	if err3 != nil {
		return false, err3
	}

	return stableCoinFinalQty > buyCoinQty, nil
}

// This function check if for a certain price we have the right liquidity
// on that level to make in a way that the market order execute properly
// without switch to the next level on the order book
func (o *Opportunity) haveLiquidity(buyCoinQty float64, price float64, quantity float64) bool {
	return buyCoinQty < price*quantity
}

// This function will simulate the exchange order in terms of money
// it will do the calculation for BUY/SELL and then it will
// subtract the fees
func (o *Opportunity) simulateExchangeOrder(side string, numeratorCoinQty float64, denominatorCoinQty float64) (float64, error) {
	// Check if the proper side is passed
	if side != "BUY" && side != "SELL" {
		return 0.0, &ce.SideError{
			Err: ce.ExchangeError{
				Context: "Opportunity.simulateExchangeOrder",
			},
		}
	}

	switch side {
	case "BUY":
		return ((numeratorCoinQty / denominatorCoinQty) * (1 - o.Fee)), nil
	case "SELL":
		return ((numeratorCoinQty * denominatorCoinQty) * (1 - o.Fee)), nil

	}

	return 0.0, nil // This should never happen
}
