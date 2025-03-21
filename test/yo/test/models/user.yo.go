// Code generated by yo. DO NOT EDIT.
// Package models contains the types.
package models

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/grpc/codes"
)

// User represents a row from 'Users'.
type User struct {
	ID        int64            `spanner:"id" json:"id"`                 // id
	Name      string           `spanner:"name" json:"name"`             // name
	IconURL   string           `spanner:"icon_url" json:"icon_url"`     // icon_url
	UserType  int64            `spanner:"user_type" json:"user_type"`   // user_type
	CreatedAt time.Time        `spanner:"created_at" json:"created_at"` // created_at
	UpdatedAt spanner.NullTime `spanner:"updated_at" json:"updated_at"` // updated_at
}

func UserPrimaryKeys() []string {
	return []string{
		"id",
	}
}

func UserColumns() []string {
	return []string{
		"id",
		"name",
		"icon_url",
		"user_type",
		"created_at",
		"updated_at",
	}
}

func UserWritableColumns() []string {
	return []string{
		"id",
		"name",
		"icon_url",
		"user_type",
		"created_at",
		"updated_at",
	}
}

func (u *User) columnsToPtrs(cols []string, customPtrs map[string]interface{}) ([]interface{}, error) {
	ret := make([]interface{}, 0, len(cols))
	for _, col := range cols {
		if val, ok := customPtrs[col]; ok {
			ret = append(ret, val)
			continue
		}

		switch col {
		case "id":
			ret = append(ret, &u.ID)
		case "name":
			ret = append(ret, &u.Name)
		case "icon_url":
			ret = append(ret, &u.IconURL)
		case "user_type":
			ret = append(ret, &u.UserType)
		case "created_at":
			ret = append(ret, &u.CreatedAt)
		case "updated_at":
			ret = append(ret, &u.UpdatedAt)
		default:
			return nil, fmt.Errorf("unknown column: %s", col)
		}
	}
	return ret, nil
}

func (u *User) columnsToValues(cols []string) ([]interface{}, error) {
	ret := make([]interface{}, 0, len(cols))
	for _, col := range cols {
		switch col {
		case "id":
			ret = append(ret, u.ID)
		case "name":
			ret = append(ret, u.Name)
		case "icon_url":
			ret = append(ret, u.IconURL)
		case "user_type":
			ret = append(ret, u.UserType)
		case "created_at":
			ret = append(ret, u.CreatedAt)
		case "updated_at":
			ret = append(ret, u.UpdatedAt)
		default:
			return nil, fmt.Errorf("unknown column: %s", col)
		}
	}

	return ret, nil
}

// newUser_Decoder returns a decoder which reads a row from *spanner.Row
// into User. The decoder is not goroutine-safe. Don't use it concurrently.
func newUser_Decoder(cols []string) func(*spanner.Row) (*User, error) {
	customPtrs := map[string]interface{}{}

	return func(row *spanner.Row) (*User, error) {
		var u User
		ptrs, err := u.columnsToPtrs(cols, customPtrs)
		if err != nil {
			return nil, err
		}

		if err := row.Columns(ptrs...); err != nil {
			return nil, err
		}

		return &u, nil
	}
}

// Insert returns a Mutation to insert a row into a table. If the row already
// exists, the write or transaction fails.
func (u *User) Insert(ctx context.Context) *spanner.Mutation {
	values, _ := u.columnsToValues(UserWritableColumns())
	return spanner.Insert("Users", UserWritableColumns(), values)
}

// Update returns a Mutation to update a row in a table. If the row does not
// already exist, the write or transaction fails.
func (u *User) Update(ctx context.Context) *spanner.Mutation {
	values, _ := u.columnsToValues(UserWritableColumns())
	return spanner.Update("Users", UserWritableColumns(), values)
}

// InsertOrUpdate returns a Mutation to insert a row into a table. If the row
// already exists, it updates it instead. Any column values not explicitly
// written are preserved.
func (u *User) InsertOrUpdate(ctx context.Context) *spanner.Mutation {
	values, _ := u.columnsToValues(UserWritableColumns())
	return spanner.InsertOrUpdate("Users", UserWritableColumns(), values)
}

// UpdateColumns returns a Mutation to update specified columns of a row in a table.
func (u *User) UpdateColumns(ctx context.Context, cols ...string) (*spanner.Mutation, error) {
	// add primary keys to columns to update by primary keys
	colsWithPKeys := append(cols, UserPrimaryKeys()...)

	values, err := u.columnsToValues(colsWithPKeys)
	if err != nil {
		return nil, newErrorWithCode(codes.InvalidArgument, "User.UpdateColumns", "Users", err)
	}

	return spanner.Update("Users", colsWithPKeys, values), nil
}

// FindUser gets a User by primary key
func FindUser(ctx context.Context, db YORODB, id int64) (*User, error) {
	key := spanner.Key{id}
	row, err := db.ReadRow(ctx, "Users", key, UserColumns())
	if err != nil {
		return nil, newError("FindUser", "Users", err)
	}

	decoder := newUser_Decoder(UserColumns())
	u, err := decoder(row)
	if err != nil {
		return nil, newErrorWithCode(codes.Internal, "FindUser", "Users", err)
	}

	return u, nil
}

// ReadUser retrieves multiples rows from User by KeySet as a slice.
func ReadUser(ctx context.Context, db YORODB, keys spanner.KeySet) ([]*User, error) {
	var res []*User

	decoder := newUser_Decoder(UserColumns())

	rows := db.Read(ctx, "Users", keys, UserColumns())
	err := rows.Do(func(row *spanner.Row) error {
		u, err := decoder(row)
		if err != nil {
			return err
		}
		res = append(res, u)

		return nil
	})
	if err != nil {
		return nil, newErrorWithCode(codes.Internal, "ReadUser", "Users", err)
	}

	return res, nil
}

// Delete deletes the User from the database.
func (u *User) Delete(ctx context.Context) *spanner.Mutation {
	values, _ := u.columnsToValues(UserPrimaryKeys())
	return spanner.Delete("Users", spanner.Key(values))
}
