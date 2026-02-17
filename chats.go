package maxapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ListChatsParams configures pagination for ListChats.
type ListChatsParams struct {
	Count  int
	Marker *int64
}

// ListChats returns chats where the bot interacted.
func (c *Client) ListChats(ctx context.Context, params *ListChatsParams) (*ChatsPage, error) {
	query := url.Values{}
	if params != nil {
		if params.Count > 0 {
			query.Set("count", fmt.Sprint(params.Count))
		}
		if params.Marker != nil {
			query.Set("marker", fmt.Sprint(*params.Marker))
		}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/chats", query, nil)
	if err != nil {
		return nil, err
	}

	var resp ChatsPage
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetChat returns a detailed chat object.
func (c *Client) GetChat(ctx context.Context, chatID int64) (*Chat, error) {
	path := fmt.Sprintf("/chats/%d", chatID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	var chat Chat
	if err := c.do(req, &chat); err != nil {
		return nil, err
	}
	return &chat, nil
}

// UpdateChatRequest mutates chat metadata.
type UpdateChatRequest struct {
	Title  string `json:"title,omitempty"`
	Icon   *Image `json:"icon,omitempty"`
	Notify *bool  `json:"notify,omitempty"`
}

// UpdateChat edits chat title/icon/notifications.
func (c *Client) UpdateChat(ctx context.Context, chatID int64, payload UpdateChatRequest) (*Chat, error) {
	path := fmt.Sprintf("/chats/%d", chatID)
	req, err := c.newRequest(ctx, http.MethodPatch, path, nil, payload)
	if err != nil {
		return nil, err
	}

	var chat Chat
	if err := c.do(req, &chat); err != nil {
		return nil, err
	}
	return &chat, nil
}

// DeleteChat removes (closes) the chat for all participants.
func (c *Client) DeleteChat(ctx context.Context, chatID int64) (*SuccessResponse, error) {
	path := fmt.Sprintf("/chats/%d", chatID)
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return nil, err
	}
	var resp SuccessResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SendChatAction posts typing_on/sending_* or mark_seen.
func (c *Client) SendChatAction(ctx context.Context, chatID int64, action ChatActionRequest) (*SuccessResponse, error) {
	path := fmt.Sprintf("/chats/%d/actions", chatID)
	req, err := c.newRequest(ctx, http.MethodPost, path, nil, action)
	if err != nil {
		return nil, err
	}
	var resp SuccessResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetChatPin returns the pinned message or nil.
func (c *Client) GetChatPin(ctx context.Context, chatID int64) (*Message, error) {
	path := fmt.Sprintf("/chats/%d/pin", chatID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	var message Message
	if err := c.do(req, &message); err != nil {
		return nil, err
	}
	if message.MessageID == "" && message.Body == nil && message.Sender == nil {
		return nil, nil
	}
	return &message, nil
}

// PinChatMessage pins an existing message by ID.
func (c *Client) PinChatMessage(ctx context.Context, chatID int64, payload PinMessageRequest) (*SuccessResponse, error) {
	path := fmt.Sprintf("/chats/%d/pin", chatID)
	req, err := c.newRequest(ctx, http.MethodPut, path, nil, payload)
	if err != nil {
		return nil, err
	}
	var resp SuccessResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UnpinChatMessage removes the pinned message.
func (c *Client) UnpinChatMessage(ctx context.Context, chatID int64) (*SuccessResponse, error) {
	path := fmt.Sprintf("/chats/%d/pin", chatID)
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return nil, err
	}
	var resp SuccessResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
