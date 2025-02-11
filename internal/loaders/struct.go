package loaders

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

	fixgen_errors "github.com/earlgray283/fixgen/internal/errors"
	"github.com/earlgray283/fixgen/internal/gen"
)

func LoadStructInfos(goFilePath string) ([]*gen.StructInfo, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, goFilePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parser.ParseFile: %+w", err)
	}

	structInfoList := make([]*gen.StructInfo, 0)
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

		if ts.Name == nil {
			log.Printf("%+v\n", st)
		}

		structInfo, err := parse(ts.Name.Name, st)
		if err != nil {
			var uerr *fixgen_errors.UnsupportedTypeError
			if errors.As(err, &uerr) {
				log.Printf("[WARNING] %v\n", uerr)
				return true
			}
			parseErr = err
			return false
		}
		structInfo.Comments = comments
		structInfoList = append(structInfoList, structInfo)

		return true
	})
	if parseErr != nil {
		return nil, parseErr
	}

	return structInfoList, nil
}

var regexpTagKeyAndValue = regexp.MustCompile(`(.+):"(.*)"`)

func extractTagKeyValue(tag string) (string, string, error) {
	matches := regexpTagKeyAndValue.FindStringSubmatch(tag)
	if len(matches) != 3 {
		return "", "", fmt.Errorf("failed to extract tag's key and value: `%s`", tag)
	}
	return matches[1], matches[2], nil
}

func resolveType(name string, exprType ast.Expr) (typ string, defaultValue string, err error) {
	switch ty := exprType.(type) {
	case *ast.Ident: // TODO: resolve struct type
		typ := ty.Name
		return typ, defaultValueMap[typ], nil
	case *ast.SelectorExpr:
		typ := fmt.Sprintf("%s.%s", ty.X.(*ast.Ident).Name, ty.Sel.Name)
		return typ, defaultValueMap[typ], nil
	case *ast.StarExpr:
		resolved, _, err := resolveType(name, ty.X)
		if err != nil {
			return "", "", err
		}
		return fmt.Sprintf("*%s", resolved), "", nil
	case *ast.ArrayType:
		resolved, _, err := resolveType(name, ty.Elt)
		if err != nil {
			return "", "", err
		}
		typ := fmt.Sprintf("[]%s", resolved)
		return typ, defaultValueMap[typ], nil
	default:
		return "", "", fixgen_errors.NewUnsupportedTypeError(name, fmt.Sprintf("%T", exprType))
	}
}

func parse(name string, st *ast.StructType) (*gen.StructInfo, error) {
	fields := make([]*gen.Field, 0, len(st.Fields.List))
	for _, f := range st.Fields.List {
		names := make([]string, 0, len(f.Names))
		for _, n := range f.Names {
			names = append(names, n.Name)
		}

		typ, defaultValue, err := resolveType(strings.Join(names, ", "), f.Type)
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
		name := typ
		if len(f.Names) != 0 {
			name = f.Names[0].Name
		}

		// means unexported field
		if unicode.IsLower(rune(name[0])) {
			continue
		}

		fields = append(fields, &gen.Field{
			Name:         name,
			Type:         typ,
			DefaultValue: defaultValue,
			Tags:         tags,
		})
	}

	return &gen.StructInfo{
		Name:   name,
		Fields: fields,
	}, nil
}

var defaultValueMap = map[string]string{
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
