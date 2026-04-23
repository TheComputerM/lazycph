package main

import (
	"context"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/thecomputerm/lazycph/internal/app"
	"github.com/urfave/cli/v3"
)

// version is set at build time via ldflags.
var version = "dev"

func main() {
	cmd := &cli.Command{
		Name:    "lazycph",
		Usage:   "A terminal UI for competitive programming",
		Version: version,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "companion",
				Aliases: []string{"c"},
				Usage:   "Enable Competitive Companion integration",
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "filepath",
				Config: cli.StringConfig{
					TrimSpace: true,
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			filepath := cmd.StringArg("filepath")
			if filepath == "" {
				var err error
				filepath, err = os.Getwd()
				if err != nil {
					return err
				}
			}

			companion := cmd.Bool("companion")
			model, err := app.New(filepath, companion)
			if err != nil {
				return err
			}
			p := tea.NewProgram(model)
			_, err = p.Run()
			return err
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
