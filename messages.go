package maxapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ListMessagesParams configures GET /messages.
type ListMessagesParams struct {
	ChatID     *int64
	MessageIDs []string
	From       *int64
	To         *int64
	Count      int
}

// ListMessages returns a slice of messages filtered by chat or IDs.
func (c *Client) ListMessages(ctx context.Context, params ListMessagesParams) (*MessageList, error) {
	query := url.Values{}
	if params.ChatID != nil {
		query.Set("chat_id", fmt.Sprint(*params.ChatID))
	}
	if len(params.MessageIDs) > 0 {
		query.Set("message_ids", strings.Join(params.MessageIDs, ","))
	}
	if params.From != nil {
		query.Set("from", fmt.Sprint(*params.From))
	}
	if params.To != nil {
		query.Set("to", fmt.Sprint(*params.To))
	}
	if params.Count > 0 {
		query.Set("count", fmt.Sprint(params.Count))
	}
	req, err := c.newRequest(ctx, http.MethodGet, "/messages", query, nil)
	if err != nil {
		return nil, err
	}
	var resp MessageList
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SendMessageParams configure recipient + flags.
type SendMessageParams struct {
	ChatID             *int64
	UserID             *int64
	DisableLinkPreview bool
}

// SendMessage sends a new message to a chat or user.
func (c *Client) SendMessage(ctx context.Context, params SendMessageParams, body NewMessageBody) (*Message, error) {
	query := url.Values{}
	if params.ChatID != nil {
		query.Set("chat_id", fmt.Sprint(*params.ChatID))
	}
	if params.UserID != nil {
		query.Set("user_id", fmt.Sprint(*params.UserID))
	}
	if params.DisableLinkPreview {
		query.Set("disable_link_preview", "true")
	}
	if len(query) == 0 {
		return nil, fmt.Errorf("chat_id or user_id must be provided")
	}
	req, err := c.newRequest(ctx, http.MethodPost, "/messages", query, body)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Message Message `json:"message"`
	}
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Message, nil
}

// EditMessage updates message content by ID.
func (c *Client) EditMessage(ctx context.Context, messageID string, body NewMessageBody) (*SuccessResponse, error) {
	if strings.TrimSpace(messageID) == "" {
		return nil, fmt.Errorf("messageID is required")
	}
	query := url.Values{}
	query.Set("message_id", messageID)
	req, err := c.newRequest(ctx, http.MethodPut, "/messages", query, body)
	if err != nil {
		return nil, err
	}
	var resp SuccessResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteMessage removes a message younger than 24h.
func (c *Client) DeleteMessage(ctx context.Context, messageID string) (*SuccessResponse, error) {
	if strings.TrimSpace(messageID) == "" {
		return nil, fmt.Errorf("messageID is required")
	}
	query := url.Values{}
	query.Set("message_id", messageID)
	req, err := c.newRequest(ctx, http.MethodDelete, "/messages", query, nil)
	if err != nil {
		return nil, err
	}
	var resp SuccessResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetMessage returns a single message by ID.
func (c *Client) GetMessage(ctx context.Context, messageID string) (*Message, error) {
	if strings.TrimSpace(messageID) == "" {
		return nil, fmt.Errorf("messageID is required")
	}
	path := fmt.Sprintf("/messages/%s", url.PathEscape(messageID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	var msg Message
	if err := c.do(req, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
