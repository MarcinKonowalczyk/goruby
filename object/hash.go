package object

import (
	"fmt"
	"hash/fnv"
	"strings"
)

var hashClass RubyClassObject = newClass(
	"Hash",
	hashMethods,
	nil,
	func(RubyClassObject, ...RubyObject) (RubyObject, error) {
		return &Hash{Map: make(map[hashKey]hashPair)}, nil
	},
)

func init() {
	CLASSES.Set("Hash", hashClass)
}

func (h *Hash) ObjectMap() map[RubyObject]RubyObject {
	hashmap := make(map[RubyObject]RubyObject)
	for _, v := range h.Map {
		hashmap[v.Key] = v.Value
	}
	return hashmap
}

type hashKey struct {
	Type  Type
	Value uint64
}

func (h hashKey) bytes() []byte {
	return append([]byte(h.Type), byte(h.Value))
}

func hash(obj RubyObject) hashKey {
	if hashable, ok := obj.(hashable); ok {
		return hashable.hashKey()
	}
	pointer := fmt.Sprintf("%p", obj)
	h := fnv.New64a()
	h.Write([]byte(pointer))
	return hashKey{Type: obj.Type(), Value: h.Sum64()}
}

type hashPair struct {
	Key   RubyObject
	Value RubyObject
}

type Hash struct {
	Map map[hashKey]hashPair
}

func (h *Hash) init() {
	if h.Map == nil {
		h.Map = make(map[hashKey]hashPair)
	}
}

func (h *Hash) Set(key, value RubyObject) RubyObject {
	h.init()
	h.Map[hash(key)] = hashPair{Key: key, Value: value}
	return value
}

func (h *Hash) Get(key RubyObject) (RubyObject, bool) {
	v, ok := h.Map[hash(key)]
	if !ok {
		return nil, false
	}
	return v.Value, true
}

func (h *Hash) Type() Type { return HASH_OBJ }

func (h *Hash) Inspect() string {
	elems := []string{}
	for _, v := range h.Map {
		elems = append(elems, fmt.Sprintf("%q => %q", v.Key.Inspect(), v.Value.Inspect()))
	}
	return "{" + strings.Join(elems, ", ") + "}"
}

func (h *Hash) Class() RubyClass { return hashClass }

func (h *Hash) hashKey() hashKey {
	hash := fnv.New64a()
	for k := range h.Map {
		hash.Write(k.bytes())
	}
	return hashKey{Type: h.Type(), Value: hash.Sum64()}
}

var hashMethods = map[string]RubyMethod{
	"has_key?": newMethod(hashHasKey),
}

func hashHasKey(context CallContext, args ...RubyObject) (RubyObject, error) {
	hash, _ := context.Receiver().(*Hash)
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError(1, len(args))
	}
	key := args[0]
	if _, ok := hash.Get(key); ok {
		return TRUE, nil
	}
	return FALSE, nil
}
