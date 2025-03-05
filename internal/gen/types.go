package gen

import (
	"github.com/earlgray283/fixgen/internal/load"
)

type Generator interface {
	PackageInfo() *PackageInfo
	Generate(si []*load.StructInfo, data map[string]any) ([]*File, error)
}

type File struct {
	Name    string
	Content []byte
}

type PackageInfo struct {
	PackagePath     string // e.g. github.com/earlgray283/pj-todo/models
	PackageAlias    string // e.g. yo_gen
	PackageLocation string
}
