package binance

import (
	"fmt"
	"testing"
	"github.com/xiaolo66/ExchangeApi"
)

var BaFuture = NewFuture(ExchangeApi.Options{
	PassPhrase: "",
	ProxyUrl:   "http://127.0.0.1:4780"}, ExchangeApi.FutureOptions{
	FutureAccountType: ExchangeApi.UsdtMargin,
	ContractType:      ExchangeApi.Swap,
})
var fmsgChan = make(ExchangeApi.MessageChan)

func handleFutureMsg(msgChan <-chan ExchangeApi.Message) {
	for {
		select {
		case msg := <-msgChan:
			switch msg.Type {
			case ExchangeApi.MsgReConnected:
				fmt.Println("reconnected..., restart subscribe")
			case ExchangeApi.MsgDisConnected:
				fmt.Println("disconnected, stop use old data, waiting reconnect....")
			case ExchangeApi.MsgClosed:
				fmt.Println("websocket closed, stop all")
				break
			case ExchangeApi.MsgError:
				fmt.Printf("error happend: %v\n", msg.Data)
				if err, ok := msg.Data.(ExchangeApi.ExError); ok {
					if err.Code == ExchangeApi.ErrInvalidDepth {
						//depth data invalid, Do some cleanup work, wait for the latest data or resubscribe
					}
				}
			case ExchangeApi.MsgOrderBook:
				orderbook, ok := msg.Data.(ExchangeApi.OrderBook)
				if !ok {
					fmt.Printf("order book data error %v", msg)
				}
				fmt.Printf("order book:%+v\n", orderbook)
			case ExchangeApi.MsgTicker:
				ticker, ok := msg.Data.(ExchangeApi.Ticker)
				if !ok {
					fmt.Printf("ticker data error %v", msg)
				}
				fmt.Printf("ticker:%+v\n", ticker)
			case ExchangeApi.MsgTrade:
				trades, ok := msg.Data.(ExchangeApi.Trade)
				if !ok {
					fmt.Printf("trade data error %v", msg)
				}
				fmt.Printf("trade:%+v\n", trades)
			case ExchangeApi.MsgMarkPrice:
				markPrice, ok := msg.Data.(ExchangeApi.MarkPrice)
				if !ok {
					fmt.Printf("markprice data error %v", msg)
				}
				fmt.Printf("markPrice:%+v\n", markPrice)
			case ExchangeApi.MsgKLine:
				klines, ok := msg.Data.(ExchangeApi.KLine)
				if !ok {
					fmt.Printf("kline data error %v", msg)
				}
				fmt.Printf("kline:%+v\n", klines)
			case ExchangeApi.MsgBalance:
				balances, ok := msg.Data.(ExchangeApi.BalanceUpdate)
				if !ok {
					fmt.Printf("balance data error %v", msg)
				}
				fmt.Printf("balance:%+v\n", balances)
			case ExchangeApi.MsgOrder:
				order, ok := msg.Data.(ExchangeApi.Order)
				if !ok {
					fmt.Printf("order data error %v", msg)
				}
				fmt.Printf("order:%+v\n", order)
			case ExchangeApi.MsgPositions:
				position, ok := msg.Data.(ExchangeApi.FuturePositonsUpdate)
				if !ok {
					fmt.Printf("order data error %v", msg)
				}
				fmt.Printf("position:%+v\n", position)
			}
		}
	}
}

func TestBinanceFutureWs_SubscribeOrderBook(t *testing.T) {

	if _, err := BaFuture.SubscribeOrderBook(symbol, 5, 0, true, fmsgChan); err == nil {
		handleFutureMsg(fmsgChan)
	}
}

func TestBinanceFutureWs_SubscribeTicker(t *testing.T) {
	if _, err := BaFuture.SubscribeTicker(symbol, fmsgChan); err == nil {
		handleFutureMsg(fmsgChan)
	}
}

func TestBinanceFutureWs_SubscribeTrades(t *testing.T) {
	if _, err := BaFuture.SubscribeTrades(symbol, fmsgChan); err == nil {
		handleFutureMsg(fmsgChan)
	}
}

func TestBinanceFutureWs_SubscribeKLine(t *testing.T) {
	if _, err := BaFuture.SubscribeKLine(symbol, ExchangeApi.KLine3Minute, fmsgChan); err == nil {
		handleFutureMsg(fmsgChan)
	}
}

func TestBinanceFutureWs_SubscribeBalance(t *testing.T) {
	if _, err := BaFuture.SubscribeBalance(symbol, fmsgChan); err == nil {
		handleFutureMsg(fmsgChan)
	}
}
func TestZbFutureWs_SubscribeOrder(t *testing.T) {
	if _, err := BaFuture.SubscribeOrder(symbol, fmsgChan); err == nil {
		handleFutureMsg(fmsgChan)
	}
}

func TestBinanceFutureWs_SubscribePositions(t *testing.T) {
	if _, err := BaFuture.SubscribePositions(symbol, fmsgChan); err == nil {
		handleFutureMsg(fmsgChan)
	}
}

func TestBinanceFutureWs_SubscribeMarkPrice(t *testing.T) {
	if _, err := BaFuture.SubscribeMarkPrice("EOS/USDT", fmsgChan); err == nil {
		handleFutureMsg(fmsgChan)
	}
}
