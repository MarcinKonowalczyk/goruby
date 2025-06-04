package assert

import (
	"strings"
	"testing"

	"github.com/MarcinKonowalczyk/assert/compare"
)

func That(t testing.TB, predicate bool, args ...any) {
	t.Helper()
	assert(t, nestedAssertParent, predicate, args)
}

func Equal[T comparable](t testing.TB, a T, b T) {
	t.Helper()
	assert(t, nestedAssertParent, a == b, []any{"expected '%v' (%T) == '%v' (%T)", a, a, b, b})
}

// func AssertDeepEqual[T any](t testing.TB, a T, b T) {
// 	t.Helper()
// 	if !reflect.DeepEqual(a, b) {
// 		file, line := getParentInfo(thisFunctionsParent)
// 		t.Errorf("expected reflect.DeepEqual(%v, %v) in %s:%d", a, b, file, line)
// 	}
// }

func NotEqual[T comparable](t testing.TB, a T, b T) {
	t.Helper()
	assert(t, nestedAssertParent, a != b, []any{"expected '%v' (%T) != '%v' (%T)", a, a, b, b})
}

// Assert that error is nil.
func NoError(t testing.TB, err error, args ...any) {
	t.Helper()
	assert_error(t, nestedAssertParent, err, nil, args)
}

// Assert that an error is not nil.
func Error(t testing.TB, err error, expected any, args ...any) {
	t.Helper()
	assert_error(t, nestedAssertParent, err, expected, args)
}

// Compare two values using a custom comparator function.
func EqualCmp[T any](t testing.TB, a T, b T, comparator func(T, T) bool) {
	t.Helper()
	equal_cmp(t, nestedAssertParent, a, b, comparator)
}

// Compare two values of any type using a custom comparator function.
// This is a more generic version of AssertEqualCmp, but it is less type-safe.
// The comparator function is responsible for type assertions.
func EqualCmpAny(t testing.TB, a any, b any, comparator func(any, any) bool) {
	t.Helper()
	equal_cmp_any(t, nestedAssertParent, a, b, comparator)
}

// Utility functions for comparing arrays. Equivalent to AssertEqualWithComparator
// where the comparator is CompareArrays.
func EqualArrays[T comparable](t testing.TB, a []T, b []T) {
	t.Helper()
	equal_cmp(t, nestedAssertParent, a, b, compare.Arrays)
}

// Utility functions for comparing maps. Equivalent to AssertEqualWithComparator
// where the comparator is CompareMaps.
func EqualMaps[T comparable, V comparable](t testing.TB, a map[T]V, b map[T]V) {
	t.Helper()
	equal_cmp(t, nestedAssertParent, a, b, compare.Maps)
}

func EqualArraysUnordered[T comparable](t testing.TB, a []T, b []T) {
	t.Helper()
	equal_cmp(t, nestedAssertParent, a, b, compare.ArraysUnordered)
}

func Type[T any](t testing.TB, obj any, args ...any) T {
	t.Helper()
	return assert_type[T](t, nestedAssertParent, obj, args...)
}

func EqualLineByLine(t testing.TB, a string, b string) {
	t.Helper()
	// EqualCmp(t, a, b, compare.LineByLine)
	// assert(t, nestedAssertParent, comparator(a, b), []any{"expected '%v' (%T) == '%v' (%T)", a, a, b, b})
	a_lines := strings.Split(a, "\n")
	b_lines := strings.Split(b, "\n")
	assert(t, nestedAssertParent, len(a_lines) == len(b_lines), []any{"expected '%d' lines, got '%d'", len(a_lines), len(b_lines)})
	if len(a_lines) != len(b_lines) {
		return // no point in checking the lines if the number of lines is different
	}
	for i := range a_lines {
		assert(t, nestedAssertParent, a_lines[i] == b_lines[i], []any{"expected line %d to be '%s', got '%s'", i + 1, a_lines[i], b_lines[i]})
	}
}

func ContainsString(t testing.TB, haystack string, needle string) {
	t.Helper()
	assert(t, nestedAssertParent, strings.Contains(haystack, needle), []any{"expected needle string '%s' to be in a haystack string '%s'", needle, haystack})
}
