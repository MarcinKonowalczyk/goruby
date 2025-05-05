package object

import (
	"reflect"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/utils"
)

func TestHashSet(t *testing.T) {
	t.Run("Set on initialized hash", func(t *testing.T) {
		hash := &Hash{Map: make(map[HashKey]hashPair)}

		key := NewString("foo")
		value := NewInteger(42)

		result := hash.Set(key, value)
		utils.AssertEqual(t, len(hash.Map), 1)

		var values []hashPair
		for _, v := range hash.Map {
			values = append(values, v)
		}

		utils.AssertEqualCmpAny(t, values[0].Key, key, CompareRubyObjectsForTests)
		utils.AssertEqualCmpAny(t, values[0].Value, value, CompareRubyObjectsForTests)
		utils.AssertEqualCmpAny(t, result, value, CompareRubyObjectsForTests)
	})
	t.Run("Set on uninitialized hash", func(t *testing.T) {
		var hash Hash

		key := NewString("foo")
		value := NewInteger(42)

		result := hash.Set(key, value)
		utils.AssertEqual(t, len(hash.Map), 1)

		var values []hashPair
		for _, v := range hash.Map {
			values = append(values, v)
		}

		utils.AssertEqualCmpAny(t, values[0].Key, key, CompareRubyObjectsForTests)
		utils.AssertEqualCmpAny(t, values[0].Value, value, CompareRubyObjectsForTests)
		utils.AssertEqualCmpAny(t, result, value, CompareRubyObjectsForTests)
	})
}

func TestHashGet(t *testing.T) {
	t.Run("value found", func(t *testing.T) {
		key := NewString("foo")
		value := NewInteger(42)

		hash := &Hash{Map: map[HashKey]hashPair{
			key.HashKey(): hashPair{Key: key, Value: value},
		}}

		result, ok := hash.Get(key)

		utils.Assert(t, ok, "Expected returned bool to be true, got false")
		utils.AssertEqualCmpAny(t, result, value, CompareRubyObjectsForTests)
	})
	t.Run("value not found", func(t *testing.T) {
		key := NewString("foo")

		hash := &Hash{Map: map[HashKey]hashPair{}}

		result, ok := hash.Get(key)

		utils.Assert(t, !ok, "Expected returned bool to be false, got true")
		utils.AssertEqual(t, result, nil)
	})
	t.Run("on uninitalized hash", func(t *testing.T) {
		key := NewString("foo")

		var hash Hash

		result, ok := hash.Get(key)

		utils.Assert(t, !ok, "Expected returned bool to be false, got true")
		utils.AssertEqual(t, result, nil)
	})
}

func TestHashMap(t *testing.T) {
	t.Run("on initialized hash", func(t *testing.T) {
		key := NewString("foo")
		value := NewInteger(42)

		hash := &Hash{Map: map[HashKey]hashPair{
			key.HashKey(): hashPair{Key: key, Value: value},
		}}

		var result map[RubyObject]RubyObject = hash.ObjectMap()

		expected := map[string]RubyObject{
			"foo": value,
		}
		actual := make(map[string]RubyObject)
		for k, v := range result {
			actual[k.Inspect()] = v
		}

		utils.Assert(t, reflect.DeepEqual(expected, actual), "Expected hash to equal\n%s\n\tgot\n%s\n", expected, actual)
	})
	t.Run("on uninitialized hash", func(t *testing.T) {
		var hash Hash

		var result map[RubyObject]RubyObject = hash.ObjectMap()

		expected := map[string]RubyObject{}
		actual := make(map[string]RubyObject)
		for k, v := range result {
			actual[k.Inspect()] = v
		}

		utils.Assert(t, reflect.DeepEqual(expected, actual), "Expected hash to equal\n%s\n\tgot\n%s\n", expected, actual)
	})
}
