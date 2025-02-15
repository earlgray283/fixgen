package config

import "text/template"

type Config struct {
	*CommonConfig
}

type CommonConfig struct {
	Fields          map[string]*FieldConfig       // fieldName=>FieldConfig
	CustomTemplates map[string]*template.Template // fileName=>Template
	Imports         []string
}

type FieldConfig struct {
	Value string `yaml:"value"`
}
