package gen

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

func FindAndReadDirByFileName(rootDir, fileName string) (string, []string, error) {
	var dirPath string
	if err := fs.WalkDir(os.DirFS(rootDir), ".", func(path string, d fs.DirEntry, err error) error {
		if d.Name() == fileName {
			dirPath = filepath.Dir(path)
			return fs.SkipAll
		}
		return nil
	}); err != nil {
		return "", nil, err
	}
	if dirPath == "" {
		return "", nil, fmt.Errorf("failed to find %s", fileName)
	}

	entries, err := os.ReadDir(filepath.Join(rootDir, dirPath))
	if err != nil {
		return "", nil, err
	}

	filepaths := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		filepaths = append(filepaths, filepath.Join(rootDir, dirPath, e.Name()))
	}

	return dirPath, filepaths, nil
}

func FindFilePath(rootDir, fileName string) (string, error) {
	var p string
	if err := fs.WalkDir(os.DirFS(rootDir), ".", func(path string, d fs.DirEntry, err error) error {
		if d.Name() == fileName {
			p = path
			return fs.SkipAll
		}
		return nil
	}); err != nil {
		return "", err
	}
	if p == "" {
		return "", fmt.Errorf("failed to find %s", fileName)
	}

	return filepath.Join(rootDir, p), nil
}

func LoadGoModulePath(dirPath string) (string, error) {
	b, err := os.ReadFile(filepath.Join(dirPath, "go.mod"))
	if err != nil {
		return "", err
	}
	f, err := modfile.Parse("", b, nil)
	if err != nil {
		return "", err
	}
	return f.Module.Mod.Path, nil
}
