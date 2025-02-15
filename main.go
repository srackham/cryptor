package main

import (
	"net/http"
	"os"
	"path"
	"time"

	"github.com/srackham/cryptor/internal/cli"
	. "github.com/srackham/cryptor/internal/global"
	"github.com/srackham/cryptor/internal/helpers"
)

func main() {
	ctx := Context{
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		DataDir:   path.Join(helpers.GetDataDir(), "cryptor"),
		CacheDir:  path.Join(helpers.GetCacheDir(), "cryptor"),
		ConfigDir: path.Join(helpers.GetConfigDir(), "cryptor"),
		Now:       func() time.Time { return time.Now() },
		HttpGet:   http.Get,
	}
	cli := cli.New(&ctx)
	if err := cli.Execute(os.Args...); err != nil {
		os.Exit(1)
	}
}
