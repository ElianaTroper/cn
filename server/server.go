// Package main reads the command line arguments and calls the proper
// sub-package.
package main


import (
	// Built ins
	"fmt"
	"os"
	// Third parties
	"github.com/urfave/cli/v2"
	// This repo
	"github.com/ElianaTroper/cn/server/app"
	"github.com/ElianaTroper/cn/server/config"
	"github.com/ElianaTroper/cn/server/ipfs"
)


// TODO: Allow launching the daemon in either mode (default to mirror)
// TODO: Add `app enable`
// TODO: Add start
// TODO: Add stop


func main() {
	app := &cli.App{
		Name:  "CN Server App",
		Usage: "Runs and configures the CN server",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Eliana Troper",
				Email: "eliana@troper.report",
			},
		},
		Copyright: "(c) 2023 Eliana Troper",
		Version:   "t.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "conf",
				Aliases: []string{"c"},
				Value:   "./config/default.json",
				Usage:   "path to the cn config to use",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "ipfs",
				Aliases: []string{"i"},
				Usage:   "Manages ipfs interactions",
				Subcommands: []*cli.Command{
					{
						Name:      "init",
						Aliases:   []string{"i"},
						Usage:     "sets configs, using ./default.json",
						ArgsUsage: "[FILE (optional, default: ./default.json)]",
						Action: func(cCtx *cli.Context) error {
							if cCtx.Args().Len() > 0 {
								return fmt.Errorf("no args are allowed when using `cn ipfs init`")
							}
							conf, err := config.Load(cCtx.String("conf"))
							if err != nil {
								return err
							}
							return ipfs.Init(conf)
						},
					},
				},
			},
			{
				Name:    "app",
				Aliases: []string{"a"},
				Usage:   "Manages app deployments",
				Subcommands: []*cli.Command{
					{
						Name:      "deploy",
						Aliases:   []string{"d"},
						Usage:     "deploys an app (currently only available for a root node)",
						ArgsUsage: "[APP (required), APP2 (optional)...]",
						Action: func(cCtx *cli.Context) error {
							if cCtx.Args().Len() == 0 {
								return fmt.Errorf("at least one argument required when using `cn app deploy`")
							}
							conf, err := config.Load(cCtx.String("conf"))
							if err != nil {
								return err
							}
							for i := 0; i < cCtx.Args().Len(); i++ {
								err = app.Deploy(cCtx.Args().Get(i), conf)
								if err != nil {
									return err
								}
							}
							return nil
						},
					},
				},
			},
		},
		EnableBashCompletion: true,
		Suggest:              true,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.exit(1)
	}
}
