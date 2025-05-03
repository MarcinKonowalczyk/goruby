package object

import (
	"os"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/utils"
)

func TestFileExpandPath(t *testing.T) {
	t.Run("one arg flavour", func(t *testing.T) {
		env := NewEnvironment()
		context := &callContext{
			receiver: fileClass,
			env:      env,
		}
		filename := &String{Value: "./fixtures/testfile.rb"}

		result, err := fileExpandPath(context, filename)

		utils.AssertNoError(t, err)

		cwd, err := os.Getwd()
		if err != nil {
			t.Skip("Cannot determine working directory")
		}
		expected := &String{Value: cwd + "/fixtures/testfile.rb"}

		utils.AssertEqualCmpAny(t, result, expected, CompareRubyObjectsForTests)
	})
	t.Run("two arg flavour", func(t *testing.T) {
		env := NewEnvironment()
		context := &callContext{
			receiver: fileClass,
			env:      env,
		}
		filename := &String{Value: "../../main.go"}
		dirname := &String{Value: "object/fixtures/"}

		result, err := fileExpandPath(context, filename, dirname)

		utils.AssertNoError(t, err)

		cwd, err := os.Getwd()
		if err != nil {
			t.Skip("Cannot determine working directory")
		}
		expected := &String{Value: cwd + "/main.go"}

		utils.AssertEqualCmpAny(t, result, expected, CompareRubyObjectsForTests)
	})
}

func TestFileDirname(t *testing.T) {
	context := &callContext{
		receiver: fileClass,
		env:      NewEnvironment(),
	}
	filename := &String{Value: "/var/log/foo.log"}

	result, err := fileDirname(context, filename)

	utils.AssertNoError(t, err)

	expected := &String{Value: "/var/log"}

	utils.AssertEqualCmpAny(t, result, expected, CompareRubyObjectsForTests)
}
