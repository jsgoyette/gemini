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
	BASE_URL    = "https://api.gemini.com"
	SANDBOX_URL = "https://api.sandbox.gemini.com"

	// public
	SYMBOLS_URL = "/v1/symbols"
	TICKER_URL  = "/v1/pubticker/"
	BOOK_URL    = "/v1/book/"
	TRADES_URL  = "/v1/trades/"

	// authenticated
	ACTIVE_ORDERS_URL  = "/v1/orders"
	PAST_TRADES_URL    = "/v1/mytrades"
	ORDER_STATUS_URL   = "/v1/order/status"
	NEW_ORDER_URL      = "/v1/order/new"
	CANCEL_ORDER_URL   = "/v1/order/cancel"
	CANCEL_ALL_URL     = "/v1/order/cancel/all"
	CANCEL_SESSION_URL = "/v1/order/cancel/session"
	HEARTBEAT_URL      = "/v1/heartbeat"
)

type GeminiAPI struct {
	url    string
	key    string
	secret string
}

func New(live bool, key, secret string) *GeminiAPI {
	var url string
	if url = SANDBOX_URL; live == true {
		url = BASE_URL
	}

	return &GeminiAPI{url: url, key: key, secret: secret}
}

type GeminiError struct {
	Reason  string
	Message string
}

func (e *GeminiError) Error() string {
	return fmt.Sprintf("[%v] %v", e.Reason, e.Message)
}

type Id string

// custom Unmarshal since Gemini returns array of ints instead of strings in
// the CancelAll response, and thus type needs to handle unmarshalling from
// both string and int json types
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
	OrderId         Id      `json:"order_id"`
	ClientOrderId   string  `json:"client_order_id"`
	Symbol          string  `json:"symbol"`
	Price           float64 `json:",string"`
	Side            string  `json:"side"`
	Type            string  `json:"type"`
	Timestamp       uint64  `json:"timestampms"`
	IsLive          bool    `json:"is_live"`
	IsCancelled     bool    `json:"is_cancelled"`
	IsHidden        bool    `json:"is_hidden"`
	WasForced       bool    `json:"was_forced"`
	ExecutedAmount  float64 `json:"executed_amount,string"`
	RemainingAmount float64 `json:"remaining_amount,string"`
	OriginalAmount  float64 `json:"original_amount,string"`
}

type Trade struct {
	OrderId       Id      `json:"order_id"`
	TradeId       Id      `json:"tid"`
	Exchange      string  `json:"exchange"`
	Price         float64 `json:",string"`
	Amount        float64 `json:",string"`
	Timestamp     uint64  `json:"timestampms"`
	Type          string  `json:"type"`
	Aggressor     bool    `json:"aggressor"`
	FeeCurrency   string  `json:"fee_currency"`
	FeeAmount     float64 `json:"fee_amount,string"`
	IsAuctionFill bool    `json:"is_auction_fill"`
}

type Ticker struct {
	Bid    float64      `json:"bid,string"`
	Ask    float64      `json:"ask,string"`
	Volume TickerVolume `json:"volume"`
}

type TickerVolume struct {
	BTC       float64 `json:"BTC,string"`
	ETH       float64 `json:"ETH,string"`
	USD       float64 `json:"USD,string"`
	Last      float64 `json:"last,string"`
	Timestamp uint64  `json:"timestamp"`
}

type Book struct {
	Bids []BookEntry
	Asks []BookEntry
}

type BookEntry struct {
	Price  float64 `json:",string"`
	Amount float64 `json:",string"`
}

type GenericResponse struct {
	Result  string
	Reason  string
	Message string
}

type CancelResult struct {
	GenericResponse
	Details CancelResultDetails `json:"details"`
}

type CancelResultDetails struct {
	CancelledOrders []Id `json:"cancelledOrders"`
	CancelRejects   []Id `json:"cancelRejects"`
}

// internal types
type requestHeaders struct {
	key       string
	payload   string
	signature string
}

type requestParams map[string]interface{}

// internal functions
func getNonce() int64 {
	return time.Now().UnixNano()
}

// internal methods
func (g *GeminiAPI) prepPayload(req *requestParams) *requestHeaders {

	reqStr, _ := json.Marshal(req)
	payload := base64.StdEncoding.EncodeToString([]byte(reqStr))

	mac := hmac.New(sha512.New384, []byte(g.secret))
	mac.Write([]byte(payload))

	signature := hex.EncodeToString(mac.Sum(nil))

	return &requestHeaders{
		key:       g.key,
		payload:   payload,
		signature: signature,
	}
}

func (g *GeminiAPI) request(verb, url string, postParams, getParams requestParams) ([]byte, error) {

	req, err := http.NewRequest(verb, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}

	if postParams != nil {
		headers := g.prepPayload(&postParams)

		req.Header.Set("X-GEMINI-APIKEY", headers.key)
		req.Header.Set("X-GEMINI-PAYLOAD", headers.payload)
		req.Header.Set("X-GEMINI-SIGNATURE", headers.signature)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// fmt.Println("response Body:", string(body))

	// check for Gemini error
	var res GenericResponse
	json.Unmarshal(body, &res)
	if res.Result == "error" {
		return nil, &GeminiError{
			Reason:  res.Reason,
			Message: res.Message,
		}
	}

	return body, nil
}
