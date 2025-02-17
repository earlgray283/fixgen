package gen

import (
	"fmt"

	"golang.org/x/tools/imports"

	"github.com/earlgray283/fixgen/internal/templates"
)

type Generator interface {
	Generate() ([]*File, error)
	GenPackageInfo() *GenPackageInfo
	IsExperimental() bool
}

type File struct {
	Name    string
	Content []byte
}

type GenPackageInfo struct {
	PackagePath  string // e.g. github.com/earlgray283/pj-todo/models
	PackageAlias string // e.g. yo_gen
}

func GenerateWithFormat[G Generator](g G, opts ...OptionFunc) ([]*File, error) {
	opt := defaultOption()
	opt.applyOptionFuncs(opts...)

	files, err := g.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generator.Generate: %+w", err)
	}
	commonFile, err := generateCommonFile()
	if err != nil {
		return nil, err
	}
	files = append(files, commonFile)

	genPkgInfo := g.GenPackageInfo()
	header, err := templates.Execute(templates.TmplHeaderFile, map[string]string{
		"PackageName": opt.packageName,
		"GenPkgAlias": genPkgInfo.PackageAlias,
		"GenPkgPath":  genPkgInfo.PackagePath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to ExecuteTemplate: %+w", err)
	}

	for _, f := range files {
		f.Content = append(header, f.Content...)
		content, err := Format(f.Content)
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
