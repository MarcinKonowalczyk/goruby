package compare_test

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/testutils/assert/compare"
)

func TestCompareArraysUnordered(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := []int{5, 4, 3, 2, 1}
	if !compare.ArraysUnordered(a, b) {
		t.Errorf("Expected arrays to be equal, but they are not: %v != %v", a, b)
	}
}

func TestCompareArraysUnordered_Duplicates(t *testing.T) {
	a := []int{1, 2, 3, 3, 5}
	b := []int{5, 3, 3, 2, 1}
	if !compare.ArraysUnordered(a, b) {
		t.Errorf("Expected arrays with duplicates to be equal, but they are not: %v != %v", a, b)
	}

}

func TestCompareArraysUnordered_DifferentLengths(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := []int{5, 4, 3, 2}
	if compare.ArraysUnordered(a, b) {
		t.Errorf("Expected arrays of different lengths to not be equal, but they are: %v == %v", a, b)
	}
}
