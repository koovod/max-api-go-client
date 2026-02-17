package maxapi

import (
	"context"
	"net/http"
)

// GetMe returns information about the current bot token.
func (c *Client) GetMe(ctx context.Context) (*BotInfo, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/me", nil, nil)
	if err != nil {
		return nil, err
	}
	var resp BotInfo
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
