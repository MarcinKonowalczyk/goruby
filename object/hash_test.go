package object

import (
	"reflect"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/object/hash"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
)

func TestHashSet(t *testing.T) {
	t.Run("Set on initialized hash", func(t *testing.T) {
		hash := &Hash{Map: make(map[hash.Key]hashPair)}

		key := NewString("foo")
		value := NewInteger(42)

		result := hash.Set(key, value)
		assert.Equal(t, len(hash.Map), 1)

		var values []hashPair
		for _, v := range hash.Map {
			values = append(values, v)
		}

		assert.EqualCmpAny(t, values[0].Key, key, CompareRubyObjectsForTests)
		assert.EqualCmpAny(t, values[0].Value, value, CompareRubyObjectsForTests)
		assert.EqualCmpAny(t, result, value, CompareRubyObjectsForTests)
	})
	t.Run("Set on uninitialized hash", func(t *testing.T) {
		var hash Hash

		key := NewString("foo")
		value := NewInteger(42)

		result := hash.Set(key, value)
		assert.Equal(t, len(hash.Map), 1)

		var values []hashPair
		for _, v := range hash.Map {
			values = append(values, v)
		}

		assert.EqualCmpAny(t, values[0].Key, key, CompareRubyObjectsForTests)
		assert.EqualCmpAny(t, values[0].Value, value, CompareRubyObjectsForTests)
		assert.EqualCmpAny(t, result, value, CompareRubyObjectsForTests)
	})
}

func TestHashGet(t *testing.T) {
	t.Run("value found", func(t *testing.T) {
		key := NewString("foo")
		value := NewInteger(42)

		hash := &Hash{Map: map[hash.Key]hashPair{
			key.HashKey(): {Key: key, Value: value},
		}}

		result, ok := hash.Get(key)

		assert.That(t, ok, "Expected returned bool to be true, got false")
		assert.EqualCmpAny(t, result, value, CompareRubyObjectsForTests)
	})
	t.Run("value not found", func(t *testing.T) {
		key := NewString("foo")

		hash := &Hash{Map: map[hash.Key]hashPair{}}

		result, ok := hash.Get(key)

		assert.That(t, !ok, "Expected returned bool to be false, got true")
		assert.Equal(t, result, nil)
	})
	t.Run("on uninitalized hash", func(t *testing.T) {
		key := NewString("foo")

		var hash Hash

		result, ok := hash.Get(key)

		assert.That(t, !ok, "Expected returned bool to be false, got true")
		assert.Equal(t, result, nil)
	})
}

func TestHashMap(t *testing.T) {
	t.Run("on initialized hash", func(t *testing.T) {
		key := NewString("foo")
		value := NewInteger(42)

		hash := &Hash{Map: map[hash.Key]hashPair{
			key.HashKey(): {Key: key, Value: value},
		}}

		var result map[ruby.Object]ruby.Object = hash.ObjectMap()

		expected := map[string]ruby.Object{
			"foo": value,
		}
		actual := make(map[string]ruby.Object)
		for k, v := range result {
			actual[k.Inspect()] = v
		}

		assert.That(t, reflect.DeepEqual(expected, actual), "Expected hash to equal\n%s\n\tgot\n%s\n", expected, actual)
	})
	t.Run("on uninitialized hash", func(t *testing.T) {
		var hash Hash

		var result map[ruby.Object]ruby.Object = hash.ObjectMap()

		expected := map[string]ruby.Object{}
		actual := make(map[string]ruby.Object)
		for k, v := range result {
			actual[k.Inspect()] = v
		}

		assert.That(t, reflect.DeepEqual(expected, actual), "Expected hash to equal\n%s\n\tgot\n%s\n", expected, actual)
	})
}
