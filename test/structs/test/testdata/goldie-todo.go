// Code generated by fixgen, DO NOT EDIT.

package fixture

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/samber/lo"
)

func CreateTodo(t *testing.T, db Inserter, m *structs_gen.Todo, opts ...func(*structs_gen.Todo)) *yo_gen.Todo {
	t.Helper()

	tbl := &structs_gen.Todo{
		ID:          rand.Int64(),
		Title:       lo.RandomString(32, lo.AlphanumericCharset),
		Description: lo.RandomString(32, lo.AlphanumericCharset),
		// Tags is slice
		CreatedAt: time.Now(),
		// UpdatedAt is unknown
		// DoneAt is unknown
	}

	if isModified(m.ID) {
		tbl.ID = m.ID
	}
	tbl.Title = m.Title // must overwrite
	if isModified(m.Description) {
		tbl.Description = m.Description
	}
	if len(m.Tags) > 0 {
		tbl.Tags = m.Tags
	}
	if isModified(m.CreatedAt) {
		tbl.CreatedAt = m.CreatedAt
	}
	if isModified(m.UpdatedAt) {
		tbl.UpdatedAt = m.UpdatedAt
	}
	if isModified(m.DoneAt) {
		tbl.DoneAt = m.DoneAt
	}
	for _, opt := range opts {
		opt(tbl)
	}

	m, err := db.Insert(context.Background(), tbl)
	if err != nil {
		t.Fatal(err)
	}

	return m
}
