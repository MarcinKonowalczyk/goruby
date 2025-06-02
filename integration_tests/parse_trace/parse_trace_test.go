package parse_trace_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/integration_tests/utils"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
)

const TEST_FILES_FOLDER = "test_files"

type testFilePair struct {
	rubyFile  string
	traceFile string
}

func findTestFiles() []testFilePair {
	this_path := utils.ThisPackagePath()
	test_files_folder := filepath.Join(this_path, TEST_FILES_FOLDER)

	ruby_files := utils.FindTestFiles(test_files_folder, ".rb")
	trace_files := utils.FindTestFiles(test_files_folder, ".trace")

	trace_files_set := make(map[string]struct{})
	for _, out_file := range trace_files {
		basename := filepath.Base(out_file)
		trace_files_set[basename] = struct{}{}
	}

	var test_pairs []testFilePair
	for i := 0; i < len(ruby_files); i++ {
		base_name := filepath.Base(ruby_files[i])
		trace_name := strings.TrimSuffix(base_name, ".rb") + ".trace"
		if _, ok := trace_files_set[trace_name]; ok {
			test_pairs = append(test_pairs, testFilePair{
				rubyFile:  ruby_files[i],
				traceFile: filepath.Join(test_files_folder, trace_name),
			})
		} else {
			fmt.Printf("No trace file for %s\n", ruby_files[i])
		}
	}
	return test_pairs
}

func readExpectedOutput(traceFile string) ([]byte, error) {
	expected, err := os.ReadFile(traceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read expected output from %s: %w", traceFile, err)
	}
	// Remove the header line if it exists
	if len(expected) > 0 && expected[0] == '#' {
		lines := strings.Split(string(expected), "\n")
		if len(lines) > 1 {
			expected = []byte(strings.Join(lines[1:], "\n"))
		} else {
			expected = []byte{}
		}
	}
	return expected, nil
}

func TestAll(t *testing.T) {
	test_files := findTestFiles()
	fmt.Printf("Found %d test files\n", len(test_files))

	grb, err := utils.InitGoRuby()
	if err != nil {
		t.Skipf("GoRuby not compiled, skipping tests: %v", err)
	} else {
		t.Logf("Using GoRuby binary: %s", grb.Version())
	}

	flags := []string{"--trace-parse=only-no-messages"}
	for _, pair := range test_files {
		out, err := grb.RunFile(pair.rubyFile, flags...)
		assert.NoError(t, err)
		expected, err := readExpectedOutput(pair.traceFile)
		assert.NoError(t, err)
		assert.EqualLineByLine(t, string(expected), out)
	}

}
