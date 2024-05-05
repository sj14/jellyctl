package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sj14/jellyctl/controller"
	"github.com/sj14/jellyfin-go/api"
	"github.com/urfave/cli/v2"
)

var (
	// will be replaced during the build process
	version = "undefined"
	// commit  = "undefined"
	// date    = "undefined"
)

func main() {
	app := &cli.App{
		Name:                 "jellyctl",
		Usage:                "Manage Jellyfin from the CLI",
		Version:              version,
		Suggest:              true,
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "url",
				Value:   "http://127.0.0.1:8096",
				EnvVars: []string{"JELLYCTL_URL"},
				Usage:   "URL of the Jellyfin server",
			},
			&cli.StringFlag{
				Name:    "token",
				Value:   " ", // https://github.com/urfave/cli/issues/1902
				EnvVars: []string{"JELLYCTL_TOKEN"},
				Usage:   "API token",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "docs",
				Hidden: true,
				Action: func(ctx *cli.Context) error {
					s, err := ctx.App.ToMarkdown()
					if err != nil {
						return err
					}
					return os.WriteFile("DOCS.md", []byte(s), os.ModePerm)
				},
			},
			{
				Name:  "activity",
				Usage: "List activities",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "start",
						Usage: "start at the given index",
					},
					&cli.IntFlag{
						Name:  "limit",
						Usage: "limit the output logs",
						Value: 15,
					},
					&cli.TimestampFlag{
						Name:     "after",
						Usage:    "only logs after the given time",
						Layout:   time.DateTime,
						Timezone: time.Local,
						Value:    cli.NewTimestamp(time.Time{}),
					},
					jsonFlag,
				},
				Action: func(ctx *cli.Context) error {
					return Exec(ctx, func(ctrl *controller.Controller) error {
						return ctrl.GetLogEntries(
							int32(ctx.Int("start")),
							int32(ctx.Int("limit")),
							*ctx.Timestamp("after"),
							ctx.Bool("json"),
						)
					})
				},
			},
			{
				Name:  "system",
				Usage: "Manage the system",
				Subcommands: []*cli.Command{
					{
						Name:  "shutdown",
						Usage: "Stop the Jellyfin process",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.SystemShutdown()
							})
						},
					},
					{
						Name:  "restart",
						Usage: "Restart the Jellyfin process",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.SystemRestart()
							})
						},
					},
					{
						Name:  "info",
						Usage: "Shows system information",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "public",
								Usage: "show public info which won't need a token",
							},
						},
						Action: func(ctx *cli.Context) error {
							if ctx.Bool("public") {
								return Exec(ctx, func(ctrl *controller.Controller) error {
									return ctrl.GetPublicSystemInfo()
								})
							} else {
								return Exec(ctx, func(ctrl *controller.Controller) error {
									return ctrl.GetSystemInfo()
								})
							}
						},
					},
					{
						Name:  "backup",
						Usage: "Export some data (EXPERIMENTAL)",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.SystemBackup()
							})
						},
					},
					{
						Name:      "restore",
						Usage:     "Import played and favourite information (based on the user name not user ID!) (EXPERIMENTAL)",
						Args:      true,
						ArgsUsage: " <PATH>",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "unplay",
								Usage: "mark media as unplayed when this is the backup state",
							},
							&cli.BoolFlag{
								Name:  "unfav",
								Usage: "unfavorite media when this is the backup state",
							},
						},
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.SystemRestore(
									ctx.Args().Get(0),
									ctx.Bool("unplay"),
									ctx.Bool("unfav"),
								)
							})
						},
					},
				},
			},
			{
				Name:  "user",
				Usage: "Manage users",
				Subcommands: []*cli.Command{
					{
						Name:      "create",
						Usage:     "Add a user",
						Args:      true,
						ArgsUsage: " <NAME> <PASSWORD>",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.UserAdd(
									ctx.Args().Get(0),
									ctx.Args().Get(1),
								)
							})
						},
					},
					{
						Name:      "delete",
						Usage:     "Remove a user",
						Args:      true,
						ArgsUsage: " <ID>",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.UserDel(ctx.Args().Get(0))
							})
						},
					},
					{
						Name:      "enable",
						Usage:     "Enable a user",
						Args:      true,
						ArgsUsage: " <ID>",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.UserEnable(ctx.Args().Get(0))
							})
						},
					},
					{
						Name:      "disable",
						Usage:     "Disable a user",
						Args:      true,
						ArgsUsage: " <ID>",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.UserDisable(ctx.Args().Get(0))
							})
						},
					},
					{
						Name:      "set-admin",
						Usage:     "Promote admin privileges",
						Args:      true,
						ArgsUsage: " <ID>",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.UserSetAdmin(ctx.Args().Get(0))
							})
						},
					},
					{
						Name:      "unset-admin",
						Usage:     "Revoke admin privileges",
						Args:      true,
						ArgsUsage: " <ID>",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.UserUnsetAdmin(ctx.Args().Get(0))
							})
						},
					},
					{
						Name:      "set-hidden",
						Usage:     "Hide user from login page",
						Args:      true,
						ArgsUsage: " <ID>",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.UserSetHidden(ctx.Args().Get(0))
							})
						},
					},
					{
						Name:      "unset-hidden",
						Usage:     "Expose user on login page",
						Args:      true,
						ArgsUsage: " <ID>",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.UserUnsetHidden(ctx.Args().Get(0))
							})
						},
					},
					{
						Name:  "list",
						Usage: "Shows users",
						Flags: []cli.Flag{
							jsonFlag,
						},
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.UserList(
									ctx.Bool("json"),
								)
							})
						},
					},
				},
			},
			{
				Name:  "library",
				Usage: "Manage media libraries",
				Subcommands: []*cli.Command{
					{
						Name:  "scan",
						Usage: "Trigger a rescan of all libraries",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.LibraryScan()
							})
						},
					},
					{
						Name:  "unscraped",
						Usage: "List entries which were not scraped",
						Flags: []cli.Flag{
							itemTypesFlag,
							jsonFlag,
						},
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.LibraryUnscraped(
									ctx.StringSlice("types"),
									ctx.Bool("json"),
								)
							})
						},
					},
					{
						Name:      "search",
						Usage:     "Search throught the library",
						Args:      true,
						ArgsUsage: " <TERM>",
						Flags: []cli.Flag{
							itemTypesFlag,
							jsonFlag,
						},
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.LibrarySearch(
									ctx.Args().Get(0),
									ctx.StringSlice("types"),
									ctx.Bool("json"),
								)
							})
						},
					},
				},
			},
			{
				Name:  "key",
				Usage: "Manage API keys",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "Show keys",
						Flags: []cli.Flag{
							jsonFlag,
						},
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.KeyList(
									ctx.Bool("json"),
								)
							})
						},
					},
					{
						Name:      "create",
						Usage:     "Add a new key",
						Args:      true,
						ArgsUsage: " <NAME>",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.KeyCreate(ctx.Args().Get(0))
							})
						},
					},
					{
						Name:      "delete",
						Usage:     "Remove a key",
						Args:      true,
						ArgsUsage: " <TOKEN>",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.KeyDelete(ctx.Args().Get(0))
							})
						},
					},
				},
			},
			{
				Name:  "task",
				Usage: "Manage schedule tasks",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "Show tasks",
						Flags: []cli.Flag{
							jsonFlag,
						},
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.TaskList(
									ctx.Bool("json"),
								)
							})
						},
					},
					{
						Name:      "start",
						Usage:     "Start task",
						Args:      true,
						ArgsUsage: " <ID>",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.TaskStart(
									ctx.Args().Get(0),
								)
							})
						},
					},
					{
						Name:      "stop",
						Usage:     "Stop task",
						Args:      true,
						ArgsUsage: " <ID>",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.TaskStop(
									ctx.Args().Get(0),
								)
							})
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Exec(ctx *cli.Context, fn func(ctrl *controller.Controller) error) error {
	config := &api.Configuration{
		Servers: api.ServerConfigurations{
			{URL: ctx.String("url")},
		},
		DefaultHeader: map[string]string{
			"Authorization": fmt.Sprintf("MediaBrowser Token=\"%s\"", ctx.String("token")),
		},
	}
	client := api.NewAPIClient(config)

	ctrl := controller.New(context.Background(), client)
	return fn(ctrl)
}

var (
	itemTypesFlag = &cli.StringSliceFlag{
		Name:  "types",
		Value: cli.NewStringSlice("Movie", "Series"),
		Usage: "filter media types",
		Action: func(ctx *cli.Context, types []string) error {
			for _, t := range types {
				_, err := api.NewBaseItemKindFromValue(t)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	jsonFlag = &cli.BoolFlag{
		Name:  "json",
		Usage: "print output as JSON",
	}
)
