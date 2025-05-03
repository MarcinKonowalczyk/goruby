package object

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/utils"
)

func TestSymbol_hashKey(t *testing.T) {
	hello1 := NewSymbol("Hello World")
	hello2 := NewSymbol("Hello World")
	diff1 := NewSymbol("My name is johnny")
	diff2 := NewSymbol("My name is johnny")

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestSymbolToS(t *testing.T) {
	context := &callContext{
		receiver: NewSymbol("foo"),
	}

	result, err := symbolToS(context)

	utils.AssertNoError(t, err)

	expected := NewString("foo")

	utils.AssertEqualCmpAny(t, result, expected, CompareRubyObjectsForTests)
}

func TestSymbolToBool(t *testing.T) {
	val, ok := SymbolToBool(TRUE)
	if !ok {
		t.Errorf("Expected true, got false")
	}
	if val != true {
		t.Errorf("Expected true, got false")
	}
	val, ok = SymbolToBool(FALSE)
	if !ok {
		t.Errorf("Expected false, got true")
	}
	if val != false {
		t.Errorf("Expected false, got true")
	}
	val, ok = SymbolToBool(NewSymbol("foo"))
	if ok {
		t.Errorf("Expected false, got true")
	}
	if val {
		t.Errorf("Expected false, got true")
	}
	val, ok = SymbolToBool(NewString("foo"))
	if ok {
		t.Errorf("Expected false, got true")
	}
	if val {
		t.Errorf("Expected false, got true")
	}
	val, ok = SymbolToBool(nil)
	if ok {
		t.Errorf("Expected false, got true")
	}
	if val {
		t.Errorf("Expected false, got true")
	}
	var b *Symbol = nil
	val, ok = SymbolToBool(b)
	if ok {
		t.Errorf("Expected false, got true")
	}
	if val {
		t.Errorf("Expected false, got true")
	}
}
