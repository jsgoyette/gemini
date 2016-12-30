package gemini

import (
	"encoding/json"
	// "fmt"
	"strconv"
)

// Active Orders
// [
// 	{
// 		"order_id":"68053346",
// 		"id":"68053346",
// 		"symbol":"btcusd",
// 		"exchange":"gemini",
// 		"avg_execution_price":"0.00",
// 		"side":"buy",
// 		"type":"exchange limit",
// 		"timestamp":"1483052547",
// 		"timestampms":1483052547532,
// 		"is_live":true,
// 		"is_cancelled":false,
// 		"is_hidden":false,
// 		"was_forced":false,
// 		"executed_amount":"0",
// 		"remaining_amount":"0.212",
// 		"price":"123.45",
// 		"original_amount":"0.212"
// 	}
// ]
func (g *GeminiAPI) ActiveOrders() ([]Order, error) {

	url := g.url + ACTIVE_ORDERS_URL
	req := RequestParams{
		"request": ACTIVE_ORDERS_URL,
		"nonce":   getNonce(),
	}

	var orders []Order

	body, err := request("POST", url, g.prepPayload(req), nil)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, &orders)

	return orders, nil
}

// Past Trades
// [
// 	{
// 		"price":"900.97",
// 		"amount":"0.00221429",
// 		"timestamp":1482907289,
// 		"timestampms":1482907289592,
// 		"type":"Buy",
// 		"aggressor":true,
// 		"fee_currency":"USD",
// 		"fee_amount":"0.00498752215325",
// 		"tid":68053255,
// 		"order_id":"68053253",
// 		"exchange":"gemini",
// 		"is_auction_fill":false,
// 	}
// ]
func (g *GeminiAPI) PastTrades(symbol string, limitTrades int, timestamp int64) ([]Trade, error) {

	url := g.url + PAST_TRADES_URL

	req := RequestParams{
		"request":      PAST_TRADES_URL,
		"nonce":        getNonce(),
		"symbol":       symbol,
		"limit_trades": limitTrades,
		"timestamp":    timestamp,
	}

	body, err := request("POST", url, g.prepPayload(req), nil)
	if err != nil {
		return nil, err
	}

	var trades []Trade
	json.Unmarshal(body, &trades)

	return trades, nil
}

// New Order
// {
// 	"order_id":"68053332",
// 	"id":"68053332",
// 	"symbol":"btcusd",
// 	"exchange":"gemini",
// 	"avg_execution_price":"0.00",
// 	"side":"buy",
// 	"type":"exchange limit",
// 	"timestamp":"1483050634",
// 	"timestampms":1483050634767,
// 	"is_live":true,
// 	"is_cancelled":false,
// 	"is_hidden":false,
// 	"was_forced":false,
// 	"executed_amount":"0",
// 	"remaining_amount":"0.12",
// 	"price":"123.45",
// 	"original_amount":"0.12"
// }
// {
// 	"result":"error",
// 	"reason":"InsufficientFunds",
// 	"message":"Failed to place buy order on symbol 'BTCUSD' for price $123.45 and quantity 120,120 BTC due to insufficient funds"
// }
func (g *GeminiAPI) NewOrder(symbol string, clientOrderId string, amount float64, price float64, side string, options []string) (Order, error) {

	url := g.url + NEW_ORDER_URL
	req := RequestParams{
		"request":         NEW_ORDER_URL,
		"nonce":           getNonce(),
		"client_order_id": clientOrderId,
		"symbol":          symbol,
		"amount":          strconv.FormatFloat(amount, 'f', -1, 64),
		"price":           strconv.FormatFloat(price, 'f', -1, 64),
		"side":            side,
		"type":            "exchange limit",
	}

	if options != nil {
		req["options"] = options
	}

	var order Order

	body, err := request("POST", url, g.prepPayload(req), nil)
	if err != nil {
		return order, err
	}

	// fmt.Println(string(body))

	json.Unmarshal(body, &order)

	return order, nil
}

// {
// 	"order_id":"68053346",
// 	"id":"68053346",
// 	"symbol":"btcusd",
// 	"exchange":"gemini",
// 	"avg_execution_price":"0.00",
// 	"side":"buy",
// 	"type":"exchange limit",
// 	"timestamp":"1483052547",
// 	"timestampms":1483052547532,
// 	"is_live":true,
// 	"is_cancelled":false,
// 	"is_hidden":false,
// 	"was_forced":false,
// 	"executed_amount":"0",
// 	"remaining_amount":"0.212",
// 	"price":"123.45",
// 	"original_amount":"0.212"
// }
// {
// 	"result":"error",
// 	"reason":"OrderNotFound",
// 	"message":"Order 68053334 not found"
// }
func (g *GeminiAPI) OrderStatus(orderId OrderId) (Order, error) {

	url := g.url + ORDER_STATUS_URL
	req := RequestParams{
		"request":  ORDER_STATUS_URL,
		"nonce":    getNonce(),
		"order_id": orderId,
	}

	var order Order

	body, err := request("POST", url, g.prepPayload(req), nil)
	if err != nil {
		return order, err
	}

	json.Unmarshal(body, &order)

	return order, nil
}

// {
// 	"result":"ok",
// 	"details":{
// 		"cancelledOrders":[68053332],
// 		"cancelRejects":[]
// 	}
// }
func (g *GeminiAPI) CancelAll() (CancelResponse, error) {
	url := g.url + CANCEL_ALL_URL
	req := RequestParams{
		"request": CANCEL_ALL_URL,
		"nonce":   getNonce(),
	}

	var res CancelResponse

	body, err := request("POST", url, g.prepPayload(req), nil)
	if err != nil {
		return res, err
	}

	json.Unmarshal(body, &res)

	return res, nil
}

// Cancel Order
// {
// 	"order_id":"68053338",
// 	"id":"68053338",
// 	"symbol":"btcusd",
// 	"exchange":"gemini",
// 	"avg_execution_price":"0.00",
// 	"side":"buy",
// 	"type":"exchange limit",
// 	"timestamp":"1483051780",
// 	"timestampms":1483051780029,
// 	"is_live":false,
// 	"is_cancelled":true,
// 	"is_hidden":false,
// 	"was_forced":false,
// 	"executed_amount":"0",
// 	"remaining_amount":"0.12",
// 	"price":"123.45",
// 	"original_amount":"0.12"
// }
// {
// 	"result":"error",
// 	"reason":"OrderNotFound",
// 	"message":"Order 6803338 not found"
// }
func (g *GeminiAPI) CancelOrder(orderId OrderId) (Order, error) {

	url := g.url + CANCEL_ORDER_URL
	req := RequestParams{
		"request":  CANCEL_ORDER_URL,
		"nonce":    getNonce(),
		"order_id": orderId,
	}

	var order Order

	body, err := request("POST", url, g.prepPayload(req), nil)
	if err != nil {
		return order, err
	}

	json.Unmarshal(body, &order)

	return order, nil
}
