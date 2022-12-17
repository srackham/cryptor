// Reads crypto currency prices in USD.
// TODO cache/save/restore all prices (cf. exchangerates package)
package price

import (
	"github.com/srackham/cryptor/internal/cache"
	"github.com/srackham/cryptor/internal/logger"
)

type IPriceAPI interface {
	GetPrice(symbol string, date string) (float64, error)
	SetCacheDir(cacheDir string)
	LoadCacheFiles() error
	SaveCacheFiles() error
}

type PriceReader struct {
	API IPriceAPI
	log *logger.Log
	cache.Cache[cache.RatesCache]
}

func NewPriceReader(reader IPriceAPI, log *logger.Log) PriceReader {
	data := make(cache.RatesCache)
	return PriceReader{
		API:   reader,
		log:   log,
		Cache: *cache.NewCache(&data),
	}
}

// TODO Add GetPrices()

// GetPrice returns the value in USD of the `symbol` crypto currency on `date`.
// If `force` is `true` then then today's price is unconditionally fetched and the cache updated.
func (r *PriceReader) GetPrice(symbol string, date string, force bool) (float64, error) {
	var val float64
	var ok bool
	var err error
	if val, ok = (*r.CacheData)[date][symbol]; !ok || force {
		val, err = r.API.GetPrice(symbol, date)
		if err != nil {
			return 0.0, err
		}
		r.log.Verbose("price request: %s %s %.2f", symbol, date, val)
		if (*r.CacheData)[date] == nil {
			(*r.CacheData)[date] = make(cache.Rates)
		}
		(*r.CacheData)[date][symbol] = val
	}
	return val, nil
}
