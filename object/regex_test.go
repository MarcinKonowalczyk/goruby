package object

import "testing"

func TestRegex_hashKey(t *testing.T) {
	hello1 := &Regex{Value: "Hello World"}
	hello2 := &Regex{Value: "Hello World"}
	diff1 := &Regex{Value: "My name is johnny"}
	diff2 := &Regex{Value: "My name is johnny"}

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

func TestRegex_Inspect(t *testing.T) {
	hello := &Regex{Value: "Hello World", Modifiers: "i"}

	if hello.Inspect() != "/Hello World/i" {
		t.Errorf("Inspect should return the Value")
	}
}
