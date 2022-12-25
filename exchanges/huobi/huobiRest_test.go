package huobi

import (
	"github.com/xiaolo66/ExchangeApi"
	"testing"
)

var huobi = New(ExchangeApi.Options{AccessKey: "20fdf4fb-5c28360a-qv2d5ctgbn-3d3a8", SecretKey: "65f50f0e-7a5f6387-7ca98911-ed85c", ProxyUrl: "http://127.0.0.1:4780"})

func TestHuobiRest_FetchTicker(t *testing.T) {
	ticker, err := huobi.FetchTicker(symbol)
	if err != nil {
		t.Error(err)
	}
	t.Log(ticker)
}

func TestHuobiRest_FetchOrderBook(t *testing.T) {
	res, err := huobi.FetchOrderBook("BTC/USDT", 0)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

func TestHuobiRest_FetchAllTicker(t *testing.T) {
	res, err := huobi.FetchAllTicker()
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

func TestHuobiRest_FetchTrade(t *testing.T) {
	res, err := huobi.FetchTrade("ETH/USDT")
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

func TestHuobiRest_FetchKLine(t *testing.T) {
	res, err := huobi.FetchKLine("ETH/USDT", ExchangeApi.KLine1Day)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

func TestHuobiRest_FetchOrder(t *testing.T) {
	res, err := huobi.FetchOrder("ETH/USDT", "320896583379611")
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

func TestHuobiRest_FetchBalance(t *testing.T) {
	res, err := huobi.FetchBalance()
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

func TestHuobiRest_CreateOrder(t *testing.T) {
	res, err := huobi.CreateOrder("EOS/USDT", 0.5, 20, ExchangeApi.Buy, ExchangeApi.LIMIT, ExchangeApi.PostOnly, false)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

func TestHuobiRest_CancelOrder(t *testing.T) {
	err := huobi.CancelOrder("ETH/USDT", "702143010867965")
	if err != nil {
		t.Error(err)
	}
}

func TestHuobiRest_CancelAllOrders(t *testing.T) {
	err := huobi.CancelAllOrders("ETH/USDT")
	if err != nil {
		t.Error(err)
	}
}

func TestHuobiRest_FetchAllOrders(t *testing.T) {
	orders, err := huobi.FetchOpenOrders("ETH/USDT", 0, 1000)
	if err != nil {
		t.Error(err)
	}
	t.Log(orders)
}
