package maxapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// CallbackNotification shows a toast/alert to the user.
type CallbackNotification struct {
	Text      string `json:"text"`
	ShowAlert bool   `json:"show_alert,omitempty"`
	URL       string `json:"url,omitempty"`
}

// AnswerCallbackRequest either updates message content or sends a notification.
type AnswerCallbackRequest struct {
	Message      *NewMessageBody       `json:"message,omitempty"`
	Notification *CallbackNotification `json:"notification,omitempty"`
}

// AnswerCallback responds to inline keyboard callbacks.
func (c *Client) AnswerCallback(ctx context.Context, callbackID string, payload AnswerCallbackRequest) (*SuccessResponse, error) {
	if strings.TrimSpace(callbackID) == "" {
		return nil, fmt.Errorf("callbackID is required")
	}
	query := url.Values{}
	query.Set("callback_id", callbackID)
	req, err := c.newRequest(ctx, http.MethodPost, "/answers", query, payload)
	if err != nil {
		return nil, err
	}
	var resp SuccessResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
