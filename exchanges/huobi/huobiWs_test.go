package huobi

import (
	"fmt"
	"github.com/xiaolo66/ExchangeApi"
	"testing"
	"time"
)

var (
	symbol  = "BTC/USDT"
	symbol1 = "FIL/USDT"
	e       = New(ExchangeApi.Options{AccessKey: "dqnh6tvdf3-d6158aae-690f264d-ef474", SecretKey: "a5f7a8f5-86220cae-b6d77900-c1a01", AutoReconnect: true, ProxyUrl: "http://127.0.0.1:4780"})
	msgChan = make(ExchangeApi.MessageChan)
)

func handleMsg(msgChan <-chan ExchangeApi.Message) {
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
			}
		}
	}
}

func TestHuobiWs_SubscribeTicker(t *testing.T) {
	if _, err := e.SubscribeTicker(symbol, msgChan); err == nil {
		handleMsg(msgChan)
	}
}

func TestHuobiWs_SubscribeOrderBook(t *testing.T) {
	if _, err := e.SubscribeOrderBook(symbol, 5, 100, false, msgChan); err == nil {
		handleMsg(msgChan)
		if err!=nil{
			fmt.Println(err)
		}
	}

}

func TestHuobiWs_SubscribeTrade(t *testing.T) {
	if _, err := e.SubscribeTrades(symbol, msgChan); err == nil {
		handleMsg(msgChan)
	}
}

func TestHuobiWs_SubscribeKline(t *testing.T) {
	if _, err := e.SubscribeKLine(symbol, ExchangeApi.KLine1Minute, msgChan); err == nil {
		handleMsg(msgChan)
	}
}

func TestHuobiWs_UnSubscribe(t *testing.T) {
	topic, err := e.SubscribeKLine(symbol, ExchangeApi.KLine1Minute, msgChan)
	if err == nil {
		go handleMsg(msgChan)
	}
	time.Sleep(time.Second * 10)
	e.UnSubscribe(topic, msgChan)
	time.Sleep(time.Second * 10)
}

func TestHuobiWs_SubscribeBalance(t *testing.T) {
	_, err := e.SubscribeBalance(symbol, msgChan)
	if err == nil {
		handleMsg(msgChan)
	}
}

func TestHuobiWs_SubscribeOrders(t *testing.T) {
	_, err := e.SubscribeOrder("BTC/USDT", msgChan)
	if err == nil {
		handleMsg(msgChan)
	}
}