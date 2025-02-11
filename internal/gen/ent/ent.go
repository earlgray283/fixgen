package ent

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"entgo.io/ent/entc/load"
	"entgo.io/ent/schema/field"
	"github.com/earlgray283/fixgen/internal"
	"github.com/earlgray283/fixgen/internal/caseconv"
	"github.com/earlgray283/fixgen/internal/gen"
)

type Generator struct {
	packageName string
	genPkgPath  string
	spec        *load.SchemaSpec
}

func NewGenerator(packageName string) (*Generator, error) {
	var generatedDirPath, schemaDirPath string
	if err := fs.WalkDir(os.DirFS("."), ".", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, "schema") {
			schemaDirPath = fmt.Sprintf("./%s", path)
		}
		if d.Name() == "client.go" {
			generatedDirPath = filepath.Dir(path)
		}
		if schemaDirPath != "" && generatedDirPath != "" {
			return fs.SkipAll
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if generatedDirPath == "" || schemaDirPath == "" {
		return nil, errors.New("failed to load ent package information")
	}

	spec, err := (&load.Config{Path: schemaDirPath}).Load()
	if err != nil {
		return nil, err
	}

	return &Generator{
		packageName: packageName,
		genPkgPath:  strings.Join([]string{spec.Module.Path, generatedDirPath}, "/"),
		spec:        spec,
	}, nil
}

func (g *Generator) Generate() ([]*gen.File, error) {
	files := make([]*gen.File, 0, len(g.spec.Schemas))
	for _, s := range g.spec.Schemas {
		dfvs, err := makeDefaultFieldValues(s.Fields)
		if err != nil {
			return nil, err
		}
		f, err := g.generate(s.Name, dfvs)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (g *Generator) IsExperimental() bool {
	return true
}

func (g *Generator) generate(tableName string, dfvs []*DefaultFieldValue) (*gen.File, error) {
	buf := &bytes.Buffer{}
	if err := internal.TmplMockEntFile.Execute(buf, map[string]any{
		"PackageName":        g.packageName,
		"GenPkgPath":         g.genPkgPath,
		"TableName":          tableName,
		"DefaultFieldValues": dfvs,
	}); err != nil {
		return nil, err
	}

	formated, err := gen.Format(buf.Bytes())
	if err != nil {
		return nil, err
	}

	return &gen.File{
		Name:    caseconv.ConvertPascalToSnake(tableName),
		Content: formated,
	}, nil
}

func makeDefaultFieldValues(schemas []*load.Field) ([]*DefaultFieldValue, error) {
	dfvs := make([]*DefaultFieldValue, 0, len(schemas))
	for _, s := range schemas {
		dfv, err := convertInfoTypeToDefaultFieldValue(s.Name, s.Info, s.Nillable)
		if err != nil {
			return nil, err
		}
		dfvs = append(dfvs, dfv)
	}
	return dfvs, nil
}

func convertInfoTypeToDefaultFieldValue(name string, typeInfo *field.TypeInfo, isNillable bool) (*DefaultFieldValue, error) {
	value, ok := defaultValueMap[typeInfo.Type]
	if !ok {
		return nil, fmt.Errorf("unsupported type: %s", typeInfo.Type.String())
	}

	if isNillable {
		value = fmt.Sprintf("lo.ToPtr(%s)", value)
	}
	return &DefaultFieldValue{
		FieldName: caseconv.ConvertSnakeToPascal(name),
		Value:     value,
		IsPointer: isNillable,
	}, nil
}

type DefaultFieldValue struct {
	FieldName string
	Value     string
	IsPointer bool
}

var defaultValueMap = map[field.Type]string{
	field.TypeInt64:  "rand.Int64()",
	field.TypeString: "lo.RandomString(32, lo.AlphanumericCharset)",
	field.TypeTime:   "time.Now()",
	field.TypeBytes:  "[]byte(lo.RandomString(32, lo.AlphanumericCharset))",
}
