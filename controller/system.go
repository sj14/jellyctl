package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"time"

	"github.com/sj14/jellyfin-go/api"
)

func (c *Controller) GetSystemInfo() error {
	result, _, err := c.client.SystemAPI.GetSystemInfo(c.ctx).Execute()
	if err != nil {
		return err
	}

	printAsJSON(result)
	return err
}

func (c *Controller) GetPublicSystemInfo() error {
	result, _, err := c.client.SystemAPI.GetPublicSystemInfo(c.ctx).Execute()
	if err != nil {
		return err
	}

	printAsJSON(result)
	return err
}

func (c *Controller) SystemShutdown() error {
	_, err := c.client.SystemAPI.ShutdownApplication(c.ctx).Execute()
	return err
}

func (c *Controller) SystemRestart() error {
	_, err := c.client.SystemAPI.RestartApplication(c.ctx).Execute()
	return err
}

func (c *Controller) SystemBackup() error {
	users, _, err := c.client.UserAPI.GetUsers(c.ctx).Execute()
	if err != nil {
		return err
	}

	basedir := filepath.Join("jellyctl-backup", fmt.Sprint(time.Now().Unix()))

	for _, user := range users {
		userdir := filepath.Join(basedir, "users", user.GetName())
		err = os.MkdirAll(userdir, os.ModePerm)
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(userdir, "user.json"), b, os.ModePerm)
		if err != nil {
			return err
		}

		items, _, err := c.client.ItemsAPI.GetItems(c.ctx).
			SearchTerm("").
			Recursive(true).
			UserId(user.GetId()). //  needed for getting the userData (favorite, played)
			Execute()
		if err != nil {
			return err
		}

		b, err = json.MarshalIndent(items.Items, "", "  ")
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(userdir, "items.json"), b, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO:
// - Complete user restore
// - Playlists
func (c *Controller) SystemRestore(backupDir string, unplayed, unfav bool) error {
	if backupDir == "" {
		return errors.New("missing path to the backup directory")
	}

	userdir := filepath.Join(backupDir, "users")
	dirEntries, err := os.ReadDir(userdir)
	if err != nil {
		return err
	}

	users, _, err := c.client.UserAPI.GetUsers(c.ctx).Execute()
	if err != nil {
		return err
	}

	var backupUsernames []string
	for _, dirEntry := range dirEntries {
		backupUsernames = append(backupUsernames, dirEntry.Name())
	}

	// create missing users
	// TODO: use user.json from backup for settings (e.g. admin, hidden, disabled, ...)
	for _, backupUser := range backupUsernames {
		found := false
		for _, systemUser := range users {
			if systemUser.GetName() == backupUser {
				found = true
				break
			}
		}
		if !found {
			pass := fmt.Sprint(rand.Int())
			fmt.Printf("creating new user %q with (unsafe) password %q\n", backupUser, pass)
			err := c.UserAdd(backupUser, pass)
			if err != nil {
				return err
			}
		}
	}

	// reload users as missing ones might have been just created
	users, _, err = c.client.UserAPI.GetUsers(c.ctx).Execute()
	if err != nil {
		return err
	}

	for _, dirEntry := range dirEntries {
		userName := dirEntry.Name()

		for _, user := range users {
			if user.GetName() != userName {
				continue
			}

			fmt.Printf("restoring data for %q\n", user.GetName())

			itemsJson, err := os.ReadFile(filepath.Join(userdir, userName, "items.json"))
			if err != nil {
				return err
			}

			var items []api.BaseItemDto
			err = json.Unmarshal(itemsJson, &items)
			if err != nil {
				return err
			}

			for _, item := range items {
				if played, ok := item.UserData.Get().GetPlayedOk(); ok {
					if *played {
						_, _, err = c.client.PlaystateAPI.MarkPlayedItem(
							c.ctx,
							item.GetId()).
							UserId(user.GetId()).
							DatePlayed(item.UserData.Get().GetLastPlayedDate()).
							Execute()
						if err != nil {
							return err
						}

						// TODO: probably not the right API, where to set the user ID?
						// Check model_playback_progress_info_item.go / UserData NullableBaseItemDtoUserData
						//
						// _, err = c.client.PlaystateAPI.ReportPlaybackProgress(c.ctx).
						// 	PlaybackProgressInfo(api.PlaybackProgressInfo{
						// 		ItemId:        item.Id,
						// 		PositionTicks: *api.NewNullableInt64(item.GetUserData().PlaybackPositionTicks),
						// 	},
						// 	).Execute()
						// if err != nil {
						// 	return err
						// }
					} else if unplayed {
						_, _, err = c.client.PlaystateAPI.MarkUnplayedItem(
							c.ctx,
							item.GetId()).
							UserId(user.GetId()).
							Execute()
						if err != nil {
							return err
						}
					}

					if fav, ok := item.UserData.Get().GetIsFavoriteOk(); ok {
						if *fav {
							_, _, err = c.client.UserLibraryAPI.MarkFavoriteItem(
								c.ctx,
								item.GetId()).
								UserId(user.GetId()).
								Execute()
							if err != nil {
								return err
							}
						} else if unfav {
							_, _, err = c.client.UserLibraryAPI.UnmarkFavoriteItem(
								c.ctx,
								item.GetId()).
								UserId(user.GetId()).
								Execute()
							if err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}

	return nil
}
