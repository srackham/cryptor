package portfolio

import (
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/srackham/cryptor/internal/cache"
	"github.com/srackham/cryptor/internal/fsx"
	"github.com/srackham/cryptor/internal/price"
	"github.com/srackham/cryptor/internal/set"
)

type Asset struct {
	Symbol     string
	Amount     float64
	USD        float64 // Value in USD at the time of valuation
	Allocation float64 // Percentage of total portfolio USD value
}

type Assets []Asset

type Portfolio struct {
	Name   string
	Notes  string
	Date   string  // The valuation date formatted "YYYY-MM-DD"
	Time   string  // The valuation time formatted "HH:MM:SS"
	USD    float64 // Combined assets value in USD
	Assets Assets
}

type Portfolios []Portfolio

func (assets Assets) SortByValue() {
	// TODO tests
	// Sort assets by descending USD value.
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].USD > assets[j].USD
	})
}

// Find searches assets slice for an asset with matching `symbol`.
// If found it return the asset index else returns -1.
func (assets Assets) Find(symbol string) int {
	// TODO tests
	for i := range assets {
		if assets[i].Symbol == symbol {
			return i
		}
	}
	return -1
}

// GetPrices gets the prices of all portfolio crypto assets.
// TODO tests
func (ps *Portfolios) GetPrices(reader price.PriceReader, date string, refresh bool) (cache.Rates, error) {
	ss := set.New[string]()
	for _, p := range *ps {
		for _, a := range p.Assets {
			ss.Add(a.Symbol)
		}
	}
	symbols := ss.Values()
	sort.Strings(symbols)
	prices := make(cache.Rates)
	for _, sym := range symbols {
		price, err := reader.GetPrice(sym, date, refresh)
		if err != nil {
			return prices, err
		}
		prices[sym] = price
	}
	return prices, nil
}

// SetUSDValues calculates the current USD value of portfolio assets and their total value.
// TODO tests
func (p *Portfolio) SetUSDValues(prices cache.Rates) {
	total := 0.0
	for i, a := range p.Assets {
		rate := prices[a.Symbol]
		val := a.Amount * rate
		p.Assets[i].USD = val
		total += val
	}
	p.USD = total
}

// SetAllocations synthesizes asset allocation as a percentage of the total portfolio USD value.
// TODO tests
func (p *Portfolio) SetAllocations() {
	for i, a := range p.Assets {
		if p.USD != 0.00 {
			p.Assets[i].Allocation = a.USD / p.USD * 100
		}
	}
}

// SetTimeStamp timestamps the portfolio.
// TODO tests
func (p *Portfolio) SetTimeStamp(date string, refresh bool) {
	if refresh {
		now := time.Now()
		p.Date = now.Format("2006-01-02")
		p.Time = now.Format("15:04:05")
	} else {
		p.Date = date
		p.Time = ""
	}
}

// See [How to deep copy a struct in Go](https://www.educative.io/answers/how-to-deep-copy-a-struct-in-go)
// TODO tests
func (p Portfolio) DeepCopy() Portfolio {
	res := p
	res.Assets = nil
	res.Assets = append(res.Assets, p.Assets...)
	return res
}

// LoadPortfoliosFile reads TOML portfolios file.
// Returns a Portfolios slice.
func LoadPortfoliosFile(filename string) (Portfolios, error) {
	res := Portfolios{}
	s, err := fsx.ReadFile(filename)
	if err != nil {
		return res, err
	}
	raw := struct {
		Portfolios []struct {
			Name   string `toml:"name"`
			Notes  string `toml:"notes"`
			Assets []struct {
				Symbol string  `toml:"symbol"`
				Amount float64 `toml:"amount"`
			} `toml:"assets"`
		} `toml:"portfolios"`
	}{}
	_, err = toml.Decode(s, &raw)
	if err != nil {
		return res, err
	}
	// Copy parsed portfolios configuration to Portfolios slice.
	for _, c := range raw.Portfolios {
		p := Portfolio{}
		p.Name = c.Name
		p.Notes = c.Notes
		p.Assets = []Asset{}
		for _, a := range c.Assets {
			asset := Asset{}
			asset.Symbol = a.Symbol
			asset.Amount = a.Amount
			p.Assets = append(p.Assets, asset)
		}
		res = append(res, p)
	}
	return res, err
}

func LoadValuationsFile(valuationsFile string) (Portfolios, error) {
	res := Portfolios{}
	if !fsx.FileExists(valuationsFile) {
		return res, nil
	}
	s, err := fsx.ReadFile(valuationsFile)
	if err == nil {
		err = json.Unmarshal([]byte(s), &res)
	}
	return res, err
}

func (ps Portfolios) SaveValuationsFile(valuationsFile string) error {
	ps.SortByDate()
	data, err := json.MarshalIndent(ps, "", "  ")
	if err == nil {
		err = fsx.WriteFile(valuationsFile, string(data))
	}
	return err
}

// TODO tests
func (ps *Portfolios) UpdateValuations(p Portfolio) {
	i := ps.FindByNameAndDate(p.Name, p.Date)
	if i == -1 {
		*ps = append(*ps, p)
	} else {
		(*ps)[i] = p
	}
}

// Aggregate returns a new portfolio that combines assets from one or more portfolios.
// Returns an aggregated portfolio with `name` and `notes`.
// Portfolio Date and Time fields are left unfilled.
// Asset.Amount and Asset.USD asset fields are aggregated (summed).
// TODO tests
func (ps Portfolios) Aggregate(name, notes string) Portfolio {
	res := Portfolio{
		Name:   name,
		Notes:  notes,
		Assets: Assets{},
	}
	for _, p := range ps {
		for _, a := range p.Assets {
			i := res.Assets.Find(a.Symbol)
			if i == -1 {
				res.Assets = append(res.Assets, Asset{Symbol: a.Symbol, Amount: a.Amount, USD: a.USD})
			} else {
				res.Assets[i].Amount += a.Amount
				res.Assets[i].USD += a.USD
			}
		}
	}
	return res
}

// FilterByDate returns a list of portfolios dated `date`.
func (ps Portfolios) FilterByDate(date string) Portfolios {
	res := []Portfolio{}
	for _, p := range ps {
		if p.Date == date {
			res = append(res, p)
		}
	}
	return res
}

// FindByNameAndDate searches portfolios slice for a portfolio whose name and date matches portfolio `p`.
// If found it return the portfolio index else returns -1.
// TODO tests
func (ps Portfolios) FindByNameAndDate(name string, date string) int {
	for i := range ps {
		if name == ps[i].Name && date == ps[i].Date {
			return i
		}
	}
	return -1
}

// Find searches portfolios slice for a portfolio whose name matches `name`.
// If found return the portfolio index else return -1.
// TODO tests
func (ps Portfolios) FindByName(name string) int {
	for i := range ps {
		if name == ps[i].Name {
			return i
		}
	}
	return -1
}

// FilterByName returns a list of named portfolios.
func (ps Portfolios) FilterByName(names ...string) Portfolios {
	res := []Portfolio{}
	for _, p := range ps {
		for _, name := range names {
			if p.Name == name {
				res = append(res, p)
			}
		}
	}
	return res
}

func (ps Portfolios) SortByDate() {
	// Sort documents by ascending date.
	sort.Slice(ps, func(i, j int) bool {
		return strings.Compare(ps[i].Date, ps[j].Date) == -1
	})
}
