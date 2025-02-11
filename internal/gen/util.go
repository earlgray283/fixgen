package gen

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

func FindAndReadDirByFileName(rootDir, fileName string) (string, []string, error) {
	return findAndReadDir(rootDir, func(d fs.DirEntry) bool {
		return !d.IsDir() && d.Name() == fileName
	})
}

func FindDir(rootDir, dirName string) (string, error) {
	return findByKey(rootDir, func(d fs.DirEntry) bool {
		return d.IsDir() && d.Name() == dirName
	})
}

func FindFile(rootDir, fileName string) (string, error) {
	return findByKey(rootDir, func(d fs.DirEntry) bool {
		return !d.IsDir() && d.Name() == fileName
	})
}

func findAndReadDir(rootDir string, keyFunc func(d fs.DirEntry) bool) (string, []string, error) {
	keyPath, err := findByKey(rootDir, keyFunc)
	if err != nil {
		return "", nil, err
	}
	dirPath := filepath.Dir(keyPath)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", nil, err
	}

	filepaths := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		filepaths = append(filepaths, filepath.Join(dirPath, e.Name()))
	}

	return dirPath, filepaths, nil
}

func findByKey(rootDir string, keyFunc func(d fs.DirEntry) bool) (string, error) {
	var p string
	if err := fs.WalkDir(os.DirFS(rootDir), ".", func(path string, d fs.DirEntry, err error) error {
		if keyFunc(d) {
			p = path
			return fs.SkipAll
		}
		return nil
	}); err != nil {
		return "", err
	}
	if p == "" {
		return "", fmt.Errorf("no such file or directory")
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
