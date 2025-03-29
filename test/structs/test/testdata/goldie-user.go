// Code generated by fixgen, DO NOT EDIT.

package fixture

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"
	"time"
)

func CreateUser(t *testing.T, db Inserter, m *structs_gen.User, opts ...func(*structs_gen.User)) *structs_gen.User {
	t.Helper()

	tbl := &structs_gen.User{
		ID:        rand.Int64(),
		Name:      "Taro Yamada",                                // Name is overwritten
		IconURL:   fmt.Sprintf("http://example.com/%d", 123456), // IconURL is overwritten
		UserType:  1,                                            // UserType is overwritten
		CreatedAt: time.Now(),
		// UpdatedAt is unknown
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
	if m.UserType != 1 {
		tbl.UserType = m.UserType
	}
	if isModified(m.CreatedAt) {
		tbl.CreatedAt = m.CreatedAt
	}
	if isModified(m.UpdatedAt) {
		tbl.UpdatedAt = m.UpdatedAt
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
