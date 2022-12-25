package factory

import (
	"github.com/xiaolo66/ExchangeApi"
	"github.com/xiaolo66/ExchangeApi/exchanges/binance"
	"github.com/xiaolo66/ExchangeApi/exchanges/okex"
)

func NewExchange(t ExchangeApi.ExchangeType, option ExchangeApi.Options) ExchangeApi.IExchange {
	switch t {
	case ExchangeApi.Binance:
		return binance.New(option)
	//case ExchangeApi.ZB:
	//	return zb.New(option)
	case ExchangeApi.Okex:
		return okex.New(option)
	//case ExchangeApi.Huobi:
	//	return huobi.New(option)
	//case ExchangeApi.GateIo:
	//	return gateio.New(option)
	}
	return nil
}

func NewFutureExchange(t ExchangeApi.ExchangeType, option ExchangeApi.Options, futureOptions ExchangeApi.FutureOptions) ExchangeApi.IFutureExchange {
	switch t {
	case ExchangeApi.Binance:
		return binance.NewFuture(option,futureOptions)
	//case ExchangeApi.ZB:
	//	return zb.NewFuture(option, futureOptions)
	}
	return nil
}
