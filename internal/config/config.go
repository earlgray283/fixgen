package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Structs Structs   `yaml:"structs"`
	Imports []*Import `yaml:"imports"`
}

type Structs map[string]*Struct

type Struct struct {
	Fields map[string]*Field `yaml:"fields"`
}

type Field struct {
	Value          any    `yaml:"value"`
	Expr           string `yaml:"expr"`
	IsModifiedCond string `yaml:"isModifiedCond"` // should be expr
	MustOverwrite  bool   `yaml:"overwrite"`
}

type Import struct {
	Alias   string `yaml:"alias"`
	Package string `yaml:"package"`
}

func Load(name string) (*Config, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return newConfig(), nil
		}

		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func newConfig() *Config {
	return &Config{
		Structs: make(map[string]*Struct),
	}
}

func (f *Field) DefaultValue() (string, bool) {
	var zero any
	if f.Expr == "" && f.Value == zero {
		return "", false
	}

	if f.Expr != "" {
		return f.Expr, true
	}

	switch v := f.Value.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v), true
	default:
		return fmt.Sprintf("%v", v), true
	}
}
