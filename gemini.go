package gemini

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	BASE_URL       = "https://api.gemini.com"
	SANDBOX_URL    = "https://api.sandbox.gemini.com"
	WS_BASE_URL    = "wss://api.gemini.com"
	WS_SANDBOX_URL = "wss://api.sandbox.gemini.com"

	// public
	SYMBOLS_URI = "/v1/symbols"
	TICKER_URI  = "/v1/pubticker/"
	BOOK_URI    = "/v1/book/"
	TRADES_URI  = "/v1/trades/"
	AUCTION_URI = "/v1/auction/"

	// authenticated
	PAST_TRADES_URI    = "/v1/mytrades"
	TRADE_VOLUME_URI   = "/v1/tradevolume"
	ACTIVE_ORDERS_URI  = "/v1/orders"
	ORDER_STATUS_URI   = "/v1/order/status"
	NEW_ORDER_URI      = "/v1/order/new"
	CANCEL_ORDER_URI   = "/v1/order/cancel"
	CANCEL_ALL_URI     = "/v1/order/cancel/all"
	CANCEL_SESSION_URI = "/v1/order/cancel/session"
	HEARTBEAT_URI      = "/v1/heartbeat"

	// fund mgmt
	BALANCES_URI            = "/v1/balances"
	NEW_DEPOSIT_ADDRESS_URI = "/v1/deposit/"
	WITHDRAW_FUNDS_URI      = "/v1/withdraw/"

	// websockets
	ORDER_EVENTS_URI = "/v1/order/events"
	MARKET_DATA_URI  = "/v1/marketdata/"
)

type Api struct {
	url    string
	key    string
	secret string
}

func New(live bool, key, secret string) *Api {
	var url string
	if url = SANDBOX_URL; live == true {
		url = BASE_URL
	}

	return &Api{url: url, key: key, secret: secret}
}

type ApiError struct {
	Reason  string
	Message string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("[%v] %v", e.Reason, e.Message)
}

type GenericResponse struct {
	Result string
	ApiError
}

type Id string

// Id has a custom Unmarshal since it needs to handle unmarshalling from both
// string and int json types. This package takes the position that throughout
// ids should be strings and converted from json into strings where needed.
func (id *Id) UnmarshalJSON(b []byte) error {

	if len(b) > 0 && b[0] == '"' {
		b = b[1:]
	}

	l := len(b)
	if l > 0 && b[l-1] == '"' {
		b = b[:l-1]
	}

	*id = Id(b)
	return nil
}

type Order struct {
	OrderId           Id     `json:"order_id"`
	ClientOrderId     string `json:"client_order_id"`
	Symbol            string
	Side              string
	Type              string
	Timestamp         int64   `json:"timestampms"`
	IsLive            bool    `json:"is_live"`
	IsCancelled       bool    `json:"is_cancelled"`
	IsHidden          bool    `json:"is_hidden"`
	WasForced         bool    `json:"was_forced"`
	Price             float64 `json:",string"`
	ExecutedAmount    float64 `json:"executed_amount,string"`
	RemainingAmount   float64 `json:"remaining_amount,string"`
	OriginalAmount    float64 `json:"original_amount,string"`
	AvgExecutionPrice float64 `json:"avg_execution_price,string"`
}

type Trade struct {
	OrderId       Id    `json:"order_id"`
	TradeId       Id    `json:"tid"`
	Timestamp     int64 `json:"timestampms"`
	Exchange      string
	Type          string
	FeeCurrency   string  `json:"fee_currency"`
	FeeAmount     float64 `json:"fee_amount,string"`
	Amount        float64 `json:",string"`
	Price         float64 `json:",string"`
	IsAuctionFill bool    `json:"is_auction_fill"`
	Aggressor     bool
	Broken        bool
	Break         string
}

type Ticker struct {
	Bid    float64 `json:",string"`
	Ask    float64 `json:",string"`
	Last   float64 `json:",string"`
	Volume TickerVolume
}

type TickerVolume struct {
	BTC       float64 `json:",string"`
	ETH       float64 `json:",string"`
	USD       float64 `json:",string"`
	Timestamp int64
}

type TradeVolume struct {
	AccountId         Id      `json:"account_id"`
	Symbol            string  `json:"symbol"`
	BaseCurrency      string  `json:"base_currency"`
	NotionalCurrency  string  `json:"notional_currency"`
	DataDate          string  `json:"data_date"`
	TotalVolumeBase   float64 `json:"total_volume_base"`
	MakeBuySellRatio  float64 `json:"maker_buy_sell_ratio"`
	BuyMakerBase      float64 `json:"buy_maker_base"`
	BuyMakerNotional  float64 `json:"buy_maker_notional"`
	BuyMakerCount     float64 `json:"buy_maker_count"`
	SellMakerBase     float64 `json:"sell_maker_base"`
	SellMakerNotional float64 `json:"sell_maker_notional"`
	SellMakerCount    float64 `json:"sell_maker_count"`
	BuyTakerBase      float64 `json:"buy_taker_base"`
	BuyTakerNotional  float64 `json:"buy_taker_notional"`
	BuyTakerCount     float64 `json:"buy_taker_count"`
	SellTakerBase     float64 `json:"sell_taker_base"`
	SellTakerNotional float64 `json:"sell_taker_notional"`
	SellTakerCount    float64 `json:"sell_taker_count"`
}

type CurrentAuction struct {
	ClosedUntil                  int64   `json:"closed_until_ms"`
	LastAuctionEid               Id      `json:"last_auction_eid"`
	LastAuctionPrice             float64 `json:"last_auction_price,string"`
	LastAuctionQuantity          float64 `json:"last_auction_quantity,string"`
	LastHighestBidPrice          float64 `json:"last_highest_bid_price,string"`
	LastLowestAskPrice           float64 `json:"last_lowest_ask_price,string"`
	MostRecentIndicativePrice    float64 `json:"most_recent_indicative_price,string"`
	MostRecentIndicativeQuantity float64 `json:"most_recent_indicative_quantity,string"`
	MostRecentHighestBidPrice    float64 `json:"most_recent_highest_bid_price,string"`
	MostRecentLowestAskPrice     float64 `json:"most_recent_lowest_ask_price,string"`
	NextUpdate                   int64   `json:"next_update_ms"`
	NextAuction                  int64   `json:"next_auction_ms"`
}

type Auction struct {
	Timestamp       int64   `json:"timestampms"`
	AuctionId       Id      `json:"auction_id"`
	Eid             Id      `json:"eid"`
	EventType       string  `json:"event_type"`
	AuctionResult   string  `json:"auction_result"`
	AuctionPrice    float64 `json:"auction_price,string"`
	AuctionQuantity float64 `json:"auction_quantity,string"`
	HighestBidPrice float64 `json:"highest_bid_price,string"`
	LowestAskPrice  float64 `json:"lowest_ask_price,string"`
}

type CancelResult struct {
	GenericResponse
	Details CancelResultDetails
}

type CancelResultDetails struct {
	CancelledOrders []Id
	CancelRejects   []Id
}

type FundBalance struct {
	Type                   string
	Currency               string
	Amount                 float64 `json:",string"`
	Available              float64 `json:",string"`
	AvailableForWithdrawal float64 `json:",string"`
}

type DepositAddress struct {
	Currency string
	Address  string
	Label    string
}

type WithdrawFundsResult struct {
	Destination string
	TxHash      string
	Amount      float64 `json:",string"`
}

// RequestHeaders contains the values to be included in POST headers, according
// to Gemini specification
type RequestHeaders struct {
	ApiKey    string
	Payload   string
	Signature string
}

// Nonce returns a generic nonce based on unix timestamp
func Nonce() int64 {
	return time.Now().UnixNano()
}

// BuildHeader handles the conversion of post parameters into headers formatted
// according to Gemini specification. Resulting headers include the API key,
// the signature and payload.
func (g *Api) BuildHeader(req *map[string]interface{}) *RequestHeaders {

	reqStr, _ := json.Marshal(req)
	payload := base64.StdEncoding.EncodeToString([]byte(reqStr))

	mac := hmac.New(sha512.New384, []byte(g.secret))
	mac.Write([]byte(payload))

	signature := hex.EncodeToString(mac.Sum(nil))

	return &RequestHeaders{
		ApiKey:    g.key,
		Payload:   payload,
		Signature: signature,
	}
}

// request makes the HTTP request to Gemini and handles any returned errors
func (g *Api) request(verb, url string, postParams, getParams map[string]interface{}) ([]byte, error) {

	req, err := http.NewRequest(verb, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}

	if postParams != nil {
		headers := g.BuildHeader(&postParams)
		req.Header.Set("X-GEMINI-APIKEY", headers.ApiKey)
		req.Header.Set("X-GEMINI-PAYLOAD", headers.Payload)
		req.Header.Set("X-GEMINI-SIGNATURE", headers.Signature)
	}

	if getParams != nil {
		q := req.URL.Query()
		for key, val := range getParams {
			q.Add(key, val.(string))
		}
		req.URL.RawQuery = q.Encode()
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// check for error from Gemini
	var res GenericResponse

	json.Unmarshal(body, &res)
	if res.Result == "error" {
		return nil, &res.ApiError
	}

	return body, nil
}
