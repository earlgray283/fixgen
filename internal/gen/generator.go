package gen

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/tools/imports"

	"github.com/earlgray283/fixgen/internal/config"
	"github.com/earlgray283/fixgen/internal/load"
	"github.com/earlgray283/fixgen/internal/templates"
)

func GenerateWithFormat[G Generator](g G, c *config.Config, opts ...OptionFunc) ([]*File, error) {
	opt := defaultOption()
	opt.applyOptionFuncs(opts...)

	defaultValueMap := defaultValueMapRandv2
	useMathv1 := false
	if c.DefaultValuePolicy != nil {
		switch c.DefaultValuePolicy.Type {
		case config.DefaultValuePolicyTypeRandv1:
			defaultValueMap = defaultValueMapRandv1
			useMathv1 = true
		case config.DefaultValuePolicyTypeRandv2:
			defaultValueMap = defaultValueMapRandv2
		case config.DefaultValuePolicyTypeZero:
			defaultValueMap = defaultValueMapZero
		case config.DefaultValuePolicyTypeCustom:
			defaultValueMap = c.DefaultValuePolicy.CustomMap
		}
	}

	loader := load.New(c.Structs)
	genPkgInfo := g.PackageInfo()

	entries, err := os.ReadDir(genPkgInfo.PackageLocation)
	if err != nil {
		return nil, err
	}
	structInfos := make([]*load.StructInfo, 0)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		siList, err := loader.Load(filepath.Join(genPkgInfo.PackageLocation, e.Name()), defaultValueMap)
		if err != nil {
			return nil, err
		}
		structInfos = append(structInfos, siList...)
	}

	files, err := g.Generate(structInfos, map[string]any{
		"UseContext":       opt.useContext,
		"UseValueModifier": opt.useValueModifier,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generator.Generate: %+w", err)
	}
	commonFile, err := generateCommonFile()
	if err != nil {
		return nil, err
	}
	files = append(files, commonFile)

	header, err := templates.Execute(templates.TmplHeaderFile, map[string]any{
		"PackageName": opt.packageName,
		"GenPkgAlias": genPkgInfo.PackageAlias,
		"GenPkgPath":  genPkgInfo.PackagePath,
		"Imports":     c.Imports,
		"UseMathv1":   useMathv1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to ExecuteTemplate: %+w", err)
	}

	for _, f := range files {
		content, err := Format(append(header, f.Content...))
		if err != nil {
			return nil, fmt.Errorf("failed to Format(%s): %+w", f.Name, err)
		}
		f.Content = content
	}

	return files, nil
}

func generateCommonFile() (*File, error) {
	content, err := templates.Execute(templates.TmplCommonFile, nil)
	if err != nil {
		return nil, err
	}
	return &File{Name: "common", Content: content}, nil
}

func Format(content []byte) ([]byte, error) {
	formated, err := imports.Process("", content, &imports.Options{Comments: true})
	if err != nil {
		return nil, fmt.Errorf("failed to apply format: %w\n%s", err, content)
	}
	return formated, nil
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
	defaultValueMapZero = map[string]string{
		"int32":     `0`,
		"int64":     `0`,
		"uint32":    `0`,
		"uint64":    `0`,
		"float32":   `0`,
		"float64":   `0`,
		"string":    `""`,
		"[]byte":    `nil`,
		"bool":      `false`,
		"time.Time": `time.Time{}`,
	}
)
