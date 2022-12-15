package portfolio

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/helpers"
	"github.com/srackham/cryptor/internal/logger"
	"github.com/srackham/cryptor/internal/mockprice"
	"github.com/srackham/cryptor/internal/price"
)

func TestPortfolios(t *testing.T) {
	ps, err := LoadHistoryFile("../../testdata/portfolios.json")
	assert.PassIf(t, err == nil, "error reading JSON file")
	assert.Equal(t, 2, len(ps))
	assert.Equal(t, 3, len(ps[0].Assets))
	assert.Equal(t, ps[0].Assets[2], Asset{
		Symbol:      "USDC",
		Amount:      100.00,
		USD:         0.0,
		Description: "On exchange",
	})
}

func TestHistory(t *testing.T) {
	h, err := LoadHistoryFile("../../testdata/history.json")
	assert.PassIf(t, err == nil, "error reading JSON file")
	assert.Equal(t, 14, len(h))
	assert.Equal(t, "2022-12-02", h[0].Date)
	assert.Equal(t, "2022-12-06", h[13].Date)
	assert.Equal(t, "12:34:56", h[13].Time)
	assert.Equal(t, 3, len(h[0].Assets))
	assert.Equal(t, h[0].Assets[0], Asset{
		Symbol:      "BTC",
		Amount:      0.5,
		USD:         5000.0,
		Description: "Cold storage",
	})
	assert.Equal(t, h[2].Assets[2], Asset{
		Symbol:      "USDC",
		Amount:      100.0,
		USD:         100.00,
		Description: "On exchange",
	})
}

func TestSaveHistoryFile(t *testing.T) {
	h, err := LoadHistoryFile("../../testdata/history.json")
	assert.PassIf(t, err == nil, "error reading JSON file")
	tmpdir, err := os.MkdirTemp("", "cryptor")
	assert.PassIf(t, err == nil, "%v", err)
	fname := filepath.Join(tmpdir, "history.json")
	err = h.SaveHistoryFile(fname)
	assert.PassIf(t, err == nil, "%v", err)
	savedHistory := h
	h, err = LoadHistoryFile(fname)
	assert.PassIf(t, err == nil, "%v", err)
	assert.PassIf(t, reflect.DeepEqual(savedHistory, h),
		"expected:\n%v\n\ngot:\n%v", savedHistory, h)
}

func TestEvaluate(t *testing.T) {
	ps, err := LoadHistoryFile("../../testdata/portfolios.json")
	assert.PassIf(t, err == nil, "error reading JSON file")
	p := ps[0]
	reader := price.NewPriceReader(&mockprice.Reader{}, &logger.Log{})
	prices, err := ps.GetPrices(reader, helpers.DateNowString(), true)
	assert.PassIf(t, err == nil, "error valuating portfolio: %v", err)
	p.SetUSDValues(prices)
	p.Assets.SortByValue()
	assert.Equal(t, 5000.0, p.Assets[0].USD)
	assert.Equal(t, 2500.0, p.Assets[1].USD)
	assert.Equal(t, 100.0, p.Assets[2].USD)
	p.Assets[0].USD = 1000.00
	p.Assets.SortByValue()
	assert.Equal(t, 2500.0, p.Assets[0].USD)
	assert.Equal(t, 1000.0, p.Assets[1].USD)
	assert.Equal(t, 100.0, p.Assets[2].USD)
}

func TestSortAndFilter(t *testing.T) {
	h, err := LoadHistoryFile("../../testdata/history.json")
	assert.PassIf(t, err == nil, "error reading JSON file")

	h2 := h.FilterByDate("2022-12-02")
	assert.Equal(t, 2, len(h2))

	h2 = h.FilterByName("personal")
	assert.Equal(t, 7, len(h2))
	h2.SortByDate()
	assert.Equal(t, "2022-12-01", h2[0].Date)
	assert.Equal(t, "2022-12-07", h2[6].Date)

	h2 = h.FilterByDate("2022-12-07").FilterByName("joint")
	assert.Equal(t, 1, len(h2))
	assert.Equal(t, "2022-12-07", h2[0].Date)

	h2 = h.FilterByName("personal", "joint")
	assert.Equal(t, 14, len(h2))
}
