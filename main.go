package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sj14/jellyctl/pkg"
	jellyapi "github.com/sj14/jellyfin-go/api"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                 "jellyctl",
		Usage:                "Manage Jellyfin from the CLI",
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
					return Exec(ctx, func(ctrl *pkg.Controller) (*http.Response, error) {
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
				Usage: "Show system information",
				Action: func(ctx *cli.Context) error {
					return Exec(ctx, func(ctrl *pkg.Controller) (*http.Response, error) {
						return ctrl.GetPublicSystemInfo()
					})
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
							return Exec(ctx, func(ctrl *pkg.Controller) (*http.Response, error) {
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
							return Exec(ctx, func(ctrl *pkg.Controller) (*http.Response, error) {
								return ctrl.UserDel(ctx.Args().Get(0))
							})
						},
					},
					{
						Name:  "list",
						Usage: "Show users",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *pkg.Controller) (*http.Response, error) {
								return ctrl.UserList()
							})
						},
					},
				},
			},
			{
				Name:  "library",
				Usage: "Manage the media libraries",
				Subcommands: []*cli.Command{
					{
						Name:  "scan",
						Usage: "Trigger a rescan of all libraries",
						Action: func(ctx *cli.Context) error {
							return Exec(ctx, func(ctrl *pkg.Controller) (*http.Response, error) {
								return ctrl.LibraryScan()
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

func Exec(ctx *cli.Context, fn func(ctrl *pkg.Controller) (*http.Response, error)) error {
	config := &jellyapi.Configuration{
		Servers:       jellyapi.ServerConfigurations{{URL: ctx.String("url")}},
		DefaultHeader: map[string]string{"Authorization": fmt.Sprintf("MediaBrowser Token=\"%s\"", ctx.String("token"))},
	}
	client := jellyapi.NewAPIClient(config)

	ctrl := pkg.NewController(context.Background(), client)
	_, err := fn(ctrl)
	// if err != nil {
	// 	log.Printf("%#v\n", resp)
	// }
	return err
}
