// Code generated by fixgen, DO NOT EDIT.

package fixture

import (
	"context"
	ent_gen "ent/ent"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/samber/lo"
)

func CreateUser(t *testing.T, db *ent_gen.Client, m *ent_gen.User, opts ...func(*ent_gen.UserCreate)) *ent_gen.User {
	t.Helper()

	tbl := &ent_gen.User{
		ID:        rand.Int64(),
		Name:      "Taro Yamada", // Name is overwritten
		Bytes:     []byte(lo.RandomString(32, lo.AlphanumericCharset)),
		CreatedAt: time.Now(),
		// UpdatedAt is nillable
	}

	builder := db.User.Create()
	if isModified(m.ID) {
		builder = builder.SetID(tbl.ID)
	}
	if isModified(m.Name) {
		builder = builder.SetName(tbl.Name)
	}
	if len(m.Bytes) > 0 {
		builder = builder.SetBytes(tbl.Bytes)
	}
	if isModified(m.CreatedAt) {
		builder = builder.SetCreatedAt(tbl.CreatedAt)
	}
	if m.UpdatedAt != nil {
		builder = builder.SetUpdatedAt(*tbl.UpdatedAt)
	}
	for _, opt := range opts {
		opt(builder)
	}

	createdTbl, err := builder.Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	return createdTbl
}
