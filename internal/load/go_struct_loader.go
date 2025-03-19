package load

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"regexp"
	"strings"
	"unicode"

	"github.com/earlgray283/fixgen/internal/config"
)

type StructInfoLoader struct {
	structConfigs config.Structs
}

func New(structConfigs config.Structs) *StructInfoLoader {
	return &StructInfoLoader{structConfigs}
}

type unsupportedTypeError struct {
	fieldName string
	typeName  string
}

func (e *unsupportedTypeError) Error() string {
	return fmt.Sprintf("unsupported type error(field: `%s`, type: `%s`)", e.fieldName, e.typeName)
}

func (l *StructInfoLoader) Load(goFilePath string, useRandv1 bool) ([]*StructInfo, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, goFilePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parser.ParseFile: %+w", err)
	}

	defaultValueMap := defaultValueMapRandv2
	if useRandv1 {
		defaultValueMap = defaultValueMapRandv1
	}

	structInfos := make([]*StructInfo, 0)

	var parseErr error
	ast.Inspect(f, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok {
			return true
		}

		comments := make([]string, 0)

		// `type T struct` の形式だったら `genDecl.Doc.List` を struct のコメントとみなす
		if len(genDecl.Specs) == 1 {
			if genDecl.Doc != nil {
				for _, c := range genDecl.Doc.List {
					comments = append(comments, c.Text)
				}
			}
		}

		// TODO: support multiple struct definitions in typeDecl
		ts, ok := genDecl.Specs[0].(*ast.TypeSpec)
		if !ok {
			return true
		}

		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		structInfo, err := parse(ts.Name.Name, st, defaultValueMap)
		if err != nil {
			var uerr *unsupportedTypeError
			if errors.As(err, &uerr) {
				log.Printf("[WARNING] %v\n", uerr)
				return true
			}
			parseErr = err
			return false
		}
		structInfo.Comments = comments

		structInfos = append(structInfos, structInfo)

		return true
	})
	if parseErr != nil {
		return nil, parseErr
	}

	for _, si := range structInfos {
		sc, ok := l.structConfigs[si.Name]
		if !ok {
			continue
		}
		for _, f := range si.Fields {
			fc, ok := sc.Fields[f.Name]
			if !ok {
				continue
			}
			defaultValue, ok := fc.DefaultValue()
			if ok {
				f.DefaultValue = defaultValue
				f.IsOverwritten = true
			}
			f.IsModifiedCond = fc.IsModifiedCond
			f.MustOverwrite = fc.MustOverwrite
		}
	}

	return structInfos, nil
}

var regexpTagKeyAndValue = regexp.MustCompile(`(.+):"(.*)"`)

func extractTagKeyValue(tag string) (string, string, error) {
	matches := regexpTagKeyAndValue.FindStringSubmatch(tag)
	if len(matches) != 3 {
		return "", "", fmt.Errorf("failed to extract tag's key and value: `%s`", tag)
	}

	return matches[1], matches[2], nil
}

func resolveType(name string, exprType ast.Expr, defaultValueMap map[string]string) (typ *Type, defaultValue string, err error) {
	switch ty := exprType.(type) {
	case *ast.Ident: // TODO: resolve struct type
		typ := ty.Name
		return &Type{Name: typ}, defaultValueMap[typ], nil
	case *ast.SelectorExpr:
		typ := fmt.Sprintf("%s.%s", ty.X.(*ast.Ident).Name, ty.Sel.Name)
		return &Type{Name: typ}, defaultValueMap[typ], nil
	case *ast.StarExpr:
		resolved, _, err := resolveType(name, ty.X, defaultValueMap)
		if err != nil {
			return nil, "", err
		}
		return &Type{Name: fmt.Sprintf("*%s", resolved.Name), IsNillable: true}, "", nil
	case *ast.ArrayType:
		resolved, _, err := resolveType(name, ty.Elt, defaultValueMap)
		if err != nil {
			return nil, "", err
		}
		typ := fmt.Sprintf("[]%s", resolved.Name)
		return &Type{Name: typ, IsSlice: true}, defaultValueMap[typ], nil
	default:
		return nil, "", &unsupportedTypeError{fieldName: name, typeName: fmt.Sprintf("%T", exprType)}
	}
}

func parse(name string, st *ast.StructType, defaultValueMap map[string]string) (*StructInfo, error) {
	fields := make([]*Field, 0, len(st.Fields.List))
	for _, f := range st.Fields.List {
		names := make([]string, 0, len(f.Names))
		for _, n := range f.Names {
			names = append(names, n.Name)
		}

		typ, defaultValue, err := resolveType(strings.Join(names, ", "), f.Type, defaultValueMap)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve type: %+w", err)
		}

		tags := make(map[string]string, 0)
		if f.Tag != nil {
			tagValues := strings.Fields(strings.Trim(f.Tag.Value, "`"))
			tags = make(map[string]string, len(tagValues))
			for _, tagValue := range tagValues {
				key, value, err := extractTagKeyValue(tagValue)
				if err != nil {
					return nil, fmt.Errorf("%+w", err)
				}
				tags[key] = value
				break
			}
		}

		// for supporting anonymous field(=type only)
		name := typ.Name
		if len(f.Names) != 0 {
			name = f.Names[0].Name
		}

		// means unexported field
		if unicode.IsLower(rune(name[0])) {
			continue
		}

		fields = append(fields, &Field{
			Name:         name,
			Type:         typ,
			DefaultValue: defaultValue,
			Tags:         tags,
		})
	}

	return &StructInfo{
		Name:   name,
		Fields: fields,
	}, nil
}

var (
	defaultValueMapRandv2 = map[string]string{
		"int32":     "rand.Int32()",
		"int64":     "rand.Int64()",
		"uint32":    "rand.Uint32()",
		"uint64":    "rand.Uint64()",
		"float32":   "rand.Float32()",
		"float64":   "rand.Float64()",
		"string":    "lo.RandomString(32, lo.AlphanumericCharset)",
		"[]byte":    "[]byte(lo.RandomString(32, lo.AlphanumericCharset))",
		"bool":      "false",
		"time.Time": "time.Now()",
	}
	defaultValueMapRandv1 = map[string]string{
		"int32":     "rand.Int31()",
		"int64":     "rand.Int63()",
		"uint32":    "rand.Uint31()",
		"uint64":    "rand.Uint63()",
		"float32":   "rand.Float32()",
		"float64":   "rand.Float64()",
		"string":    "lo.RandomString(32, lo.AlphanumericCharset)",
		"[]byte":    "[]byte(lo.RandomString(32, lo.AlphanumericCharset))",
		"bool":      "false",
		"time.Time": "time.Now()",
	}
)
