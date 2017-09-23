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
	OrderId           Id      `json:"order_id"`
	ClientOrderId     string  `json:"client_order_id"`
	Symbol            string  `json:"symbol"`
	Side              string  `json:"side"`
	Type              string  `json:"type"`
	Timestamp         int64   `json:"timestampms"`
	IsLive            bool    `json:"is_live"`
	IsCancelled       bool    `json:"is_cancelled"`
	IsHidden          bool    `json:"is_hidden"`
	WasForced         bool    `json:"was_forced"`
	Price             float64 `json:"price,string"`
	ExecutedAmount    float64 `json:"executed_amount,string"`
	RemainingAmount   float64 `json:"remaining_amount,string"`
	OriginalAmount    float64 `json:"original_amount,string"`
	AvgExecutionPrice float64 `json:"avg_execution_price,string"`
}

type Trade struct {
	OrderId       Id      `json:"order_id"`
	TradeId       Id      `json:"tid"`
	Timestamp     int64   `json:"timestampms"`
	Exchange      string  `json:"exchange"`
	Type          string  `json:"type"`
	FeeCurrency   string  `json:"fee_currency"`
	FeeAmount     float64 `json:"fee_amount,string"`
	Amount        float64 `json:"amount,string"`
	Price         float64 `json:"price,string"`
	IsAuctionFill bool    `json:"is_auction_fill"`
	Aggressor     bool    `json:"aggressor"`
	Broken        bool    `json:"broken"`
	Break         string  `json:"break"`
}

type Ticker struct {
	Bid    float64      `json:"bid,string"`
	Ask    float64      `json:"ask,string"`
	Last   float64      `json:"last,string"`
	Volume TickerVolume `json:"volume"`
}

type TickerVolume struct {
	BTC       float64 `json:",string"`
	ETH       float64 `json:",string"`
	USD       float64 `json:",string"`
	Timestamp int64   `json:"timestamp"`
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
	Type                   string  `json:"type"`
	Currency               string  `json:"currency"`
	Amount                 float64 `json:"amount,string"`
	Available              float64 `json:"available,string"`
	AvailableForWithdrawal float64 `json:"availableForWithdrawal,string"`
}

type DepositAddress struct {
	Currency string `json:"currency"`
	Address  string `json:"address"`
	Label    string `json:"label"`
}

type WithdrawFundsResult struct {
	Destination string  `json:"destination"`
	TxHash      string  `json:"txHash"`
	Amount      float64 `json:"amount,string"`
}

// Nonce returns a generic nonce based on unix timestamp
func Nonce() int64 {
	return time.Now().UnixNano()
}

// BuildHeader handles the conversion of post parameters into headers formatted
// according to Gemini specification. Resulting headers include the API key,
// the payload and the signature.
func (api *Api) BuildHeader(req *map[string]interface{}) http.Header {

	reqStr, _ := json.Marshal(req)
	payload := base64.StdEncoding.EncodeToString([]byte(reqStr))

	mac := hmac.New(sha512.New384, []byte(api.secret))
	mac.Write([]byte(payload))

	signature := hex.EncodeToString(mac.Sum(nil))

	header := http.Header{}
	header.Set("X-GEMINI-APIKEY", api.key)
	header.Set("X-GEMINI-PAYLOAD", payload)
	header.Set("X-GEMINI-SIGNATURE", signature)

	return header
}

// request makes the HTTP request to Gemini and handles any returned errors
func (api *Api) request(verb, url string, params map[string]interface{}) ([]byte, error) {

	req, err := http.NewRequest(verb, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}

	if params != nil {
		if verb == "GET" {
			q := req.URL.Query()
			for key, val := range params {
				q.Add(key, val.(string))
			}
			req.URL.RawQuery = q.Encode()
		} else {
			req.Header = api.BuildHeader(&params)
		}
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
