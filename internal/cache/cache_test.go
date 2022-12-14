package cache

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/fsx"
)

func TestRates(t *testing.T) {
	data := make(RatesCache)
	r := NewCache(&data)
	tmpdir, err := os.MkdirTemp("", "cryptor")
	assert.PassIf(t, err == nil, "%v", err)
	r.CacheFile = filepath.Join(tmpdir, "history.json")

	err = r.SaveCacheFile()
	assert.PassIf(t, err == nil, "%v", err)
	savedCache := (*r.CacheData)
	err = r.LoadCacheFile()
	assert.PassIf(t, err == nil, "%v", err)
	assert.PassIf(t, reflect.DeepEqual(savedCache, (*r.CacheData)), "expected:\n%v\n\ngot:\n%v", savedCache, (*r.CacheData))

	(*r.CacheData)["2022-06-01"] = make(Rates)
	(*r.CacheData)["2022-06-01"]["BTC"] = 10000.00
	err = r.SaveCacheFile()
	assert.PassIf(t, err == nil, "%v", err)
	savedCache = (*r.CacheData)
	err = r.LoadCacheFile()
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, (*r.CacheData)["2022-06-01"]["BTC"], 10000.00)
	assert.PassIf(t, reflect.DeepEqual(savedCache, (*r.CacheData)), "expected:\n%v\n\ngot:\n%v", savedCache, (*r.CacheData))

	s, err := fsx.ReadFile(r.CacheFile)
	sha := sha256.Sum256([]byte(s))
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, sha, r.sha256)
	(*r.CacheData)["2022-06-01"]["BTC"] = 0.00
	err = r.SaveCacheFile()
	assert.PassIf(t, err == nil, "%v", err)
	assert.NotEqual(t, sha, r.sha256)
	sha = r.sha256
	err = r.LoadCacheFile()
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, sha, r.sha256)
}
