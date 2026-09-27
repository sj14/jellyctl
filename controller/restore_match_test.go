package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sj14/jellyfin-go/api"
)

func TestFindRestoreItem(t *testing.T) {
	decode := func(data string) api.BaseItemDto {
		t.Helper()
		var item api.BaseItemDto
		if err := json.Unmarshal([]byte(data), &item); err != nil {
			t.Fatal(err)
		}
		return item
	}
	backup := decode(`{"Name":"Pilot","Type":"Episode","ProductionYear":2020,"SeriesName":"Series A","ParentIndexNumber":1,"IndexNumber":1,"ProviderIds":{"tvdb":"123"}}`)
	candidates := []api.BaseItemDto{
		decode(`{"Id":"prefix","Name":"Pilot Part 2","Type":"Episode","ProductionYear":2020,"SeriesName":"Series A","ParentIndexNumber":1,"IndexNumber":1}`),
		decode(`{"Id":"type","Name":"Pilot","Type":"Movie","ProductionYear":2020,"SeriesName":"Series A","ParentIndexNumber":1,"IndexNumber":1}`),
		decode(`{"Id":"year","Name":"Pilot","Type":"Episode","ProductionYear":2021,"SeriesName":"Series A","ParentIndexNumber":1,"IndexNumber":1}`),
		decode(`{"Id":"series","Name":"Pilot","Type":"Episode","ProductionYear":2020,"SeriesName":"Series B","ParentIndexNumber":1,"IndexNumber":1}`),
		decode(`{"Id":"episode","Name":"Pilot","Type":"Episode","ProductionYear":2020,"SeriesName":"Series A","ParentIndexNumber":1,"IndexNumber":2}`),
		decode(`{"Id":"provider","Name":"Pilot","Type":"Episode","ProductionYear":2020,"SeriesName":"Series A","ParentIndexNumber":1,"IndexNumber":1,"ProviderIds":{"tvdb":"456"}}`),
		decode(`{"Id":"right","Name":"Pilot","Type":"Episode","ProductionYear":2020,"SeriesName":"Series A","ParentIndexNumber":1,"IndexNumber":1,"ProviderIds":{"tvdb":"123"}}`),
	}

	match, count := findRestoreItem(backup, candidates)
	if count != 1 || match == nil || match.GetId() != "right" {
		t.Fatalf("match=%v, count=%d; want right", match, count)
	}
	match, count = findRestoreItem(backup, append(candidates, candidates[len(candidates)-1]))
	if match != nil || count != 2 {
		t.Fatalf("ambiguous match=%v, count=%d; want nil and 2", match, count)
	}
	match, count = findRestoreItem(backup, candidates[:len(candidates)-1])
	if match != nil || count != 0 {
		t.Fatalf("missing match=%v, count=%d; want nil and 0", match, count)
	}
}

func TestSystemRestoreUsesOnlyUniqueExactMatch(t *testing.T) {
	for _, tt := range []struct {
		name       string
		candidates string
		wantPlayed string
		wantError  bool
	}{
		{
			name: "unique match",
			candidates: `[{"Id":"prefix","Name":"Pilot Part 2","Type":"Episode","ProductionYear":2020},` +
				`{"Id":"right","Name":"Pilot","Type":"Episode","ProductionYear":2020}]`,
			wantPlayed: "right",
		},
		{
			name: "ambiguous match",
			candidates: `[{"Id":"one","Name":"Pilot","Type":"Episode","ProductionYear":2020},` +
				`{"Id":"two","Name":"Pilot","Type":"Episode","ProductionYear":2020}]`,
			wantError: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			backupDir := t.TempDir()
			userDir := filepath.Join(backupDir, "users", "alice")
			if err := os.MkdirAll(userDir, 0700); err != nil {
				t.Fatal(err)
			}
			backup := `[{"Name":"Pilot","Type":"Episode","ProductionYear":2020,"UserData":{"Key":"pilot","Played":true}}]`
			if err := os.WriteFile(filepath.Join(userDir, "items.json"), []byte(backup), 0600); err != nil {
				t.Fatal(err)
			}

			var played []string
			itemsRequests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case r.Method == http.MethodGet && r.URL.Path == "/Users":
					_, _ = w.Write([]byte(`[{"Id":"user","Name":"alice"}]`))
				case r.Method == http.MethodGet && r.URL.Path == "/Items":
					itemsRequests++
					if r.URL.Query().Get("recursive") != "true" {
						t.Errorf("unexpected item query: %s", r.URL.RawQuery)
					}
					_, _ = w.Write([]byte(`{"Items":` + tt.candidates + `,"TotalRecordCount":2}`))
				case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/UserPlayedItems/"):
					played = append(played, strings.TrimPrefix(r.URL.Path, "/UserPlayedItems/"))
					_, _ = w.Write([]byte(`{"Key":"pilot"}`))
				default:
					http.NotFound(w, r)
				}
			}))
			defer server.Close()

			client := api.NewAPIClient(&api.Configuration{Servers: api.ServerConfigurations{{URL: server.URL}}})
			err := New(context.Background(), client).SystemRestore(backupDir, false, false)
			if tt.wantError {
				if err == nil || !strings.Contains(err.Error(), "restore incomplete: 1") {
					t.Fatalf("expected incomplete restore error, got %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			var wantPlayed []string
			if tt.wantPlayed != "" {
				wantPlayed = []string{tt.wantPlayed}
			}
			if !reflect.DeepEqual(played, wantPlayed) {
				t.Fatalf("played IDs = %v, want %v", played, wantPlayed)
			}
			if itemsRequests != 1 {
				t.Errorf("library fetched %d times, want once", itemsRequests)
			}
		})
	}
}
