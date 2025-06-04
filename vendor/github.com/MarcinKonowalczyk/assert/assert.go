package assert

import (
	"fmt"
	"regexp"
	"runtime"
	"testing"

	"github.com/MarcinKonowalczyk/assert/compare"
)

func getParentInfo(N int) (string, int) {
	parent, _, _, _ := runtime.Caller(1 + N)
	return runtime.FuncForPC(parent).FileLine(parent)
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

const nestedAssertParent = 2

func assert(t testing.TB, N int, predicate bool, args []any) {
	t.Helper()
	if !predicate {
		file, line := getParentInfo(N)
		msg := argsToMessage(func() string { return "assertion failed" }, args)
		t.Errorf(msg+" in %s:%d", file, line)
	}
}

func assert_error(t testing.TB, N int, err error, expected any, args []any) {
	t.Helper()
	var msg_fun func() string
	switch expected := expected.(type) {
	case string:
		if err == nil {
			if expected != "" {
				msg_fun = func() string {
					return fmt.Sprintf("expected no error, got '%s'", expected)
				}
			}
		} else {
			re := regexp.MustCompile(expected)
			if !re.MatchString(err.Error()) {
				msg_fun = func() string {
					return fmt.Sprintf("expected error to match '%s', got '%v' (%T)", expected, err, err)
				}
			}
		}

	case error:
		if expected == nil {
			if err != nil {
				msg_fun = func() string {
					return fmt.Sprintf("expected no error, got '%v' (%T)", err, err)
				}
			}
		} else {
			if err == nil {
				msg_fun = func() string {
					return fmt.Sprintf("expected error '%v' (%T), got no error (nil)", expected, expected)
				}
			} else {
				if !compare.Errors(err, expected) {
					msg_fun = func() string {
						return fmt.Sprintf("expected error '%v' (%T), got '%v' (%T)", expected, expected, err, err)
					}
				}
			}
		}
	case nil:
		if err != nil {
			// msg = fmt.Sprintf("expected no error, got '%v' (%T)", err, err)
			msg_fun = func() string {
				return fmt.Sprintf("expected no error, got '%v' (%T)", err, err)
			}
		}
	case *regexp.Regexp:
		if err == nil {
			msg_fun = func() string {
				return fmt.Sprintf("expected error '%v' (%T), got no error (nil)", expected, expected)
			}
		} else {
			re := regexp.MustCompile(expected.String())
			if !re.MatchString(err.Error()) {
				msg_fun = func() string {
					return fmt.Sprintf("expected error to match '%s', got '%v' (%T)", expected, err, err)
				}
			}
		}
	default:
		panic("expected type is not an error or string")

	}

	if msg_fun != nil {
		msg := argsToMessage(msg_fun, args)
		file, line := getParentInfo(N)
		t.Errorf(msg+" in %s:%d", file, line)
	}
}

func equal_cmp[T any](t testing.TB, N int, a T, b T, comparator func(T, T) bool) {
	t.Helper()
	assert(t, N+1, comparator(a, b), []any{"expected '%v' (%T) == '%v' (%T)", a, a, b, b})
}

func equal_cmp_any(t testing.TB, N int, a any, b any, comparator func(any, any) bool) {
	defer func() {
		if r := recover(); r != nil {
			// If the comparator panics, we want to catch it and report it as a test failure.
			file, line := getParentInfo(4)
			t.Errorf("Comparator panicked: %v in %s:%d", r, file, line)
		}
	}()
	t.Helper()
	assert(t, N+1, comparator(a, b), []any{"expected '%v' (%T) == '%v' (%T)", a, a, b, b})
}

// Check that the type of obj is T.
func assert_type[T any](t testing.TB, N int, obj any, args ...any) T {
	t.Helper()
	if obj_T, ok := obj.(T); ok {
		return obj_T
	} else {
		file, line := getParentInfo(N)
		msg := argsToMessage(func() string {
			return fmt.Sprintf("expected type %T, got %T", (*T)(nil), obj)
		}, args)
		t.Errorf(msg+" in %s:%d", file, line)
	}
	return *new(T)
}
