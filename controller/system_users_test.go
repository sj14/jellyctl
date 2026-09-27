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

func TestRestoreUserSettingsKeepsServerSpecificFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "user.json")
	backup := `{"Name":"alice","Policy":{"IsAdministrator":true,"IsHidden":true,"IsDisabled":false,` +
		`"AuthenticationProviderId":"old","PasswordResetProviderId":"old-reset"},` +
		`"Configuration":{"AudioLanguagePreference":"de","RememberSubtitleSelections":true,"OrderedViews":["old-view"]}}`
	if err := os.WriteFile(path, []byte(backup), 0600); err != nil {
		t.Fatal(err)
	}

	var policy api.UserPolicy
	var config api.UserConfiguration
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/Users/user":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"Name":"alice","Policy":{"IsAdministrator":false,"IsHidden":false,"IsDisabled":true,` +
				`"EnableMediaPlayback":true,"AuthenticationProviderId":"new","PasswordResetProviderId":"reset"},` +
				`"Configuration":{"AudioLanguagePreference":"en","OrderedViews":["new-view"]}}`))
		case r.Method == http.MethodPost && r.URL.Path == "/Users/user/Policy":
			if err := json.NewDecoder(r.Body).Decode(&policy); err != nil {
				t.Errorf("decode policy: %v", err)
			}
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPost && r.URL.Path == "/Users/Configuration":
			if r.URL.Query().Get("userId") != "user" {
				t.Errorf("configuration userId = %q", r.URL.Query().Get("userId"))
			}
			if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
				t.Errorf("decode configuration: %v", err)
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewAPIClient(&api.Configuration{Servers: api.ServerConfigurations{{URL: server.URL}}})
	if err := New(context.Background(), client).restoreUserSettings("user", "alice", path); err != nil {
		t.Fatal(err)
	}
	if !policy.GetIsAdministrator() || !policy.GetIsHidden() || policy.GetIsDisabled() || !policy.GetEnableMediaPlayback() || policy.AuthenticationProviderId != "new" {
		t.Errorf("policy = %+v", policy)
	}
	if config.GetAudioLanguagePreference() != "de" || !config.GetRememberSubtitleSelections() || !reflect.DeepEqual(config.OrderedViews, []string{"new-view"}) {
		t.Errorf("configuration = %+v", config)
	}
}

func TestRestoreUserSettingsRejectsWrongBackupUser(t *testing.T) {
	path := filepath.Join(t.TempDir(), "user.json")
	if err := os.WriteFile(path, []byte(`{"Name":"bob"}`), 0600); err != nil {
		t.Fatal(err)
	}
	err := New(context.Background(), nil).restoreUserSettings("user", "alice", path)
	if err == nil || !strings.Contains(err.Error(), "belongs to") {
		t.Fatalf("expected a user name mismatch, got %v", err)
	}
}
