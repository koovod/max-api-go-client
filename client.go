package maxapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultBaseURL = "https://platform-api.max.ru"

// Client is a typed HTTP wrapper around the MAX Platform bot API.
type Client struct {
	token      string
	baseURL    *url.URL
	httpClient *http.Client
}

// Option changes the default client behaviour.
type Option func(*Client) error

// WithHTTPClient overrides the HTTP client used for requests.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) error {
		if hc == nil {
			return fmt.Errorf("http client cannot be nil")
		}
		c.httpClient = hc
		return nil
	}
}

// WithBaseURL overrides the API host. Useful for tests.
func WithBaseURL(raw string) Option {
	return func(c *Client) error {
		if raw == "" {
			return fmt.Errorf("base URL cannot be empty")
		}
		u, err := url.Parse(raw)
		if err != nil {
			return fmt.Errorf("parse base url: %w", err)
		}
		if u.Scheme == "" {
			u.Scheme = "https"
		}
		c.baseURL = u
		return nil
	}
}

// NewClient constructs a Client with the provided access token.
func NewClient(token string, opts ...Option) (*Client, error) {
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("token is required")
	}

	base, err := url.Parse(defaultBaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse default base url: %w", err)
	}

	c := &Client{
		token:      token,
		baseURL:    base,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, query url.Values, body any) (*http.Request, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if c.baseURL == nil {
		return nil, fmt.Errorf("client base URL is not configured")
	}

	u := *c.baseURL
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u.Path = strings.TrimSuffix(c.baseURL.Path, "/") + path
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, fmt.Errorf("encode request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *Client) do(req *http.Request, out any) error {
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("perform request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return decodeAPIError(resp.StatusCode, body)
	}

	if out == nil {
		return nil
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

// APIError represents a non-2xx response from the MAX API.
type APIError struct {
	StatusCode int
	Message    string
	Body       []byte
}

func (e *APIError) Error() string {
	msg := e.Message
	if msg == "" {
		msg = http.StatusText(e.StatusCode)
	}
	return fmt.Sprintf("max api error: status=%d message=%s", e.StatusCode, msg)
}

func decodeAPIError(code int, body []byte) error {
	errResp := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Error   string `json:"error"`
	}{}
	_ = json.Unmarshal(body, &errResp)

	msg := strings.TrimSpace(errResp.Message)
	if msg == "" {
		msg = strings.TrimSpace(errResp.Error)
	}
	if msg == "" && len(body) > 0 {
		msg = string(body)
	}
	return &APIError{StatusCode: code, Message: msg, Body: body}
}
