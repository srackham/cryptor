// Binance pricereader.IPriceAPI interface implementation.
package binance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/srackham/cryptor/internal/cache"
	. "github.com/srackham/cryptor/internal/global"
)

// Cache data types.
type Rates map[string]float64 // Key = currency symbol; value = value in USD.

type TickerPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type PriceReader struct {
	*Context
	cache.Cache[Rates]
}

func (r *PriceReader) LoadCache() error { return nil }
func (r *PriceReader) SaveCache() error { return nil }

func NewPriceReader(ctx *Context) PriceReader {
	data := make(Rates)
	result := PriceReader{
		ctx,
		*cache.New(&data),
	}
	return result
}

func (r *PriceReader) getPrice(symbol string) (float64, error) {
	if symbol == "USDT" {
		return 1.0, nil // Because "USDTUSDT" is an illegal trading pair.
	}
	url := PRICE_QUERY + symbol + "USDT"
	var resp *http.Response
	var err error
	resp, err = r.HttpGet(url)
	if err != nil {
		return 0, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return 0, fmt.Errorf("invalid trading pair: %sUSDT", symbol)
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected HTTP response status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response: %v", err)
	}
	var ticker TickerPrice
	err = json.Unmarshal(body, &ticker)
	if err != nil {
		return 0, fmt.Errorf("error parsing JSON: %v", err)
	}
	price, err := strconv.ParseFloat(ticker.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting price to float: %v", err)
	}
	return price, nil
}

func (r *PriceReader) GetCachedPrice(symbol string) (price float64, err error) {
	var ok bool
	if price, ok = (*r.CacheData)[strings.ToUpper(symbol)]; !ok {
		price, err = r.getPrice(symbol)
		if err != nil {
			return 0.0, err
		}
		(*r.CacheData)[strings.ToUpper(symbol)] = price
	}
	return
}
