package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/earlgray283/fixgen/internal/config"
	"github.com/earlgray283/fixgen/internal/gen"
	gen_ent "github.com/earlgray283/fixgen/internal/gen/ent"
	gen_yo "github.com/earlgray283/fixgen/internal/gen/yo"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:  "fixgen <toolname> [flags]",
		Long: "Generate your fixture",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load(opts.config)
			if err != nil {
				return fmt.Errorf("failed to load config: %+w", err)
			}

			generatorType := GeneratorType(args[0])
			if err := generatorType.Validate(); err != nil {
				return fmt.Errorf("%+w", err)
			}

			generator, err := loadGenerator(generatorType, ".")
			if err != nil {
				return fmt.Errorf("failed to load generator: %+w", err)
			}

			if !opts.skipConfirm && !yesNo("Proceed?[y/N]") {
				return nil
			}

			files, err := gen.GenerateWithFormat(generator, c, opts.buildGenOptions()...)
			if err != nil {
				return fmt.Errorf("failed to generate: %+w", err)
			}

			packageDir := opts.packageDir()
			if err := createDirIfNotExists(packageDir); err != nil {
				return fmt.Errorf("failed to create dir: %+w", err)
			}
			for _, f := range files {
				if err := saveFile(opts.fullFilePath(f.Name), f.Content); err != nil {
					_ = os.RemoveAll(packageDir)
					return fmt.Errorf("failed to save file: %+w", err)
				}
			}

			return nil
		},
	}

	fs := cmd.Flags()
	fs.StringVarP(&opts.prefix, "prefix", "p", "", "prefix for the generated files")
	fs.StringVar(&opts.ext, "ext", ".gen.go", "extension for the generated files")
	fs.StringVar(&opts.packageName, "package", "fixture", "package name for the generated files")
	fs.StringVarP(&opts.destDir, "output", "o", ".", "output directory for the generated files")
	fs.BoolVarP(&opts.skipConfirm, "skip-confirm", "y", false, "skip confirmation")
	fs.BoolVar(&opts.useContext, "use-context", false, "add `context.Context` argument for the generated functions")
	fs.BoolVar(&opts.useValueModifier, "use-value-modifier", false, "use value modifier for the generated functions")
	fs.StringVarP(&opts.config, "config", "c", "fixgen.yaml", "config file path")

	return cmd
}

func loadGenerator(typ GeneratorType, workDir string) (gen.Generator, error) {
	switch typ {
	case GeneratorTypeEnt:
		return gen_ent.NewGenerator(workDir)
	case GeneratorTypeYo:
		return gen_yo.NewGenerator(workDir)
	default:
		return nil, fmt.Errorf("unrecognized generator type: %s", typ)
	}
}

func eprintf(format string, v ...any) {
	fmt.Fprintf(os.Stderr, format, v...)
}

func yesNo(prompt string) bool {
	var yn string
	eprintf("%s\n", prompt)
	_, _ = fmt.Scanln(&yn)
	return strings.ToLower(yn) == "y"
}

func saveFile(name string, content []byte) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, bytes.NewReader(content)); err != nil {
		return err
	}

	return nil
}

func createDirIfNotExists(p string) error {
	if _, err := os.Stat(p); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err := os.MkdirAll(p, 0o755); err != nil {
			return err
		}
	}
	return nil
}
