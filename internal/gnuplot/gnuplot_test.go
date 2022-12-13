package gnuplot

import (
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/portfolio"
)

func TestPlotHistory(t *testing.T) {
	portfolios, err := portfolio.LoadHistoryFile("../../testdata/portfolios.json")
	assert.PassIf(t, err == nil, "error reading JSON file")
	p := &Plotter{}
	p.PlotHistory(portfolios)
}

func TestPlotAllocation(t *testing.T) {
	// TODO
	portfolios, err := portfolio.LoadHistoryFile("../../testdata/portfolios.json")
	assert.PassIf(t, err == nil, "error reading JSON file")
	p := &Plotter{}
	p.PlotAllocation(portfolios[0])
}
