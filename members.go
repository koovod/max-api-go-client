package maxapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GetMyChatMember returns member info for the bot in a chat.
func (c *Client) GetMyChatMember(ctx context.Context, chatID int64) (*ChatMember, error) {
	path := fmt.Sprintf("/chats/%d/members/me", chatID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	var member ChatMember
	if err := c.do(req, &member); err != nil {
		return nil, err
	}
	return &member, nil
}

// LeaveChat removes the bot from the chat.
func (c *Client) LeaveChat(ctx context.Context, chatID int64) (*SuccessResponse, error) {
	path := fmt.Sprintf("/chats/%d/members/me", chatID)
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

// ListChatAdmins returns admin members for chat.
func (c *Client) ListChatAdmins(ctx context.Context, chatID int64, marker *int64) (*ChatMembersPage, error) {
	query := url.Values{}
	if marker != nil {
		query.Set("marker", fmt.Sprint(*marker))
	}
	path := fmt.Sprintf("/chats/%d/members/admins", chatID)
	req, err := c.newRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, err
	}
	var page ChatMembersPage
	if err := c.do(req, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// SetChatAdmins assigns admins with permission sets.
func (c *Client) SetChatAdmins(ctx context.Context, chatID int64, payload SetAdminsRequest) (*SuccessResponse, error) {
	path := fmt.Sprintf("/chats/%d/members/admins", chatID)
	req, err := c.newRequest(ctx, http.MethodPost, path, nil, payload)
	if err != nil {
		return nil, err
	}
	var resp SuccessResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RemoveChatAdmin revokes admin rights.
func (c *Client) RemoveChatAdmin(ctx context.Context, chatID int64, userID int64) (*SuccessResponse, error) {
	path := fmt.Sprintf("/chats/%d/members/admins/%d", chatID, userID)
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

// ListMembersParams configures GET /members.
type ListMembersParams struct {
	UserIDs []int64
	Marker  *int64
	Count   int
}

// ListChatMembers returns page of members or filtered user IDs.
func (c *Client) ListChatMembers(ctx context.Context, chatID int64, params *ListMembersParams) (*ChatMembersPage, error) {
	query := url.Values{}
	if params != nil {
		if len(params.UserIDs) > 0 {
			var ids []string
			for _, id := range params.UserIDs {
				ids = append(ids, fmt.Sprint(id))
			}
			query.Set("user_ids", strings.Join(ids, ","))
		}
		if params.Marker != nil {
			query.Set("marker", fmt.Sprint(*params.Marker))
		}
		if params.Count > 0 {
			query.Set("count", fmt.Sprint(params.Count))
		}
	}
	path := fmt.Sprintf("/chats/%d/members", chatID)
	req, err := c.newRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, err
	}
	var page ChatMembersPage
	if err := c.do(req, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// AddChatMembers adds members by ID.
func (c *Client) AddChatMembers(ctx context.Context, chatID int64, userIDs []int64) (*SuccessResponse, error) {
	if len(userIDs) == 0 {
		return nil, fmt.Errorf("userIDs cannot be empty")
	}
	payload := AddMembersRequest{UserIDs: userIDs}
	path := fmt.Sprintf("/chats/%d/members", chatID)
	req, err := c.newRequest(ctx, http.MethodPost, path, nil, payload)
	if err != nil {
		return nil, err
	}
	var resp SuccessResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RemoveChatMember removes or blocks a user from the chat.
func (c *Client) RemoveChatMember(ctx context.Context, chatID int64, userID int64, block bool) (*SuccessResponse, error) {
	path := fmt.Sprintf("/chats/%d/members", chatID)
	query := url.Values{}
	query.Set("user_id", fmt.Sprint(userID))
	if block {
		query.Set("block", "true")
	}
	req, err := c.newRequest(ctx, http.MethodDelete, path, query, nil)
	if err != nil {
		return nil, err
	}
	var resp SuccessResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
