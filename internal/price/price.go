// Reads crypto currency prices in USD.
// TODO cache/save/restore all prices (cf. exchangerates package)
package price

import (
	"github.com/srackham/cryptor/internal/cache"
	"github.com/srackham/cryptor/internal/helpers"
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
	return PriceReader{
		API: reader,
		log: log,
		Cache: cache.Cache[cache.RatesCache]{
			CacheData: make(cache.RatesCache),
		},
	}
}

// TODO Add GetPrices()

// GetPrice returns the value in USD of the `symbol` crypto currency on `date`.
// If `date` is `"latest"` then then today's price is unconditionally fetched and the cache updated.
// If `date` is a "YYYY-MM-DD" date string the price is fetched only if it is not cached.
func (r *PriceReader) GetPrice(symbol string, date string) (float64, error) {
	var val float64
	var ok bool
	var err error
	if val, ok = r.CacheData[date][symbol]; !ok {
		// fmt.Printf("UPDATING CACHE: date=%s\n", date)
		val, err = r.API.GetPrice(symbol, date)
		if err != nil {
			return 0.0, err
		}
		r.log.Verbose("price request: %s %s %.2f", symbol, date, val)
		if date == "latest" {
			date = helpers.DateNowString()
		}
		if r.CacheData[date] == nil {
			r.CacheData[date] = make(cache.Rates)
		}
		r.CacheData[date][symbol] = val
	}
	return val, nil
}
