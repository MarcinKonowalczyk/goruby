package object

import (
	"context"
	"os"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/env"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
)

func TestFileExpandPath(t *testing.T) {
	t.Run("one arg flavour", func(t *testing.T) {
		env := env.NewEnvironment[ruby.Object]()
		ctx := call.NewContext[ruby.Object](context.Background(), fileClass, env)
		filename := NewString("./fixtures/testfile.rb")

		result, err := fileExpandPath(ctx, filename)

		assert.NoError(t, err)

		cwd, err := os.Getwd()
		if err != nil {
			t.Skip("Cannot determine working directory")
		}
		expected := NewString(cwd + "/fixtures/testfile.rb")

		assert.EqualCmpAny(t, result, expected, CompareRubyObjectsForTests)
	})
	t.Run("two arg flavour", func(t *testing.T) {
		env := env.NewEnvironment[ruby.Object]()
		ctx := call.NewContext[ruby.Object](context.Background(), fileClass, env)
		filename := NewString("../../main.go")
		dirname := NewString("object/fixtures/")

		result, err := fileExpandPath(ctx, filename, dirname)
		assert.NoError(t, err)

		cwd, err := os.Getwd()
		if err != nil {
			t.Skip("Cannot determine working directory")
		}
		expected := NewString(cwd + "/main.go")
		assert.EqualCmpAny(t, result, expected, CompareRubyObjectsForTests)
	})
}

func TestFileDirname(t *testing.T) {
	ctx := call.NewContext[ruby.Object](context.Background(), fileClass, NewMainEnvironment())
	filename := NewString("/var/log/foo.log")
	result, err := fileDirname(ctx, filename)
	assert.NoError(t, err)
	expected := NewString("/var/log")
	assert.EqualCmpAny(t, result, expected, CompareRubyObjectsForTests)
}
