package util

import (
	"encoding/json"
	"strings"
)

type Cursor[V any] struct {
	Value V
	ID    string
}

func CursorFromString[V any](cursor *string) *Cursor[V] {
	if cursor == nil {
		return nil
	}

	parts := strings.Split(*cursor, "::")
	if len(parts) != 2 {
		return nil
	}

	var value V
	err := json.Unmarshal([]byte(parts[0]), &value)
	if err != nil {
		return nil
	}

	return &Cursor[V]{
		Value: value,
		ID:    parts[1],
	}
}

func CursorToString[V any](cursor *Cursor[V]) *string {
	if cursor == nil {
		return nil
	}

	value, err := json.Marshal(cursor.Value)
	if err != nil {
		return nil
	}

	str := string(value) + "::" + cursor.ID

	return &str
}

type Page[I any, C any] struct {
	Items []I
	Next  *Cursor[C]
}

type Pagination[C any] struct {
	Limit int
	From  *Cursor[C]
}
