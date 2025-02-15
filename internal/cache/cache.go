package cache

import (
	"crypto/sha256"
	"encoding/json"

	"github.com/srackham/cryptor/internal/fsx"
)

// Cache is designed to be embedded, it implements file-based data persistance with load and save functions.
// The cache data is external to the Cache struct and is accessed via the CacheData pointer.
type Cache[T any] struct {
	CacheData *T
	sha256    [32]byte // Cache file checksum.
}

func New[T any](data *T) *Cache[T] {
	return &Cache[T]{
		CacheData: data,
	}
}

func (c *Cache[T]) Load(cacheFile string) error {
	var err error
	if fsx.FileExists(cacheFile) {
		s, err := fsx.ReadFile(cacheFile)
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
func (c *Cache[T]) Save(cacheFile string) error {
	json, err := json.MarshalIndent(*c.CacheData, "", "  ")
	if err == nil {
		sha := sha256.Sum256(json)
		if c.sha256 != sha {
			err = fsx.WriteFile(cacheFile, string(json))
		}
		c.sha256 = sha
	}
	return err
}
