package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sj14/jellyfin-go/api"
)

func (c *Controller) restorePlaylists(path, userID string, backupItems map[string]api.BaseItemDto, serverItems map[string][]api.BaseItemDto, userIDsByName map[string]string, restored map[string]string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	var playlists []playlistBackup
	if err := json.Unmarshal(data, &playlists); err != nil {
		return 0, fmt.Errorf("parse playlists.json: %w", err)
	}

	skipped := 0
	for _, playlist := range playlists {
		if playlist.PlaylistID == "" || playlist.PlaylistName == "" {
			fmt.Printf("skipping playlist %q: missing name or ID\n", playlist.PlaylistName)
			skipped++
			continue
		}
		if restored[playlist.PlaylistID] != "" {
			continue // A shared playlist was already restored for another user.
		}

		backupPlaylist, ok := backupItems[playlist.PlaylistID]
		if !ok {
			fmt.Printf("skipping playlist %q: missing metadata in items.json\n", playlist.PlaylistName)
			skipped++
			continue
		}
		if backupPlaylist.GetName() != playlist.PlaylistName || backupPlaylist.GetType() != api.BASEITEMKIND_PLAYLIST {
			fmt.Printf("skipping playlist %q: inconsistent backup metadata\n", playlist.PlaylistName)
			skipped++
			continue
		}

		ids := make([]string, 0, len(playlist.Items))
		complete := true
		for _, entry := range playlist.Items {
			backupItem, ok := backupItems[entry.ID]
			if !ok || backupItem.GetName() != entry.Name {
				fmt.Printf("skipping playlist %q: no backup metadata for item %q\n", playlist.PlaylistName, entry.Name)
				complete = false
				break
			}
			match, count := findRestoreItem(backupItem, serverItems[backupItem.GetName()])
			if match == nil {
				fmt.Printf("skipping playlist %q: item %q has %d exact matches\n", playlist.PlaylistName, entry.Name, count)
				complete = false
				break
			}
			ids = append(ids, match.GetId())
		}
		if !complete {
			skipped++
			continue
		}
		permissions := make([]api.PlaylistUserPermissions, 0, len(playlist.Users))
		for _, saved := range playlist.Users {
			id := userIDsByName[strings.ToLower(saved.Name)]
			if id == "" {
				fmt.Printf("skipping playlist %q: shared user %q is missing\n", playlist.PlaylistName, saved.Name)
				complete = false
				break
			}
			permission := api.NewPlaylistUserPermissions()
			permission.SetUserId(id)
			if saved.CanEdit != nil {
				permission.SetCanEdit(*saved.CanEdit)
			}
			permissions = append(permissions, *permission)
		}
		if !complete {
			skipped++
			continue
		}

		existing, count := findRestoreItem(backupPlaylist, serverItems[playlist.PlaylistName])
		if count > 1 {
			fmt.Printf("skipping playlist %q: %d exact playlists on the server\n", playlist.PlaylistName, count)
			skipped++
			continue
		}
		if existing != nil {
			if err := c.reconcilePlaylist(existing.GetId(), userID, ids); err != nil {
				return skipped, fmt.Errorf("reconcile playlist %q: %w", playlist.PlaylistName, err)
			}
			if playlist.IsPublic != nil || playlist.Users != nil {
				update := api.NewUpdatePlaylistDto()
				if playlist.IsPublic != nil {
					update.SetIsPublic(*playlist.IsPublic)
				}
				if playlist.Users != nil {
					update.SetUsers(permissions)
				}
				if _, err := c.client.PlaylistAPI.UpdatePlaylist(c.ctx, existing.GetId()).UpdatePlaylistDto(*update).Execute(); err != nil {
					return skipped, fmt.Errorf("update playlist %q: %w", playlist.PlaylistName, err)
				}
			}
			restored[playlist.PlaylistID] = existing.GetId()
			continue
		}

		create := api.NewCreatePlaylistDto(playlist.PlaylistName)
		create.SetUserId(userID)
		create.SetIds(ids)
		if playlist.IsPublic != nil {
			create.SetIsPublic(*playlist.IsPublic)
		}
		if playlist.Users != nil {
			create.SetUsers(permissions)
		}
		result, _, err := c.client.PlaylistAPI.CreatePlaylist(c.ctx).CreatePlaylistDto(*create).Execute()
		if err != nil {
			return skipped, fmt.Errorf("create playlist %q: %w", playlist.PlaylistName, err)
		}
		if result == nil || result.GetId() == "" {
			return skipped, fmt.Errorf("create playlist %q: server returned no ID", playlist.PlaylistName)
		}
		restored[playlist.PlaylistID] = result.GetId()
	}
	return skipped, nil
}

// Keep items that already exist only on the server, while restoring the backup
// entries in their saved order at the beginning of the playlist.
func (c *Controller) reconcilePlaylist(playlistID, userID string, want []string) error {
	response, err := getAllRequestItems(c.client.PlaylistAPI.GetPlaylistItems(c.ctx, playlistID).UserId(userID))
	if err != nil {
		return err
	}
	current := response.Items
	for i, wantedID := range want {
		if i < len(current) && current[i].GetId() == wantedID {
			continue
		}
		found := -1
		for j := i + 1; j < len(current); j++ {
			if current[j].GetId() == wantedID {
				found = j
				break
			}
		}
		if found >= 0 {
			entryID := current[found].GetPlaylistItemId()
			if entryID == "" {
				return fmt.Errorf("item %q has no playlist entry ID", wantedID)
			}
			if _, err := c.client.PlaylistAPI.MoveItem(c.ctx, playlistID, entryID, int32(i)).Execute(); err != nil {
				return err
			}
			item := current[found]
			copy(current[i+1:found+1], current[i:found])
			current[i] = item
			continue
		}
		if _, err := c.client.PlaylistAPI.AddItemToPlaylist(c.ctx, playlistID).
			Ids([]string{wantedID}).
			Position(int32(i)).
			UserId(userID).
			Execute(); err != nil {
			return err
		}
		current = append(current, api.BaseItemDto{})
		copy(current[i+1:], current[i:])
		current[i] = api.BaseItemDto{Id: pointer(wantedID)}
	}
	return nil
}
