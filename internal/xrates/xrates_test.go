package xrates

import (
	"os"
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/logger"
)

func TestExchangeRates(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skip on Github Actions because this test requires HTTP access")
	}

	x := NewExchangeRates(&logger.Log{})

	rate, err := x.GetRate("USD", false)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 1.00, rate)

	rate, err = x.GetRate("NZD", false)
	assert.PassIf(t, err == nil, "%v", err)
	assert.PassIf(t, rate > 0, "invalid NZD rate: %f", rate)

	rate, err = x.GetRate("AUD", false)
	assert.PassIf(t, err == nil, "%v", err)
	assert.PassIf(t, rate > 0, "invalid AUD rate: %f", rate)

	_, err = x.GetRate("FOOBAR", false)
	assert.PassIf(t, err != nil, "should have return error for FOOBAR currency")
	assert.Equal(t, "unknown currency: FOOBAR", err.Error())
}
