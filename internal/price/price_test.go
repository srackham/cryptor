package price

import (
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/helpers"
	"github.com/srackham/cryptor/internal/logger"
	"github.com/srackham/cryptor/internal/mockprice"
)

func TestPrice(t *testing.T) {
	r := NewPriceReader(&mockprice.Reader{}, &logger.Log{})
	date := helpers.TodaysDate()

	amt, err := r.GetPrice("BTC", date, true)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 10000, amt)

	amt, err = r.GetPrice("ETH", date, true)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 1000, amt)

	_, err = r.GetPrice("UNDEFINED_SYMBOL", date, true)
	assert.Equal(t, "unknown symbol: UNDEFINED_SYMBOL", err.Error())
}
