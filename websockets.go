package gemini

type MarketData struct {
	Type    string        `json:"type"`
	EventId Id            `json:"eventId"`
	Events  []MarketEvent `json:"events"`
}

type MarketEvent struct {
	Type  string  `json:"type"`
	Price float64 `json:"price,string"`

	// change event
	Side      string  `json:"side"`
	Remaining float64 `json:"remaining,string"`
	Delta     float64 `json:"delta,string"`
	Reason    string  `json:"reason"`

	// trade event
	TradeId   Id      `json:"tid"`
	Amount    float64 `json:"amount,string"`
	MakerSide string  `json:"makerSide"`

	// auction open event
	AuctionOpen     int64 `json:"auction_open_ms"`
	AuctionTime     int64 `json:"auction_time_ms"`
	FirstIndicative int64 `json:"first_indicative_ms"`
	LastCancelTime  int64 `json:"last_cancel_time_ms"`

	// auction indicative price event
	// auction outcome event
	Eid                Id
	AuctionResult      string  `json:"auction_result"`
	EventTime          int64   `json:"event_time_ms"`
	HighestBidPrice    float64 `json:"highest_bid_price,string"`
	LowestAskPrice     float64 `json:"lowest_ask_price,string"`
	IndicativePrice    float64 `json:"indicative_price,string"`
	IndicativeQuantity float64 `json:"indicative_quantity,string"`
}

type OrderEvent struct {
	Type              string  `json:"type"`
	OrderId           Id      `json:"order_id"`
	EventId           Id      `json:"event_id"`
	ClientOrderId     string  `json:"client_order_id"`
	ApiSession        string  `json:"api_session"`
	Symbol            string  `json:"symbol"`
	Side              string  `json:"side"`
	Behavior          string  `json:"behavior"`
	OrderType         string  `json:"order_type"`
	Timestamp         int64   `json:"timestampms"`
	IsLive            bool    `json:"is_live"`
	IsCancelled       bool    `json:"is_cancelled"`
	IsHidden          bool    `json:"is_hidden"`
	Price             float64 `json:"price,string"`
	ExecutedAmount    float64 `json:"executed_amount,string"`
	RemainingAmount   float64 `json:"remaining_amount,string"`
	OriginalAmount    float64 `json:"original_amount,string"`
	AvgExecutionPrice float64 `json:"avg_execution_price,string"`
	TotalSpend        float64 `json:"total_spend,string"`

	// subscription acknowledgement
	AccountId        Id       `json:"accountId"`
	SymbolFilter     []string `json:"symbolFilter"`
	ApiSessionFilter []string `json:"apiSessionFilter"`
	EventTypeFilter  []string `json:"eventTypeFilter"`

	// heartbeat
	Sequence int    `json:"sequence"`
	TraceId  string `json:"trace_id"`

	// fill
	Fill OrderFill `json:"fill"`

	// reject / cancel
	Reason string `json:"reason"`

	// cancel
	CancelCommandId string `json:"cancel_command_id"`
}

type OrderFill struct {
	TradeId     Id      `json:"trade_id"`
	Liquidity   string  `json:"liquidity"`
	Price       float64 `json:"price,string"`
	Amount      float64 `json:"amount,string"`
	Fee         float64 `json:"fee,string"`
	FeeCurrency string  `json:"fee_currency"`
}
