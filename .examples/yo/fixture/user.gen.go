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

func CreateUser(t *testing.T, db *spanner.Client, m *yo_gen.User, opts ...func(*yo_gen.User)) *yo_gen.User {
	t.Helper()

	tbl := &yo_gen.User{
		ID:        rand.Int64(),
		Name:      lo.RandomString(32, lo.AlphanumericCharset),
		IconURL:   lo.RandomString(32, lo.AlphanumericCharset),
		CreatedAt: spanner.CommitTimestamp,
		// UpdatedAt is nullable
	}

	if isModified(m.ID) {
		tbl.ID = m.ID
	}
	if isModified(m.Name) {
		tbl.Name = m.Name
	}
	if isModified(m.IconURL) {
		tbl.IconURL = m.IconURL
	}
	if isModified(m.CreatedAt) {
		t.Log("CreatedAt: spanner.CommitTimestamp should be used")
	}
	if !m.UpdatedAt.IsNull() {
		tbl.UpdatedAt = m.UpdatedAt
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
