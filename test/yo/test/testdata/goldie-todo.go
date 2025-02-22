// Code generated by fixgen, DO NOT EDIT.

package fixture

import (
	"context"
	"math/rand/v2"
	"testing"
	yo_gen "yo/models"

	"cloud.google.com/go/spanner"
	"github.com/samber/lo"
)

func CreateTodo(t *testing.T, db *spanner.Client, m *yo_gen.Todo, opts ...func(*yo_gen.Todo)) *yo_gen.Todo {
	t.Helper()

	tbl := &yo_gen.Todo{
		ID:          rand.Int64(),
		Title:       lo.RandomString(32, lo.AlphanumericCharset),
		Description: lo.RandomString(32, lo.AlphanumericCharset),
		// Tags is slice
		CreatedAt: spanner.CommitTimestamp,
		// UpdatedAt is nullable
		// DoneAt is nullable
	}

	if isModified(m.ID) {
		tbl.ID = m.ID
	}
	if isModified(m.Title) {
		tbl.Title = m.Title
	}
	if isModified(m.Description) {
		tbl.Description = m.Description
	}
	if len(m.Tags) > 0 {
		tbl.Tags = m.Tags
	}
	if isModified(m.CreatedAt) {
		t.Log("CreatedAt: spanner.CommitTimestamp should be used")
	}
	if !m.UpdatedAt.IsNull() {
		tbl.UpdatedAt = m.UpdatedAt
	}
	if !m.DoneAt.IsNull() {
		tbl.DoneAt = m.DoneAt
	}
	for _, opt := range opts {
		opt(tbl)
	}

	_, err := db.ReadWriteTransaction(context.Background(), func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		return tx.BufferWrite([]*spanner.Mutation{tbl.Insert(ctx)})
	})
	if err != nil {
		t.Fatal(err)
	}

	return tbl
}
