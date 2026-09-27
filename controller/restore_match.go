package controller

import "github.com/sj14/jellyfin-go/api"

// findRestoreItem returns a server item only when exactly one candidate matches
// the identity fields stored in the backup.
func findRestoreItem(backup api.BaseItemDto, candidates []api.BaseItemDto) (*api.BaseItemDto, int) {
	var match *api.BaseItemDto
	count := 0
	for i := range candidates {
		if !sameRestoreItem(backup, candidates[i]) {
			continue
		}
		count++
		match = &candidates[i]
	}
	if count != 1 {
		return nil, count
	}
	return match, count
}

func sameRestoreItem(backup, candidate api.BaseItemDto) bool {
	if backup.GetName() == "" || backup.GetType() == "" || candidate.GetId() == "" ||
		backup.GetName() != candidate.GetName() ||
		backup.GetType() != candidate.GetType() ||
		backup.GetProductionYear() != candidate.GetProductionYear() ||
		backup.GetSeriesName() != candidate.GetSeriesName() ||
		backup.GetParentIndexNumber() != candidate.GetParentIndexNumber() ||
		backup.GetIndexNumber() != candidate.GetIndexNumber() {
		return false
	}

	// Provider IDs are stable across servers when both sides have them.
	for provider, backupID := range backup.GetProviderIds() {
		candidateID := candidate.GetProviderIds()[provider]
		if backupID != nil && candidateID != nil && *backupID != "" && *candidateID != "" && *backupID != *candidateID {
			return false
		}
	}
	return true
}
