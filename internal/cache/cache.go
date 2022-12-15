package cache

import (
	"crypto/sha256"
	"encoding/json"

	"github.com/srackham/cryptor/internal/fsx"
)

// Cache data types.
type Rates map[string]float64    // Key = currency symbol; value = value in USD.
type RatesCache map[string]Rates // Key = date string "YYYY-MM-DD".

// Cache is designed to be embedded, it implements file-based data persistance with load and save functions.
// The cache data is external to the Cache struct and is accessed via the CacheData pointer.
type Cache[T any] struct {
	CacheData *T
	CacheFile string
	sha256    [32]byte // Cache file checksum.
}

func NewCache[T any](data *T) *Cache[T] {
	return &Cache[T]{
		CacheData: data,
	}
}

func (c *Cache[T]) Load() error {
	var err error
	if fsx.FileExists(c.CacheFile) {
		s, err := fsx.ReadFile(c.CacheFile)
		if err == nil {
			err = json.Unmarshal([]byte(s), c.CacheData)
			if err == nil {
				c.sha256 = sha256.Sum256([]byte(s))
			}
		}
	}
	return err
}

// Save writes the cache to disk if it has been modified.
func (c *Cache[T]) Save() error {
	json, err := json.MarshalIndent(*c.CacheData, "", "  ")
	if err == nil {
		sha := sha256.Sum256(json)
		if c.sha256 != sha {
			err = fsx.WriteFile(c.CacheFile, string(json))
		}
		c.sha256 = sha
	}
	return err
}
