package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

func this_package_path(N int) string {
	_, file, _, ok := runtime.Caller(1 + N)
	if !ok {
		return ""
	}
	return filepath.Dir(file)
}

// Return the path of this package
func ThisPackagePath() string {
	return this_package_path(0)
}

// Find test files in a given folder, relative to the package path.
// If extension is non-empty, only files with that extension will be returned.
// Extension must have the leading dot, e.g. ".rb".
func FindTestFiles(
	folder string,
	extension string,
) []string {
	test_files := []string{}
	this_path := this_package_path(1)
	test_files_folder := filepath.Join(this_path, folder)

	err := filepath.Walk(test_files_folder, func(path string, info os.FileInfo, err error) error {
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
