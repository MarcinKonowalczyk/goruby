package function_lift

import (
	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/parser"
)

type context_spec string

type replacement_spec struct {
	statement ast.Statement
	name      string
}

var context_call_expr_replacements map[context_spec]replacement_spec = map[context_spec]replacement_spec{
	"find_all": {
		name: "__builtin_find_all",
		statement: parser.MustParseStatement(`
__builtin_find_all = -> (range, fun) {
    out = []
    i = 0
    loop {
        break if i >= range.size
        if fun[i]
            out.push(i)
        end
        i += 1
    }
    out
}`),
	},
}
