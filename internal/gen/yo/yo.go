package yo

import (
	"errors"
	"fmt"
	"maps"
	"regexp"
	"strings"

	yo_loaders "go.mercari.io/yo/loaders"
	"go.mercari.io/yo/models"

	"github.com/earlgray283/fixgen/internal/caseconv"
	"github.com/earlgray283/fixgen/internal/gen"
	"github.com/earlgray283/fixgen/internal/load"
	"github.com/earlgray283/fixgen/internal/templates"
)

type Generator struct {
	yoPackagePath string
	genDirPath    string
	filepaths     []string
	tables        Tables
}

var _ gen.Generator = (*Generator)(nil)

// NewGenerator is a constructor for the struct Generator.
func NewGenerator(workDir string) (*Generator, error) {
	goModulePath, err := gen.LoadGoModulePath(workDir)
	if err != nil {
		return nil, fmt.Errorf(": %+w", err)
	}

	genDirPath, filepaths, err := gen.FindAndReadDirByFileName(workDir, "yo_db.yo.go")
	if err != nil {
		return nil, fmt.Errorf("failed to util.FindAndReadDirByFileName: %+w", err)
	}

	ddlPath, err := gen.FindFile(workDir, "schema.sql")
	if err != nil {
		return nil, fmt.Errorf("failed to util.FindFilePath: %+w", err)
	}
	tables, err := loadYoTables(ddlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load yo tables: %+w", err)
	}

	return &Generator{
		yoPackagePath: strings.Join([]string{goModulePath, genDirPath}, "/"),
		genDirPath:    genDirPath,
		filepaths:     filepaths,
		tables:        tables,
	}, nil
}

func (g *Generator) PackageInfo() *gen.PackageInfo {
	return &gen.PackageInfo{
		PackagePath:     g.yoPackagePath,
		PackageAlias:    "yo_gen",
		PackageLocation: g.genDirPath,
	}
}

type (
	Tables  map[string]Columns
	Columns map[string]*models.Column
)

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

func (g *Generator) Generate(structInfos []*load.StructInfo, data map[string]any) ([]*gen.File, error) {
	files := make([]*gen.File, 0, len(structInfos))

	yoStructInfos := make([]*structInfo, 0)
	for _, si := range structInfos {
		yosi, err := g.parse(si)
		if err != nil {
			if errors.Is(err, gen.ErrNotTargetStruct) {
				continue
			}
			return nil, err
		}
		if yosi == nil {
			continue
		}
		yoStructInfos = append(yoStructInfos, yosi)
	}

	for _, si := range yoStructInfos {
		file, err := g.execute(si, data)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

type structInfo struct {
	tableName string
	fields    []*field
}

type field struct {
	*load.Field
	IsSpannerNullType    bool // the type is `spanner.Null{type}`
	AllowCommitTimestamp bool
}

func (g *Generator) parse(si *load.StructInfo) (*structInfo, error) {
	sqlTableName := extractSQLTableNameFromComments(si.Comments)
	if sqlTableName == "" {
		return nil, gen.ErrNotTargetStruct
	}

	columns, ok := g.tables[sqlTableName]
	if !ok {
		return nil, fmt.Errorf("failed to get the table information `%s`", sqlTableName)
	}

	fields := make([]*field, 0, len(si.Fields))
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

		fields = append(fields, &field{
			Field:                f,
			IsSpannerNullType:    strings.HasPrefix(f.Type.Name, "spanner.Null"),
			AllowCommitTimestamp: column.AllowCommitTimestamp,
		})
	}

	return &structInfo{
		tableName: si.Name,
		fields:    fields,
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

func (g *Generator) execute(si *structInfo, data map[string]any) (*gen.File, error) {
	yoData := map[string]any{
		"TableName": si.tableName,
		"Fields":    si.fields,
	}
	maps.Copy(yoData, data)

	content, err := templates.Execute(templates.TmplYoFile, yoData)
	if err != nil {
		return nil, err
	}

	return &gen.File{
		Name:    caseconv.ConvertPascalToSnake(si.tableName),
		Content: content,
	}, nil
}
