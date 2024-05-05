package controller

import (
	"encoding/json"
	"errors"
	"fmt"
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
		userdir := filepath.Join(basedir, "users", user.GetId())
		err = os.MkdirAll(userdir, os.ModePerm)
		if err != nil {
			return err
		}

		b, err := json.Marshal(user)
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

		b, err = json.Marshal(items.Items)
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

func (c *Controller) SystemRestore(backupDir string, unplayed, unfav bool) error {
	if backupDir == "" {
		return errors.New("missing path to the backup directory")
	}
	// TODO:
	// - Complete user restore
	// - Playlists

	userdir := filepath.Join(backupDir, "users")
	entries, err := os.ReadDir(userdir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		userID := entry.Name()

		itemsPath := filepath.Join(userdir, userID, "items.json")

		itemsJson, err := os.ReadFile(itemsPath)
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
						userID,
						item.GetId()).
						DatePlayed(item.UserData.Get().GetLastPlayedDate()).
						Execute()
					if err != nil {
						return err
					}
				} else if unplayed {
					_, _, err = c.client.PlaystateAPI.MarkUnplayedItem(
						c.ctx,
						userID,
						item.GetId()).
						Execute()
					if err != nil {
						return err
					}
				}

				if fav, ok := item.UserData.Get().GetIsFavoriteOk(); ok {
					if *fav {
						_, _, err = c.client.UserLibraryAPI.MarkFavoriteItem(
							c.ctx,
							userID,
							item.GetId()).
							Execute()
						if err != nil {
							return err
						}
					} else if unfav {
						_, _, err = c.client.UserLibraryAPI.UnmarkFavoriteItem(
							c.ctx,
							userID,
							item.GetId()).
							Execute()
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}
