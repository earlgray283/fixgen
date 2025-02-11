package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/earlgray283/fixgen/internal/gen"
	gen_ent "github.com/earlgray283/fixgen/internal/gen/ent"
	gen_yo "github.com/earlgray283/fixgen/internal/gen/yo"
)

type Flags struct {
	GoFilePrefix          string
	PackageName           string
	DestDir               string
	CleanIfFailed         bool
	ConfirmIfExperimental bool
}

func parseFlags() *Flags {
	f := &Flags{}

	flag.StringVar(&f.GoFilePrefix, "go-file-prefix", "mock_", "")
	flag.StringVar(&f.PackageName, "pkgname", "fixture", "")
	flag.StringVar(&f.DestDir, "dest-dir", ".", "the path the destination directory is created")
	flag.BoolVar(&f.CleanIfFailed, "clean-if-failed", false, "clean the directory and files if failed")
	flag.BoolVar(&f.ConfirmIfExperimental, "confirm-if-experimental", true, "confirm before generation if the generator is experimental")

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

	generator, err := loadGenerator(generatorType)
	if err != nil {
		eprintf("failed to loadGenerator: %v\n", err)
		os.Exit(1)
	}

	if flgs.ConfirmIfExperimental && generator.IsExperimental() {
		if !yesNo(fmt.Sprintf("The generator \"%s\" is experimental.\nProceed?[y/N]", generatorType)) {
			os.Exit(1)
		}
	}

	files, err := gen.GenerateWithFormat(generator)
	if err != nil {
		eprintf("failed to generator.Generate: %v\n", err)
		os.Exit(1)
	}

	packageDirPath := filepath.Join(flgs.DestDir, flgs.PackageName)
	if err := gen.CreateDirIfNotExists(packageDirPath); err != nil {
		eprintf("failed to CreateDirIfNotExists: %v\n", err)
		os.Exit(1)
	}
	for _, f := range files {
		fileName := buildFileName(flgs.DestDir, flgs.PackageName, flgs.GoFilePrefix, f.Name)
		if err := gen.SaveFile(fileName, f.Content); err != nil {
			if flgs.CleanIfFailed {
				_ = os.RemoveAll(packageDirPath)
			}
			eprintf("failed to FormatAndSaveFile: %v\n", err)
			os.Exit(1)
		}
	}
}

func loadGenerator(typ string) (gen.Generator, error) {
	switch typ {
	case "ent":
		return gen_ent.NewGenerator()
	case "yo":
		return gen_yo.NewGenerator()
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

func buildFileName(destDir, packageName, goFilePrefix, name string) string {
	return filepath.Join(destDir, packageName, fmt.Sprintf("%s%s.go", goFilePrefix, name))
}
