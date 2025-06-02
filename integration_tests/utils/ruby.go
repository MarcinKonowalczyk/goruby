package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/testutils/ruby"
)

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
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", err
	} else if len(output) > 0 {
	}

	if _, err := os.Stat(goruby_path); err != nil {
		return "", fmt.Errorf("failed to compile goruby: %w", err)
	}

	return goruby_path, nil
}

func InitGoRuby() (ruby.Ruby, error) {
	path, err := CompileGoRuby(true)
	if err != nil {
		return nil, fmt.Errorf("failed to compile GoRuby: %w", err)
	}

	rb, err := ruby.FindRuby(path)
	if err != nil {
		return nil, fmt.Errorf("failed to find GoRuby binary: %w", err)
	}

	// make sure we are actually running GoRuby
	if !strings.Contains(strings.ToLower(rb.Version()), "goruby") {
		return nil, fmt.Errorf("GoRuby binary does not contain 'goruby' in version string: %s", rb.Version())
	}
	return rb, nil
}
