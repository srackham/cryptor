// Mock AOI implements pricereader.IPriceAPI interface.
package mockprice

import (
	"fmt"
	"regexp"
)

var prices = map[string]float64{
	"BTC":  10000.00,
	"ETH":  1000.00,
	"USDC": 1.00,
}

type Reader struct{}

func (r *Reader) LoadCacheFiles() error       { return nil }
func (r *Reader) SaveCacheFiles() error       { return nil }
func (r *Reader) SetCacheDir(cacheDir string) {}

func (r *Reader) GetPrice(symbol string, date string) (float64, error) {
	if amt, ok := prices[symbol]; ok {
		return amt, nil
	}
	re := regexp.MustCompile("^[a-zA-Z]{3,4}$")
	if re.MatchString(symbol) {
		return 1000.00, nil
	}
	return 0.00, fmt.Errorf("unknown symbol: %s", symbol)
}
