package utils_test

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/utils"
)

func TestTesting_CompareArraysUnordered(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := []int{5, 4, 3, 2, 1}
	utils.Assert(t, utils.CompareArraysUnordered(a, b), "Arrays are not equal")
}

func TestTesting_CompareArraysUnordered_Duplicates(t *testing.T) {
	a := []int{1, 2, 3, 3, 5}
	b := []int{5, 3, 3, 2, 1}
	utils.Assert(t, utils.CompareArraysUnordered(a, b), "Arrays are not equal")
}

func TestTesting_CompareArraysUnordered_DifferentLengths(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := []int{5, 4, 3, 2}
	utils.Assert(t, !utils.CompareArraysUnordered(a, b), "Arrays are equal")
}
