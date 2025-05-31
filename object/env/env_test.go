package env_test

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/object"
	"github.com/MarcinKonowalczyk/goruby/object/env"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
)

// TODO: we don't actually need object/ruby.Object here. we can test with whatever.

func TestEnvironmentSet(t *testing.T) {
	e := env.NewEnvironment[ruby.Object]()
	e.Set("foo", object.NIL)

	val, ok := e.ShallowGet("foo")
	assert.That(t, ok, "Expected store to contain 'foo'")
	assert.EqualCmpAny(t, val, object.NIL, object.CompareRubyObjectsForTests)
}

func TestEnvironmentSetGlobal(t *testing.T) {
	t.Run("toplevel env", func(t *testing.T) {
		e := env.NewEnvironment[ruby.Object]()
		e.SetGlobal("$foo", object.NIL)

		val, ok := e.ShallowGet("$foo")
		assert.That(t, ok, "Expected store to contain '$foo'")
		assert.EqualCmpAny(t, val, object.NIL, object.CompareRubyObjectsForTests)
	})
	t.Run("inner env one level", func(t *testing.T) {
		outer := env.NewEnvironment[ruby.Object]()
		e := env.NewEnclosedEnvironment(outer)
		e.SetGlobal("$foo", object.NIL)

		_, ok := e.ShallowGet("$foo")
		assert.That(t, !ok, "Expected env store to not contain '$foo'")

		val, ok := outer.ShallowGet("$foo")

		assert.That(t, ok, "Expected outer store to contain '$foo'")
		assert.EqualCmpAny(t, val, object.NIL, object.CompareRubyObjectsForTests)
	})
	t.Run("inner env two level", func(t *testing.T) {
		root := env.NewEnvironment[ruby.Object]()
		outer := env.NewEnclosedEnvironment(root)
		e := env.NewEnclosedEnvironment(outer)
		e.SetGlobal("$foo", object.NIL)

		_, ok := e.ShallowGet("$foo")
		assert.That(t, !ok, "Expected env store to not contain '$foo'")

		_, ok = outer.ShallowGet("$foo")
		assert.That(t, !ok, "Expected outer store to not contain '$foo'")

		val, ok := root.ShallowGet("$foo")
		assert.That(t, ok, "Expected root store to contain '$foo'")
		assert.EqualCmpAny(t, val, object.NIL, object.CompareRubyObjectsForTests)
	})
}

func TestEnvironmentGet(t *testing.T) {
	t.Run("toplevel env", func(t *testing.T) {
		// env := &environment{store: map[string]ruby.Object{"foo": object.TRUE}}
		env := env.NewEnvironment[ruby.Object]()
		env.Set("foo", object.TRUE)

		val, ok := env.Get("foo")
		assert.That(t, ok, "Expected store to contain 'foo'")
		assert.EqualCmpAny(t, val, object.TRUE, object.CompareRubyObjectsForTests)
	})
	t.Run("inner env one level", func(t *testing.T) {
		outer := env.NewEnvironment[ruby.Object]()
		outer.Set("foo", object.TRUE)
		env := env.NewEnclosedEnvironment(outer)

		val, ok := env.Get("foo")
		assert.That(t, ok, "Expected store to contain 'foo'")
		assert.EqualCmpAny(t, val, object.TRUE, object.CompareRubyObjectsForTests)
	})
	t.Run("inner env two level", func(t *testing.T) {
		root := env.NewEnvironment[ruby.Object]()
		root.Set("foo", object.TRUE)
		outer := env.NewEnclosedEnvironment(root)
		env := env.NewEnclosedEnvironment(outer)

		val, ok := env.Get("foo")
		assert.That(t, ok, "Expected store to contain 'foo'")
		assert.EqualCmpAny(t, val, object.TRUE, object.CompareRubyObjectsForTests)
	})
	t.Run("inner env two level overridden key", func(t *testing.T) {
		root := env.NewEnvironment[ruby.Object]()
		root.Set("foo", object.FALSE)
		outer := env.NewEnclosedEnvironment(root)
		outer.Set("foo", object.TRUE)
		env := env.NewEnclosedEnvironment(outer)

		val, ok := env.Get("foo")
		assert.That(t, ok, "Expected store to contain 'foo'")
		assert.EqualCmpAny(t, val, object.TRUE, object.CompareRubyObjectsForTests)
	})
}
