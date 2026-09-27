package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sj14/jellyfin-go/api"
)

func TestUserPolicyMissing(t *testing.T) {
	updates := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/Users/user":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"Id":"user"}`))
		case r.Method == http.MethodPost && r.URL.Path == "/Users/user/Policy":
			updates++
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewAPIClient(&api.Configuration{Servers: api.ServerConfigurations{{URL: server.URL}}})
	if err := New(context.Background(), client).UserDisable("user"); err == nil || !strings.Contains(err.Error(), "has no policy") {
		t.Fatalf("expected a missing-policy error, got %v", err)
	}
	if updates != 0 {
		t.Fatalf("sent %d policy updates without a policy", updates)
	}
}

func TestUserPolicyUpdatePreservesExistingFields(t *testing.T) {
	var updated api.UserPolicy
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/Users/user":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"Id":"user","Policy":{"IsDisabled":true,"IsAdministrator":false,"AuthenticationProviderId":"provider","PasswordResetProviderId":"reset"}}`))
		case r.Method == http.MethodPost && r.URL.Path == "/Users/user/Policy":
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
	if err := New(context.Background(), client).UserSetAdmin("user"); err != nil {
		t.Fatal(err)
	}
	if !updated.GetIsAdministrator() || !updated.GetIsDisabled() || updated.AuthenticationProviderId != "provider" || updated.PasswordResetProviderId != "reset" {
		t.Errorf("policy fields changed unexpectedly: %+v", updated)
	}
}
