package mock

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	. "github.com/srackham/cryptor/internal/global"
	"github.com/srackham/go-utils/assert"
	"github.com/srackham/go-utils/fsx"
)

func MkdirTemp(t *testing.T) string {
	tmpdir := path.Join(os.TempDir(), "cryptor-temp")
	err := os.RemoveAll(tmpdir)
	assert.PassIf(t, err == nil, "%v", err)
	if !fsx.DirExists(tmpdir) {
		err := os.Mkdir(tmpdir, 0o755)
		assert.PassIf(t, err == nil, "%v", err)
	}
	return tmpdir
}

func NewContext() Context {
	return Context{
		Stdout:    new(bytes.Buffer),
		Stderr:    new(bytes.Buffer),
		DataDir:   "../../testdata/data",
		CacheDir:  "../../testdata/cache",
		ConfigDir: "../../testdata",
		Now:       now,
		HttpGet:   httpGet,
	}
}

func now() time.Time {
	layout := "2006-01-02 15:04:05"
	mockTime := "2000-12-01 12:30:00"
	t, _ := time.Parse(layout, mockTime)
	return t
}

func httpGet(url string) (resp *http.Response, err error) {
	switch url {
	case PRICE_QUERY + "BTC" + "USDT":
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"symbol":"BTCUSDT","price":"100000.00000000"}`)),
		}, nil
	case PRICE_QUERY + "ETH" + "USDT":
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"symbol":"ETHUSDT","price":"1000.00"}`)),
		}, nil
	case PRICE_QUERY + "USDC" + "USDT":
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"symbol":"USDCUSDT","price":"1.00"}`)),
		}, nil
	case PRICE_QUERY + "BAD_JSON" + "USDT":
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(``)),
		}, nil
	case PRICE_QUERY + "INVALID_SYMBOL" + "USDT":
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(strings.NewReader(``)),
		}, nil
	case XRATES_QUERY + "1234":
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(strings.NewReader(`{
  "rates": {
    "AUD": 1.6,
    "NZD": 1.5,
    "USD": 1
  }
}`)),
		}, nil
	default:
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader(`not found`)),
		}, nil
	}
}
