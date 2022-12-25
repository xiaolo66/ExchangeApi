package okex
import (
	"github.com/xiaolo66/ExchangeApi"
)
type Okex struct {
	OkexRest
	OkexWs
}

func New(options ExchangeApi.Options) *Okex {
	instance := &Okex{}
	instance.OkexRest.Init(options)
	instance.OkexWs.Init(options)

	if len(options.Markets) == 0 {
		instance.OkexWs.Option.Markets, _ = instance.FetchMarkets()
	}

	return instance
}
