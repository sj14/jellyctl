package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sj14/jellyfin-go/api"
)

func (c *Controller) restoreUserSettings(userID, userName, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var backup api.UserDto
	if err := json.Unmarshal(data, &backup); err != nil {
		return fmt.Errorf("parse user.json: %w", err)
	}
	if backup.GetName() != "" && !strings.EqualFold(backup.GetName(), userName) {
		return fmt.Errorf("user.json belongs to %q, not %q", backup.GetName(), userName)
	}
	if backup.Policy.Get() == nil && backup.Configuration.Get() == nil {
		return nil
	}

	current, _, err := c.client.UserAPI.GetUserById(c.ctx, userID).Execute()
	if err != nil {
		return fmt.Errorf("get current user: %w", err)
	}
	if current == nil {
		return fmt.Errorf("user %q was not returned by the server", userID)
	}
	if !strings.EqualFold(current.GetName(), userName) {
		return fmt.Errorf("server user %q does not match backup user %q", current.GetName(), userName)
	}

	if source := backup.Policy.Get(); source != nil {
		target := current.Policy.Get()
		if target == nil {
			return fmt.Errorf("user %q has no current policy", userID)
		}
		policy := *target
		changed := false
		if source.IsAdministrator != nil {
			policy.IsAdministrator = source.IsAdministrator
			changed = true
		}
		if source.IsHidden != nil {
			policy.IsHidden = source.IsHidden
			changed = true
		}
		if source.IsDisabled != nil {
			policy.IsDisabled = source.IsDisabled
			changed = true
		}
		if changed {
			if err := c.userUpdatePolicy(userID, policy); err != nil {
				return fmt.Errorf("update policy: %w", err)
			}
		}
	}

	if source := backup.Configuration.Get(); source != nil {
		target := current.Configuration.Get()
		if target == nil {
			return fmt.Errorf("user %q has no current configuration", userID)
		}
		config := *target
		changed := copyPortableUserConfiguration(&config, source)
		if changed {
			_, err := c.client.UserAPI.UpdateUserConfiguration(c.ctx).
				UserId(userID).
				UserConfiguration(config).
				Execute()
			if err != nil {
				return fmt.Errorf("update configuration: %w", err)
			}
		}
	}
	return nil
}

// Library view and cast receiver IDs refer to the old server, so only portable
// preferences are copied from the backup.
func copyPortableUserConfiguration(target, source *api.UserConfiguration) bool {
	changed := false
	if source.AudioLanguagePreference.IsSet() {
		target.AudioLanguagePreference = source.AudioLanguagePreference
		changed = true
	}
	if source.PlayDefaultAudioTrack != nil {
		target.PlayDefaultAudioTrack = source.PlayDefaultAudioTrack
		changed = true
	}
	if source.SubtitleLanguagePreference.IsSet() {
		target.SubtitleLanguagePreference = source.SubtitleLanguagePreference
		changed = true
	}
	if source.SubtitleMode != nil {
		target.SubtitleMode = source.SubtitleMode
		changed = true
	}
	if source.DisplayMissingEpisodes != nil {
		target.DisplayMissingEpisodes = source.DisplayMissingEpisodes
		changed = true
	}
	if source.DisplayCollectionsView != nil {
		target.DisplayCollectionsView = source.DisplayCollectionsView
		changed = true
	}
	if source.HidePlayedInLatest != nil {
		target.HidePlayedInLatest = source.HidePlayedInLatest
		changed = true
	}
	if source.RememberAudioSelections != nil {
		target.RememberAudioSelections = source.RememberAudioSelections
		changed = true
	}
	if source.RememberSubtitleSelections != nil {
		target.RememberSubtitleSelections = source.RememberSubtitleSelections
		changed = true
	}
	if source.EnableNextEpisodeAutoPlay != nil {
		target.EnableNextEpisodeAutoPlay = source.EnableNextEpisodeAutoPlay
		changed = true
	}
	return changed
}
