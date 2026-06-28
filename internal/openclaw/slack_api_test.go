package openclaw

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNormalizeSlackChannelName(t *testing.T) {
	if got := normalizeSlackChannelName("#General"); got != "general" {
		t.Fatalf("normalizeSlackChannelName() = %q, want general", got)
	}
}

func TestResolveSlackTestTargetFromConfig(t *testing.T) {
	target, ok := resolveSlackTestTargetFromConfig(slackConfigView{testTarget: "C111"})
	if !ok || target != "C111" {
		t.Fatalf("resolveSlackTestTargetFromConfig() = (%q, %v)", target, ok)
	}

	target, ok = resolveSlackTestTargetFromConfig(slackConfigView{channelIDs: []string{"*", "C222"}})
	if !ok || target != "C222" {
		t.Fatalf("resolveSlackTestTargetFromConfig() = (%q, %v)", target, ok)
	}

	_, ok = resolveSlackTestTargetFromConfig(slackConfigView{})
	if ok {
		t.Fatal("resolveSlackTestTargetFromConfig() ok = true, want false")
	}
}

func TestLookupSlackChannelByName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/conversations.list" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("types"); got != "public_channel" {
			t.Fatalf("types = %q, want public_channel", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok": true,
			"channels": []map[string]any{
				{"id": "C111", "name": "random", "is_archived": false, "is_member": true},
				{"id": "C222", "name": "general", "is_archived": false, "is_member": true, "is_general": true},
			},
			"response_metadata": map[string]any{"next_cursor": ""},
		})
	}))
	defer server.Close()

	origURL := slackConversationsListURL
	origClient := slackHTTPClient
	t.Cleanup(func() {
		slackConversationsListURL = origURL
		slackHTTPClient = origClient
	})
	slackConversationsListURL = server.URL + "/conversations.list"
	slackHTTPClient = server.Client()

	channel, err := lookupSlackChannelByName("xoxb-test", "#general")
	if err != nil {
		t.Fatalf("lookupSlackChannelByName() error = %v", err)
	}
	if channel.ID != "C222" || channel.Name != "general" {
		t.Fatalf("channel = %+v", channel)
	}
}

func TestSendSlackChannelMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat.postMessage" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload["channel"] != "C222" || !strings.Contains(payload["text"], slackTestMessagePrefix) {
			t.Fatalf("payload = %#v", payload)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "channel": "C222", "ts": "1"})
	}))
	defer server.Close()

	origURL := slackChatPostMessageURL
	origClient := slackHTTPClient
	t.Cleanup(func() {
		slackChatPostMessageURL = origURL
		slackHTTPClient = origClient
	})
	slackChatPostMessageURL = server.URL + "/chat.postMessage"
	slackHTTPClient = server.Client()

	if err := sendSlackChannelMessage("xoxb-test", "C222", slackTestMessagePrefix); err != nil {
		t.Fatalf("sendSlackChannelMessage() error = %v", err)
	}
}

func TestResolveSlackTestTargetDefaultsToGeneral(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok": true,
			"channels": []map[string]any{
				{"id": "C999", "name": "general", "is_archived": false, "is_member": true, "is_general": true},
			},
			"response_metadata": map[string]any{"next_cursor": ""},
		})
	}))
	defer server.Close()

	origURL := slackConversationsListURL
	origClient := slackHTTPClient
	t.Cleanup(func() {
		slackConversationsListURL = origURL
		slackHTTPClient = origClient
	})
	slackConversationsListURL = server.URL + "/conversations.list"
	slackHTTPClient = server.Client()

	targetID, label, err := resolveSlackTestTarget(slackConfigView{}, "xoxb-test", "")
	if err != nil {
		t.Fatalf("resolveSlackTestTarget() error = %v", err)
	}
	if targetID != "C999" || label != "#general" {
		t.Fatalf("resolveSlackTestTarget() = (%q, %q)", targetID, label)
	}
}

func TestResolveSlackTestTargetSlackChannelOverride(t *testing.T) {
	targetID, label, err := resolveSlackTestTarget(slackConfigView{testTarget: "C111"}, "xoxb-test", "C222")
	if err != nil {
		t.Fatalf("resolveSlackTestTarget() error = %v", err)
	}
	if targetID != "C222" || label != "C222" {
		t.Fatalf("resolveSlackTestTarget() = (%q, %q)", targetID, label)
	}
}

func TestFormatSlackAPIErrorMissingScope(t *testing.T) {
	err := formatSlackAPIError("conversations.list", slackAPIResponse{
		Error:    "missing_scope",
		Needed:   "groups:read",
		Provided: "channels:read,chat:write",
	})
	msg := err.Error()
	if !strings.Contains(msg, "groups:read") || !strings.Contains(msg, "channels:read,chat:write") {
		t.Fatalf("error = %q", msg)
	}
}