package assert

import (
	"fmt"
	"regexp"
	"runtime"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/testutils/assert/compare"
)

func getParentInfo(N int) (string, int) {
	parent, _, _, _ := runtime.Caller(1 + N)
	info := runtime.FuncForPC(parent)
	file, line := info.FileLine(parent)
	return file, line
}

// convert 'args ...any' to the assertion message
// internal utility so we don't use variadics to make the calls a bit more consistent
func argsToMessage(default_func func() string, args []any) string {
	var msg string
	if len(args) == 0 {
		msg = default_func()
	} else {
		switch args[0].(type) {
		case string:
			msg = args[0].(string)
			msg = fmt.Sprintf(msg, args[1:]...)
		default:
			msg = fmt.Sprintf("%v", args)
		}
	}
	return msg
}

const thisFunctionsParent = 1
const nestedAssertParent = 2

func assert(t *testing.T, N int, predicate bool, args []any) {
	t.Helper()
	if !predicate {
		file, line := getParentInfo(N)
		msg := argsToMessage(func() string { return "assertion failed" }, args)
		t.Errorf(msg+" in %s:%d", file, line)
	}
}

func That(t *testing.T, predicate bool, args ...any) {
	t.Helper()
	assert(t, nestedAssertParent, predicate, args)
}

func Equal[T comparable](t *testing.T, a T, b T) {
	t.Helper()
	assert(t, nestedAssertParent, a == b, []any{"expected '%v' (%T) == '%v' (%T)", a, a, b, b})
}

// func AssertDeepEqual[T any](t *testing.T, a T, b T) {
// 	t.Helper()
// 	if !reflect.DeepEqual(a, b) {
// 		file, line := getParentInfo(thisFunctionsParent)
// 		t.Errorf("expected reflect.DeepEqual(%v, %v) in %s:%d", a, b, file, line)
// 	}
// }

func NotEqual[T comparable](t *testing.T, a T, b T) {
	t.Helper()
	assert(t, nestedAssertParent, a != b, []any{"expected '%v' (%T) != '%v' (%T)", a, a, b, b})
}

// Assert that error is nil.
func NoError(t *testing.T, err error) {
	t.Helper()
	Error(t, err, nil)
}

// Assert that an error is not nil.
func Error(t *testing.T, err error, expected any) {
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
				if !compare.Errors(err, expected) {
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
		file, line := getParentInfo(thisFunctionsParent)
		t.Errorf(msg+" in %s:%d", file, line)
	}
}

// Compare two values using a custom comparator function.
func EqualCmp[T any](t *testing.T, a T, b T, comparator func(T, T) bool) {
	t.Helper()
	assert(t, nestedAssertParent, comparator(a, b), []any{"expected '%v' (%T) == '%v' (%T)", a, a, b, b})
}

// Compare two values of any type using a custom comparator function.
// This is a more generic version of AssertEqualCmp, but it is less type-safe.
// The comparator function is responsible for type assertions.
func EqualCmpAny(t *testing.T, a any, b any, comparator func(any, any) bool) {
	defer func() {
		if r := recover(); r != nil {
			// If the comparator panics, we want to catch it and report it as a test failure.
			file, line := getParentInfo(4)
			t.Errorf("Comparator panicked: %v in %s:%d", r, file, line)
		}
	}()
	t.Helper()
	assert(t, nestedAssertParent, comparator(a, b), []any{"expected '%v' (%T) == '%v' (%T)", a, a, b, b})
}

// Utility functions for comparing arrays. Equivalent to AssertEqualWithComparator
// where the comparator is CompareArrays.
func EqualArrays[T comparable](t *testing.T, a []T, b []T) {
	t.Helper()
	EqualCmp(t, a, b, compare.Arrays)
}

// Utility functions for comparing maps. Equivalent to AssertEqualWithComparator
// where the comparator is CompareMaps.
func EqualMaps[T comparable, V comparable](t *testing.T, a map[T]V, b map[T]V) {
	t.Helper()
	EqualCmp(t, a, b, compare.Maps)
}

func EqualArraysUnordered[T comparable](t *testing.T, a []T, b []T) {
	t.Helper()
	EqualCmp(t, a, b, compare.ArraysUnordered)
}

// Check that the type of obj is T.
func assertType[T any](t *testing.T, N int, obj any, args ...any) T {
	t.Helper()
	if obj_T, ok := obj.(T); ok {
		return obj_T
	} else {
		file, line := getParentInfo(N)
		// t.Errorf("expected type %T, got %T in %s:%d", (*T)(nil), obj, file, line)
		msg := argsToMessage(func() string {
			return fmt.Sprintf("expected type %T, got %T", (*T)(nil), obj)
		}, args)
		t.Errorf(msg+" in %s:%d", file, line)
	}
	return *new(T)
}

func Type[T any](t *testing.T, obj any, args ...any) T {
	t.Helper()
	return assertType[T](t, nestedAssertParent, obj, args...)
}

// func EqualLineByLine(t *testing.T, a string, b string) {
// 	t.Helper()
// 	// EqualCmp(t, a, b, compare.LineByLine)
// 	// assert(t, nestedAssertParent, comparator(a, b), []any{"expected '%v' (%T) == '%v' (%T)", a, a, b, b})
// 	a_lines := strings.Split(a, "\n")
// 	b_lines := strings.Split(b, "\n")
// 	assert(t, nestedAssertParent, len(a_lines) == len(b_lines), []any{"expected '%d' lines, got '%d'", len(a_lines), len(b_lines)})
// 	if len(a_lines) != len(b_lines) {
// 		return // no point in checking the lines if the number of lines is different
// 	}
// 	for i := range a_lines {
// 		assert(t, nestedAssertParent, a_lines[i] == b_lines[i], []any{"expected line %d to be '%s', got '%s'", i + 1, a_lines[i], b_lines[i]})
// 	}
// }
