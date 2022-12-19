package cgprice

import (
	"os"
	"reflect"
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/helpers"
)

func TestPrice(t *testing.T) {
	cg := NewReader()
	tmpdir, err := os.MkdirTemp("", "cryptor")
	if err != nil {
		return
	}
	cg.SetCacheDir(tmpdir)
	cg.LoadCacheFiles()

	today := helpers.TodaysDate()
	amt, err := cg.GetPrice("BTC", today)
	assert.PassIf(t, err == nil, "%#v", err)
	assert.PassIf(t, amt > 0.00, "expected BTC price to be greater than zero")

	amt, err = cg.GetPrice("ETH", "2022-06-01")
	assert.PassIf(t, err == nil, "%#v", err)
	assert.PassIf(t, amt > 0.00, "expected 2022-06-01 ETH price to be greater than zero")

	_, err = cg.GetPrice("undefined_symbol", today)
	assert.Equal(t, "unsupported coin: undefined_symbol", err.Error())

	savedCache := cg.coinList.CacheData
	err = cg.SaveCacheFiles()
	assert.PassIf(t, err == nil, "%v", err)
	err = cg.LoadCacheFiles()
	assert.PassIf(t, err == nil, "%v", err)
	assert.PassIf(t, reflect.DeepEqual(&savedCache, &cg.coinList.CacheData), "expected:\n%v\n\ngot:\n%v", &savedCache, &cg.coinList.CacheData)
}
