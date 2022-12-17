package exchangerates

import (
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/helpers"
	"github.com/srackham/cryptor/internal/logger"
)

func TestExchangeRates(t *testing.T) {
	x := NewExchangeRates(&logger.Log{})

	today := helpers.DateNowString()
	rate, err := x.GetRate("USD", today, false)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 1.00, rate)

	rate, err = x.GetRate("NZD", today, false)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 1, len(*x.CacheData))
	assert.PassIf(t, rate > 0, "invalid NZD rate: %f", rate)

	rate, err = x.GetRate("AUD", today, false)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 1, len(*x.CacheData))
	assert.PassIf(t, rate > 0, "invalid AUD rate: %f", rate)

	_, err = x.GetRate("FOOBAR", today, false)
	assert.PassIf(t, err != nil, "should have return error for FOOBAR currency")
	assert.Equal(t, "unknown currency: FOOBAR", err.Error())
}
