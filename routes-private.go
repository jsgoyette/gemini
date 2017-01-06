package gemini

import (
	"encoding/json"
	"strconv"
)

// Past Trades
func (api *Api) PastTrades(symbol string, limitTrades int, timestamp int64) ([]Trade, error) {

	url := api.url + PAST_TRADES_URI

	params := map[string]interface{}{
		"request":      PAST_TRADES_URI,
		"nonce":        Nonce(),
		"symbol":       symbol,
		"limit_trades": limitTrades,
		"timestamp":    timestamp,
	}

	var trades []Trade

	body, err := api.request("POST", url, params, nil)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, &trades)

	return trades, nil
}

// Trade Volume
func (api *Api) TradeVolume() ([][]TradeVolume, error) {

	url := api.url + TRADE_VOLUME_URI
	params := map[string]interface{}{
		"request": TRADE_VOLUME_URI,
		"nonce":   Nonce(),
	}

	var volumes [][]TradeVolume

	body, err := api.request("POST", url, params, nil)
	if err != nil {
		return volumes, err
	}

	json.Unmarshal(body, &volumes)

	return volumes, nil
}

// Active Orders
func (api *Api) ActiveOrders() ([]Order, error) {

	url := api.url + ACTIVE_ORDERS_URI
	params := map[string]interface{}{
		"request": ACTIVE_ORDERS_URI,
		"nonce":   Nonce(),
	}

	var orders []Order

	body, err := api.request("POST", url, params, nil)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, &orders)

	return orders, nil
}

// Order Status
func (api *Api) OrderStatus(orderId string) (Order, error) {

	url := api.url + ORDER_STATUS_URI
	params := map[string]interface{}{
		"request":  ORDER_STATUS_URI,
		"nonce":    Nonce(),
		"order_id": orderId,
	}

	var order Order

	body, err := api.request("POST", url, params, nil)
	if err != nil {
		return order, err
	}

	json.Unmarshal(body, &order)

	return order, nil
}

// New Order
func (api *Api) NewOrder(symbol, clientOrderId string, amount, price float64, side string, options []string) (Order, error) {

	url := api.url + NEW_ORDER_URI
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

	body, err := api.request("POST", url, params, nil)
	if err != nil {
		return order, err
	}

	json.Unmarshal(body, &order)

	return order, nil
}

// Cancel Order
func (api *Api) CancelOrder(orderId string) (Order, error) {

	url := api.url + CANCEL_ORDER_URI
	params := map[string]interface{}{
		"request":  CANCEL_ORDER_URI,
		"nonce":    Nonce(),
		"order_id": orderId,
	}

	var order Order

	body, err := api.request("POST", url, params, nil)
	if err != nil {
		return order, err
	}

	json.Unmarshal(body, &order)

	return order, nil
}

// Cancel All
func (api *Api) CancelAll() (CancelResult, error) {

	url := api.url + CANCEL_ALL_URI
	params := map[string]interface{}{
		"request": CANCEL_ALL_URI,
		"nonce":   Nonce(),
	}

	var res CancelResult

	body, err := api.request("POST", url, params, nil)
	if err != nil {
		return res, err
	}

	json.Unmarshal(body, &res)

	return res, nil
}

// Cancel Session
func (api *Api) CancelSession() (GenericResponse, error) {

	url := api.url + CANCEL_SESSION_URI
	params := map[string]interface{}{
		"request": CANCEL_SESSION_URI,
		"nonce":   Nonce(),
	}

	var res GenericResponse

	body, err := api.request("POST", url, params, nil)
	if err != nil {
		return res, err
	}

	json.Unmarshal(body, &res)

	return res, nil
}

// Heartbeat
func (api *Api) Heartbeat() (GenericResponse, error) {

	url := api.url + HEARTBEAT_URI
	params := map[string]interface{}{
		"request": HEARTBEAT_URI,
		"nonce":   Nonce(),
	}

	var res GenericResponse

	body, err := api.request("POST", url, params, nil)
	if err != nil {
		return res, err
	}

	json.Unmarshal(body, &res)

	return res, nil
}

// Balances
func (api *Api) Balances() ([]FundBalance, error) {

	url := api.url + BALANCES_URI
	params := map[string]interface{}{
		"request": BALANCES_URI,
		"nonce":   Nonce(),
	}

	var balances []FundBalance

	body, err := api.request("POST", url, params, nil)
	if err != nil {
		return balances, err
	}

	json.Unmarshal(body, &balances)

	return balances, nil
}

// New Deposit Address
func (api *Api) NewDepositAddress(currency, label string) (DepositAddress, error) {

	path := NEW_DEPOSIT_ADDRESS_URI + currency + "/newAddress"
	url := api.url + path
	params := map[string]interface{}{
		"request": path,
		"nonce":   Nonce(),
		"label":   label,
	}

	var res DepositAddress

	body, err := api.request("POST", url, params, nil)
	if err != nil {
		return res, err
	}

	json.Unmarshal(body, &res)

	return res, nil
}

// Withdraw Crypto Funds
func (api *Api) WithdrawFunds(currency, address string, amount float64) (WithdrawFundsResult, error) {

	path := WITHDRAW_FUNDS_URI + currency
	url := api.url + path
	params := map[string]interface{}{
		"request": path,
		"nonce":   Nonce(),
		"address": address,
		"amount":  strconv.FormatFloat(amount, 'f', -1, 64),
	}

	var res WithdrawFundsResult

	body, err := api.request("POST", url, params, nil)
	if err != nil {
		return res, err
	}

	json.Unmarshal(body, &res)

	return res, nil
}
