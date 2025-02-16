package ent

import (
	"errors"
	"fmt"
	"io/fs"
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
	structInfos    []*load.StructInfo
	tables         Tables
}

func NewGenerator(workDir string) (*Generator, error) {
	schemaDirPath, genDirPath, err := findEntDirs(workDir)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(genDirPath)
	if err != nil {
		return nil, err
	}

	targetStructInfos := make([]*load.StructInfo, 0)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		p := filepath.Join(genDirPath, e.Name())
		structInfos, err := load.LoadStructInfos(p)
		if err != nil {
			return nil, fmt.Errorf("failed to load structInfos from %s: %+w", p, err)
		}
		targetStructInfos = append(targetStructInfos, structInfos...)
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
		structInfos:    targetStructInfos,
		tables:         tables,
	}, nil
}

func findEntDirs(workDir string) (schemaDirPath, genDirPath string, err error) {
	if err := fs.WalkDir(os.DirFS(workDir), ".", func(path string, d fs.DirEntry, err error) error {
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

func (g *Generator) Generate() ([]*gen.File, error) {
	files := make([]*gen.File, 0)

	for _, si := range g.structInfos {
		f, err := g.generate(si)
		if err != nil {
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

func (g *Generator) GenPackageInfo() *gen.GenPackageInfo {
	return &gen.GenPackageInfo{
		PackagePath:  g.entPackagePath,
		PackageAlias: "ent_gen",
	}
}

func (g *Generator) IsExperimental() bool {
	return true
}

func (g *Generator) generate(si *load.StructInfo) (*gen.File, error) {
	columns, ok := g.tables[si.Name]
	if !ok {
		return nil, nil
	}

	fields := make([]*Field, 0, len(si.Fields))
	for _, f := range si.Fields {
		sqlColumnName := extractSQLColumnName(f.Tags["json"])
		if sqlColumnName == "" {
			return nil, fmt.Errorf("failed to extractSQLColumnName(fieldName: %s)", f.Name)
		}
		column, ok := columns[sqlColumnName]
		if !ok {
			continue
		}
		fields = append(fields, &Field{
			Field:              f,
			Immutable:          column.Immutable,
			Nillable:           column.Nillable,
			HasDefaultOnCreate: column.Default,
			Ignore:             column.Immutable || column.Nillable || column.Default || f.DefaultValue == "",
		})
	}
	fields[len(fields)-1].IsLast = true

	content, err := templates.Execute(templates.TmplEntFile, map[string]any{
		"TableName": si.Name,
		"Fields":    fields,
	})
	if err != nil {
		return nil, err
	}

	return &gen.File{
		Name:    caseconv.ConvertPascalToSnake(si.Name),
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

type Tables map[string]Columns
type Columns map[string]*ent_load.Field

type Field struct {
	*load.Field
	Immutable          bool
	Nillable           bool
	HasDefaultOnCreate bool
	Ignore             bool
	IsLast             bool
}
