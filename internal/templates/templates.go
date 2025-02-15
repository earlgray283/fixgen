package templates

import (
	"embed"
	"text/template"
)

var (
	//go:embed *
	embedFs         embed.FS
	TmplMockEntFile *template.Template
	TmplMockYoFile  *template.Template
	TmplCommonFile  *template.Template
	TmplHeaderFile  *template.Template
)

func init() {
	TmplMockEntFile = parseFS("ent.go.tmpl")
	TmplMockYoFile = parseFS("yo.go.tmpl")
	TmplCommonFile = parseFS("common.go.tmpl")
	TmplHeaderFile = parseFS("header.go.tmpl")
}

func parseFS(pattern string) *template.Template {
	tmpl, err := template.ParseFS(embedFs, pattern)
	if err != nil {
		panic(err)
	}
	return tmpl
}
