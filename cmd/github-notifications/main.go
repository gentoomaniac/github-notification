package main

import (
	"github.com/alecthomas/kong"

	"github.com/gentoomaniac/github-notifications/pkg/app"
	gocli "github.com/gentoomaniac/github-notifications/pkg/cli"
	"github.com/gentoomaniac/github-notifications/pkg/logging"
)

var (
	version = "unknown"
	commit  = "unknown"
	binName = "unknown"
	builtBy = "unknown"
	date    = "unknown"
)

var cli struct {
	logging.LoggingConfig

	GithubToken string `help:"Github PAT (classic) that has access to notifications" default:"fancy" env:"GH_TOKEN"`

	Version gocli.VersionFlag `short:"V" help:"Display version."`
}

func main() {
	ctx := kong.Parse(&cli, kong.UsageOnError(), kong.Vars{
		"version": version,
		"commit":  commit,
		"binName": binName,
		"builtBy": builtBy,
		"date":    date,
	})
	logging.Setup(&cli.LoggingConfig)

	a := app.New(cli.GithubToken)

	a.Run()

	ctx.Exit(0)
}
