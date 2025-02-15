package global

import (
	"io"
	"net/http"
	"time"
)

// Build ldflags.
var (
	// VERS is the latest cryptor version tag. Set by linker -ldflags "-X main.VERS=..."
	VERS = "v0.3.0"
	// OS is the target operating system and architecture. Set by linker -ldflags "-X main.OS=..."
	OS = "-"
	// BUILD is the date the executable was built.
	BUILT = "-"
	// COMMIT is the Git commit hash.
	COMMIT = "-"

	PRICE_QUERY  = "https://api.binance.com/api/v1/ticker/price?symbol="
	XRATES_QUERY = "https://openexchangerates.org/api/latest.json?app_id="
)

// Application dependency injection container
type Context struct {
	Stdout    io.Writer
	Stderr    io.Writer
	DataDir   string
	CacheDir  string
	ConfigDir string
	Now       func() time.Time
	HttpGet   func(url string) (*http.Response, error)
}
