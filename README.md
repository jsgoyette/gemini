# gemini-api for golang
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


// Get the first ask price from the order book
orderBook, err := api.OrderBook("btcusd", 1, 1)

var askPrice float64
if len(orderBook.Asks) > 0 {
	askPrice = orderBook.Asks[0].Price
}


// Fetch your active orders
activeOrders, err := api.ActiveOrders()


// Place an order
clientOrderId := "20161229-838492"
btcAmount := 4.75
askPrice := 925.5
order, err := api.NewOrder("btcusd", clientOrderId, btcAmount, askPrice, "buy", []string{"immediate-or-cancel"})


// see code for other available methods
```
