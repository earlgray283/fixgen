package ent

import (
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"strings"

	ent_load "entgo.io/ent/entc/load"

	"github.com/earlgray283/fixgen/internal/caseconv"
	"github.com/earlgray283/fixgen/internal/gen"
	"github.com/earlgray283/fixgen/internal/load"
	"github.com/earlgray283/fixgen/internal/templates"
)

type Generator struct {
	entPackagePath string
	genDirPath     string
	tables         Tables
}

var _ gen.Generator = (*Generator)(nil)

func NewGenerator(workDir string) (*Generator, error) {
	schemaDirPath, genDirPath, err := findEntDirs(workDir)
	if err != nil {
		return nil, err
	}

	spec, err := (&ent_load.Config{Path: schemaDirPath}).Load()
	if err != nil {
		return nil, fmt.Errorf("failed to ent load: %+w", err)
	}

	tables := make(Tables, len(spec.Schemas))
	for _, s := range spec.Schemas {
		cols := make(Columns, len(s.Fields))
		for _, f := range s.Fields {
			cols[f.Name] = f
		}
		tables[s.Name] = cols
	}

	return &Generator{
		entPackagePath: strings.Join([]string{spec.Module.Path, genDirPath}, "/"),
		genDirPath:     genDirPath,
		tables:         tables,
	}, nil
}

func findEntDirs(workDir string) (schemaDirPath, genDirPath string, err error) {
	if err := fs.WalkDir(os.DirFS(workDir), ".", func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() && d.Name() == "schema" {
			schemaDirPath = filepath.Join(workDir, path)
		}
		if !d.IsDir() && d.Name() == "client.go" {
			genDirPath = filepath.Join(workDir, filepath.Dir(path))
		}
		if schemaDirPath != "" && genDirPath != "" {
			return fs.SkipAll
		}
		return nil
	}); err != nil {
		return "", "", err
	}
	if schemaDirPath == "" {
		return "", "", errors.New("the `schema` directory was not found")
	}
	if genDirPath == "" {
		return "", "", errors.New("the directory which `client.go` exists was not found")
	}
	if !filepath.IsAbs(schemaDirPath) {
		schemaDirPath = fmt.Sprintf("./%s", schemaDirPath)
	}
	return schemaDirPath, genDirPath, nil
}

func (g *Generator) PackageInfo() *gen.PackageInfo {
	return &gen.PackageInfo{
		PackagePath:     g.entPackagePath,
		PackageAlias:    "ent_gen",
		PackageLocation: g.genDirPath,
	}
}

func (g *Generator) Generate(structInfos []*load.StructInfo, data map[string]any) ([]*gen.File, error) {
	files := make([]*gen.File, 0)

	for _, si := range structInfos {
		f, err := g.generate(si, data)
		if err != nil {
			if errors.Is(err, gen.ErrNotTargetStruct) {
				continue
			}
			return nil, err
		}
		// not target
		if f == nil {
			continue
		}
		files = append(files, f)
	}
	return files, nil
}

func (g *Generator) generate(si *load.StructInfo, data map[string]any) (*gen.File, error) {
	columns, ok := g.tables[si.Name]
	if !ok {
		return nil, gen.ErrNotTargetStruct
	}

	fields := make([]*field, 0, len(si.Fields))
	for _, f := range si.Fields {
		sqlColumnName := extractSQLColumnName(f.Tags["json"])
		if sqlColumnName == "" {
			return nil, fmt.Errorf("failed to extractSQLColumnName(fieldName: %s)", f.Name)
		}
		column, ok := columns[sqlColumnName]
		if !ok {
			continue
		}
		fields = append(fields, &field{
			Field:              f,
			IsNillable:         column.Nillable,
			HasDefaultOnCreate: column.Default,
		})
	}
	fields[len(fields)-1].IsLast = true

	f, err := g.execute(&structInfo{
		tableName: si.Name,
		fields:    fields,
	}, data)
	if err != nil {
		return nil, err
	}

	return f, nil
}

type structInfo struct {
	tableName string
	fields    []*field
}

func (g *Generator) execute(si *structInfo, data map[string]any) (*gen.File, error) {
	entData := map[string]any{
		"TableName": si.tableName,
		"Fields":    si.fields,
	}
	maps.Copy(entData, data)

	content, err := templates.Execute(templates.TmplEntFile, entData)
	if err != nil {
		return nil, err
	}

	return &gen.File{
		Name:    caseconv.ConvertPascalToSnake(si.tableName),
		Content: content,
	}, nil
}

func extractSQLColumnName(tagValue string) string {
	values := strings.Split(tagValue, ",")
	if len(values) == 0 {
		return ""
	}
	return strings.TrimSpace(values[0])
}

type (
	Tables  map[string]Columns
	Columns map[string]*ent_load.Field
)

type field struct {
	*load.Field
	IsNillable         bool
	HasDefaultOnCreate bool
	IsLast             bool
}
