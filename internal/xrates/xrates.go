package xrates

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
	fetched bool
	url     string
}

func NewExchangeRates(url string, log *logger.Log) ExchangeRates {
	data := make(RatesCache)
	return ExchangeRates{
		log:     log,
		Cache:   *NewCache(&data),
		fetched: false,
		url:     url,
	}
}

// getRates fetches a list of currency exchange rates against the USD
// TODO getRates should be an IXRatesAPI interface cf. prices.IPriceAPI.
func getRates(url string) (Rates, error) {
	if helpers.GithubActions() {
		// getRates() requires HTTP access and should never execute from Github Actions.
		mockRates := Rates{"usd": 1.0, "nzd": 1.5, "aud": 1.6}
		return mockRates, nil
	}
	rates := make(Rates)
	if url == "" {
		return rates, fmt.Errorf("exchange rate request: URL has not been specified")
	}
	client := http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return rates, fmt.Errorf("exchange rate request: %s: %s", url, err.Error())
	}
	resp, err := client.Do(request)
	if err != nil {
		return rates, fmt.Errorf("exchange rate request: %s: %s", url, err.Error())
	}
	// See https://www.sohamkamani.com/golang/json/#decoding-json-to-maps---unstructured-data
	var m map[string]any
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return rates, fmt.Errorf("exchange rate decode: %s", err.Error())
	}
	_, exists := m["rates"]
	if !exists {
		jsonData, err := json.Marshal(m)
		if err != nil {
			return rates, fmt.Errorf("invalid exchange rate response: %s: %v", url, m)
		}
		return rates, fmt.Errorf("invalid exchange rate response: %s: %s", url, string(jsonData))
	}
	m = m["rates"].(map[string]any)
	for k, v := range m {
		rates[strings.ToLower(k)] = v.(float64)
	}
	return rates, nil
}

// GetRate returns the amount of `currency` that $1 USD would buy at today's rates.
// `currency` is a currency symbol.
// If `force` is `true` then then today's rates are unconditionally fetched and the cache updated.
// TODO tests
func (x *ExchangeRates) GetRate(currency string, force bool) (float64, error) {
	if currency == "USD" {
		return 1.00, nil
	}
	var rate float64
	var ok bool
	today := helpers.TodaysDate()
	force = force && !x.fetched // Don't force if rates have previously been fetched during this session.
	if rate, ok = (*x.CacheData)[today][strings.ToLower(currency)]; !ok || force {
		x.log.Note("exchange rates request: " + helpers.TodaysDate())
		rates, err := getRates(x.url)
		if err != nil {
			return 0.0, err
		}
		x.CacheData = &RatesCache{today: rates}
		if rate, ok = (*x.CacheData)[today][strings.ToLower(currency)]; !ok {
			return 0.0, fmt.Errorf("unknown currency: %s", currency)
		}
		x.fetched = true
	}
	return rate, nil
}
