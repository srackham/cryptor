package portfolio

import (
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/mock"
)

func TestIsValidPortfolioName(t *testing.T) {
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
	ctx := mock.NewContext()
	test := func(format string) {
		valuationsFile := path.Join(ctx.DataDir, "valuations."+format)
		valuations, err := LoadValuations(valuationsFile)
		assert.PassIf(t, err == nil, "error reading valuations file")
		assert.Equal(t, 15, len(valuations))
		assert.Equal(t, "2022-12-02", valuations[0].Date)
		assert.Equal(t, "2022-12-06", valuations[13].Date)
		assert.Equal(t, 3, len(valuations[0].Assets))
		assert.Equal(t, valuations[0].Assets[0], Asset{
			Symbol: "BTC",
			Amount: 0.5,
			Value:  5000.0,
		})
		assert.Equal(t, valuations[2].Assets[2], Asset{
			Symbol: "USDC",
			Amount: 100.0,
			Value:  100.00,
		})
	}
	test("json")
	test("yaml")
}

func TestSaveValuationsFile(t *testing.T) {
	ctx := mock.NewContext()
	tmpdir := mock.MkdirTemp(t)
	test := func(format string) {
		// Read valuations from ../../testdata/data/ directory
		valuationsFile := path.Join(ctx.DataDir, "valuations."+format)
		valuations, err := LoadValuations(valuationsFile)
		assert.PassIf(t, err == nil, "error reading valuations file")
		assert.PassIf(t, len(valuations) == 15, "valuations file: \"%v\": expected 15 portfolios, got %d", valuationsFile, len(valuations))
		// Save valuations to temp directory
		fname := filepath.Join(tmpdir, "valuations."+format)
		err = valuations.SaveValuations(fname)
		assert.PassIf(t, err == nil, "%v", err)
		savedValuations := valuations
		valuations, err = LoadValuations(fname)
		assert.PassIf(t, err == nil, "%v", err)
		assert.PassIf(t, reflect.DeepEqual(savedValuations, valuations),
			"valuations file: \"%v\": expected:\n%v\n\ngot:\n%v", valuationsFile, savedValuations, valuations)
	}
	test("json")
	test("yaml")
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

func TestSortAndFilter(t *testing.T) {
	ctx := mock.NewContext()
	valuationsFile := path.Join(ctx.DataDir, "valuations.yaml")
	valuations, err := LoadValuations(valuationsFile)
	assert.PassIf(t, err == nil, "error reading valuations file")

	filteredValuations := valuations.FilterByDate("2022-12-02")
	assert.Equal(t, 3, len(filteredValuations))

	filteredValuations = valuations.FilterByName("personal")
	assert.Equal(t, 7, len(filteredValuations))
	filteredValuations.Sort()
	assert.Equal(t, "2022-12-01", filteredValuations[0].Date)
	assert.Equal(t, "2022-12-02", filteredValuations[1].Date)
	assert.Equal(t, "10:30:00", filteredValuations[1].Time)
	assert.Equal(t, "2022-12-02", filteredValuations[2].Date)
	assert.Equal(t, "12:30:00", filteredValuations[2].Time)
	assert.Equal(t, "2022-12-07", filteredValuations[6].Date)

	filteredValuations = valuations.FilterByDate("2022-12-07").FilterByName("joint")
	assert.Equal(t, 1, len(filteredValuations))
	assert.Equal(t, "2022-12-07", filteredValuations[0].Date)

	filteredValuations = valuations.FilterByName("personal", "joint")
	assert.Equal(t, 14, len(filteredValuations))
}
