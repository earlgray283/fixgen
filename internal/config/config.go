package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	RandPackage        string              `yaml:"randPackage,omitempty"`
	DefaultValuePolicy *DefaultValuePolicy `yaml:"defaultValuePolicy,omitempty"`
	Structs            Structs             `yaml:"structs,omitempty"`
	Imports            []*Import           `yaml:"imports,omitempty"`
}

type DefaultValuePolicy struct {
	Type      DefaultValuePolicyType `yaml:"type,omitempty"`
	CustomMap map[string]string      `yaml:"customMap,omitempty"`
}

type DefaultValuePolicyType string

const (
	DefaultValuePolicyTypeRandv2 DefaultValuePolicyType = "rand" // default
	DefaultValuePolicyTypeRandv1 DefaultValuePolicyType = "randlegacy"
	DefaultValuePolicyTypeZero   DefaultValuePolicyType = "zero"
	DefaultValuePolicyTypeCustom DefaultValuePolicyType = "custom"
)

type Structs map[string]*Struct

type Struct struct {
	Fields map[string]*Field `yaml:"fields,omitempty"`
}

type Field struct {
	Value          any    `yaml:"value,omitempty"`
	Expr           string `yaml:"expr,omitempty"`
	IsModifiedCond string `yaml:"isModifiedCond,omitempty"` // should be expr
	MustOverwrite  bool   `yaml:"overwrite,omitempty"`
}

type Import struct {
	Alias   string `yaml:"alias,omitempty"`
	Package string `yaml:"package,omitempty"`
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
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

const (
	pkgMathRandV1 = "math/rand"
	pkgMathRandV2 = "math/rand/v2"
)

func validateConfig(c *Config) error {
	if c.RandPackage == "" {
		c.RandPackage = pkgMathRandV2
	} else {
		switch c.RandPackage {
		case pkgMathRandV1, pkgMathRandV2:
		default:
			return errors.New("key `randPackage` should be \"math/rand\" or \"math/rand/v2\"")
		}
	}

	if c.DefaultValuePolicy == nil {
		c.DefaultValuePolicy = &DefaultValuePolicy{Type: DefaultValuePolicyTypeRandv2}
	} else {
		switch c.DefaultValuePolicy.Type {
		case DefaultValuePolicyTypeRandv1, DefaultValuePolicyTypeRandv2, DefaultValuePolicyTypeZero, DefaultValuePolicyTypeCustom:
		default:
			return errors.New("key `defaultValuePolicy.type` should be `rand` or `randlegacy` or `zero` or `custom")
		}
	}

	return nil
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
