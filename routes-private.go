package gemini

import (
	"encoding/json"
	"strconv"
)

// Past Trades
func (g *GeminiApi) PastTrades(symbol string, limitTrades int, timestamp int64) ([]Trade, error) {

	url := g.url + PAST_TRADES_URI

	params := requestParams{
		"request":      PAST_TRADES_URI,
		"nonce":        getNonce(),
		"symbol":       symbol,
		"limit_trades": limitTrades,
		"timestamp":    timestamp,
	}

	var trades []Trade

	body, err := g.request("POST", url, params, nil)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, &trades)

	return trades, nil
}

// Trade Volume
func (g *GeminiApi) TradeVolume() ([][]TradeVolume, error) {

	url := g.url + TRADE_VOLUME_URI
	params := requestParams{
		"request": TRADE_VOLUME_URI,
		"nonce":   getNonce(),
	}

	var volumes [][]TradeVolume

	body, err := g.request("POST", url, params, nil)
	if err != nil {
		return volumes, err
	}

	json.Unmarshal(body, &volumes)

	return volumes, nil
}

// Active Orders
func (g *GeminiApi) ActiveOrders() ([]Order, error) {

	url := g.url + ACTIVE_ORDERS_URI
	params := requestParams{
		"request": ACTIVE_ORDERS_URI,
		"nonce":   getNonce(),
	}

	var orders []Order

	body, err := g.request("POST", url, params, nil)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, &orders)

	return orders, nil
}

// Order Status
func (g *GeminiApi) OrderStatus(orderId string) (Order, error) {

	url := g.url + ORDER_STATUS_URI
	params := requestParams{
		"request":  ORDER_STATUS_URI,
		"nonce":    getNonce(),
		"order_id": orderId,
	}

	var order Order

	body, err := g.request("POST", url, params, nil)
	if err != nil {
		return order, err
	}

	json.Unmarshal(body, &order)

	return order, nil
}

// New Order
func (g *GeminiApi) NewOrder(symbol, clientOrderId string, amount, price float64, side string, options []string) (Order, error) {

	url := g.url + NEW_ORDER_URI
	params := requestParams{
		"request":         NEW_ORDER_URI,
		"nonce":           getNonce(),
		"client_order_id": clientOrderId,
		"symbol":          symbol,
		"amount":          strconv.FormatFloat(amount, 'f', -1, 64),
		"price":           strconv.FormatFloat(price, 'f', -1, 64),
		"side":            side,
		"type":            "exchange limit",
	}

	if options != nil {
		params["options"] = options
	}

	var order Order

	body, err := g.request("POST", url, params, nil)
	if err != nil {
		return order, err
	}

	json.Unmarshal(body, &order)

	return order, nil
}

// Cancel Order
func (g *GeminiApi) CancelOrder(orderId string) (Order, error) {

	url := g.url + CANCEL_ORDER_URI
	params := requestParams{
		"request":  CANCEL_ORDER_URI,
		"nonce":    getNonce(),
		"order_id": orderId,
	}

	var order Order

	body, err := g.request("POST", url, params, nil)
	if err != nil {
		return order, err
	}

	json.Unmarshal(body, &order)

	return order, nil
}

// Cancel All
func (g *GeminiApi) CancelAll() (CancelResult, error) {

	url := g.url + CANCEL_ALL_URI
	params := requestParams{
		"request": CANCEL_ALL_URI,
		"nonce":   getNonce(),
	}

	var res CancelResult

	body, err := g.request("POST", url, params, nil)
	if err != nil {
		return res, err
	}

	json.Unmarshal(body, &res)

	return res, nil
}

// Cancel Session
func (g *GeminiApi) CancelSession() (GenericResponse, error) {

	url := g.url + CANCEL_SESSION_URI
	params := requestParams{
		"request": CANCEL_SESSION_URI,
		"nonce":   getNonce(),
	}

	var res GenericResponse

	body, err := g.request("POST", url, params, nil)
	if err != nil {
		return res, err
	}

	json.Unmarshal(body, &res)

	return res, nil
}

// Heartbeat
func (g *GeminiApi) Heartbeat() (GenericResponse, error) {

	url := g.url + HEARTBEAT_URI
	params := requestParams{
		"request": HEARTBEAT_URI,
		"nonce":   getNonce(),
	}

	var res GenericResponse

	body, err := g.request("POST", url, params, nil)
	if err != nil {
		return res, err
	}

	json.Unmarshal(body, &res)

	return res, nil
}

// Balances
func (g *GeminiApi) Balances() ([]FundBalance, error) {

	url := g.url + BALANCES_URI
	params := requestParams{
		"request": BALANCES_URI,
		"nonce":   getNonce(),
	}

	var balances []FundBalance

	body, err := g.request("POST", url, params, nil)
	if err != nil {
		return balances, err
	}

	json.Unmarshal(body, &balances)

	return balances, nil
}

// New Deposit Address
func (g *GeminiApi) NewDepositAddress(currency, label string) (DepositAddress, error) {

	path := NEW_DEPOSIT_ADDRESS_URI + currency + "/newAddress"
	url := g.url + path
	params := requestParams{
		"request": path,
		"nonce":   getNonce(),
		"label":   label,
	}

	var res DepositAddress

	body, err := g.request("POST", url, params, nil)
	if err != nil {
		return res, err
	}

	json.Unmarshal(body, &res)

	return res, nil
}

// Withdraw Crypto Funds
func (g *GeminiApi) WithdrawFunds(currency, address string, amount float64) (WithdrawFundsResult, error) {

	path := WITHDRAW_FUNDS_URI + currency
	url := g.url + path
	params := requestParams{
		"request": path,
		"nonce":   getNonce(),
		"address": address,
		"amount":  strconv.FormatFloat(amount, 'f', -1, 64),
	}

	var res WithdrawFundsResult

	body, err := g.request("POST", url, params, nil)
	if err != nil {
		return res, err
	}

	json.Unmarshal(body, &res)

	return res, nil
}
