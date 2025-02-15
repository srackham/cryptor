package cache

import (
	"crypto/sha256"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/fsx"
	"github.com/srackham/cryptor/internal/mock"
)

// Cache data types.
type Rates map[string]float64 // Key = currency symbol; value = value in USD.

func TestRates(t *testing.T) {
	data := make(Rates)
	r := New(&data)
	tmpdir := mock.MkdirTemp(t)
	valuationsFile := filepath.Join(tmpdir, "valuations.json")

	err := r.Save(valuationsFile)
	assert.PassIf(t, err == nil, "%v", err)
	savedCache := (*r.CacheData)
	err = r.Load(valuationsFile)
	assert.PassIf(t, err == nil, "%v", err)
	assert.PassIf(t, reflect.DeepEqual(savedCache, (*r.CacheData)), "expected:\n%v\n\ngot:\n%v", savedCache, (*r.CacheData))

	(*r.CacheData) = make(Rates)
	(*r.CacheData)["USD"] = 1.00
	err = r.Save(valuationsFile)
	assert.PassIf(t, err == nil, "%v", err)
	savedCache = (*r.CacheData)
	err = r.Load(valuationsFile)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, (*r.CacheData)["USD"], 1.00)
	assert.PassIf(t, reflect.DeepEqual(savedCache, (*r.CacheData)), "expected:\n%v\n\ngot:\n%v", savedCache, (*r.CacheData))

	s, err := fsx.ReadFile(valuationsFile)
	sha := sha256.Sum256([]byte(s))
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, sha, r.sha256)
	(*r.CacheData)["USD"] = 0.00
	err = r.Save(valuationsFile)
	assert.PassIf(t, err == nil, "%v", err)
	assert.NotEqual(t, sha, r.sha256)
	sha = r.sha256
	err = r.Load(valuationsFile)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, sha, r.sha256)
}
