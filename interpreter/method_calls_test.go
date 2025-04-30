package interpreter_test

import (
	"reflect"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/interpreter"
	"github.com/MarcinKonowalczyk/goruby/object"
)

func TestKernelTap(t *testing.T) {
	input := `
		x = []
		x.tap {|z|z.push(true)}
	`
	i := interpreter.NewInterpreter()

	evaluated, err := i.Interpret("", input)
	if err != nil {
		t.Logf("Expected no error, got %T:%v", err, err)
		t.FailNow()
	}

	expected := &object.Array{Elements: []object.RubyObject{object.TRUE}}
	actual := evaluated

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected result to equal\n%+#v\n\tgot\n%+#v\n", expected, actual)
		t.Fail()
	}
}
