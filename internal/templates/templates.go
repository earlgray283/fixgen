package templates

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

var (
	//go:embed *
	embedFs               embed.FS
	TmplEntFile           *template.Template
	TmplYoFile            *template.Template
	TmplCommonFile        *template.Template
	TmplHeaderFile        *template.Template
	TmplStructsFile       *template.Template
	TmplStructsCommonFile *template.Template
)

func init() {
	TmplEntFile = parseFS("ent.go.tmpl")
	TmplYoFile = parseFS("yo.go.tmpl")
	TmplCommonFile = parseFS("common.go.tmpl")
	TmplHeaderFile = parseFS("header.go.tmpl")
	TmplStructsFile = parseFS("structs.go.tmpl")
	TmplStructsCommonFile = parseFS("structs_common.go.tmpl")
}

func parseFS(pattern string) *template.Template {
	tmpl, err := template.New(pattern).ParseFS(embedFs, pattern)
	if err != nil {
		panic(err)
	}
	return tmpl
}

func Execute(tmpl *template.Template, data map[string]any) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template `%s`: %+w", tmpl.Name(), err)
	}

	return buf.Bytes(), nil
}
