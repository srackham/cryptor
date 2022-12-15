package cli

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/helpers"
	"github.com/srackham/cryptor/internal/mockprice"
	"github.com/srackham/cryptor/internal/portfolio"
)

func mockCli() *cli {
	return New(&mockprice.Reader{})
}

func TestParseArgs(t *testing.T) {
	var cli cli
	var err error
	parse := func(cmd string) {
		args := strings.Split(cmd, " ")
		cli = *mockCli()
		err = cli.parseArgs(args)
	}

	parse("cryptor")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, "help", cli.command)
	parse("cryptor illegal-command")
	assert.Equal(t, `illegal command: "illegal-command"`, err.Error())
}

func exec(cli *cli, cmd string) (out string, err error) {
	tmpdir, err := os.MkdirTemp("", "cryptor")
	if err != nil {
		return
	}
	cmd = fmt.Sprintf("%s -conf ../../testdata/portfolios.toml -confdir %s", cmd, tmpdir)
	args := strings.Split(cmd, " ")
	cli.log.Out = make(chan string, 100)
	err = cli.Execute(args)
	close(cli.log.Out)
	if err != nil {
		return
	}
	for line := range cli.log.Out {
		out += line + "\n"
	}
	out = strings.Replace(out, `\`, `/`, -1) // Normalize MS Windows path separators.
	return
}

func TestEvaluateCmd(t *testing.T) {
	out, err := exec(mockCli(), "cryptor valuate")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Contains(t, out, "Name:        personal\nDescription: Personal holdings\n")
	assert.NotContains(t, out, "price request:")
	fmt.Println(out)

	today := helpers.DateNowString()
	out, err = exec(mockCli(), "cryptor valuate -date "+today+" -refresh -v")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Contains(t, out, "price request: BTC "+today+" 10000.00")
	assert.Contains(t, out, "price request: ETH "+today+" 1000.00")
	assert.Contains(t, out, "price request: USDC "+today+" 1.00")
	fmt.Println(out)
}

func TestHelpCmd(t *testing.T) {
	out, err := exec(mockCli(), "cryptor help")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Contains(t, out, "Cryptor reports crypto currency portfolio statistics")

	out, err = exec(mockCli(), "cryptor --help")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Contains(t, out, "Cryptor reports crypto currency portfolio statistics")

	out, err = exec(mockCli(), "cryptor -h")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Contains(t, out, "Cryptor reports crypto currency portfolio statistics")

}

func TestParseConfig(t *testing.T) {
	cli := mockCli()
	cli.configFile = "../../testdata/portfolios.toml"
	println(os.Getwd())
	err := cli.loadPortfolios()
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 2, len(cli.portfolios))

	assert.Equal(t, 3, len(cli.portfolios[0].Assets))
	assert.Equal(t, "personal", cli.portfolios[0].Name)
	assert.Equal(t, cli.portfolios[0].Assets[0], portfolio.Asset{
		Symbol:      "BTC",
		Amount:      0.5,
		USD:         0.0,
		Description: "Cold storage",
	})

	assert.Equal(t, 2, len(cli.portfolios[1].Assets))
	assert.Equal(t, "joint", cli.portfolios[1].Name)
	assert.Equal(t, cli.portfolios[1].Assets[1], portfolio.Asset{
		Symbol:      "ETH",
		Amount:      2.5,
		USD:         0.0,
		Description: "",
	})
}
