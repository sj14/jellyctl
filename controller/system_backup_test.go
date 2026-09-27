package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/sj14/jellyfin-go/api"
)

func TestSystemBackupPreservesPlaylistSharing(t *testing.T) {
	t.Chdir(t.TempDir())
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/Users":
			_, _ = w.Write([]byte(`[{"Id":"user","Name":"alice"},{"Id":"shared","Name":"bob"}]`))
		case "/Items":
			if r.URL.Query().Get("enableUserData") != "true" {
				t.Errorf("user data was not requested: %s", r.URL.RawQuery)
			}
			if r.URL.Query().Get("userId") == "user" {
				_, _ = w.Write([]byte(`{"Items":[{"Id":"old-list","Name":"Favorites","Type":"Playlist"},{"Id":"old-media","Name":"First","Type":"Movie"}],"TotalRecordCount":2}`))
			} else {
				_, _ = w.Write([]byte(`{"Items":[],"TotalRecordCount":0}`))
			}
		case "/Playlists/old-list/Items":
			_, _ = w.Write([]byte(`{"Items":[{"Id":"old-media","Name":"First"}],"TotalRecordCount":1}`))
		case "/Playlists/old-list":
			_, _ = w.Write([]byte(`{"OpenAccess":true,"Shares":[{"UserId":"shared","CanEdit":true}],"ItemIds":["old-media"]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewAPIClient(&api.Configuration{Servers: api.ServerConfigurations{{URL: server.URL}}})
	if err := New(context.Background(), client).SystemBackup(); err != nil {
		t.Fatal(err)
	}
	backups, err := os.ReadDir("jellyctl-backup")
	if err != nil || len(backups) != 1 {
		t.Fatalf("backup directories = %v, err=%v", backups, err)
	}
	data, err := os.ReadFile(filepath.Join("jellyctl-backup", backups[0].Name(), "users", "alice", "playlists.json"))
	if err != nil {
		t.Fatal(err)
	}
	var playlists []playlistBackup
	if err := json.Unmarshal(data, &playlists); err != nil {
		t.Fatal(err)
	}
	if len(playlists) != 1 || playlists[0].IsPublic == nil || !*playlists[0].IsPublic || len(playlists[0].Users) != 1 || playlists[0].Users[0].Name != "bob" || playlists[0].Users[0].CanEdit == nil || !*playlists[0].Users[0].CanEdit {
		t.Fatalf("playlist backup = %+v", playlists)
	}
}
