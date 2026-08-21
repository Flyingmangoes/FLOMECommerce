package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type Page[T any] struct {
	Items 		[]T 	`json:"items"`
	NextCursor 	*string	`json:"nextCursor" binding:"omitempty"`
	HasMore 	bool 	`json:"hasMore"`
	Total 		int 	`json:"total"`
}

type PagCursor struct {
	CreatedAt time.Time `json:"createdAt"`
	TargetID string `json:"targetId"`
}

type PagFilter struct {
	Cursor 	*PagCursor
	Limit 	int
}

type CursorExtractor[T any] func(item T)(createdAt time.Time, targetID string)

func(f *PagFilter) Normalize() {
	const maxLimit = 100; const defaultLimit = 5
	if f.Limit <= 0 || f.Limit >= maxLimit {
		f.Limit = defaultLimit
	}
}	

func(c *PagCursor) EncodeCursor() (string, error){
	b, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("Pagination Failed to encode cursor: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func DecodeCursor(s string) (*PagCursor, error){
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("Pagination Failed Invalid cursor encoding: %w", err)
	}

	var c PagCursor
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("Pagination Failed Invalid cursor format: %w", err)
	}

	return &c, nil
}

func(f *PagFilter) CursorValue()(createdAt interface{}, id interface{}){
	if f.Cursor == nil {
		return nil, nil
	}

	return f.Cursor.CreatedAt, f.Cursor.TargetID
}

func Build[T any](items []T, limit int, extract CursorExtractor[T]) (*Page[T], error) {
	page := &Page[T]{}

    if len(items) > limit {
        items = items[:limit]

        last := items[limit-1]
        createdAt, id := extract(last)
        
        cursor := &PagCursor{CreatedAt: createdAt, TargetID: id}
        encoded, err := cursor.EncodeCursor()
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