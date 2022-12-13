package exchangerates

import (
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/helpers"
	"github.com/srackham/cryptor/internal/logger"
)

func TestExchangeRates(t *testing.T) {
	x := NewExchangeRates(&logger.Log{})

	rate, err := x.GetRate("USD", "latest")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 1.00, rate)

	rate, err = x.GetRate("NZD", "latest")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 1, len(x.CacheData))
	assert.PassIf(t, rate > 0, "invalid NZD rate: %f", rate)

	today := helpers.DateNowString()
	rate, err = x.GetRate("AUD", today)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 1, len(x.CacheData))
	assert.PassIf(t, rate > 0, "invalid AUD rate: %f", rate)

	rate, err = x.GetRate("AUD", "2022-06-01")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 2, len(x.CacheData))
	assert.PassIf(t, rate > 0, "invalid AUD rate: %f", rate)

	_, err = x.GetRate("FOOBAR", today)
	assert.PassIf(t, err != nil, "should have return error for FOOBAR currency")
	assert.Equal(t, "unknown currency: FOOBAR", err.Error())
}
