package gemini

type MarketData struct {
	Type    string
	EventId Id
	Events  []MarketEvent
}

type MarketEvent struct {
	Type  string
	Price float64 `json:",string"`

	// change event
	Side      string
	Remaining float64 `json:",string"`
	Delta     float64 `json:",string"`
	Reason    string

	// trade event
	TradeId   Id      `json:"tid"`
	Amount    float64 `json:",string"`
	MakerSide string

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
	Type              string
	OrderId           Id     `json:"order_id"`
	EventId           Id     `json:"event_id"`
	ClientOrderId     string `json:"client_order_id"`
	ApiSession        string `'json:"api_session"`
	Symbol            string
	Side              string
	Behavior          string
	OrderType         string  `json:"order_type"`
	Timestamp         int64   `json:"timestampms"`
	IsLive            bool    `json:"is_live"`
	IsCancelled       bool    `json:"is_cancelled"`
	IsHidden          bool    `json:"is_hidden"`
	Price             float64 `json:",string"`
	ExecutedAmount    float64 `json:"executed_amount,string"`
	RemainingAmount   float64 `json:"remaining_amount,string"`
	OriginalAmount    float64 `json:"original_amount,string"`
	AvgExecutionPrice float64 `json:"avg_execution_price,string"`
	TotalSpend        float64 `json:"total_spend,string"`

	// subscription acknowledgement
	AccountId        Id
	SymbolFilter     []string
	ApiSessionFilter []string
	EventTypeFilter  []string

	// heartbeat
	Sequence int
	TraceId  string `json:"trace_id"`

	// fill
	Fill OrderFill

	// reject / cancel
	Reason string

	// cancel
	CancelCommandId string `json:"cancel_command_id"`
}

type OrderFill struct {
	TradeId     Id `json:"trade_id"`
	Liquidity   string
	Price       float64 `json:",string"`
	Amount      float64 `json:",string"`
	Fee         float64 `json:",string"`
	FeeCurrency string  `json:"fee_currency"`
}
