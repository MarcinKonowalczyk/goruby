package ruby

import (
	"os"
	"os/exec"
)

type Ruby interface {
	RunCode(code string) (string, error)
	RunFile(filename string) (string, error)
}

type ruby_interpreter struct {
	ruby_path string
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

func FindRuby() (Ruby, error) {
	ruby_path := os.Getenv("RUBY_PATH")
	if ruby_path == "" {
		ruby_path = "ruby" // default to 'ruby' in PATH
	}
	cmd := exec.Command(ruby_path, "--version")
	if err := cmd.Run(); err != nil {
		return &ruby_interpreter{}, err // Ruby not found
	}
	return &ruby_interpreter{ruby_path: ruby_path}, nil
}
