// Code generated by fixgen, DO NOT EDIT.

package fixture

import "github.com/samber/lo"

func isOverWritten[T comparable](v T) bool {
  zero := lo.Empty[T]()
  return v != zero
}
