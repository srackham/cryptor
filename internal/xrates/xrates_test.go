package xrates

import (
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/fsx"
	"github.com/srackham/cryptor/internal/mock"
)

func TestExchangeRates(t *testing.T) {
	ctx := mock.NewContext()
	tmpdir := mock.MkdirTemp(t)
	ctx.CacheDir = tmpdir
	x := New(&ctx)

	rate, err := x.GetCachedRate("USD", false)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 1.00, rate)

	rate, err = x.GetCachedRate("NZD", false)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 1.5, rate)

	rate, err = x.GetCachedRate("AUD", false)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 1.6, rate)

	_, err = x.GetCachedRate("", false)
	assert.PassIf(t, err != nil, "should have return error for blank currency")
	assert.Equal(t, "no currency specified", err.Error())

	_, err = x.GetCachedRate("FOOBAR", false)
	assert.PassIf(t, err != nil, "should have returned error for FOOBAR currency")
	assert.Equal(t, "unknown currency: FOOBAR", err.Error())

	err = x.Save(x.CacheFile())
	assert.PassIf(t, err == nil, "error writing exchange rates cache: \"%v\": %v", x.CacheFile(), err)
	got, err := fsx.ReadFile(x.CacheFile())
	assert.PassIf(t, err == nil, "error reading exchange rates cache: \"%v\": %v", x.CacheFile(), err)
	assert.Equal(t, `{
  "2000-12-01": {
    "AUD": 1.6,
    "NZD": 1.5,
    "USD": 1
  }
}`, got)

	rates := *x.CacheData
	rates["1999-12-31"] = rates["2000-12-01"]
	assert.Equal(t, len(*(x.CacheData)), 2)
	_, err = x.GetCachedRate("AUD", true) // Replaces rates with only today's rates
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, len(*(x.CacheData)), 1)
	err = x.Save(x.CacheFile())
	assert.PassIf(t, err == nil, "error writing exchange rates cache: \"%v\": %v", x.CacheFile(), err)
	got, err = fsx.ReadFile(x.CacheFile())
	assert.PassIf(t, err == nil, "error reading exchange rates cache: \"%v\": %v", x.CacheFile(), err)
	assert.Equal(t, `{
  "2000-12-01": {
    "AUD": 1.6,
    "NZD": 1.5,
    "USD": 1
  }
}`, got)
}
