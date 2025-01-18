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

	GithubClassicToken     string `help:"Github PAT (classic) that has access to notifications" env:"GH_CLASSIC_TOKEN"`
	GithubFinegrainedToken string `help:"Github PAT (finegrained) that has access to PRs" env:"GH_FINEGRAINED_TOKEN"`

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

	a := app.New(cli.GithubClassicToken, cli.GithubFinegrainedToken)

	a.Run()

	// log.Debug().Str("classic", cli.GithubClassicToken).Str("fine", cli.GithubFinegrainedToken).Msg("tokens")

	// g := gh.New(cli.GithubClassicToken, cli.GithubFinegrainedToken)

	// pr, err := g.GetPr("gentoomaniac", "github-notification", 1)
	// if err != nil {
	// 	log.Error().Err(err).Msg("")
	// }

	// json, _ := json.Marshal(pr)
	// fmt.Println(string(json))

	ctx.Exit(0)
}
