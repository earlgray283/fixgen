// Code generated by fixgen, DO NOT EDIT.

package fixture

import (
	"context"
	"math/rand/v2"
	"testing"

	ent_gen "ent/ent"

	"github.com/samber/lo"
)

func CreateUser(t *testing.T, db *ent_gen.Client, m ent_gen.User, opts ...func(*ent_gen.UserCreate)) *ent_gen.User {
	t.Helper()

	tbl := &ent_gen.User{
		ID:    rand.Int64(),
		Name:  lo.RandomString(32, lo.AlphanumericCharset),
		Bytes: []byte(lo.RandomString(32, lo.AlphanumericCharset)),
	}

	if isModified(m.ID) {
		tbl.ID = m.ID
	}
	if isModified(m.Name) {
		tbl.Name = m.Name
	}
	if len(m.Bytes) > 0 {
		tbl.Bytes = m.Bytes
	}
	if isModified(m.CreatedAt) {
		tbl.CreatedAt = m.CreatedAt
	}
	if m.UpdatedAt != nil {
		tbl.UpdatedAt = m.UpdatedAt
	}

	builder := db.User.Create().
		SetID(tbl.ID).
		SetName(tbl.Name).
		SetBytes(tbl.Bytes).
		SetUpdatedAt(*tbl.UpdatedAt)
	for _, opt := range opts {
		opt(builder)
	}

	createdTbl, err := builder.Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	return createdTbl
}
