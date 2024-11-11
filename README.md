# triangular_bot

Traingular bot is a bot that checks for market discrepancies anche simulate buy a crypto, in a triangular way from the same exchange
and check for eventual opportunity.

## How it works

The idea behind this bot is to check via a websocket stream the price (In this specific scenario) the current pairs:

- BTCUSDT
- ETHBTC
- ETHUSDT

And understand if some type of discrepancy happen.

This project uses a package that I developer as MVP (so it needs refector), to communicate with Binance (Used for this experiment).

When a discrepancy is detected, I check for the liquidity of the market on the price of the moment to understand if the liquidity is enough to proceed with the orders.

Keep in mind that this is an example e needs some optimization that will not happen shortly.

## Spot opportunity

Spot opportunity it's the challenge here, so the things that can make the differences are:

- Choosen pairs
- Market velocity

This bot use the following:

- BTCUSDT
- ETHBTC
- ETHUSDT

All these market are really in rapid change even now you are reading.

I consider this bot as a way to experiment, the opportunity are hard to find but most likely with the proper
adjustment can be a kind let's say income.

Last update to the project: 11/11/2024
