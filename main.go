package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/earlgray283/fixgen/internal/gen"
	gen_ent "github.com/earlgray283/fixgen/internal/gen/ent"
	gen_yo "github.com/earlgray283/fixgen/internal/gen/yo"
)

type Flags struct {
	Prefix                string
	Ext                   string
	PackageName           string
	DestDir               string
	CleanIfFailed         bool
	ConfirmIfExperimental bool
	UseContext            bool
	UsePointerModifier    bool
}

func parseFlags() *Flags {
	f := &Flags{}

	flag.StringVar(&f.Prefix, "prefix", "", "")
	flag.StringVar(&f.Ext, "ext", ".gen.go", "")
	flag.StringVar(&f.PackageName, "pkgname", "fixture", "")
	flag.StringVar(&f.DestDir, "dest-dir", ".", "the path the destination directory is created")
	flag.BoolVar(&f.CleanIfFailed, "clean-if-failed", false, "clean the directory and files if failed")
	flag.BoolVar(&f.ConfirmIfExperimental, "confirm-if-experimental", true, "confirm before generation if the generator is experimental")
	flag.BoolVar(&f.UseContext, "use-context", false, "provide `context.Context` argument")
	flag.BoolVar(&f.UsePointerModifier, "use-pointer-modifier", true, "")

	flag.Parse()

	return f
}

func main() {
	flgs := parseFlags()

	generatorType := flag.Arg(0)
	if generatorType == "" {
		eprintf("target must be specified(ent, yo)\n")
		os.Exit(1)
	}

	generator, err := loadGenerator(generatorType, ".", flgs.UseContext, flgs.UsePointerModifier)
	if err != nil {
		eprintf("failed to loadGenerator: %v\n", err)
		os.Exit(1)
	}

	if flgs.ConfirmIfExperimental && generator.IsExperimental() {
		if !yesNo(fmt.Sprintf("The generator \"%s\" is experimental.\nProceed?[y/N]", generatorType)) {
			os.Exit(1)
		}
	}

	files, err := gen.GenerateWithFormat(generator, gen.WithPackageName(flgs.PackageName))
	if err != nil {
		eprintf("failed to generator.Generate: %v\n", err)
		os.Exit(1)
	}

	packageDirPath := filepath.Join(flgs.DestDir, flgs.PackageName)
	if err := createDirIfNotExists(packageDirPath); err != nil {
		eprintf("failed to CreateDirIfNotExists: %v\n", err)
		os.Exit(1)
	}
	for _, f := range files {
		fileName := buildFileName(flgs.DestDir, flgs.PackageName, flgs.Prefix, f.Name, flgs.Ext)
		if err := saveFile(fileName, f.Content); err != nil {
			if flgs.CleanIfFailed {
				_ = os.RemoveAll(packageDirPath)
			}
			eprintf("failed to FormatAndSaveFile: %v\n", err)
			os.Exit(1)
		}
	}
}

func loadGenerator(typ, workDir string, useContext, usePointerModifier bool) (gen.Generator, error) {
	switch typ {
	case "ent":
		return gen_ent.NewGenerator(workDir)
	case "yo":
		return gen_yo.NewGenerator(workDir, useContext, usePointerModifier)
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

func buildFileName(destDir, packageName, prefix, name, ext string) string {
	return filepath.Join(destDir, packageName, fmt.Sprintf("%s%s%s", prefix, name, ext))
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
		if err := os.MkdirAll(p, 0755); err != nil {
			return err
		}
	}
	return nil
}
