package main

import (
	"os"

	"github.com/srackham/cryptor/internal/cgprice"
	"github.com/srackham/cryptor/internal/cli"
)

func main() {
	// TODO implement -test option to use testapi instead of cgapi.
	cli := cli.New(&cgprice.Reader{})
	if err := cli.Execute(os.Args); err != nil {
		os.Exit(1)
	}
}
