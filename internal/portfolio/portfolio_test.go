package portfolio

import (
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/binance"
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
	tests := []struct {
		input   string
		wantAmt float64
		wantSym string
		wantErr bool
		errMsg  string
	}{
		{input: "$5,000.00 NZD", wantAmt: 5000.00, wantSym: "NZD", wantErr: false},
		{input: "1000aud", wantAmt: 1000.0, wantSym: "AUD", wantErr: false},
		{input: ".5", wantAmt: 0.5, wantSym: "USD", wantErr: false},
		{input: "123.45EUR", wantAmt: 123.45, wantSym: "EUR", wantErr: false},
		{input: "123.45 eur", wantAmt: 123.45, wantSym: "EUR", wantErr: false},
		{input: "12345", wantAmt: 12345.0, wantSym: "USD", wantErr: false},
		{input: "$1,234.56", wantAmt: 1234.56, wantSym: "USD", wantErr: false},
		{input: "abc", wantErr: true, errMsg: "invalid currency value: \"abc\""},
		{input: "", wantErr: true, errMsg: "invalid currency value: \"\""},
		{input: "   ", wantErr: true, errMsg: "invalid currency value: \"   \""},
		{input: "123.45abc123", wantErr: true, errMsg: "invalid currency value: \"123.45abc123\""},
		{input: "123.45.67USD", wantErr: true, errMsg: "invalid currency value: \"123.45.67USD\""},
	}
	for _, tt := range tests {
		gotVal, gotSym, gotErr := ParseCurrency(tt.input)
		if (gotErr != nil) != tt.wantErr {
			t.Errorf("ParseCurrency(%q) error = %v, wantErr %v", tt.input, gotErr, tt.wantErr)
			continue
		}
		if tt.wantErr && gotErr != nil && gotErr.Error() != tt.errMsg {
			t.Errorf("ParseCurrency(%q) error message = %v, want message %v", tt.input, gotErr.Error(), tt.errMsg)
			continue
		}
		if gotVal != tt.wantAmt {
			t.Errorf("ParseCurrency(%q) gotVal = %v, want %v", tt.input, gotVal, tt.wantAmt)
		}
		if gotSym != tt.wantSym {
			t.Errorf("ParseCurrency(%q) gotSym = %v, want %v", tt.input, gotSym, tt.wantSym)
		}
	}
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

func TestAssets_Sort(t *testing.T) {
	tests := []struct {
		name     string
		assets   Assets
		expected Assets
	}{
		{
			name: "Sort assets by descending value",
			assets: Assets{
				{Symbol: "BTC", Value: 100},
				{Symbol: "ETH", Value: 200},
				{Symbol: "XRP", Value: 50},
			},
			expected: Assets{
				{Symbol: "ETH", Value: 200},
				{Symbol: "BTC", Value: 100},
				{Symbol: "XRP", Value: 50},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assets.Sort()
			if !reflect.DeepEqual(tt.assets, tt.expected) {
				t.Errorf("Assets.Sort() = %v, want %v", tt.assets, tt.expected)
			}
		})
	}
}

func TestAssets_Find(t *testing.T) {
	assets := Assets{
		{Symbol: "BTC"},
		{Symbol: "ETH"},
		{Symbol: "XRP"},
	}

	tests := []struct {
		name     string
		symbol   string
		expected int
	}{
		{"Find existing asset", "ETH", 1},
		{"Find non-existing asset", "LTC", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := assets.Find(tt.symbol); got != tt.expected {
				t.Errorf("Assets.Find() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPortfolio_SetUSDValues(t *testing.T) {
	ctx := mock.NewContext()
	mockReader := binance.NewPriceReader(&ctx)
	p := &Portfolio{
		Assets: Assets{
			{Symbol: "BTC", Amount: 2},
			{Symbol: "ETH", Amount: 10},
		},
	}
	err := p.SetUSDValues(&mockReader)
	if err != nil {
		t.Fatalf("SetUSDValues() error = %v", err)
	}
	expectedTotal := 210_000.0
	if p.Value != expectedTotal {
		t.Errorf("SetUSDValues() total = %v, want %v", p.Value, expectedTotal)
	}
	expectedAssets := Assets{
		{Symbol: "BTC", Amount: 2, Price: 100_000, Value: 200_000},
		{Symbol: "ETH", Amount: 10, Price: 1000, Value: 10_000},
	}
	for i, asset := range p.Assets {
		if asset.Value != expectedAssets[i].Value || asset.Price != expectedAssets[i].Price {
			t.Errorf("SetUSDValues() asset %s = %+v, want %+v", asset.Symbol, asset, expectedAssets[i])
		}
	}
}

func TestPortfolio_SetAllocations(t *testing.T) {
	p := &Portfolio{
		Value: 100000,
		Assets: Assets{
			{Symbol: "BTC", Value: 60000},
			{Symbol: "ETH", Value: 40000},
		},
	}
	p.SetAllocations()
	expected := Assets{
		{Symbol: "BTC", Value: 60000, Allocation: 60},
		{Symbol: "ETH", Value: 40000, Allocation: 40},
	}
	for i, asset := range p.Assets {
		if asset.Allocation != expected[i].Allocation {
			t.Errorf("SetAllocations() asset %s allocation = %v, want %v", asset.Symbol, asset.Allocation, expected[i].Allocation)
		}
	}
}

func TestPortfolio_DeepCopy(t *testing.T) {
	original := Portfolio{
		Name: "Test Portfolio",
		Assets: Assets{
			{Symbol: "BTC", Amount: 1},
			{Symbol: "ETH", Amount: 10},
		},
	}
	copy := original.DeepCopy()
	if !reflect.DeepEqual(original, copy) {
		t.Errorf("DeepCopy() = %v, want %v", copy, original)
	}
	// Modify the copy to ensure it doesn't affect the original
	copy.Name = "Modified Portfolio"
	copy.Assets[0].Amount = 2
	if reflect.DeepEqual(original, copy) {
		t.Errorf("DeepCopy() did not create a separate copy")
	}
}

func TestPortfolios_Aggregate(t *testing.T) {
	portfolios := Portfolios{
		{
			Name:  "Portfolio 1",
			Value: 100000,
			Cost:  90000,
			Assets: Assets{
				{Symbol: "BTC", Amount: 1, Value: 60000},
				{Symbol: "ETH", Amount: 20, Value: 40000},
			},
		},
		{
			Name:  "Portfolio 2",
			Value: 50000,
			Cost:  45000,
			Assets: Assets{
				{Symbol: "BTC", Amount: 0.5, Value: 30000},
				{Symbol: "XRP", Amount: 1000, Value: 20000},
			},
		},
	}
	aggregated := portfolios.Aggregate("Aggregated Portfolio")
	expectedValue := 150000.0
	if aggregated.Value != expectedValue {
		t.Errorf("Aggregate() value = %v, want %v", aggregated.Value, expectedValue)
	}
	expectedCost := 135000.0
	if aggregated.Cost != expectedCost {
		t.Errorf("Aggregate() cost = %v, want %v", aggregated.Cost, expectedCost)
	}
	expectedAssets := Assets{
		{Symbol: "BTC", Amount: 1.5, Value: 90000},
		{Symbol: "ETH", Amount: 20, Value: 40000},
		{Symbol: "XRP", Amount: 1000, Value: 20000},
	}
	if len(aggregated.Assets) != len(expectedAssets) {
		t.Errorf("Aggregate() asset count = %d, want %d", len(aggregated.Assets), len(expectedAssets))
	}
	for _, expectedAsset := range expectedAssets {
		found := false
		for _, asset := range aggregated.Assets {
			if asset.Symbol == expectedAsset.Symbol && asset.Amount == expectedAsset.Amount && asset.Value == expectedAsset.Value {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Aggregate() missing or incorrect asset: %+v", expectedAsset)
		}
	}
}

func TestPortfolios_FindByNameAndDate(t *testing.T) {
	portfolios := Portfolios{
		{Name: "Portfolio 1", Date: "2025-03-15"},
		{Name: "Portfolio 2", Date: "2025-03-15"},
		{Name: "Portfolio 1", Date: "2025-03-16"},
	}
	tests := []struct {
		name     string
		pName    string
		date     string
		expected int
	}{
		{"Find existing portfolio", "Portfolio 1", "2025-03-15", 0},
		{"Find non-existing portfolio", "Portfolio 3", "2025-03-15", -1},
		{"Find existing portfolio with different date", "Portfolio 1", "2025-03-16", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := portfolios.FindByNameAndDate(tt.pName, tt.date); got != tt.expected {
				t.Errorf("Portfolios.FindByNameAndDate() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPortfolios_FindByName(t *testing.T) {
	portfolios := Portfolios{
		{Name: "Portfolio 1"},
		{Name: "Portfolio 2"},
		{Name: "Portfolio 3"},
	}
	tests := []struct {
		name     string
		pName    string
		expected int
	}{
		{"Find existing portfolio", "Portfolio 2", 1},
		{"Find non-existing portfolio", "Portfolio 4", -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := portfolios.FindByName(tt.pName); got != tt.expected {
				t.Errorf("Portfolios.FindByName() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPortfolios_SetAssetPrice(t *testing.T) {
	portfolios := Portfolios{
		{
			Name: "Portfolio 1",
			Assets: Assets{
				{Symbol: "BTC", Price: 50000},
				{Symbol: "ETH", Price: 3000},
			},
		},
		{
			Name: "Portfolio 2",
			Assets: Assets{
				{Symbol: "BTC", Price: 51000},
				{Symbol: "XRP", Price: 1},
			},
		},
	}
	tests := []struct {
		name      string
		assetName string
		price     float64
		wantErr   bool
	}{
		{"Set existing asset price", "BTC", 55000, false},
		{"Set non-existing asset price", "LTC", 200, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := portfolios.SetAssetPrice(tt.assetName, tt.price)
			if (err != nil) != tt.wantErr {
				t.Errorf("Portfolios.SetAssetPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				for _, p := range portfolios {
					i := p.Assets.Find(tt.assetName)
					if i != -1 && p.Assets[i].Price != tt.price {
						t.Errorf("Portfolios.SetAssetPrice() did not set price correctly for %s in %s", tt.assetName, p.Name)
					}
				}
			}
		})
	}
}

func TestPortfolios_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ps      Portfolios
		nodups  bool
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid portfolios",
			ps: Portfolios{
				{Name: "Portfolio1", Assets: Assets{{Symbol: "BTC"}, {Symbol: "ETH"}}},
				{Name: "Portfolio2", Assets: Assets{{Symbol: "XRP"}}},
			},
			nodups:  true,
			wantErr: false,
		},
		{
			name: "Invalid portfolio name",
			ps: Portfolios{
				{Name: "Invalid Name", Assets: Assets{{Symbol: "BTC"}}},
			},
			nodups:  true,
			wantErr: true,
			errMsg:  "invalid portfolio name: \"Invalid Name\"",
		},
		{
			name: "Invalid asset name",
			ps: Portfolios{
				{Name: "Portfolio1", Assets: Assets{{Symbol: "Invalid Asset"}}},
			},
			nodups:  true,
			wantErr: true,
			errMsg:  "invalid portfolio asset name: \"Invalid Asset\"",
		},
		{
			name: "Duplicate asset name",
			ps: Portfolios{
				{Name: "Portfolio1", Assets: Assets{{Symbol: "BTC"}, {Symbol: "BTC"}}},
			},
			nodups:  true,
			wantErr: true,
			errMsg:  "duplicate asset name: \"BTC\"",
		},
		{
			name: "Duplicate portfolio name with nodups true",
			ps: Portfolios{
				{Name: "Portfolio1", Assets: Assets{{Symbol: "BTC"}}},
				{Name: "Portfolio1", Assets: Assets{{Symbol: "ETH"}}},
			},
			nodups:  true,
			wantErr: true,
			errMsg:  "duplicate portfolio name: \"Portfolio1\"",
		},
		{
			name: "Duplicate portfolio name with nodups false",
			ps: Portfolios{
				{Name: "Portfolio1", Assets: Assets{{Symbol: "BTC"}}},
				{Name: "Portfolio1", Assets: Assets{{Symbol: "ETH"}}},
			},
			nodups:  false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ps.Validate(tt.nodups)
			if (err != nil) != tt.wantErr {
				t.Errorf("Portfolios.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Portfolios.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
