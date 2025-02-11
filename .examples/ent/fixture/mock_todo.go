package fixture

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

	ent_gen "ent-tutorial/ent"

	"github.com/samber/lo"
)

func CreateTodo(t *testing.T, db *ent_gen.Client, ovrTbl ent_gen.Todo) *ent_gen.Todo {
	t.Helper()

	tbl := &ent_gen.Todo{
		ID:          rand.Int64(),
		Title:       lo.RandomString(32, lo.AlphanumericCharset),
		Description: lo.RandomString(32, lo.AlphanumericCharset),
		CreatedAt:   time.Now(),
		UpdatedAt:   lo.ToPtr(time.Now()),
		DoneAt:      lo.ToPtr(time.Now()),
	}

	createdTbl, err := db.Todo.Create().
		SetID(tbl.ID).
		SetTitle(tbl.Title).
		SetDescription(tbl.Description).
		SetCreatedAt(tbl.CreatedAt).
		SetUpdatedAt(*tbl.UpdatedAt).
		SetDoneAt(*tbl.DoneAt).
		Save(context.Background())
	if err != nil {
		panic(err)
	}

	return createdTbl
}
