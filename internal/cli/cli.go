package cli

import (
	"context"
	"os"

	"github.com/gabefiori/gsp/internal/app"
	"github.com/gabefiori/gsp/internal/config"
	"github.com/urfave/cli/v3"
)

// flags
var (
	flagConfig = &cli.StringFlag{
		Name:      "config",
		Aliases:   []string{"c"},
		Usage:     "Load configuration from the specified `file`",
		Value:     "~/.config/gsp/config.json",
		TakesFile: true,
	}

	flagList = &cli.BoolFlag{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "Print entries to stdout",
		Value:   false,
	}

	flagMeasure = &cli.BoolFlag{
		Name:    "measure",
		Aliases: []string{"m"},
		Usage:   "Measure performance (time taken and number of entries processed)",
		Value:   false,
	}

	flagSelector = &cli.StringFlag{
		Name:    "selector",
		Aliases: []string{"sl"},
		Usage:   "Selector for displaying entries (available options: 'fzf', 'fzy', 'sk')",
	}

	flagSort = &cli.StringFlag{
		Name:    "sort",
		Aliases: []string{"s"},
		Usage:   "Specify the sort order for displaying entries (available options: 'asc', 'desc', 'nosort')",
		Value: "nosort",
	}

	flagUnique = &cli.BoolFlag{
		Name:    "unique",
		Aliases: []string{"u"},
		Usage:   "Display only unique entries",
		Value:   false,
	}

	flagExpand = &cli.BoolFlag{
		Name:    "expand-output",
		Aliases: []string{"eo"},
		Usage:   "Expand selection output",
		Value:   true,
	}
)

// Run initializes and executes the command-line interface (CLI) application.
func Run(version string) error {
	cmd := cli.Command{
		Name:    "gsp",
		Usage:   "Select projects.",
		Version: version,
		Action:  action,
		Flags: []cli.Flag{
			flagConfig,
			flagList,
			flagMeasure,
			flagSelector,
			flagSort,
			flagUnique,
			flagExpand,
		},
	}

	return cmd.Run(context.Background(), os.Args)
}

func action(ctx context.Context, c *cli.Command) error {
	params := &config.LoadParams{
		Path:     c.String(flagConfig.Name),
		Measure:  c.Bool(flagMeasure.Name),
		List:     c.Bool(flagList.Name),
		Selector: c.String(flagSelector.Name),
		Sort:     c.String(flagSort.Name),
	}

	params.Unique = optionalBoolFlag(flagUnique, c)
	params.ExpandOutput = optionalBoolFlag(flagExpand, c)

	cfg, err := config.Load(params)
	if err != nil {
		return err
	}

	a, err := app.New(cfg)
	if err != nil {
		return err
	}

	return a.Run()
}

func optionalBoolFlag(f *cli.BoolFlag, c *cli.Command) int8 {
	if !c.IsSet(f.Name) {
		return 0
	}

	if c.Bool(flagExpand.Name) {
		return 1
	}

	return -1
}
