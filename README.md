# Gemini API wrapper for Golang
Utility types and methods to facilitate the use of the Gemini bitcoin exchange REST API

## Usage

```
// Initialize
const (
	GEMINI_API_KEY    = ""
	GEMINI_API_SECRET = ""
	LIVE              = false
)

api := gemini.New(LIVE, GEMINI_API_KEY, GEMINI_API_SECRET)

// Get the first ask and first bid from the order book
limitBids := 1
limitAsks := 1
orderBook, err := api.OrderBook("btcusd", limitBids, limitAsks)

// Fetch your active orders
activeOrders, err := api.ActiveOrders()

// Place an order
clientOrderId := "20161229-838492"
btcAmount := 4.75
askPrice := 925.5
order, err := api.NewOrder("btcusd", clientOrderId, btcAmount, askPrice, "buy", []string{"immediate-or-cancel"})


// see code for other available methods
```
