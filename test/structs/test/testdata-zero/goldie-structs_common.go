// Code generated by fixgen, DO NOT EDIT.

package fixture

import "context"

type Inserter[M any] interface {
	Insert(ctx context.Context, m M) (M, error)
}
