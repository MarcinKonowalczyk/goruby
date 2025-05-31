package object

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/hash"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/trace"
)

var hashClass ruby.ClassObject = newClass(
	"Hash",
	hashMethods,
	nil,
	func(ruby.ClassObject, ...ruby.Object) (ruby.Object, error) {
		return &Hash{Map: make(map[hash.Key]hashPair)}, nil
	},
)

func init() {
	CLASSES.Set("Hash", hashClass)
}

func (h *Hash) ObjectMap() map[ruby.Object]ruby.Object {
	hashmap := make(map[ruby.Object]ruby.Object)
	for _, v := range h.Map {
		hashmap[v.Key] = v.Value
	}
	return hashmap
}

type hashPair struct {
	Key   ruby.Object
	Value ruby.Object
}

type Hash struct {
	Map map[hash.Key]hashPair
}

func (h *Hash) init() {
	if h.Map == nil {
		h.Map = make(map[hash.Key]hashPair)
	}
}

func (h *Hash) Set(key, value ruby.Object) ruby.Object {
	h.init()
	h.Map[key.HashKey()] = hashPair{Key: key, Value: value}
	return value
}

func (h *Hash) Get(key ruby.Object) (ruby.Object, bool) {
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

func (h *Hash) Class() ruby.Class { return hashClass }

func (h *Hash) HashKey() hash.Key {
	hsh := fnv.New64a()
	for k := range h.Map {
		hsh.Write(k.Bytes())
	}
	return hash.Key(hsh.Sum64())
}

var hashMethods = map[string]ruby.Method{
	"has_key?": ruby.NewMethod(hashHasKey),
}

func hashHasKey(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx, trace.HereCtx(ctx))()
	hash, _ := ctx.Receiver().(*Hash)
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError(1, len(args))
	}
	key := args[0]
	if _, ok := hash.Get(key); ok {
		return TRUE, nil
	}
	return FALSE, nil
}
