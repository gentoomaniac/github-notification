package main

import (
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"

	"github.com/gentoomaniac/github-notifications/pkg/app"
	gocli "github.com/gentoomaniac/github-notifications/pkg/cli"
	"github.com/gentoomaniac/github-notifications/pkg/config"
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

	GithubNotificationToken string `help:"Github PAT (classic) that has access to notifications" env:"GH_TOKEN"`

	ConfigPath string `help:"path to config file" default:"~/.config/ghn.json" type:"existingfile"`

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

	conf, err := config.FromFile(cli.ConfigPath)
	if err != nil {
		log.Error().Err(err).Msg("failed loading config")
	}

	if cli.GithubNotificationToken != "" {
		conf.NotificationToken = cli.GithubNotificationToken
	}

	a := app.New(conf)

	a.Run()

	ctx.Exit(0)
}
