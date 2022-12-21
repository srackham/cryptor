package portfolio

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/srackham/cryptor/internal/fsx"
	"github.com/srackham/cryptor/internal/price"
	"gopkg.in/yaml.v3"
)

type Asset struct {
	Symbol     string  `json:"symbol"`
	Amount     float64 `json:"amount"`
	Value      float64 `json:"value"`      // Value in USD at the time of valuation
	Allocation float64 `json:"allocation"` // Percentage of total portfolio USD value
}

type Assets []Asset

type Portfolio struct {
	Name    string  `json:"name"`
	Notes   string  `json:"notes"`
	Date    string  `json:"date"`    // The valuation date formatted "YYYY-MM-DD"
	Value   float64 `json:"value"`   // Combined assets value in USD
	Cost    string  `json:"cost"`    // Combined assets cost, format = "<amount> <currency>"
	USDCost float64 `json:"usdcost"` // Calculated cost in USD.
	Assets  Assets  `json:"assets"`
}

type Portfolios []Portfolio

// ParseCurrency extracts the amount and currency symbol from a string formatted like "<amount>[ <currency>]".
func ParseCurrency(cstr string) (value float64, currency string, err error) {
	s := strings.ReplaceAll(cstr, "$", "")
	s = strings.ReplaceAll(s, ",", "")
	split := regexp.MustCompile(`\s+`).Split(s, -1)
	switch len(split) {
	case 1:
		currency = "USD"
	case 2:
		currency = strings.ToUpper(split[1])
	default:
		err = fmt.Errorf("invalid currency value: %q", cstr)
		return
	}
	value, err = strconv.ParseFloat(split[0], 64)
	if err != nil {
		err = fmt.Errorf("invalid currency value: %q", cstr)
		return
	}
	return
}

func (assets Assets) SortByValue() {
	// TODO tests
	// Sort assets by descending value.
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].Value > assets[j].Value
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

// SetUSDValues calculates the current USD value of portfolio assets and their total value.
// TODO tests
func (p *Portfolio) SetUSDValues(reader price.PriceReader, date string, force bool) error {
	total := 0.0
	for i, a := range p.Assets {
		rate, err := reader.GetPrice(a.Symbol, date, force)
		if err != nil {
			return err
		}
		val := a.Amount * rate
		p.Assets[i].Value = val
		total += val
	}
	p.Value = total
	return nil
}

// SetAllocations synthesizes asset allocation as a percentage of the total portfolio value.
// TODO tests
func (p *Portfolio) SetAllocations() {
	for i, a := range p.Assets {
		if p.Value != 0.00 {
			p.Assets[i].Allocation = a.Value / p.Value * 100
		}
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

// LoadPortfoliosFile reads portfolios configuration file.
// Returns a Portfolios slice.
func LoadPortfoliosFile(filename string) (Portfolios, error) {
	res := Portfolios{}
	s, err := fsx.ReadFile(filename)
	if err != nil {
		return res, err
	}
	config := []struct {
		Name   string             `yaml:"name"`
		Notes  string             `yaml:"notes"`
		Cost   string             `yaml:"cost"`
		Assets map[string]float64 `yaml:"assets"`
	}{}
	err = yaml.Unmarshal([]byte(s), &config)
	if err != nil {
		return res, err
	}
	// Copy parsed portfolios configuration to Portfolios slice.
	for _, c := range config {
		p := Portfolio{}
		p.Name = c.Name
		p.Notes = c.Notes
		p.Cost = c.Cost
		p.Assets = []Asset{}
		for k, v := range c.Assets {
			asset := Asset{}
			asset.Symbol = strings.ToUpper(k)
			asset.Amount = v
			p.Assets = append(p.Assets, asset)
		}
		res = append(res, p)
	}
	// Check for duplicate portfolio names.
	for i := range res {
		for j := range res {
			if i != j && res[i].Name == res[j].Name {
				return res, fmt.Errorf("duplicate portfolio name: %q", res[j].Name)
			}
		}
	}
	// Synthesise missing portfolio names.
	for i := range res {
		if res[i].Name == "" {
			for j := 1; ; j++ {
				name := fmt.Sprintf("portfolio%d", j)
				if res.FindByName(name) == -1 {
					res[i].Name = name
					break
				}
			}
		}
	}
	return res, err
}

func LoadValuationsFile(valuationsFile string) (Portfolios, error) {
	res := Portfolios{}
	s, err := fsx.ReadFile(valuationsFile)
	if err == nil {
		err = json.Unmarshal([]byte(s), &res)
	}
	return res, err
}

func (ps Portfolios) SaveValuationsFile(valuationsFile string) error {
	ps.SortByDateAndName()
	data, err := json.MarshalIndent(ps, "", "  ")
	if err == nil {
		err = fsx.WriteFile(valuationsFile, string(data))
	}
	return err
}

// Aggregate returns a new portfolio that combines the receiver portfolios.
// Portfolio Notes field is assigned the list of combined portfolios.
// Aggregated costs are valid only if all portfolios are costed.
// TODO tests
func (ps Portfolios) Aggregate(name string, date string) Portfolio {
	res := Portfolio{
		Name:   name,
		Date:   date,
		Assets: Assets{},
	}
	var notes string
	missingCosts := false
	for _, p := range ps {
		notes += fmt.Sprintf("%s, ", p.Name)
		res.Value += p.Value
		if p.USDCost == 0 {
			missingCosts = true
		}
		res.USDCost += p.USDCost
		for _, a := range p.Assets {
			i := res.Assets.Find(a.Symbol)
			if i == -1 {
				res.Assets = append(res.Assets, Asset{Symbol: a.Symbol, Amount: a.Amount, Value: a.Value})
			} else {
				res.Assets[i].Amount += a.Amount
				res.Assets[i].Value += a.Value
			}
		}
	}
	res.Notes = strings.TrimSuffix(notes, ", ")
	res.SetAllocations()
	res.Assets.SortByValue()
	if missingCosts {
		res.USDCost = 0.00
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

func (ps Portfolios) SortByDateAndName() {
	// Sort documents by ascending date and name.
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].Date != ps[j].Date {
			return ps[i].Date < ps[j].Date
		}
		return ps[i].Name < ps[j].Name
	})
}
