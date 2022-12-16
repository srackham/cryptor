package cli

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/fsx"
	"github.com/srackham/cryptor/internal/helpers"
	"github.com/srackham/cryptor/internal/mockprice"
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
	err = fsx.CopyFile("../../testdata/portfolios.toml", path.Join(tmpdir, "portfolios.toml"))
	if err != nil {
		return
	}
	cmd = fmt.Sprintf("%s -confdir %s", cmd, tmpdir)
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
	today := helpers.DateNowString()
	cli := mockCli()
	out, err := exec(cli, "cryptor valuate")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Contains(t, out, "NAME:      personal\nNOTES:     Personal holdings\n")
	assert.NotContains(t, out, "price request:")
	fmt.Println(out)
	assert.PassIf(t, cli.valuations.FindByNameAndDate("personal", today) != -1, "missing valuation: %v", today)

	cli = mockCli()
	out, err = exec(cli, "cryptor valuate -date "+today+" -refresh -v")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Contains(t, out, "price request: BTC "+today+" 10000.00")
	assert.Contains(t, out, "price request: ETH "+today+" 1000.00")
	assert.Contains(t, out, "price request: USDC "+today+" 1.00")
	fmt.Println(out)

	cli = mockCli()
	date := "2022-06-01"
	out, err = exec(cli, "cryptor valuate -v -date "+date)
	assert.Contains(t, out, "price request: BTC "+date+" 10000.00")
	assert.Contains(t, out, "price request: ETH "+date+" 1000.00")
	assert.Contains(t, out, "price request: USDC "+date+" 1.00")
	assert.PassIf(t, err == nil, "%v", err)
	assert.PassIf(t, cli.valuations.FindByNameAndDate("personal", date) == -1, "past valuation should not be saved: %v", date)
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
