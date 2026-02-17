package maxapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetMe(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/me" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "token" {
			t.Fatalf("missing auth header, got %q", got)
		}
		_ = json.NewEncoder(w).Encode(BotInfo{UserID: 42, FirstName: "Echo", IsBot: true})
	}))
	defer srv.Close()

	client, err := NewClient("token", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	info, err := client.GetMe(context.Background())
	if err != nil {
		t.Fatalf("GetMe: %v", err)
	}
	if info.UserID != 42 || info.FirstName != "Echo" {
		t.Fatalf("unexpected response: %+v", info)
	}
}

func TestSendMessage(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method %s", r.Method)
		}
		values, _ := url.ParseQuery(r.URL.RawQuery)
		if values.Get("chat_id") != "99" {
			t.Fatalf("expected chat_id=99, got %s", values.Get("chat_id"))
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"message": map[string]any{"message_id": "mid"}})
	}))
	defer srv.Close()

	client, err := NewClient("token", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	chatID := int64(99)
	msg, err := client.SendMessage(context.Background(), SendMessageParams{ChatID: &chatID}, NewMessageBody{Text: "hi"})
	if err != nil {
		t.Fatalf("SendMessage: %v", err)
	}
	if msg.MessageID != "mid" {
		t.Fatalf("unexpected message id %s", msg.MessageID)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "bad things"})
	}))
	defer srv.Close()

	client, err := NewClient("token", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	_, err = client.GetMe(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("unexpected error type %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest || apiErr.Message != "bad things" {
		t.Fatalf("unexpected api error %+v", apiErr)
	}
}
