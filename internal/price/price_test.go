package price

import (
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/logger"
	"github.com/srackham/cryptor/internal/mockprice"
)

func TestPrice(t *testing.T) {
	r := NewPriceReader(&mockprice.Reader{}, &logger.Log{})

	amt, err := r.GetPrice("BTC", "latest")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 10000, amt)

	amt, err = r.GetPrice("ETH", "latest")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 1000, amt)

	_, err = r.GetPrice("UNDEFINED_SYMBOL", "latest")
	assert.Equal(t, "unknown symbol: UNDEFINED_SYMBOL", err.Error())
}
