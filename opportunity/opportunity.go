package opportunity

import (
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
	SecondCoinPrice float64 // This will be the price of ETH in BTC
	ThirdCoinPrice  float64 // This will be the price of ETH in USDT

	// Since we are interacting with an exchange we need also to consider
	// the fees that the exange takes when executing an order.
	// In this case this will value will be 0.1000% since the exchange for now is binance
	Fee float64
}

// This function will check if actually the opportunity is profitable
// usdt will be the quantity we want to use for that opportunity
func (o *Opportunity) IsProfitable(stableCoinQty float64) (bool, error) {
	// 1) Buy the first coin, we will use this formula to have the coins bought
	firstCoinQty, err1 := o.simulateExchangeOrder("BUY", stableCoinQty, o.FirstCoinPrice) // This will be the quantity of BTC(Following the example)
	if err1 != nil {
		return false, err1
	}

	// 2) Buy the second coin with the coin bought in the previous step
	secondCoinQty, err2 := o.simulateExchangeOrder("BUY", firstCoinQty, o.SecondCoinPrice) // This will be the quantity of ETH(Following the example)
	if err2 != nil {
		return false, err2
	}

	// 3) Sell the coin with
	stableCoinFinalQty, err3 := o.simulateExchangeOrder("SELL", secondCoinQty, o.ThirdCoinPrice) // This will be the final amount of stable coin after the exchanges
	if err3 != nil {
		return false, err3
	}

	return stableCoinFinalQty > stableCoinQty, nil
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
