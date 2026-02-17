package maxapi

import (
	"context"
	"fmt"
	"net/http"
)

// UploadType enumerates supported media categories.
type UploadType string

const (
	UploadTypeImage UploadType = "image"
	UploadTypeVideo UploadType = "video"
	UploadTypeAudio UploadType = "audio"
	UploadTypeFile  UploadType = "file"
)

// uploadRequest is used internally.
type uploadRequest struct {
	Type UploadType `json:"type"`
}

// CreateUpload creates an upload session for attachments.
func (c *Client) CreateUpload(ctx context.Context, uploadType UploadType) (*UploadResponse, error) {
	if uploadType == "" {
		return nil, fmt.Errorf("upload type is required")
	}
	payload := uploadRequest{Type: uploadType}
	req, err := c.newRequest(ctx, http.MethodPost, "/uploads", nil, payload)
	if err != nil {
		return nil, err
	}
	var resp UploadResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
