package controller

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

type playlistBackup struct {
	PlaylistID   string
	PlaylistName string
	Items        []NameID
	IsPublic     *bool
	Users        []playlistUserBackup
}

type playlistUserBackup struct {
	Name    string
	CanEdit *bool
}

type NameID struct {
	Name string
	ID   string
}

func (c *Controller) SystemBackup() error {
	users, _, err := c.client.UserAPI.GetUsers(c.ctx).Execute()
	if err != nil {
		return err
	}
	userNamesByID := make(map[string]string, len(users))
	for _, user := range users {
		userNamesByID[user.GetId()] = user.GetName()
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

		items, err := getAllRequestItems(c.client.LibraryAPI.GetItems(c.ctx).
			SearchTerm("").
			Recursive(true).
			Fields([]api.ItemFields{api.ITEMFIELDS_PROVIDER_IDS}).
			EnableUserData(true).
			UserId(user.GetId())) // needed for getting the userData (favorite, played)
		if err != nil {
			return err
		}

		var playlists []playlistBackup
		for _, item := range items.Items {
			// If item is a playlist, get all items of the playlist.
			// Otherwise, we won't have a link between playlist and its content.
			if item.GetType() == api.BASEITEMKIND_PLAYLIST {
				playlistItems, err := getAllRequestItems(c.client.PlaylistAPI.GetPlaylistItems(c.ctx, item.GetId()).
					EnableUserData(false).
					UserId(user.GetId()))
				if err != nil {
					return err
				}

				backup := playlistBackup{
					PlaylistID:   item.GetId(),
					PlaylistName: item.GetName(),
				}
				playlistInfo, _, err := c.client.PlaylistAPI.GetPlaylist(c.ctx, item.GetId()).Execute()
				if err != nil {
					return fmt.Errorf("get playlist %q: %w", item.GetName(), err)
				}
				if playlistInfo == nil {
					return fmt.Errorf("get playlist %q: server returned no details", item.GetName())
				}
				backup.IsPublic = playlistInfo.OpenAccess
				backup.Users = make([]playlistUserBackup, 0, len(playlistInfo.Shares))
				for _, permission := range playlistInfo.Shares {
					name := userNamesByID[permission.GetUserId()]
					if name == "" {
						return fmt.Errorf("playlist %q refers to unknown user %q", item.GetName(), permission.GetUserId())
					}
					backup.Users = append(backup.Users, playlistUserBackup{Name: name, CanEdit: permission.CanEdit})
				}

				for _, playlistItem := range playlistItems.GetItems() {
					backup.Items = append(backup.Items, NameID{
						Name: playlistItem.GetName(),
						ID:   playlistItem.GetId(),
					})
				}

				playlists = append(playlists, backup)
			}
		}

		b, err = json.MarshalIndent(playlists, "", "  ")
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(userdir, "playlists.json"), b, os.ModePerm)
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
		if !dirEntry.IsDir() {
			return fmt.Errorf("unexpected file %q in backup users directory", dirEntry.Name())
		}
		backupUsernames = append(backupUsernames, dirEntry.Name())
	}

	// Create missing users before restoring their settings and item data.
	for _, backupUser := range backupUsernames {
		found := false
		for _, systemUser := range users {
			if strings.EqualFold(systemUser.GetName(), backupUser) {
				found = true
				break
			}
		}
		if !found {
			pass := rand.Text()
			fmt.Printf("creating new user %q with temporary password %q\n", backupUser, pass)
			err := c.UserAdd(backupUser, pass)
			if err != nil {
				return fmt.Errorf("add user: %w", err)
			}
		}
	}

	// reload users as missing ones might have been just created
	users, _, err = c.client.UserAPI.GetUsers(c.ctx).Execute()
	if err != nil {
		return err
	}
	userIDsByName := make(map[string]string, len(users))
	for _, user := range users {
		userIDsByName[strings.ToLower(user.GetName())] = user.GetId()
	}

	serverItems, err := getAllRequestItems(c.client.LibraryAPI.GetItems(c.ctx).
		Recursive(true).
		Fields([]api.ItemFields{api.ITEMFIELDS_PROVIDER_IDS}))
	if err != nil {
		return fmt.Errorf("get server items: %w", err)
	}
	itemsByName := make(map[string][]api.BaseItemDto)
	for _, item := range serverItems.Items {
		itemsByName[item.GetName()] = append(itemsByName[item.GetName()], item)
	}

	unmatched := 0
	restoredPlaylists := make(map[string]string)
	for _, dirEntry := range dirEntries {
		userName := dirEntry.Name()
		foundUser := false

		for _, user := range users {
			if !strings.EqualFold(user.GetName(), userName) {
				continue
			}
			foundUser = true

			fmt.Printf("restoring data for %q\n", user.GetName())
			if err := c.restoreUserSettings(user.GetId(), userName, filepath.Join(userdir, userName, "user.json")); err != nil {
				return fmt.Errorf("restore settings for %q: %w", userName, err)
			}

			itemsJson, err := os.ReadFile(filepath.Join(userdir, userName, "items.json"))
			if err != nil {
				return fmt.Errorf("read items.json: %w", err)
			}

			var items []api.BaseItemDto
			err = json.Unmarshal(itemsJson, &items)
			if err != nil {
				return fmt.Errorf("unmarshal items.json: %w", err)
			}
			backupItemsByID := make(map[string]api.BaseItemDto, len(items))
			for _, item := range items {
				if item.GetId() != "" {
					backupItemsByID[item.GetId()] = item
				}
			}

			for _, backupItem := range items {
				if backupItem.GetType() == api.BASEITEMKIND_PLAYLIST {
					continue // Playlists and their contents are restored below.
				}
				// We have to find the same item on the server again, as the IDs won't match when the server changed.
				if backupItem.GetName() == "" {
					fmt.Printf("skipping item with empty name for user %q\n", userName)
					unmatched++
					continue
				}
				serverItem, matches := findRestoreItem(backupItem, itemsByName[backupItem.GetName()])
				if serverItem == nil {
					fmt.Printf("skipping %q (%s) for user %q: %d exact matches\n", backupItem.GetName(), backupItem.GetType(), userName, matches)
					unmatched++
					continue
				}

				userData := backupItem.UserData.Get()
				if userData == nil {
					fmt.Printf("skipping %q (%s) for user %q: no user data in backup\n", backupItem.GetName(), backupItem.GetType(), userName)
					unmatched++
					continue
				}

				if played, ok := userData.GetPlayedOk(); ok {
					if *played {
						_, _, err = c.client.UserDataAPI.MarkPlayedItem(
							c.ctx,
							serverItem.GetId(),
						).
							UserId(user.GetId()).
							DatePlayed(userData.GetLastPlayedDate()).
							Execute()
						if err != nil {
							return fmt.Errorf("mark played item: %w", err)
						}

					} else if unplayed {
						_, _, err = c.client.UserDataAPI.MarkUnplayedItem(
							c.ctx,
							serverItem.GetId(),
						).
							UserId(user.GetId()).
							Execute()
						if err != nil {
							return err
						}
					}
				}

				if fav, ok := userData.GetIsFavoriteOk(); ok {
					if *fav {
						_, _, err = c.client.UserDataAPI.MarkFavoriteItem(
							c.ctx,
							serverItem.GetId(),
						).
							UserId(user.GetId()).
							Execute()
						if err != nil {
							return fmt.Errorf("mark favourite item: %w", err)
						}
					} else if unfav {
						_, _, err = c.client.UserDataAPI.UnmarkFavoriteItem(
							c.ctx,
							serverItem.GetId(),
						).
							UserId(user.GetId()).
							Execute()
						if err != nil {
							return fmt.Errorf("unmark favourite item: %w", err)
						}
					}
				}

				if ticks := userData.GetPlaybackPositionTicks(); ticks > 0 {
					progress := api.NewUpdateUserItemDataDto()
					progress.SetPlaybackPositionTicks(ticks)
					_, _, err = c.client.UserDataAPI.UpdateItemUserData(c.ctx, serverItem.GetId()).
						UserId(user.GetId()).
						UpdateUserItemDataDto(*progress).
						Execute()
					if err != nil {
						return fmt.Errorf("restore playback position for %q: %w", backupItem.GetName(), err)
					}
				}
			}
			skippedPlaylists, err := c.restorePlaylists(filepath.Join(userdir, userName, "playlists.json"), user.GetId(), backupItemsByID, itemsByName, userIDsByName, restoredPlaylists)
			if err != nil {
				return fmt.Errorf("restore playlists for %q: %w", userName, err)
			}
			unmatched += skippedPlaylists
			break
		}
		if !foundUser {
			return fmt.Errorf("user %q was not returned after creation", userName)
		}
	}

	if unmatched > 0 {
		return fmt.Errorf("restore incomplete: %d items or playlists could not be restored", unmatched)
	}
	return nil
}
