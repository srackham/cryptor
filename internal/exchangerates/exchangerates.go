package exchangerates

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	. "github.com/srackham/cryptor/internal/cache"
	"github.com/srackham/cryptor/internal/helpers"
	"github.com/srackham/cryptor/internal/logger"
)

type ExchangeRates struct {
	log *logger.Log
	Cache[RatesCache]
}

func NewExchangeRates(log *logger.Log) ExchangeRates {
	data := make(RatesCache)
	return ExchangeRates{
		log:   log,
		Cache: *NewCache(&data),
	}
}

// getRates fetches a list of currency exchange rates against the USD from https://exchangerate.host/
// TODO getRates should be an interface and moved to exchangerateapi cf. prices.IPriceAPI.
func getRates(date string) (Rates, error) {
	rates := make(Rates)
	client := http.Client{}
	if date == helpers.DateNowString() {
		date = "latest"
	}
	url := fmt.Sprintf("https://api.exchangerate.host/%s", date)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return rates, err
	}
	resp, err := client.Do(request)
	if err != nil {
		return rates, err
	}
	// See https://www.sohamkamani.com/golang/json/#decoding-json-to-maps---unstructured-data
	var m map[string]any
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return rates, err
	}
	if m["success"] == false {
		return rates, fmt.Errorf("rate server query failed: %s", url)
	}
	m = m["rates"].(map[string]any)
	for k, v := range m {
		rates[strings.ToUpper(k)] = v.(float64)
	}
	return rates, nil
}

// updateCache checks the cache for rates on `date`.
// If the cache does not contain the rates for `date`  or `force` is `true`
// they are fetched and the cache is updated.
func (x *ExchangeRates) updateCache(date string, force bool) error {
	var rates Rates
	var ok bool
	var err error
	if _, ok = (*x.CacheData)[date]; !ok || force {
		x.log.Verbose("exchange rates request: %s", date)
		rates, err = getRates(date)
		if err != nil {
			return err
		}
		(*x.CacheData)[date] = rates
	}
	return nil
}

// GetRate returns the amount of `currency` that $1 USD would fetch on `date`.
// `currency` is a currency symbol.
// If `force` is `true` then then today's rates are unconditionally fetched and the cache updated.
func (x *ExchangeRates) GetRate(currency string, date string, force bool) (float64, error) {
	var err error
	if currency == "USD" {
		return 1.00, nil
	}
	err = x.updateCache(date, force)
	if err != nil {
		return 0.0, err
	}
	if rate, ok := (*x.CacheData)[date][strings.ToUpper(currency)]; !ok {
		return 0.0, fmt.Errorf("unknown currency: %s", currency)
	} else {
		return rate, nil
	}
}
