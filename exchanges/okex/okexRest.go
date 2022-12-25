package okex

import (
	"encoding/json"
	"fmt"
	"github.com/xiaolo66/ExchangeApi"
	"github.com/xiaolo66/ExchangeApi/exchanges"
	"github.com/xiaolo66/ExchangeApi/utils"
	. "github.com/xiaolo66/ExchangeApi/utils"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type OkexRest struct {
	exchanges.BaseExchange
	errors map[string]int
}

func (e *OkexRest) Init(option ExchangeApi.Options) {
	e.Option = option
	e.errors = map[string]int{
		"30009": ExchangeApi.ErrExchangeSystem,
		"36216": ExchangeApi.ErrOrderNotFound,
		"33014": ExchangeApi.ErrOrderNotFound,
		"33017": ExchangeApi.ErrInsufficientFunds,
	}

	if e.Option.RestHost == "" {
		e.Option.RestHost = "https://www.okex.com"
	}
}

func (e *OkexRest) FetchOrderBook(symbol string, size int) (orderBook ExchangeApi.OrderBook, err error) {
	market, err := e.GetMarket(symbol)
	if err != nil {
		return
	}
	params := url.Values{}
	params.Set("size", strconv.Itoa(size))

	function := fmt.Sprintf("/api/spot/v3/instruments/%s/book", market.SymbolID)
	res, err := e.Fetch(e, exchanges.Public, exchanges.GET, function, params, http.Header{})
	if err != nil {
		return
	}

	var data DepthData
	if err = json.Unmarshal(res, &data); err != nil {
		err = ExchangeApi.ExError{Code: ExchangeApi.ErrDataParse, Message: err.Error()}
		return
	}

	var ob OrderBook
	ob.update(data)
	return ob.OrderBook, nil
}

func (e *OkexRest) FetchTicker(symbol string) (ticker ExchangeApi.Ticker, err error) {
	market, err := e.GetMarket(symbol)
	if err != nil {
		return
	}
	params := url.Values{}
	function := fmt.Sprintf("/api/spot/v3/instruments/%s/ticker", market.SymbolID)
	res, err := e.Fetch(e, exchanges.Public, exchanges.GET, function, params, http.Header{})
	if err != nil {
		return
	}

	var data Ticker
	if err = json.Unmarshal(res, &data); err != nil {
		err = ExchangeApi.ExError{Code: ExchangeApi.ErrDataParse, Message: err.Error()}
		return
	}
	ticker = data.parseTicker(market.Symbol)
	return
}

func (e *OkexRest) FetchAllTicker() (tickers map[string]ExchangeApi.Ticker, err error) {
	params := url.Values{}
	res, err := e.Fetch(e, exchanges.Public, exchanges.GET, "/api/spot/v3/instruments/ticker", params, http.Header{})
	if err != nil {
		return
	}

	var data []Ticker
	if err = json.Unmarshal(res, &data); err != nil {
		err = ExchangeApi.ExError{Code: ExchangeApi.ErrDataParse, Message: err.Error()}
		return
	}

	tickers = make(map[string]ExchangeApi.Ticker, 0)
	for _, t := range data {
		market, err := e.GetMarketByID(t.Symbol)
		if err != nil {
			continue
		}
		tickers[market.Symbol] = t.parseTicker(market.Symbol)
	}

	return
}

func (e *OkexRest) FetchTrade(symbol string) (trades []ExchangeApi.Trade, err error) {
	market, err := e.GetMarket(symbol)
	if err != nil {
		return
	}
	params := url.Values{}
	function := fmt.Sprintf("/api/spot/v3/instruments/%s/trades", market.SymbolID)
	res, err := e.Fetch(e, exchanges.Public, exchanges.GET, function, params, http.Header{})
	if err != nil {
		return
	}

	var data []Trade
	if err = json.Unmarshal(res, &data); err != nil {
		err = ExchangeApi.ExError{Code: ExchangeApi.ErrDataParse, Message: err.Error()}
		return
	}

	for _, t := range data {
		trades = append(trades, t.parseTrade(market.Symbol))
	}
	return
}

func (e *OkexRest) FetchKLine(symbol string, t ExchangeApi.KLineType) (klines []ExchangeApi.KLine, err error) {
	market, err := e.GetMarket(symbol)
	if err != nil {
		return
	}
	params := url.Values{}
	kLineType := ""
	switch t {
	case ExchangeApi.KLine1Minute:
		kLineType = "60"
	case ExchangeApi.KLine3Minute:
		kLineType = "180"
	case ExchangeApi.KLine5Minute:
		kLineType = "300"
	case ExchangeApi.KLine15Minute:
		kLineType = "900"
	case ExchangeApi.KLine30Minute:
		kLineType = "1800"
	case ExchangeApi.KLine1Hour:
		kLineType = "3600"
	case ExchangeApi.KLine2Hour:
		kLineType = "7200"
	case ExchangeApi.KLine4Hour:
		kLineType = "14400"
	case ExchangeApi.KLine6Hour:
		kLineType = "21600"
	case ExchangeApi.KLine12Hour:
		kLineType = "43200"
	case ExchangeApi.KLine1Day:
		kLineType = "86400"
	case ExchangeApi.KLine1Week:
		kLineType = "604800"
	}
	function := fmt.Sprintf("/api/spot/v3/instruments/%s/candles?granularity=%s", market.SymbolID, kLineType)
	res, err := e.Fetch(e, exchanges.Public, exchanges.GET, function, params, http.Header{})
	if err != nil {
		return
	}

	var data []KLine
	if err = json.Unmarshal(res, &data); err != nil {
		err = ExchangeApi.ExError{Code: ExchangeApi.ErrDataParse, Message: err.Error()}
		return
	}

	for _, t := range data {
		klines = append(klines, t.parseKLine(market.Symbol))
	}
	return
}

func (e *OkexRest) FetchMarkets() (map[string]ExchangeApi.Market, error) {
	if len(e.Option.Markets) > 0 {
		return e.Option.Markets, nil
	}
	res, err := e.Fetch(e, exchanges.Public, exchanges.GET, "/api/spot/v3/instruments", url.Values{}, http.Header{})
	if err != nil {
		return e.Option.Markets, err
	}

	var response []struct {
		InstrumentId  string  `json:"instrument_id"`
		BaseCurrency  string  `json:"base_currency"`
		QuoteCurrency string  `json:"quote_currency"`
		MinSize       float64 `json:"min_size,string"`
		SizeIncrement string  `json:"size_increment"`
		TickSize      string  `json:"tick_size"`
	}
	if err = json.Unmarshal(res, &response); err != nil {
		err = ExchangeApi.ExError{Code: ExchangeApi.ErrDataParse, Message: err.Error()}
		return e.Option.Markets, err
	}

	e.Option.Markets = make(map[string]ExchangeApi.Market, 0)
	for _, v := range response {
		market := ExchangeApi.Market{
			SymbolID: strings.ToUpper(v.InstrumentId),
			Symbol:   strings.ToUpper(fmt.Sprintf("%s/%s", v.BaseCurrency, v.QuoteCurrency)),
			BaseID:   strings.ToUpper(v.BaseCurrency),
			QuoteID:  strings.ToUpper(v.QuoteCurrency),
			Lot:      v.MinSize,
		}
		pres := strings.Split(v.TickSize, ".")
		if len(pres) == 1 {
			market.PricePrecision = 0
		} else {
			market.PricePrecision = len(pres[1])
		}

		pres = strings.Split(v.SizeIncrement, ".")
		if len(pres) == 1 {
			market.AmountPrecision = 0
		} else {
			market.AmountPrecision = len(pres[1])
		}
		e.Option.Markets[market.Symbol] = market
	}
	return e.Option.Markets, nil
}

func (e *OkexRest) FetchBalance() (balances map[string]ExchangeApi.Balance, err error) {
	res, err := e.Fetch(e, exchanges.Private, exchanges.GET, "/api/spot/v3/accounts", url.Values{}, http.Header{})
	if err != nil {
		return
	}

	var data []Balance
	if err = json.Unmarshal(res, &data); err != nil {
		err = ExchangeApi.ExError{Code: ExchangeApi.ErrDataParse, Message: err.Error()}
		return
	}

	balances = make(map[string]ExchangeApi.Balance)
	for _, b := range data {
		balance := b.parseBalance()
		balances[balance.Asset] = balance
	}
	return
}

func (e *OkexRest) CreateOrder(symbol string, price, amount float64, side ExchangeApi.Side, tradeType ExchangeApi.TradeType, orderType ExchangeApi.OrderType, useClientID bool) (order ExchangeApi.Order, err error) {
	market, err := e.GetMarket(symbol)
	if err != nil {
		return
	}
	params := url.Values{}
	params.Set("instrument_id", market.SymbolID)
	params.Set("price", utils.Round(price, market.PricePrecision, false))
	params.Set("size", utils.Round(amount, market.AmountPrecision, false))
	if side == ExchangeApi.Sell {
		params.Set("side", "sell")
	} else if side == ExchangeApi.Buy {
		params.Set("side", "buy")
	}
	switch orderType {
	case ExchangeApi.PostOnly:
		params.Set("order_type", "1")
	case ExchangeApi.FOK:
		params.Set("order_type", "2")
	case ExchangeApi.IOC:
		params.Set("order_type", "3")
	}
	switch tradeType {
	case ExchangeApi.MARKET:
		params.Set("type", "market")
		params.Set("notional", fmt.Sprintf("%v", price*amount))
	default:
		params.Set("type", "limit")
	}
	if useClientID {
		params.Set("client_oid", GenerateOrderClientId(e.Option.ClientOrderIDPrefix, 32))
	}
	res, err := e.Fetch(e, exchanges.Private, exchanges.POST, "/api/spot/v3/orders", params, http.Header{})
	if err != nil {
		return
	}

	type response struct {
		ID  string `json:"order_id"`
		CID string `json:"client_oid"`
	}
	data := response{}
	if err = json.Unmarshal(res, &data); err != nil {
		err = ExchangeApi.ExError{Code: ExchangeApi.ErrDataParse, Message: err.Error()}
		return
	}
	order.ID = data.ID
	order.ClientID = data.CID
	return
}

func (e *OkexRest) CancelOrder(symbol, orderID string) (err error) {
	market, err := e.GetMarket(symbol)
	if err != nil {
		return
	}
	params := url.Values{}
	params.Set("instrument_id", market.SymbolID)
	function := "/api/spot/v3/cancel_orders/" + orderID
	_, err = e.Fetch(e, exchanges.Private, exchanges.POST, function, params, http.Header{})

	return err
}

func (e *OkexRest) CancelAllOrders(symbol string) (err error) {
	for {
		orders, err := e.FetchOpenOrders(symbol, 1, 10)
		if err != nil || len(orders) == 0 {
			break
		}
		for _, order := range orders {
			_ = e.CancelOrder(order.ID, symbol)
			time.Sleep(time.Millisecond * 200)
		}
	}
	return
}

//FetchOrder : 获取订单详情
func (e *OkexRest) FetchOrder(symbol, orderID string) (order ExchangeApi.Order, err error) {
	market, err := e.GetMarket(symbol)
	if err != nil {
		return
	}
	params := url.Values{}
	params.Set("instrument_id", market.SymbolID)
	function := "/api/spot/v3/orders/" + orderID
	res, err := e.Fetch(e, exchanges.Private, exchanges.GET, function, params, http.Header{})
	if err != nil {
		return
	}

	var data Order
	if err = json.Unmarshal(res, &data); err != nil {
		err = ExchangeApi.ExError{Code: ExchangeApi.ErrDataParse, Message: err.Error()}
		return
	}

	order = data.parseOrder(market.Symbol)
	return
}

//FetchOpenOrders :
func (e *OkexRest) FetchOpenOrders(symbol string, pageIndex, pageSize int) (orders []ExchangeApi.Order, err error) {
	market, err := e.GetMarket(symbol)
	if err != nil {
		return
	}
	params := url.Values{}
	params.Set("instrument_id", market.SymbolID)
	function := "/api/spot/v3/orders_pending/"
	res, err := e.Fetch(e, exchanges.Private, exchanges.GET, function, params, http.Header{})
	if err != nil {
		return
	}
	var data = make([]Order, 0)
	if err = json.Unmarshal(res, &data); err != nil {
		err = ExchangeApi.ExError{Code: ExchangeApi.ErrDataParse, Message: err.Error()}
		return
	}
	orders = make([]ExchangeApi.Order, len(data))
	for i, order := range data {
		orders[i] = order.parseOrder(market.Symbol)
	}
	return
}

func (e *OkexRest) Sign(access, method, function string, param url.Values, header http.Header) (request exchanges.Request) {
	request.Method = method
	request.Headers = header
	path := function
	if access == exchanges.Public {
		if len(param) > 0 {
			path = path + "?" + param.Encode()
		}
		request.Url = e.Option.RestHost + path
	} else {
		request.Headers.Set("OK-ACCESS-KEY", e.Option.AccessKey)
		request.Headers.Set("OK-ACCESS-PASSPHRASE", e.Option.PassPhrase)
		timestamp := IsoTime()
		request.Headers.Set("OK-ACCESS-TIMESTAMP", timestamp)
		auth := timestamp + method
		if method == exchanges.GET {
			if len(param) > 0 {
				path = path + "?" + param.Encode()
			}
			auth += path
		} else {
			request.Body = UrlValuesToJson(param)
			auth = auth + path + request.Body
		}
		request.Url = e.Option.RestHost + path
		signature, err := HmacSign(SHA256, auth, e.Option.SecretKey, true)
		if err != nil {
			return
		}
		request.Headers.Set("OK-ACCESS-SIGN", signature)
		request.Headers.Set("Content-Type", "application/json")
	}
	return request
}

func (e *OkexRest) HandleError(request exchanges.Request, response []byte) error {
	type Result struct {
		Code    string `json:"error_code"`
		Message string `json:"error_message"`
	}
	var result Result
	if err := json.Unmarshal(response, &result); err != nil {
		return nil
	}

	if result.Code == "0" || result.Code == "" {
		return nil
	}
	errCode, ok := e.errors[result.Code]
	if ok {
		return ExchangeApi.ExError{Code: errCode, Message: result.Message}
	} else {
		return ExchangeApi.ExError{Code: ExchangeApi.UnHandleError, Message: fmt.Sprintf("code:%v msg:%v", result.Code, result.Message)}
	}
}

