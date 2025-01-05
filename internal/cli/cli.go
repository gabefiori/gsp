package cli

import (
	"os"

	"github.com/gabefiori/gsp/internal/app"
	"github.com/gabefiori/gsp/internal/config"
	"github.com/urfave/cli/v2"
)

// Run initializes and executes the command-line interface (CLI) application.
func Run(version string) error {
	var (
		path         string
		selector     string
		sort         string
		expandOutput bool
		measure      bool
		unique       bool
		list         bool
	)

	cliApp := &cli.App{
		Name:        "Select Projects",
		HelpName:    "gsp",
		Usage:       "Select projects",
		Description: "A simple tool for quickly selecting projects.",
		Version:     version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Load configuration from `file`",
				Value:       "~/.config/gsp/config.json",
				TakesFile:   true,
				Destination: &path,
			},
			&cli.StringFlag{
				Name:        "selector",
				Aliases:     []string{"sl"},
				Usage:       "Selector for displaying projects (available options: 'fzf', 'fzy')",
				Value:       "fzf",
				Destination: &selector,
			},
			&cli.StringFlag{
				Name:        "sort",
				Aliases:     []string{"s"},
				Usage:       "Specify the sort order (available options: 'asc', 'desc')",
				Value:       "",
				Destination: &sort,
			},
			&cli.BoolFlag{
				Name:        "unique",
				Aliases:     []string{"u"},
				Usage:       "Display only unique projects",
				Value:       false,
				Destination: &unique,
			},
			&cli.BoolFlag{
				Name:        "expand-output",
				Aliases:     []string{"eo"},
				Usage:       "Expand the output",
				Value:       true,
				Destination: &expandOutput,
			},
			&cli.BoolFlag{
				Name:        "list",
				Aliases:     []string{"l"},
				Usage:       "List projects to stdout",
				Value:       false,
				Destination: &list,
			},
			&cli.BoolFlag{
				Name:        "measure",
				Aliases:     []string{"m"},
				Usage:       "Measure performance (time taken and number of items processed)",
				Value:       false,
				Destination: &measure,
			},
		},

		Action: func(ctx *cli.Context) error {
			params := &config.LoadParams{
				Path:    path,
				Measure: measure,
				List:    list,
			}

			if ctx.IsSet("expand-output") {
				params.ExpandOutput = &expandOutput
			}

			if ctx.IsSet("selector") {
				params.Selector = selector
			}

			if ctx.IsSet("sort") {
				params.Sort = sort
			}

			if ctx.IsSet("unique") {
				params.Unique = &unique
			}

			cfg, err := config.Load(params)
			if err != nil {
				return err
			}

			a, err := app.New(cfg)
			if err != nil {
				return err
			}

			return a.Run()
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		return err
	}

	return nil
}
