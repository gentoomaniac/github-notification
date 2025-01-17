package main

import (
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"

	gocli "github.com/gentoomaniac/github-notifications/pkg/cli"
	"github.com/gentoomaniac/github-notifications/pkg/gh"
	"github.com/gentoomaniac/github-notifications/pkg/logging"
	"github.com/gentoomaniac/github-notifications/pkg/ui"
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

	github := gh.New(cli.GithubToken)

	notifications, err := github.GetNotifications()
	if err != nil {
		log.Error().Err(err).Msg("failed getting notifications")
	}

	u := ui.Ui{}

	u.Run(notifications)

	ctx.Exit(0)
}
