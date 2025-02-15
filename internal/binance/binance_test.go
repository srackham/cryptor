package binance

import (
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/mock"
)

func TestPrice(t *testing.T) {
	ctx := mock.NewContext()
	reader := NewPriceReader(&ctx)

	price, err := reader.GetCachedPrice("BTC")
	wanted := 100_000.0
	assert.PassIf(t, err == nil, "%#v", err)
	assert.Equal(t, wanted, price)

	price, err = reader.GetCachedPrice("ETH")
	wanted = 1000.0
	assert.PassIf(t, err == nil, "%#v", err)
	assert.Equal(t, wanted, price)

	_, err = reader.GetCachedPrice("INVALID_SYMBOL")
	assert.Equal(t, "invalid trading pair: INVALID_SYMBOLUSDT", err.Error())
}
