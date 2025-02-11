package fixture

import (
	"context"
	"math/rand/v2"
	"testing"

	yo_gen "yo/models"

	"cloud.google.com/go/spanner"
	"github.com/samber/lo"
)

func CreateUser(t *testing.T, db *spanner.Client, ovrTbl *yo_gen.User) *yo_gen.User {
	t.Helper()

	tbl := &yo_gen.User{
		ID:        rand.Int64(),
		Name:      lo.RandomString(32, lo.AlphanumericCharset),
		CreatedAt: spanner.CommitTimestamp,
	}

	if isOverWritten(ovrTbl.ID) {
		tbl.ID = ovrTbl.ID
	}
	if isOverWritten(ovrTbl.Name) {
		tbl.Name = ovrTbl.Name
	}
	if isOverWritten(ovrTbl.CreatedAt) {
		t.Fatal("spanner.CommitTimestamp should be used")
	}
	if !ovrTbl.UpdatedAt.IsNull() {
		tbl.UpdatedAt = ovrTbl.UpdatedAt
	}

	_, err := db.ReadWriteTransaction(context.Background(), func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		return tx.BufferWrite([]*spanner.Mutation{tbl.Insert(ctx)})
	})
	if err != nil {
		t.Fatal(err)
	}

	return tbl
}
