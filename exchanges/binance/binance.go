package binance

import "github.com/xiaolo66/ExchangeApi"
type Binance struct {
	BinanceRest
	BinanceWs
}

func New(options ExchangeApi.Options) *Binance {
	instance := &Binance{}
	instance.BinanceRest.Init(options)
	instance.BinanceWs.Init(options)

	if len(options.Markets) == 0 {
		instance.BinanceWs.Option.Markets, _ = instance.FetchMarkets()
	}

	return instance
}
