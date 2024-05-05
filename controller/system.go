package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
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

		b, err = json.Marshal(items)
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

func (c *Controller) SystemRestore(backupDir string) error {
	// TODO: Playlists
	// c.client.PlaystateAPI.MarkPlayedItem(c.ctx, "TODO", "TODO").DatePlayed(time.Time{}).Execute()
	// c.client.UserLibraryAPI.MarkFavoriteItem(c.ctx, "TODO", "TODO").Execute()
	return nil
}
