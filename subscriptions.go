package maxapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// CreateSubscriptionRequest registers a webhook endpoint.
type CreateSubscriptionRequest struct {
	URL          string   `json:"url"`
	UpdateTypes  []string `json:"update_types,omitempty"`
	Secret       string   `json:"secret,omitempty"`
	AllowedPorts []int    `json:"allowed_ports,omitempty"`
}

// ListSubscriptions returns configured webhooks.
func (c *Client) ListSubscriptions(ctx context.Context) (*SubscriptionsResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/subscriptions", nil, nil)
	if err != nil {
		return nil, err
	}
	var resp SubscriptionsResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateSubscription registers a webhook.
func (c *Client) CreateSubscription(ctx context.Context, payload CreateSubscriptionRequest) (*SuccessResponse, error) {
	if strings.TrimSpace(payload.URL) == "" {
		return nil, fmt.Errorf("url is required")
	}
	req, err := c.newRequest(ctx, http.MethodPost, "/subscriptions", nil, payload)
	if err != nil {
		return nil, err
	}
	var resp SuccessResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteSubscription removes webhook by URL.
func (c *Client) DeleteSubscription(ctx context.Context, webhookURL string) (*SuccessResponse, error) {
	if strings.TrimSpace(webhookURL) == "" {
		return nil, fmt.Errorf("url is required")
	}
	query := url.Values{}
	query.Set("url", webhookURL)
	req, err := c.newRequest(ctx, http.MethodDelete, "/subscriptions", query, nil)
	if err != nil {
		return nil, err
	}
	var resp SuccessResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdatesParams configures GET /updates.
type UpdatesParams struct {
	Limit   int
	Timeout int
	Marker  *int64
	Types   []string
}

// GetUpdates performs long polling.
func (c *Client) GetUpdates(ctx context.Context, params *UpdatesParams) (*UpdatePage, error) {
	query := url.Values{}
	if params != nil {
		if params.Limit > 0 {
			query.Set("limit", fmt.Sprint(params.Limit))
		}
		if params.Timeout > 0 {
			query.Set("timeout", fmt.Sprint(params.Timeout))
		}
		if params.Marker != nil {
			query.Set("marker", fmt.Sprint(*params.Marker))
		}
		if len(params.Types) > 0 {
			query.Set("types", strings.Join(params.Types, ","))
		}
	}
	req, err := c.newRequest(ctx, http.MethodGet, "/updates", query, nil)
	if err != nil {
		return nil, err
	}
	var page UpdatePage
	if err := c.do(req, &page); err != nil {
		return nil, err
	}
	return &page, nil
}
