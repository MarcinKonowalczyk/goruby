package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func this_package_path(N int) string {
	_, file, _, ok := runtime.Caller(1 + N)
	if !ok {
		return ""
	}
	return filepath.Dir(file)
}

func go_mod_dir() string {
	// Find the directory of the go.mod file
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	dir := filepath.Dir(file)
	for dir != "" {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir { // reached root
			break
		}
		dir = parent
	}
	return ""
}

// Return the path of this package
//
//go:noinline
func ThisPackagePath() string {
	return this_package_path(1)
}

// check whether path ends with the given suffix.
func TrimPathSuffix(path string, suffix string) (string, error) {
	path_parts := strings.Split(path, string(filepath.Separator))
	suffix_parts := strings.Split(suffix, string(filepath.Separator))
	if len(path_parts) < len(suffix_parts) {
		return "", fmt.Errorf(
			"path '%s' is shorter than suffix '%s'",
			path, suffix,
		)
	}
	for i := 0; i < len(suffix_parts); i++ {
		a := path_parts[len(path_parts)-1-i]
		b := suffix_parts[len(suffix_parts)-1-i]
		if a != b {
			return "", fmt.Errorf(
				"path '%s' does not end with suffix '%s'",
				path, suffix,
			)
		}
	}
	return strings.Join(path_parts[:len(path_parts)-len(suffix_parts)], string(filepath.Separator)), nil
}

// Find test files in a given folder, relative to the package path.
// If extension is non-empty, only files with that extension will be returned.
// Extension must have the leading dot, e.g. ".rb".
//
//go:noinline
func FindTestFiles(
	folder string,
	extension string,
) []string {
	test_files := []string{}

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if extension != "" && filepath.Ext(path) != extension {
			return nil
		}
		test_files = append(test_files, path)
		return nil
	})

	if err != nil {
		panic(err)
	}
	return test_files
}
