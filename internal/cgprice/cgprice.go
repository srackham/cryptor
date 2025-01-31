// CoinGecko API (implements pricereader.IPriceAPI interface).
// Uses CoinGecko API Client: https://github.com/superoo7/go-gecko
package cgprice

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/srackham/cryptor/internal/cache"
	"github.com/srackham/cryptor/internal/helpers"
	gecko "github.com/superoo7/go-gecko/v3"
	"github.com/superoo7/go-gecko/v3/types"
)

type Reader struct {
	coinList cache.Cache[types.CoinList]
}

func NewReader() *Reader {
	data := types.CoinList{}
	return &Reader{
		coinList: *cache.NewCache(&data),
	}
}

func (r *Reader) LoadCacheFiles() error {
	return r.coinList.Load()
}

func (r *Reader) SaveCacheFiles() error {
	return r.coinList.Save()
}

func (r *Reader) SetCacheDir(cacheDir string) {
	r.coinList.CacheFile = filepath.Join(cacheDir, "gecko-coins-list.json")
}

// getId returns the CoinGecko API ID of the crypto symbol.
func (r *Reader) GetCoinId(symbol string) (string, error) {
	symbol = strings.ToLower(symbol)
	// FIXME: Kludge to get correct ID for well known coins.
	switch symbol {
	case "btc":
		return "bitcoin", nil
	case "eth":
		return "ethereum", nil
	case "sol":
		return "solana", nil
	case "avax":
		return "avalanche-2", nil
	}
	if len(*(r.coinList.CacheData)) == 0 {
		cl, err := getCoinsList()
		if err != nil {
			return "", err
		}
		r.coinList.CacheData = cl
	}
	for _, c := range *r.coinList.CacheData {
		if c.Symbol == symbol {
			return c.ID, nil
		}
	}
	return "", fmt.Errorf("unsupported coin: %s", symbol)
}

func (r *Reader) GetPrice(symbol string, date string) (float64, error) {
	id, err := r.GetCoinId(symbol)
	if err != nil {
		return 0.00, err
	}
	var pd PriceData
	vc := "usd" // Must be lowercase.
	// TODO if there is a rate limit error back off and retry for 3x10s.
	// TODO if there is an error retry using all IDs that map to symbol.
	if date == helpers.TodaysDate() {
		pd, err = getCurrentPriceData(id, vc)
	} else {
		date = date[8:10] + "-" + date[5:7] + "-" + date[0:4] // Convert YYYY-MM-DD date to DD-MM-YYYY
		pd, err = getHistoricalPriceData(id, vc, date)
	}
	if err != nil {
		return 0.00, err
	}
	return pd.Price, nil
}

/*
The remaining code taken from package github.com/hakochaz/crypto-price-cli
*/

type PriceData struct {
	Coin   string
	VC     string
	Price  float64
	Date   string
	Amount float64
	Value  float64
}

func getCoinsList() (*types.CoinList, error) {
	cg := gecko.NewClient(nil)
	cl, err := cg.CoinsList()
	if err != nil {
		return nil, err
	}
	return cl, nil
}

func getCurrentPriceData(id, vc string) (PriceData, error) {
	pd := PriceData{}
	cg := gecko.NewClient(nil)
	// TODO use SimplePrice() to handle multiple ids
	sp, err := cg.SimpleSinglePrice(id, vc)
	if err != nil {
		return pd, err
	}
	c := (*sp)
	pd.Coin = id
	pd.VC = vc
	pd.Price = float64(c.MarketPrice)
	return pd, nil
}

// GetHistoricalPriceData gets historical price data for a coin
// versus another currency
func getHistoricalPriceData(id, vc, d string) (PriceData, error) {
	pd := PriceData{}
	cg := gecko.NewClient(nil)
	sp, err := cg.CoinsIDHistory(id, d, true)
	if err != nil {
		return pd, err
	}
	c := (*sp)
	if c.MarketData == nil {
		return pd, fmt.Errorf("missing historical price data: id: %q: versus currency: %q: date: %q", id, vc, d)
	}
	if c.MarketData.CurrentPrice[vc] == 0 {
		return pd, fmt.Errorf("incompatible versus currency: %q", vc)
	}
	pd.Coin = id
	pd.VC = vc
	pd.Price = c.MarketData.CurrentPrice[vc]
	pd.Date = d
	return pd, nil
}
