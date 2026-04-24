package grpcutil

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// PageRequest represents a generic page request.
// Mirrors the pagination package but for gRPC context.
type PageRequest struct {
	// Page is 1-based page number (for offset pagination).
	Page int32

	// PageSize is the number of items per page.
	PageSize int32

	// PageToken is an opaque cursor token (for cursor-based pagination).
	PageToken string
}

// PageResponse wraps paging metadata alongside a list response.
type PageResponse struct {
	// TotalItems is the total count of items across all pages.
	TotalItems int64

	// TotalPages is the total number of pages (offset pagination).
	TotalPages int32

	// CurrentPage is the current page number (offset pagination).
	CurrentPage int32

	// PageSize is the number of items returned.
	PageSize int32

	// NextPageToken is set when there are more results (cursor pagination).
	NextPageToken string

	// PrevPageToken is set when there is a previous page (cursor pagination).
	PrevPageToken string

	// HasMore indicates whether more pages exist.
	HasMore bool
}

// Offset returns the SQL/DB offset for the given page request.
func (r *PageRequest) Offset() int32 {
	if r.Page <= 1 {
		return 0
	}
	return (r.Page - 1) * r.PageSize
}

// Limit returns the effective page size, clamped to sane defaults.
func (r *PageRequest) Limit(maxSize int32) int32 {
	if r.PageSize <= 0 {
		return 20
	}
	if r.PageSize > maxSize {
		return maxSize
	}
	return r.PageSize
}

// NewPageResponse builds a PageResponse from offset-pagination params.
func NewPageResponse(page, pageSize int32, totalItems int64) PageResponse {
	totalPages := int32(0)
	if pageSize > 0 {
		total := int32(totalItems)
		totalPages = total / pageSize
		if total%pageSize != 0 {
			totalPages++
		}
	}
	return PageResponse{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
		HasMore:     int64(page)*int64(pageSize) < totalItems,
	}
}

// ---- Cursor token helpers ----

// cursorPayload is the data encoded inside a page token.
type cursorPayload struct {
	ID     interface{} `json:"id,omitempty"`
	Offset int64       `json:"offset,omitempty"`
	Extra  interface{} `json:"extra,omitempty"`
}

// EncodeCursor encodes an arbitrary value into a base64 page token.
func EncodeCursor(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("grpcutil: encode cursor: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// DecodeCursor decodes a base64 page token into the target value.
func DecodeCursor(token string, target interface{}) error {
	b, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return fmt.Errorf("grpcutil: decode cursor: %w", err)
	}
	if err := json.Unmarshal(b, target); err != nil {
		return fmt.Errorf("grpcutil: unmarshal cursor: %w", err)
	}
	return nil
}

// EncodeOffsetCursor encodes an integer offset into a page token.
func EncodeOffsetCursor(offset int64) (string, error) {
	return EncodeCursor(cursorPayload{Offset: offset})
}

// DecodeOffsetCursor decodes an offset-based page token.
func DecodeOffsetCursor(token string) (int64, error) {
	if token == "" {
		return 0, nil
	}
	var p cursorPayload
	if err := DecodeCursor(token, &p); err != nil {
		return 0, err
	}
	return p.Offset, nil
}
