package transformer_tests

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/integration_tests/utils"
	transformer_pipeline "github.com/MarcinKonowalczyk/goruby/pipelines/transformer"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
	"github.com/MarcinKonowalczyk/goruby/testutils/combinatorics"
	"github.com/MarcinKonowalczyk/goruby/testutils/ruby"
	"github.com/MarcinKonowalczyk/goruby/transformer"
)

// test files folder relative to this package
const TEST_FILES_FOLDER = "test_files"
const SAVE_INTERMEDIATE_FILES = true
const TRANSFORMED_SUFFIX = "transformed"

func findTestFiles() []string {
	test_files := utils.FindTestFiles(TEST_FILES_FOLDER, ".rb")

	// Filter _transformed.rb files
	filtered_test_files := []string{}
	for _, file := range test_files {
		if strings.HasSuffix(file, fmt.Sprintf("_%s.rb", TRANSFORMED_SUFFIX)) {
			continue
		}
		filtered_test_files = append(filtered_test_files, file)
	}
	return filtered_test_files
}

func runTest(
	t *testing.T,
	rb ruby.Ruby,
	test_file string,
	intermediate_suffix string,
	stages []transformer.Stage,
) {
	t.Helper()
	// run the original Ruby file
	before, err := rb.RunFile(test_file)
	assert.NoError(t, err)

	// Transform
	fmt.Println(stages)
	transformed, err := transformer_pipeline.TransformStages(test_file, nil, stages)
	assert.NoError(t, err)

	if intermediate_suffix != "" {
		// Save the transformed code to a file
		base := path.Base(test_file)
		transformed_file := path.Join(TEST_FILES_FOLDER, base[:len(base)-3]+intermediate_suffix+".rb")
		err = os.WriteFile(transformed_file, []byte(transformed), 0644)
		assert.NoError(t, err)
	}
	// Run the transformed code
	after, err := rb.RunCode(transformed)
	assert.NoError(t, err)

	// Make sure the output is the same
	assert.Equal(t, before, after)
}

func TestAll(t *testing.T) {
	test_files := findTestFiles()
	rb, err := ruby.FindRuby()
	if err != nil {
		t.Skip("Ruby interpreter not found, skipping tests")
	}
	for _, test_file := range test_files {
		base := path.Base(test_file)
		t.Run(base, func(t *testing.T) {
			t.Parallel() // Run each test in parallel
			runTest(t, rb, test_file, "", transformer.ALL_STAGES)
		})
	}
}

func combinationsToString(combination []transformer.Stage) string {
	var sb strings.Builder
	for i, stage := range combination {
		if i > 0 {
			sb.WriteString("-")
		}
		var stage_string string
		stage_string = string(stage)
		stage_string = strings.Replace(stage_string, " ", "-", -1)
		stage_string = strings.Replace(stage_string, "_", "-", -1)
		sb.WriteString(stage_string)
	}
	return sb.String()
}
func TestAllStageCombinations(t *testing.T) {
	test_files := findTestFiles()
	rb, err := ruby.FindRuby()
	if err != nil {
		t.Skip("Ruby interpreter not found, skipping tests")
	}
	all_stages := transformer.ALL_STAGES
	for i := 1; i <= len(all_stages); i++ {
		combinations := combinatorics.CombinationsSorted(all_stages, i, func(a, b transformer.Stage) bool { return a < b })
		fmt.Printf("Found %d combinations of %d stages\n", len(combinations), i)
		for _, combination := range combinations {
			combination_string := combinationsToString(combination)
			for _, test_file := range test_files {
				base := path.Base(test_file)
				t.Run(
					strings.Join([]string{base, fmt.Sprintf("%d", i), combination_string}, "/"),
					func(t *testing.T) {
						fmt.Printf("Running test %s with stages %s\n", base, combination_string)
						intermediate_suffix := fmt.Sprintf("_%s_%s", combination_string, TRANSFORMED_SUFFIX)
						runTest(t, rb, test_file, intermediate_suffix, combination)
					})
			}
		}
	}

	// for _, test_file := range test_files {
	// 	base := path.Base(test_file)
	// 	t.Run(base, func(t *testing.T) {
	// 		t.Parallel() // Run each test in parallel
	// 		runTest(t, rb, test_file, transformer.ALL_STAGES)
	// 	})
	// }
}
