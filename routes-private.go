package gemini

import (
	"encoding/json"
	"strconv"
)

// Past Trades
func (g *Api) PastTrades(symbol string, limitTrades int, timestamp int64) ([]Trade, error) {

	url := g.url + PAST_TRADES_URI

	params := map[string]interface{}{
		"request":      PAST_TRADES_URI,
		"nonce":        Nonce(),
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
func (g *Api) TradeVolume() ([][]TradeVolume, error) {

	url := g.url + TRADE_VOLUME_URI
	params := map[string]interface{}{
		"request": TRADE_VOLUME_URI,
		"nonce":   Nonce(),
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
func (g *Api) ActiveOrders() ([]Order, error) {

	url := g.url + ACTIVE_ORDERS_URI
	params := map[string]interface{}{
		"request": ACTIVE_ORDERS_URI,
		"nonce":   Nonce(),
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
func (g *Api) OrderStatus(orderId string) (Order, error) {

	url := g.url + ORDER_STATUS_URI
	params := map[string]interface{}{
		"request":  ORDER_STATUS_URI,
		"nonce":    Nonce(),
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
func (g *Api) NewOrder(symbol, clientOrderId string, amount, price float64, side string, options []string) (Order, error) {

	url := g.url + NEW_ORDER_URI
	params := map[string]interface{}{
		"request":         NEW_ORDER_URI,
		"nonce":           Nonce(),
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
func (g *Api) CancelOrder(orderId string) (Order, error) {

	url := g.url + CANCEL_ORDER_URI
	params := map[string]interface{}{
		"request":  CANCEL_ORDER_URI,
		"nonce":    Nonce(),
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
func (g *Api) CancelAll() (CancelResult, error) {

	url := g.url + CANCEL_ALL_URI
	params := map[string]interface{}{
		"request": CANCEL_ALL_URI,
		"nonce":   Nonce(),
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
func (g *Api) CancelSession() (GenericResponse, error) {

	url := g.url + CANCEL_SESSION_URI
	params := map[string]interface{}{
		"request": CANCEL_SESSION_URI,
		"nonce":   Nonce(),
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
func (g *Api) Heartbeat() (GenericResponse, error) {

	url := g.url + HEARTBEAT_URI
	params := map[string]interface{}{
		"request": HEARTBEAT_URI,
		"nonce":   Nonce(),
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
func (g *Api) Balances() ([]FundBalance, error) {

	url := g.url + BALANCES_URI
	params := map[string]interface{}{
		"request": BALANCES_URI,
		"nonce":   Nonce(),
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
func (g *Api) NewDepositAddress(currency, label string) (DepositAddress, error) {

	path := NEW_DEPOSIT_ADDRESS_URI + currency + "/newAddress"
	url := g.url + path
	params := map[string]interface{}{
		"request": path,
		"nonce":   Nonce(),
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
func (g *Api) WithdrawFunds(currency, address string, amount float64) (WithdrawFundsResult, error) {

	path := WITHDRAW_FUNDS_URI + currency
	url := g.url + path
	params := map[string]interface{}{
		"request": path,
		"nonce":   Nonce(),
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
