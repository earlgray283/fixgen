package yo

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/earlgray283/fixgen/internal"
	"github.com/earlgray283/fixgen/internal/caseconv"
	"github.com/earlgray283/fixgen/internal/gen"
	"github.com/earlgray283/fixgen/internal/loaders"
	yo_loaders "go.mercari.io/yo/loaders"
	"go.mercari.io/yo/models"
)

type Generator struct {
	opt *gen.Option

	genDirPath    string
	yoPackagePath string
	structInfos   []*StructInfo
	tables        Tables
}

// NewGenerator is a constructor for the struct Generator
func NewGenerator(opts ...gen.OptionFunc) (*Generator, error) {
	opt := gen.DefaultOption()
	opt.ApplyOptionFuncs(opts...)

	goModulePath, err := gen.LoadGoModulePath(opt.WorkDir)
	if err != nil {
		return nil, fmt.Errorf(": %+w\n", err)
	}

	genDirPath, filepathList, err := gen.FindAndReadDirByFileName(opt.WorkDir, "yo_db.yo.go")
	if err != nil {
		return nil, fmt.Errorf("failed to util.FindAndReadDirByFileName: %+w", err)
	}

	targetStructInfos := make([]*StructInfo, 0)
	for _, p := range filepathList {
		structInfos, err := loaders.LoadStructInfos(p)
		if err != nil {
			return nil, fmt.Errorf("failed to load structInfos from %s: %+w", p, err)
		}
		for _, si := range structInfos {
			sqlTableName := extractSQLTableNameFromComments(si.Comments)
			if sqlTableName == "" {
				continue
			}
			targetStructInfos = append(targetStructInfos, &StructInfo{
				StructInfo:   si,
				SQLTableName: sqlTableName,
			})
		}
	}

	ddlPath, err := gen.FindFilePath(opt.WorkDir, "schema.sql")
	if err != nil {
		return nil, fmt.Errorf("failed to util.FindFilePath: %+w", err)
	}
	tables, err := loadYoTables(ddlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load yo tables: %+w", err)
	}

	return &Generator{
		opt:           opt,
		genDirPath:    genDirPath,
		yoPackagePath: strings.Join([]string{goModulePath, genDirPath}, "/"),
		structInfos:   targetStructInfos,
		tables:        tables,
	}, nil
}

var regexpYoStructComment = regexp.MustCompile(`(.+) represents a row from '(.+)'\.`)

func extractSQLTableNameFromComments(comments []string) string {
	for _, comment := range comments {
		matches := regexpYoStructComment.FindStringSubmatch(comment)
		if len(matches) < 3 || matches[2] == "" {
			continue
		}
		return matches[2]
	}
	return ""
}

func (g *Generator) Generate() ([]*gen.File, error) {
	files := make([]*gen.File, 0, len(g.structInfos)+1)

	commonFile, err := g.generateCommonFile()
	if err != nil {
		return nil, err
	}
	files = append(files, commonFile)

	for _, si := range g.structInfos {
		file, err := g.generate(si)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func (g *Generator) IsExperimental() bool {
	return true
}

type Tables map[string]Columns
type Columns map[string]*models.Column

func loadYoTables(ddlPath string) (Tables, error) {
	yoLoader, err := yo_loaders.NewSpannerLoaderFromDDL(ddlPath)
	if err != nil {
		return nil, err
	}

	tbls, err := yoLoader.TableList()
	if err != nil {
		return nil, err
	}

	tables := make(Tables, len(tbls))
	for _, tbl := range tbls {
		cols, err := yoLoader.ColumnList(tbl.TableName)
		if err != nil {
			return nil, err
		}
		columns := make(Columns, len(cols))
		for _, col := range cols {
			columns[col.ColumnName] = col
		}
		tables[tbl.TableName] = columns
	}

	return tables, nil
}

func (g *Generator) generateCommonFile() (*gen.File, error) {
	buf := &bytes.Buffer{}
	if err := internal.TmplCommonFile.Execute(buf, map[string]string{
		"PackageName": g.opt.PackageName,
	}); err != nil {
		return nil, err
	}
	return &gen.File{Name: "common", Content: buf.Bytes()}, nil
}

type StructInfo struct {
	*gen.StructInfo
	SQLTableName string
}

type Field struct {
	*gen.Field
	IsSpannerNullType    bool // the type is `spanner.Null{type}`
	AllowCommitTimestamp bool
}

func (g *Generator) generate(si *StructInfo) (*gen.File, error) {
	columns, ok := g.tables[si.SQLTableName]
	if !ok {
		return nil, fmt.Errorf("failed to get the table information `%s`", si.SQLTableName)
	}

	fields := make([]*Field, 0, len(si.Fields))
	for _, f := range si.Fields {
		sqlColumnName, ok := f.Tags["spanner"]
		if !ok {
			return nil, fmt.Errorf("failed to extract SQLColumnName from tag(%v)", f.Tags)
		}
		column, ok := columns[sqlColumnName]
		if !ok {
			return nil, fmt.Errorf("failed to get column: %v", sqlColumnName)
		}
		if column.AllowCommitTimestamp {
			f.DefaultValue = "spanner.CommitTimestamp"
		}
		fields = append(fields, &Field{
			Field:                f,
			IsSpannerNullType:    strings.HasPrefix(f.Type, "spanner.Null"),
			AllowCommitTimestamp: column.AllowCommitTimestamp,
		})
	}

	buf := &bytes.Buffer{}
	if err := internal.TmplMockYoFile.Execute(buf, map[string]any{
		"PackageName": g.opt.PackageName,
		"GenPkgPath":  g.yoPackagePath,
		"TableName":   si.Name,
		"Fields":      fields,
	}); err != nil {
		return nil, err
	}

	return &gen.File{
		Name:    caseconv.ConvertPascalToSnake(si.Name),
		Content: buf.Bytes(),
	}, nil
}
