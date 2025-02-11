package fixture

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

	ent_gen "ent-tutorial/ent"

	"github.com/samber/lo"
)

func CreateUser(t *testing.T, db *ent_gen.Client, ovrTbl ent_gen.User) *ent_gen.User {
	t.Helper()

	tbl := &ent_gen.User{
		ID:        rand.Int64(),
		Name:      lo.RandomString(32, lo.AlphanumericCharset),
		CreatedAt: time.Now(),
		UpdatedAt: lo.ToPtr(time.Now()),
	}

	createdTbl, err := db.User.Create().
		SetID(tbl.ID).
		SetName(tbl.Name).
		SetCreatedAt(tbl.CreatedAt).
		SetUpdatedAt(*tbl.UpdatedAt).
		Save(context.Background())
	if err != nil {
		panic(err)
	}

	return createdTbl
}
