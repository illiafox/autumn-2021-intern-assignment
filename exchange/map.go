package exchange

import (
	"sync"
)

var exchangeMap sync.Map

func Add(abbreviation string, exchange float64) {
	exchangeMap.Store(abbreviation, exchange)
}

func GetExchange(abbreviation string) (float64, bool) {
	rate, ok := exchangeMap.Load(abbreviation)
	if !ok {
		return 0, ok
	}

	c, ok := rate.(float64)

	return c, ok
}
