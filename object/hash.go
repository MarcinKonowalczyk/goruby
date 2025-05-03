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
		return &Hash{Map: make(map[HashKey]hashPair)}, nil
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

type HashKey uint64

func (h HashKey) bytes() []byte {
	bytes := [4]byte{}
	bytes[0] = byte(h >> 24)
	bytes[1] = byte(h >> 16)
	bytes[2] = byte(h >> 8)
	bytes[3] = byte(h)
	return bytes[:]
}

// func hash(obj RubyObject) HashKey {
// 	if hashable, ok := obj.(hashable); ok {
// 		return hashable.HashKey()
// 	}
// 	pointer := fmt.Sprintf("%p", obj)
// 	h := fnv.New64a()
// 	h.Write([]byte(pointer))
// 	return HashKey(h.Sum64())
// }

type hashPair struct {
	Key   RubyObject
	Value RubyObject
}

type Hash struct {
	Map map[HashKey]hashPair
}

func (h *Hash) init() {
	if h.Map == nil {
		h.Map = make(map[HashKey]hashPair)
	}
}

func (h *Hash) Set(key, value RubyObject) RubyObject {
	h.init()
	h.Map[key.HashKey()] = hashPair{Key: key, Value: value}
	return value
}

func (h *Hash) Get(key RubyObject) (RubyObject, bool) {
	v, ok := h.Map[key.HashKey()]
	if !ok {
		return nil, false
	}
	return v.Value, true
}

func (h *Hash) Inspect() string {
	elems := []string{}
	for _, v := range h.Map {
		elems = append(elems, fmt.Sprintf("%q => %q", v.Key.Inspect(), v.Value.Inspect()))
	}
	return "{" + strings.Join(elems, ", ") + "}"
}

func (h *Hash) Class() RubyClass { return hashClass }

func (h *Hash) HashKey() HashKey {
	hash := fnv.New64a()
	for k := range h.Map {
		hash.Write(k.bytes())
	}
	return HashKey(hash.Sum64())
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
