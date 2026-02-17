package maxapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GetVideoMetadata fetches information about a video attachment token.
func (c *Client) GetVideoMetadata(ctx context.Context, token string) (*VideoMetadata, error) {
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("token is required")
	}
	path := fmt.Sprintf("/videos/%s", url.PathEscape(token))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	var meta VideoMetadata
	if err := c.do(req, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}
