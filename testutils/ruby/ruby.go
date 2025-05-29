package ruby

import (
	"os"
	"os/exec"
)

type Ruby interface {
	Version() string
	RunCode(code string) (string, error)
	RunFile(filename string) (string, error)
}

type ruby_interpreter struct {
	ruby_path string
	version   string
}

func (r *ruby_interpreter) Version() string {
	return r.version
}

func (r *ruby_interpreter) RunCode(code string) (string, error) {
	tmpfile, err := os.CreateTemp("", "goruby_test_*.rb")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.WriteString(code); err != nil {
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		return "", err
	}

	return r.RunFile(tmpfile.Name())
}

func (r *ruby_interpreter) RunFile(filename string) (string, error) {
	cmd := exec.Command(r.ruby_path, filename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

var _ Ruby = (*ruby_interpreter)(nil)

func FindRuby(args ...string) (Ruby, error) {
	var ruby_path string
	if len(args) > 0 {
		ruby_path = args[0]
	}
	if ruby_path == "" {
		ruby_path = os.Getenv("RUBY_PATH")
	}
	if ruby_path == "" {
		ruby_path = "ruby" // default to 'ruby' in PATH
	}

	cmd := exec.Command(ruby_path, "--version")
	output, err := cmd.Output()
	if err != nil {
		return nil, err // Ruby not found
	}
	version := string(output)

	return &ruby_interpreter{
		ruby_path: ruby_path,
		version:   version,
	}, nil
}
