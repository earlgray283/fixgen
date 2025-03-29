package datastore

import (
	"fmt"
	"maps"
	"strings"

	"github.com/earlgray283/fixgen/internal/caseconv"
	"github.com/earlgray283/fixgen/internal/gen"
	"github.com/earlgray283/fixgen/internal/load"
	"github.com/earlgray283/fixgen/internal/templates"
)

type Generator struct {
	packagePath string
	dirPath     string
	filepaths   []string
}

var _ gen.Generator = (*Generator)(nil)

func NewGenerator(workDir, packageDirPath string) (*Generator, error) {
	goModulePath, err := gen.LoadGoModulePath(workDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load go module path: %+w", err)
	}

	filepaths, err := gen.ReadDir(packageDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %+w", err)
	}

	return &Generator{
		packagePath: strings.Join([]string{goModulePath, packageDirPath}, "/"),
		dirPath:     packageDirPath,
		filepaths:   filepaths,
	}, nil
}

// Generate implements gen.Generator.
func (g *Generator) Generate(structInfos []*load.StructInfo, data map[string]any) ([]*gen.File, error) {
	files := make([]*gen.File, 0, len(structInfos)+1)

	content, err := templates.Execute(templates.TmplStructsCommonFile, nil)
	if err != nil {
		return nil, err
	}
	files = append(files, &gen.File{
		Name:    "structs_common",
		Content: content,
	})

	for _, si := range structInfos {
		file, err := g.execute(si, data)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func (g *Generator) execute(si *load.StructInfo, data map[string]any) (*gen.File, error) {
	newData := map[string]any{
		"TableName": si.Name,
		"Fields":    si.Fields,
	}
	maps.Copy(newData, data)

	content, err := templates.Execute(templates.TmplStructsFile, newData)
	if err != nil {
		return nil, err
	}

	return &gen.File{
		Name:    caseconv.ConvertPascalToSnake(si.Name),
		Content: content,
	}, nil
}

// PackageInfo implements gen.Generator.
func (g *Generator) PackageInfo() *gen.PackageInfo {
	return &gen.PackageInfo{
		PackagePath:     g.packagePath,
		PackageAlias:    "structs",
		PackageLocation: g.dirPath,
	}
}
