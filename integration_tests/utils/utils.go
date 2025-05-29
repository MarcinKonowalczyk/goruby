package utils

import (
	"fmt"
	"os"
	"os/exec"
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

const GORUBY_BIN = "goruby"

// Run 'go build -o goruby ./cmd/goruby/main.go -C <package_root>'
func CompileGoRuby(force bool) (string, error) {
	package_root := go_mod_dir()
	goruby_path := filepath.Join(package_root, GORUBY_BIN)
	if _, err := os.Stat(goruby_path); err == nil {
		// goruby binary already exists, no need to compile
		if !force {
			return goruby_path, nil
		} else {
			// remove existing binary
			if err := os.Remove(goruby_path); err != nil {
				return "", err
			}
		}
	}

	cmd := exec.Command("go", "build", "-C", package_root, "-o", GORUBY_BIN, "./cmd/goruby/main.go")
	// fmt.Println("Running command:", cmd.String())
	if output, err := cmd.CombinedOutput(); err != nil {
		// fmt.Println("Error running go build:", err)
		return "", err
	} else if len(output) > 0 {
		// fmt.Println("Output from go build:", string(output))
	}

	if _, err := os.Stat(goruby_path); err != nil {
		return "", fmt.Errorf("failed to compile goruby: %w", err)
	}

	return goruby_path, nil
}
