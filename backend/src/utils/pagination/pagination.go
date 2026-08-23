package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type PagFilter struct {
	Limit  int
	Cursor *PagCursor
}

type Page[T any] struct {
	Items      []T     `json:"items"`
	NextCursor *string `json:"next_cursor,omitempty"`
	HasMore    bool    `json:"has_more"`
	TotalCount int     `json:"total_count"`
}

type PagCursor struct {
	TargetID  string `json:"target_id"`
	CreatedAt int64  `json:"created_at"`
}

type CursorExtractor[T any] func(item T) (createdAt time.Time, targetID string)

func (f *PagFilter) Normalize() {
	const defaultLimit = 20
	const defaultMaxLimit = 100
	if f.Limit <= 0 || f.Limit > defaultMaxLimit {
		f.Limit = defaultLimit
	}
}

func (f *PagFilter) CursorValue() (createdAt interface{}, targetID interface{}) {
	if f.Cursor == nil {
		return nil, nil
	}

	return f.Cursor.CreatedAt, f.Cursor.TargetID
}

func (c *PagCursor) EncodeCursor() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func DecodeCursor(encoded string) (*PagCursor, error) {
	b, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("Pagination Failed Invalid cursor encoding: %w", err)
	}

	var c PagCursor
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("Pagination Failed Invalid cursor format: %w", err)
	}

	return &c, nil
}

func Build[T any](items []T, limit int, extractor CursorExtractor[T]) (Page[T], error) {
	page := Page[T]{}

	return page, nil
}
