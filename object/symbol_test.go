package object

import "testing"

func TestSymbol_hashKey(t *testing.T) {
	hello1 := &Symbol{Value: "Hello World"}
	hello2 := &Symbol{Value: "Hello World"}
	diff1 := &Symbol{Value: "My name is johnny"}
	diff2 := &Symbol{Value: "My name is johnny"}

	if hello1.hashKey() != hello2.hashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.hashKey() != diff2.hashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.hashKey() == diff1.hashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestSymbolToS(t *testing.T) {
	context := &callContext{
		receiver: &Symbol{Value: "foo"},
	}

	result, err := symbolToS(context)

	checkError(t, err, nil)

	expected := &String{Value: "foo"}

	checkResult(t, result, expected)
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
	val, ok = SymbolToBool(&Symbol{Value: "foo"})
	if ok {
		t.Errorf("Expected false, got true")
	}
	if val {
		t.Errorf("Expected false, got true")
	}
	val, ok = SymbolToBool(&String{Value: "foo"})
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
