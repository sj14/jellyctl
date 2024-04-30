package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sj14/jellyctl/controller"
	jellyapi "github.com/sj14/jellyfin-go/api"
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
				Value:   "",
				EnvVars: []string{"JELLYCTL_TOKEN"},
				Usage:   "API token",
			},
		},
		Commands: []*cli.Command{
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
				},
				Action: func(ctx *cli.Context) error {
					return Exec(ctx, func(ctrl *controller.Controller) error {
						return ctrl.GetLogEntries(
							int32(ctx.Int("start")),
							int32(ctx.Int("limit")),
							*ctx.Timestamp("after"),
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
						Usage:     "Protmote admin privileges",
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
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.UserList()
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
							&cli.BoolFlag{
								Name:  "movies",
								Value: true,
								Usage: "show unscraped movies",
							},
							&cli.BoolFlag{
								Name:  "shows",
								Value: true,
								Usage: "show unscraped shows",
							},
							&cli.BoolFlag{
								Name:  "seasons",
								Value: false,
								Usage: "show unscraped seasons",
							},
							&cli.BoolFlag{
								Name:  "episodes",
								Value: false,
								Usage: "show unscraped episodes",
							},
						},
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *controller.Controller) error {
								return ctrl.LibraryUnscraped(
									ctx.Bool("movies"),
									ctx.Bool("shows"),
									ctx.Bool("seasons"),
									ctx.Bool("episodes"),
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
	config := &jellyapi.Configuration{
		Servers:       jellyapi.ServerConfigurations{{URL: ctx.String("url")}},
		DefaultHeader: map[string]string{"Authorization": fmt.Sprintf("MediaBrowser Token=\"%s\"", ctx.String("token"))},
	}
	client := jellyapi.NewAPIClient(config)

	ctrl := controller.New(context.Background(), client)
	return fn(ctrl)
}
