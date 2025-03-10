package xrates

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/srackham/cryptor/internal/cache"
	"github.com/srackham/cryptor/internal/config"
	. "github.com/srackham/cryptor/internal/global"
)

// Cache data types.
type Rates map[string]float64        // Key = currency symbol; value = value in USD.
type RatesCacheData map[string]Rates // Key = date string "YYYY-MM-DD".

type ExchangeRates struct {
	*Context
	*cache.Cache[RatesCacheData]
}

func New(ctx *Context) ExchangeRates {
	result := ExchangeRates{}
	result.Context = ctx
	data := make(RatesCacheData)
	result.Cache = cache.New(&data)
	return result
}

func (x *ExchangeRates) ConfigFile() string {
	return filepath.Join(x.ConfigDir, "config.yaml")
}

func (x *ExchangeRates) CacheFile() string {
	return filepath.Join(x.CacheDir, "exchange-rates.json")
}

// getRates fetches a list of currency exchange rates against the USD
func (x *ExchangeRates) getRates() (Rates, error) {
	rates := make(Rates)
	conf, err := config.LoadConfig(x.ConfigFile())
	if err != nil {
		return rates, err
	}
	url := XRATES_QUERY + conf.XratesAppId
	resp, err := x.HttpGet(url)
	if err != nil {
		return rates, fmt.Errorf("exchange rate request: %s: %s", url, err.Error())
	}
	defer resp.Body.Close()

	// See https://www.sohamkamani.com/golang/json/#decoding-json-to-maps---unstructured-data
	var m map[string]any
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return rates, fmt.Errorf("exchange rates decode: %s", err.Error())
	}
	_, ok := m["rates"]
	if !ok {
		return rates, fmt.Errorf("invalid exchange rate response: %s: %v", url, m)
	}
	for k, v := range m["rates"].(map[string]any) {
		rates[strings.ToUpper(k)] = v.(float64)
	}
	return rates, nil
}

// GetCachedRate returns the amount of `currency` that $1 USD would buy at today's rates.
// `symbol` is a currency symbol.
// If `force` is `true` then then today's rates are unconditionally fetched and the cache updated.
// TODO tests
func (x *ExchangeRates) GetCachedRate(currency string, force bool) (float64, error) {
	if currency == "" {
		return 0.0, fmt.Errorf("no currency specified")
	}
	if currency == "USD" {
		return 1.00, nil
	}
	today := x.Now().Format("2006-01-02")
	rate, ok := (*x.CacheData)[today][strings.ToUpper(currency)]
	if !ok || force {
		rates, err := x.getRates()
		if err != nil {
			return 0.0, err
		}
		x.CacheData = &(RatesCacheData{today: rates})
		if rate, ok = (*x.CacheData)[today][strings.ToUpper(currency)]; !ok {
			return 0.0, fmt.Errorf("unknown currency: %s", currency)
		}
	}
	return rate, nil
}
