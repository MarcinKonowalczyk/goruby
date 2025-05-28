package transformer_tests

import (
	"path"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/integration_tests/utils"
	transformer_pipeline "github.com/MarcinKonowalczyk/goruby/pipelines/transformer"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
	"github.com/MarcinKonowalczyk/goruby/testutils/ruby"
)

// test files folder relative to this package
const TEST_FILES_FOLDER = "test_files"

func TestAll(t *testing.T) {
	test_files := utils.FindTestFiles(TEST_FILES_FOLDER, ".rb")
	ruby, err := ruby.FindRuby()
	if err != nil {
		t.Skip("Ruby interpreter not found, skipping tests")
	}
	for _, test_file := range test_files {
		base := path.Base(test_file)
		t.Run(base, func(t *testing.T) {
			// run the original Ruby file
			before, err := ruby.RunFile(test_file)
			assert.NoError(t, err)

			// Transform
			transformed, err := transformer_pipeline.Transform(test_file, nil)
			assert.NoError(t, err)

			// Run the transformed code
			after, err := ruby.RunCode(transformed)
			assert.NoError(t, err)

			// Make sure the output is the same
			assert.Equal(t, before, after)
		})
	}

}
