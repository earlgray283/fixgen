package internal

import (
	"embed"
	"text/template"
)

var (
	//go:embed templates
	embedFs         embed.FS
	TmplMockEntFile *template.Template
	TmplMockYoFile  *template.Template
	TmplCommonFile  *template.Template
)

func init() {
	TmplMockEntFile = parseFS("templates/ent.go.tmpl")
	TmplMockYoFile = parseFS("templates/yo.go.tmpl")
	TmplCommonFile = parseFS("templates/common.go.tmpl")
}

func parseFS(pattern string) *template.Template {
	tmpl, err := template.ParseFS(embedFs, pattern)
	if err != nil {
		panic(err)
	}
	return tmpl
}
