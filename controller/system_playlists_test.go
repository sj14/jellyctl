package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sj14/jellyfin-go/api"
)

func TestRestorePlaylistsCreatesPlaylistWithMappedItems(t *testing.T) {
	path := filepath.Join(t.TempDir(), "playlists.json")
	data := `[{"PlaylistID":"old-list","PlaylistName":"Favorites","IsPublic":true,"Users":[{"Name":"bob","CanEdit":true}],"Items":[{"Name":"First","ID":"old-first"},{"Name":"Second","ID":"old-second"}]}]`
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	decode := func(data string) api.BaseItemDto {
		t.Helper()
		var item api.BaseItemDto
		if err := json.Unmarshal([]byte(data), &item); err != nil {
			t.Fatal(err)
		}
		return item
	}
	backupItems := map[string]api.BaseItemDto{
		"old-list":   decode(`{"Id":"old-list","Name":"Favorites","Type":"Playlist"}`),
		"old-first":  decode(`{"Id":"old-first","Name":"First","Type":"Movie","ProductionYear":2020}`),
		"old-second": decode(`{"Id":"old-second","Name":"Second","Type":"Movie","ProductionYear":2021}`),
	}
	serverItems := map[string][]api.BaseItemDto{
		"First":  {decode(`{"Id":"new-first","Name":"First","Type":"Movie","ProductionYear":2020}`)},
		"Second": {decode(`{"Id":"new-second","Name":"Second","Type":"Movie","ProductionYear":2021}`)},
	}

	var created api.CreatePlaylistDto
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/Playlists" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&created); err != nil {
			t.Errorf("decode playlist: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"Id":"new-list"}`))
	}))
	defer server.Close()

	client := api.NewAPIClient(&api.Configuration{Servers: api.ServerConfigurations{{URL: server.URL}}})
	restored := make(map[string]string)
	skipped, err := New(context.Background(), client).restorePlaylists(path, "user", backupItems, serverItems, map[string]string{"bob": "new-bob"}, restored)
	if err != nil || skipped != 0 {
		t.Fatalf("skipped=%d, err=%v", skipped, err)
	}
	if created.GetName() != "Favorites" || created.GetUserId() != "user" || !created.GetIsPublic() || !reflect.DeepEqual(created.GetIds(), []string{"new-first", "new-second"}) || restored["old-list"] != "new-list" || len(created.GetUsers()) != 1 || created.GetUsers()[0].GetUserId() != "new-bob" || !created.GetUsers()[0].GetCanEdit() {
		t.Errorf("created=%+v, restored=%v", created, restored)
	}
}

func TestReconcilePlaylistMovesAndAddsItems(t *testing.T) {
	var operations []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/Playlists/list/Items":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"Items":[{"Id":"B","PlaylistItemId":"entry-B"},{"Id":"A","PlaylistItemId":"entry-A"}],"TotalRecordCount":2}`))
		case r.Method == http.MethodPost && r.URL.Path == "/Playlists/list/Items/entry-A/Move/0":
			operations = append(operations, "move A")
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPost && r.URL.Path == "/Playlists/list/Items":
			if r.URL.Query().Get("position") != "2" || r.URL.Query().Get("userId") != "user" || r.URL.Query().Get("ids") != "C" {
				t.Errorf("unexpected add query: %s", r.URL.RawQuery)
			}
			operations = append(operations, "add C")
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewAPIClient(&api.Configuration{Servers: api.ServerConfigurations{{URL: server.URL}}})
	if err := New(context.Background(), client).reconcilePlaylist("list", "user", []string{"A", "B", "C"}); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(operations, []string{"move A", "add C"}) {
		t.Errorf("operations = %v", operations)
	}
}

func TestRestorePlaylistsUpdatesExistingSharing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "playlists.json")
	data := `[{"PlaylistID":"old-list","PlaylistName":"Favorites","IsPublic":false,"Users":[{"Name":"bob","CanEdit":false}],"Items":[]}]`
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	var old, current api.BaseItemDto
	if err := json.Unmarshal([]byte(`{"Id":"old-list","Name":"Favorites","Type":"Playlist"}`), &old); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(`{"Id":"current-list","Name":"Favorites","Type":"Playlist"}`), &current); err != nil {
		t.Fatal(err)
	}
	var updated api.UpdatePlaylistDto
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/Playlists/current-list/Items":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"Items":[],"TotalRecordCount":0}`))
		case r.Method == http.MethodPost && r.URL.Path == "/Playlists/current-list":
			if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
				t.Errorf("decode update: %v", err)
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewAPIClient(&api.Configuration{Servers: api.ServerConfigurations{{URL: server.URL}}})
	restored := make(map[string]string)
	skipped, err := New(context.Background(), client).restorePlaylists(path, "user", map[string]api.BaseItemDto{"old-list": old}, map[string][]api.BaseItemDto{"Favorites": {current}}, map[string]string{"bob": "new-bob"}, restored)
	if err != nil || skipped != 0 {
		t.Fatalf("skipped=%d, err=%v", skipped, err)
	}
	if updated.GetIsPublic() || !updated.IsPublic.IsSet() || len(updated.GetUsers()) != 1 || updated.GetUsers()[0].GetUserId() != "new-bob" || updated.GetUsers()[0].GetCanEdit() || restored["old-list"] != "current-list" {
		t.Errorf("updated=%+v, restored=%v", updated, restored)
	}
}
