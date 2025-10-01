package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sj14/jellyctl/controller"
	"github.com/sj14/jellyfin-go/api"
	docs "github.com/urfave/cli-docs/v3"
	"github.com/urfave/cli/v3"
)

var (
	// will be replaced during the build process
	version = "undefined"
	// commit  = "undefined"
	// date    = "undefined"
)

func main() {
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

var app = &cli.Command{
	Name:                  "jellyctl",
	Usage:                 "Manage Jellyfin from the CLI",
	Version:               version,
	Suggest:               true,
	EnableShellCompletion: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "url",
			Value:   "http://127.0.0.1:8096",
			Sources: cli.EnvVars("JELLYCTL_URL"),
			Usage:   "URL of the Jellyfin server",
		},
		&cli.StringFlag{
			Name:    "token",
			Sources: cli.EnvVars("JELLYCTL_TOKEN"),
			Usage:   "API token",
		},
	},
	Commands: []*cli.Command{
		{
			Name:   "docs",
			Hidden: true,
			Action: func(ctx context.Context, cmd *cli.Command) error {
				s, err := docs.ToMarkdown(cmd.Root())
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
					Name:  "after",
					Usage: "only logs after the given time",
					Config: cli.TimestampConfig{
						Layouts:  []string{time.DateTime},
						Timezone: time.Local,
					},
					Value: time.Time{},
				},
				jsonFlag,
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return Exec(cmd, func(ctrl *controller.Controller) error {
					return ctrl.GetLogEntries(
						int32(cmd.Int("start")),
						int32(cmd.Int("limit")),
						cmd.Timestamp("after"),
						cmd.Bool("json"),
					)
				})
			},
		},
		{
			Name:  "system",
			Usage: "Manage the system",
			Commands: []*cli.Command{
				{
					Name:  "shutdown",
					Usage: "Stop the Jellyfin process",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.SystemShutdown()
						})
					},
				},
				{
					Name:  "restart",
					Usage: "Restart the Jellyfin process",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
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
					Action: func(ctx context.Context, cmd *cli.Command) error {
						if cmd.Bool("public") {
							return Exec(cmd, func(ctrl *controller.Controller) error {
								return ctrl.GetPublicSystemInfo()
							})
						} else {
							return Exec(cmd, func(ctrl *controller.Controller) error {
								return ctrl.GetSystemInfo()
							})
						}
					},
				},
				{
					Name:  "backup",
					Usage: "Export some data (EXPERIMENTAL)",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.SystemBackup()
						})
					},
				},
				{
					Name:      "restore",
					Usage:     "Import played and favourite information (based on the user name not user ID!) (EXPERIMENTAL)",
					ArgsUsage: "<PATH>",
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
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.SystemRestore(
								cmd.Args().Get(0),
								cmd.Bool("unplay"),
								cmd.Bool("unfav"),
							)
						})
					},
				},
			},
		},
		{
			Name:  "user",
			Usage: "Manage users",
			Commands: []*cli.Command{
				{
					Name:      "create",
					Usage:     "Add a user",
					ArgsUsage: "<NAME> <PASSWORD>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.UserAdd(
								cmd.Args().Get(0),
								cmd.Args().Get(1),
							)
						})
					},
				},
				{
					Name:      "delete",
					Usage:     "Remove a user",
					ArgsUsage: "<ID>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.UserDel(cmd.Args().Get(0))
						})
					},
				},
				{
					Name:      "enable",
					Usage:     "Enable a user",
					ArgsUsage: "<ID>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.UserEnable(cmd.Args().Get(0))
						})
					},
				},
				{
					Name:      "disable",
					Usage:     "Disable a user",
					ArgsUsage: "<ID>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.UserDisable(cmd.Args().Get(0))
						})
					},
				},
				{
					Name:      "set-admin",
					Usage:     "Promote admin privileges",
					ArgsUsage: "<ID>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.UserSetAdmin(cmd.Args().Get(0))
						})
					},
				},
				{
					Name:      "unset-admin",
					Usage:     "Revoke admin privileges",
					ArgsUsage: "<ID>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.UserUnsetAdmin(cmd.Args().Get(0))
						})
					},
				},
				{
					Name:      "set-hidden",
					Usage:     "Hide user from login page",
					ArgsUsage: "<ID>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.UserSetHidden(cmd.Args().Get(0))
						})
					},
				},
				{
					Name:      "unset-hidden",
					Usage:     "Expose user on login page",
					ArgsUsage: "<ID>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.UserUnsetHidden(cmd.Args().Get(0))
						})
					},
				},
				{
					Name:  "list",
					Usage: "Shows users",
					Flags: []cli.Flag{
						jsonFlag,
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.UserList(
								cmd.Bool("json"),
							)
						})
					},
				},
			},
		},
		{
			Name:  "library",
			Usage: "Manage media libraries",
			Commands: []*cli.Command{
				{
					Name:  "scan",
					Usage: "Trigger a rescan of all libraries",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
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
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.LibraryUnscraped(
								cmd.StringSlice("types"),
								cmd.Bool("json"),
							)
						})
					},
				},
				{
					Name:      "search",
					Usage:     "Search throught the library",
					ArgsUsage: "<TERM>",
					Flags: []cli.Flag{
						itemTypesFlag,
						jsonFlag,
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.LibrarySearch(
								cmd.Args().Get(0),
								cmd.StringSlice("types"),
								cmd.Bool("json"),
							)
						})
					},
				},
				{
					Name:      "duplicates",
					Usage:     "List duplicates in the library",
					ArgsUsage: "<TERM>",
					Flags: []cli.Flag{
						itemTypesFlag,
						jsonFlag,
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.LibraryDuplicates(
								cmd.Args().Get(0),
								cmd.StringSlice("types"),
								cmd.Bool("json"),
							)
						})
					},
				},
			},
		},
		{
			Name:  "key",
			Usage: "Manage API keys",
			Commands: []*cli.Command{
				{
					Name:  "list",
					Usage: "Show keys",
					Flags: []cli.Flag{
						jsonFlag,
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.KeyList(
								cmd.Bool("json"),
							)
						})
					},
				},
				{
					Name:      "create",
					Usage:     "Add a new key",
					ArgsUsage: "<NAME>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.KeyCreate(cmd.Args().Get(0))
						})
					},
				},
				{
					Name:      "delete",
					Usage:     "Remove a key",
					ArgsUsage: "<TOKEN>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.KeyDelete(cmd.Args().Get(0))
						})
					},
				},
			},
		},
		{
			Name:  "task",
			Usage: "Manage schedule tasks",
			Commands: []*cli.Command{
				{
					Name:  "list",
					Usage: "Show tasks",
					Flags: []cli.Flag{
						jsonFlag,
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.TaskList(
								cmd.Bool("json"),
							)
						})
					},
				},
				{
					Name:      "start",
					Usage:     "Start task",
					ArgsUsage: "<ID>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.TaskStart(
								cmd.Args().Get(0),
							)
						})
					},
				},
				{
					Name:      "stop",
					Usage:     "Stop task",
					ArgsUsage: "<ID>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.TaskStop(
								cmd.Args().Get(0),
							)
						})
					},
				},
			},
		},
	},
}

func Exec(cmd *cli.Command, fn func(ctrl *controller.Controller) error) error {
	config := &api.Configuration{
		Servers: api.ServerConfigurations{
			{URL: cmd.String("url")},
		},
		DefaultHeader: map[string]string{
			"Authorization": fmt.Sprintf("MediaBrowser Token=\"%s\"", cmd.String("token")),
		},
	}
	client := api.NewAPIClient(config)

	ctrl := controller.New(context.Background(), client)
	return fn(ctrl)
}

var (
	itemTypesFlag = &cli.StringSliceFlag{
		Name:  "types",
		Value: []string{"Movie", "Series"},
		Usage: "filter media types",
		Action: func(ctx context.Context, cmd *cli.Command, types []string) error {
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
		Name:    "json",
		Aliases: []string{"j"},
		Usage:   "print output as JSON",
	}
)
