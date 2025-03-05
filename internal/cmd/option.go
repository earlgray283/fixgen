package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/earlgray283/fixgen/internal/gen"
)

type Options struct {
	prefix           string
	ext              string
	packageName      string
	destDir          string
	skipConfirm      bool
	useContext       bool
	useValueModifier bool
	config           string
}

func (o *Options) packageDir() string {
	return filepath.Join(o.destDir, o.packageName)
}

func (o *Options) fullFilePath(name string) string {
	return filepath.Join(o.destDir, o.packageName, fmt.Sprintf("%s%s%s", o.prefix, name, o.ext))
}

func (o *Options) buildGenOptions() []gen.OptionFunc {
	opts := make([]gen.OptionFunc, 0)

	opts = append(opts, gen.WithPackageName(o.packageName))
	if o.useContext {
		opts = append(opts, gen.UseContext())
	}
	if o.useValueModifier {
		opts = append(opts, gen.UseValueModifier())
	}

	return opts
}

type GeneratorType string

const (
	GeneratorTypeYo  GeneratorType = "yo"
	GeneratorTypeEnt GeneratorType = "ent"
)

var generatorTypes = []GeneratorType{GeneratorTypeYo, GeneratorTypeEnt}

func (t GeneratorType) Validate() error {
	if t == "" {
		return errors.New("generatorType is empty")
	}
	if slices.Contains(generatorTypes, t) {
		return fmt.Errorf("generatorType `%s` is invalid(supported types: %v)", t, generatorTypes)
	}
	return nil
}
