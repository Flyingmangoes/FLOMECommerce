package utils

import (
	"encoding/base64"
	"encoding/json"
	"time"
	"fmt"
)

type PagCursor struct {
    CreatedAt time.Time `json:"created_at"`
    ID        string    `json:"id"`
}

const (
    DefaultLimit = 20
    MaxLimit     = 100
)

type PagFilter struct {
    Cursor *PagCursor
    Limit  int
}

type CursorExtractor[T any] func(item T) (createdAt time.Time, id string)

type Page[T any] struct {
    Items      []T     `json:"items"`
    NextCursor *string `json:"next_cursor,omitempty"` // nil = last page
    HasMore    bool    `json:"has_more"`
    Total      int     `json:"total"` // count of items in this page
}

func (c *PagCursor) Encode() (string, error) {
    b, err := json.Marshal(c)
    if err != nil {
        return "", fmt.Errorf("pagination: failed to encode cursor: %w", err)
    }
    return base64.URLEncoding.EncodeToString(b), nil
}

func DecodeCursor(s string) (*PagCursor, error) {
    b, err := base64.URLEncoding.DecodeString(s)
    if err != nil {
        return nil, fmt.Errorf("pagination: invalid cursor encoding: %w", err)
    }

    var c PagCursor
    if err := json.Unmarshal(b, &c); err != nil {
        return nil, fmt.Errorf("pagination: invalid cursor format: %w", err)
    }

    return &c, nil
}

func (f *PagFilter) Normalize() {
    if f.Limit <= 0 || f.Limit > MaxLimit {
        f.Limit = DefaultLimit
    }
}

func (f *PagFilter) CursorValues() (createdAt interface{}, id interface{}) {
    if f.Cursor == nil {
        return nil, nil
    }
    return f.Cursor.CreatedAt, f.Cursor.ID
}

func Build[T any](items []T, limit int, extract CursorExtractor[T]) (*Page[T], error) {
    page := &Page[T]{}

    if len(items) > limit {
        items = items[:limit]

        last := items[limit-1]
        createdAt, id := extract(last)
        
        cursor := &PagCursor{CreatedAt: createdAt, ID: id}
        encoded, err := cursor.Encode()
        if err != nil {
            return nil, err
        }
        
        page.NextCursor = &encoded
        page.HasMore = true
    }

    page.Items = items
    page.Total = len(items)
    return page, nil
}