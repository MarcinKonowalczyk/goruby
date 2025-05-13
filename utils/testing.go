package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"testing"
)

func getParentInfo(N int) (string, int) {
	parent, _, _, _ := runtime.Caller(2 + N)
	info := runtime.FuncForPC(parent)
	file, line := info.FileLine(parent)
	return file, line
}

func Assert(t *testing.T, predicate bool, args ...any) {
	t.Helper()
	if !predicate {
		file, line := getParentInfo(0)
		var msg string
		if len(args) == 0 {
			msg = "assertion failed"
		}
		switch args[0].(type) {
		case string:
			msg = args[0].(string)
			msg = fmt.Sprintf(msg, args...)
		default:
			msg = fmt.Sprintf("%v", args)
		}
		t.Errorf(msg+" in %s:%d", file, line)
	}
}

func AssertEqual[T comparable](t *testing.T, a T, b T) {
	t.Helper()
	if a != b {
		file, line := getParentInfo(0)
		t.Errorf("expected '%v' (%T) == '%v' (%T) in %s:%d", a, a, b, b, file, line)
	}
}

// func AssertDeepEqual[T any](t *testing.T, a T, b T) {
// 	t.Helper()
// 	if !reflect.DeepEqual(a, b) {
// 		file, line := getParentInfo(0)
// 		t.Errorf("expected reflect.DeepEqual(%v, %v) in %s:%d", a, b, file, line)
// 	}
// }

func AssertNotEqual[T comparable](t *testing.T, a T, b T) {
	t.Helper()
	if a == b {
		file, line := getParentInfo(0)
		t.Errorf("expected %v (%T) != %v (%T) in %s:%d", a, a, b, b, file, line)
	}
}

// Assert that error is nil.
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	AssertError(t, err, nil)
}

func CompareErrors(a error, b error) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Error() != b.Error() {
		return false
	}
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}
	return true
}

// Assert that an error is not nil.
func AssertError(t *testing.T, err error, expected any) {
	t.Helper()
	var msg = ""
	switch expected := expected.(type) {
	case string:
		if err == nil {
			if expected != "" {
				msg = fmt.Sprintf("expected no error, got '%s'", expected)
			}
		} else {
			re := regexp.MustCompile(expected)
			if !re.MatchString(err.Error()) {
				msg = fmt.Sprintf("expected error to match '%s', got '%v' (%T)", expected, err, err)
			}
		}

	case error:
		if expected == nil {
			if err != nil {
				msg = fmt.Sprintf("expected no error, got '%v' (%T)", err, err)
			}
		} else {
			if err == nil {
				msg = fmt.Sprintf("expected error '%v' (%T), got no error (nil)", expected, expected)
			} else {
				if !CompareErrors(err, expected) {
					msg = fmt.Sprintf("expected error '%v' (%T), got '%v' (%T)", expected, expected, err, err)
				}
			}
		}
	case nil:
		if err != nil {
			msg = fmt.Sprintf("expected no error, got '%v' (%T)", err, err)
		}
	case *regexp.Regexp:
		if err == nil {
			msg = fmt.Sprintf("expected error '%v' (%T), got no error (nil)", expected, expected)
		} else {
			re := regexp.MustCompile(expected.String())
			if !re.MatchString(err.Error()) {
				msg = fmt.Sprintf("expected error to match '%s', got '%v' (%T)", expected, err, err)
			}
		}
	default:
		panic("expected type is not an error or string")

	}

	if msg != "" {
		file, line := getParentInfo(0)
		t.Errorf(msg+" in %s:%d", file, line)
	}
}

// Compare two values using a custom comparator function.
func AssertEqualCmp[T any](t *testing.T, a T, b T, comparator func(T, T) bool) {
	t.Helper()
	if !comparator(a, b) {
		file, line := getParentInfo(0)
		t.Errorf("expected %v (%T) == %v (%T) in %s:%d", a, a, b, b, file, line)
	}
}

// Compare two values of any type using a custom comparator function.
// This is a more generic version of AssertEqualCmp, but it is less type-safe.
// The comparator function is responsible for type assertions.
func AssertEqualCmpAny(t *testing.T, a any, b any, comparator func(any, any) bool) {
	defer func() {
		if r := recover(); r != nil {
			// If the comparator panics, we want to catch it and report it as a test failure.
			file, line := getParentInfo(3)
			t.Errorf("Comparator panicked: %v in %s:%d", r, file, line)
		}
	}()
	t.Helper()
	if !comparator(a, b) {
		file, line := getParentInfo(0)
		t.Errorf("expected %v (%T) == %v (%T) in %s:%d", a, a, b, b, file, line)
	}
}

func CompareArrays[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func CompareMaps[T comparable, V comparable](a map[T]V, b map[T]V) bool {
	if len(a) != len(b) {
		return false
	}
	var vb V
	var ok bool

	// NOTE: the range on a map is in random order
	for k, va := range a {
		// Check if key exists in b
		if vb, ok = b[k]; !ok {
			return false
		}
		// Check if value is the same
		if va != vb {
			return false
		}
	}

	// All keys of a exist in b, and a and b have the same length, hence they
	// must have the same keys
	return true
}

// Utility functions for comparing arrays. Equivalent to AssertEqualWithComparator
// where the comparator is CompareArrays.
func AssertEqualArrays[T comparable](t *testing.T, a []T, b []T) {
	t.Helper()
	AssertEqualCmp(t, a, b, CompareArrays)
}

// Utility functions for comparing maps. Equivalent to AssertEqualWithComparator
// where the comparator is CompareMaps.
func AssertEqualMaps[T comparable, V comparable](t *testing.T, a map[T]V, b map[T]V) {
	t.Helper()
	AssertEqualCmp(t, a, b, CompareMaps)
}

// Check if two arrays are equal, regardless of the order of the elements.
func CompareArraysUnordered[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	am := make(map[T]int) // map from element to count
	for _, e := range a {
		am[e]++
	}
	// Iterate over b, decrementing the count of each element in am.
	for _, e := range b {
		if am[e] == 0 {
			return false
		}
		am[e]--
	}
	return true
}

func AssertEqualArraysUnordered[T comparable](t *testing.T, a []T, b []T) {
	t.Helper()
	AssertEqualCmp(t, a, b, CompareArraysUnordered)
}
