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
		siList, err := loader.Load(filepath.Join(genPkgInfo.PackageLocation, e.Name()))
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
	})
	if err != nil {
		return nil, fmt.Errorf("failed to ExecuteTemplate: %+w", err)
	}

	for _, f := range files {
		content, err := Format(append(header, f.Content...))
		if err != nil {
			return nil, fmt.Errorf("failed to Format: %+w", err)
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
