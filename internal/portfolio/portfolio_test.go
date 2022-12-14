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

func TestIsValidName(t *testing.T) {
	type test struct {
		input string
		want  bool
	}
	tests := []test{
		{input: "", want: false},
		{input: "-foo", want: false},
		{input: "foo bar", want: false},
		{input: "foo", want: true},
	}
	for _, tc := range tests {
		got := IsValidName(tc.input)
		assert.PassIf(t, tc.want == got, "input: %q: wanted: %v: got: %v", tc.input, tc.want, got)
	}
}

func TestLoadValuationsFile(t *testing.T) {
	h, err := LoadValuationsFile("../../testdata/valuations.json")
	assert.PassIf(t, err == nil, "error reading JSON file")
	assert.Equal(t, 14, len(h))
	assert.Equal(t, "2022-12-02", h[0].Date)
	assert.Equal(t, "2022-12-06", h[13].Date)
	assert.Equal(t, 3, len(h[0].Assets))
	assert.Equal(t, h[0].Assets[0], Asset{
		Symbol: "BTC",
		Amount: 0.5,
		Value:  5000.0,
	})
	assert.Equal(t, h[2].Assets[2], Asset{
		Symbol: "USDC",
		Amount: 100.0,
		Value:  100.00,
	})
}

func TestSaveValuationsFile(t *testing.T) {
	h, err := LoadValuationsFile("../../testdata/valuations.json")
	assert.PassIf(t, err == nil, "error reading JSON file")
	tmpdir, err := os.MkdirTemp("", "cryptor")
	assert.PassIf(t, err == nil, "%v", err)
	fname := filepath.Join(tmpdir, "valuations.json")
	err = h.SaveValuationsFile(fname)
	assert.PassIf(t, err == nil, "%v", err)
	savedValuations := h
	h, err = LoadValuationsFile(fname)
	assert.PassIf(t, err == nil, "%v", err)
	assert.PassIf(t, reflect.DeepEqual(savedValuations, h),
		"expected:\n%v\n\ngot:\n%v", savedValuations, h)
}

func TestParseCurrency(t *testing.T) {
	value, currency, err := ParseCurrency("$1,234.56")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 1234.56, value)
	assert.Equal(t, "USD", currency)

	value, currency, err = ParseCurrency("$1,234.56 NZD")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 1234.56, value)
	assert.Equal(t, "NZD", currency)

	_, _, err = ParseCurrency("")
	assert.PassIf(t, err != nil, "%v", err)
	assert.Equal(t, "invalid currency value: \"\"", err.Error())

	_, _, err = ParseCurrency("foo")
	assert.PassIf(t, err != nil, "%v", err)
	assert.Equal(t, "invalid currency value: \"foo\"", err.Error())

	_, _, err = ParseCurrency("$1,234.56 NZD qux")
	assert.PassIf(t, err != nil, "%v", err)
	assert.Equal(t, "invalid currency value: \"$1,234.56 NZD qux\"", err.Error())
}

func TestEvaluate(t *testing.T) {
	ps, err := LoadPortfoliosFile("../../testdata/portfolios.yaml")
	assert.PassIf(t, err == nil, "error reading portfolios file")
	p := ps[0]
	reader := price.NewPriceReader(&mockprice.Reader{}, &logger.Log{})
	err = p.SetUSDValues(reader, helpers.TodaysDate(), true)
	assert.PassIf(t, err == nil, "error pricing portfolio: %v", err)
	p.Assets.SortByValue()
	assert.Equal(t, 5000.0, p.Assets[0].Value)
	assert.Equal(t, 2500.0, p.Assets[1].Value)
	assert.Equal(t, 100.0, p.Assets[2].Value)
	p.Assets[0].Value = 1000.00
	p.Assets.SortByValue()
	assert.Equal(t, 2500.0, p.Assets[0].Value)
	assert.Equal(t, 1000.0, p.Assets[1].Value)
	assert.Equal(t, 100.0, p.Assets[2].Value)
}

func TestLoadPortfoliosFile(t *testing.T) {
	ps, err := LoadPortfoliosFile("../../testdata/portfolios.yaml")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 3, len(ps))

	assert.Equal(t, 3, len(ps[0].Assets))
	assert.Equal(t, "personal", ps[0].Name)
	i := ps[0].Assets.Find("BTC")
	assert.PassIf(t, i != -1, "missing asset: BTC")
	assert.Equal(t, Asset{
		Symbol: "BTC",
		Amount: 0.5,
		Value:  0.0,
	},
		ps[0].Assets[i])

	i = ps[1].Assets.Find("ETH")
	assert.PassIf(t, i != -1, "missing asset: ETH")
	assert.Equal(t, 2, len(ps[1].Assets))
	assert.Equal(t, "joint", ps[1].Name)
	assert.Equal(t, Asset{
		Symbol: "ETH",
		Amount: 2.5,
		Value:  0.0,
	},
		ps[1].Assets[i])
}
func TestSortAndFilter(t *testing.T) {
	h, err := LoadValuationsFile("../../testdata/valuations.json")
	assert.PassIf(t, err == nil, "error reading JSON file")

	h2 := h.FilterByDate("2022-12-02")
	assert.Equal(t, 2, len(h2))

	h2 = h.FilterByName("personal")
	assert.Equal(t, 7, len(h2))
	h2.SortByDateAndName()
	assert.Equal(t, "2022-12-01", h2[0].Date)
	assert.Equal(t, "2022-12-07", h2[6].Date)

	h2 = h.FilterByDate("2022-12-07").FilterByName("joint")
	assert.Equal(t, 1, len(h2))
	assert.Equal(t, "2022-12-07", h2[0].Date)

	h2 = h.FilterByName("personal", "joint")
	assert.Equal(t, 14, len(h2))
}
