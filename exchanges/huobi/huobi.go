package huobi

import "github.com/xiaolo66/ExchangeApi"

type Huobi struct {
	HuobiRest
	HuobiWs
}

func New(options ExchangeApi.Options) *Huobi {
	instance := &Huobi{}
	instance.HuobiRest.Init(options)
	instance.HuobiWs.Init(options)

	if len(options.Markets) == 0 {
		instance.HuobiWs.Option.Markets, _ = instance.FetchMarkets()
	}
	return instance
}
