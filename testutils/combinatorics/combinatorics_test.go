package combinatorics_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/testutils/combinatorics"
)

func assert(t *testing.T, condition bool, msg string) {
	t.Helper()
	if !condition {
		t.Errorf("Assertion failed: %s", msg)
	}
}

func TestCombinations(t *testing.T) {
	myArray := []int{1, 2, 3, 4}
	t.Run("cmb 0", func(t *testing.T) {
		combinations := combinatorics.Combinations(myArray, 0)
		expected := [][]int{{}}
		assert(t, len(combinations) == 1, "Expected 1 combination for 4 elements taken 0 at a time")
		assert(t, reflect.DeepEqual(combinations, expected), "Combinations do not match expected values")
		fmt.Printf("Combinations of 4 elements taken 0 at a time: %v\n", combinations)
	})
	t.Run("cmb 1", func(t *testing.T) {
		combinations := combinatorics.Combinations(myArray, 1)
		expected := [][]int{{1}, {2}, {3}, {4}}
		assert(t, len(combinations) == 4, "Expected 4 combinations for 4 elements taken 1 at a time")
		assert(t, reflect.DeepEqual(combinations, expected), "Combinations do not match expected values")
	})
	t.Run("cmb 2", func(t *testing.T) {
		combinations := combinatorics.Combinations(myArray, 2)
		expected := [][]int{
			{1, 2}, {1, 3}, {1, 4},
			{2, 3}, {2, 4}, {3, 4},
		}
		assert(t, len(combinations) == 6, "Expected 6 combinations for 4 elements taken 2 at a time")
		assert(t, reflect.DeepEqual(combinations, expected), "Combinations do not match expected values")
	})
	t.Run("cmb 3", func(t *testing.T) {
		combinations := combinatorics.Combinations(myArray, 3)
		expected := [][]int{
			{1, 2, 3}, {1, 2, 4}, {1, 3, 4}, {2, 3, 4},
		}
		assert(t, len(combinations) == 4, "Expected 4 combinations for 4 elements taken 3 at a time")
		assert(t, reflect.DeepEqual(combinations, expected), "Combinations do not match expected values")
	})
	t.Run("cmb 4", func(t *testing.T) {
		combinations := combinatorics.Combinations(myArray, 4)
		expected := [][]int{{1, 2, 3, 4}}
		assert(t, len(combinations) == 1, "Expected 1 combination for 4 elements taken 4 at a time")
		assert(t, reflect.DeepEqual(combinations, expected), "Combinations do not match expected values")
	})
	t.Run("cmb 5", func(t *testing.T) {
		combinations := combinatorics.Combinations(myArray, 5)
		assert(t, len(combinations) == 0, "Expected no combinations for 4 elements taken 5 at a time")
	})

	// Testing combinations with more elements than available

	t.Run("cmb 6", func(t *testing.T) {
		combinations := combinatorics.Combinations(myArray, 6)
		assert(t, len(combinations) == 0, "Expected no combinations for 4 elements taken 6 at a time")
	})
	t.Run("cmb 7", func(t *testing.T) {
		combinations := combinatorics.Combinations(myArray, 7)
		assert(t, len(combinations) == 0, "Expected no combinations for 4 elements taken 7 at a time")
	})

	// Testing combinations with negative numbers

	t.Run("cmb -1", func(t *testing.T) {
		combinations := combinatorics.Combinations(myArray, -1)
		assert(t, len(combinations) == 0, "Expected no combinations for 4 elements taken -1 at a time")
	})
}

func TestCombinationsEmptyArray(t *testing.T) {
	myArray := []int{}
	t.Run("cmb 0", func(t *testing.T) {
		combinations := combinatorics.Combinations(myArray, 0)
		expected := [][]int{{}}
		assert(t, len(combinations) == 1, "Expected 1 combination for empty array taken 0 at a time")
		assert(t, reflect.DeepEqual(combinations, expected), "Combinations do not match expected values")
	})
	t.Run("cmb 1", func(t *testing.T) {
		combinations := combinatorics.Combinations(myArray, 1)
		assert(t, len(combinations) == 0, "Expected no combinations for empty array taken 1 at a time")
	})
	t.Run("cmb -1", func(t *testing.T) {
		combinations := combinatorics.Combinations(myArray, -1)
		assert(t, len(combinations) == 0, "Expected no combinations for empty array taken -1 at a time")
	})
}

func TestCombinationsOtherTypes(t *testing.T) {
	t.Run("cmb strings", func(t *testing.T) {
		myArray := []string{"a", "b", "c"}
		combinations := combinatorics.Combinations(myArray, 0)
		expected := [][]string{{}}
		assert(t, len(combinations) == 1, "Expected 1 combination for string array taken 0 at a time")
		assert(t, reflect.DeepEqual(combinations, expected), "Combinations do not match expected values")
	})

	t.Run("cmb structs", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		myArray := []Person{
			{"Alice", 30},
			{"Bob", 25},
			{"Charlie", 35},
		}
		combinations := combinatorics.Combinations(myArray, 2)
		expected := [][]Person{
			{{"Alice", 30}, {"Bob", 25}},
			{{"Alice", 30}, {"Charlie", 35}},
			{{"Bob", 25}, {"Charlie", 35}},
		}
		assert(t, len(combinations) == 3, "Expected 3 combinations for struct array taken 2 at a time")
		assert(t, reflect.DeepEqual(combinations, expected), "Combinations do not match expected values")
	})
}

func TestCombinationsSorted(t *testing.T) {
	myArray := []int{4, 1, 3, 2}
	lessFunc := func(a, b int) bool {
		return a < b
	}
	t.Run("cmb 0", func(t *testing.T) {
		combinations := combinatorics.CombinationsSorted(myArray, 0, lessFunc)
		expected := [][]int{{}}
		assert(t, len(combinations) == 1, "Expected 1 combination for 4 elements taken 0 at a time")
		assert(t, reflect.DeepEqual(combinations, expected), "Combinations do not match expected values")
		fmt.Printf("Combinations of 4 elements taken 0 at a time: %v\n", combinations)
	})
	t.Run("cmb 1", func(t *testing.T) {
		combinations := combinatorics.CombinationsSorted(myArray, 1, lessFunc)
		expected := [][]int{{1}, {2}, {3}, {4}}
		assert(t, len(combinations) == 4, "Expected 4 combinations for 4 elements taken 1 at a time")
		assert(t, reflect.DeepEqual(combinations, expected), "Combinations do not match expected values")
	})
	t.Run("cmb 2", func(t *testing.T) {
		combinations := combinatorics.CombinationsSorted(myArray, 2, lessFunc)
		expected := [][]int{
			{1, 2}, {1, 3}, {1, 4},
			{2, 3}, {2, 4}, {3, 4},
		}
		fmt.Println(combinations)
		fmt.Println(expected)
		assert(t, len(combinations) == 6, "Expected 6 combinations for 4 elements taken 2 at a time")
		assert(t, reflect.DeepEqual(combinations, expected), "Combinations do not match expected values")
	})
	t.Run("cmb 3", func(t *testing.T) {
		combinations := combinatorics.CombinationsSorted(myArray, 3, lessFunc)
		expected := [][]int{
			{1, 2, 3}, {1, 2, 4}, {1, 3, 4}, {2, 3, 4},
		}
		assert(t, len(combinations) == 4, "Expected 4 combinations for 4 elements taken 3 at a time")
		assert(t, reflect.DeepEqual(combinations, expected), "Combinations do not match expected values")
	})
	t.Run("cmb 4", func(t *testing.T) {
		combinations := combinatorics.CombinationsSorted(myArray, 4, lessFunc)
		expected := [][]int{{1, 2, 3, 4}}
		assert(t, len(combinations) == 1, "Expected 1 combination for 4 elements taken 4 at a time")
		assert(t, reflect.DeepEqual(combinations, expected), "Combinations do not match expected values")
	})
	t.Run("cmb 5", func(t *testing.T) {
		combinations := combinatorics.CombinationsSorted(myArray, 5, lessFunc)
		assert(t, len(combinations) == 0, "Expected no combinations for 4 elements taken 5 at a time")
	})
}
