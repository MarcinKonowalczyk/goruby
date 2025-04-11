package lexer

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/token"
)

func TestLexerNextToken(t *testing.T) {
	input := `five = 5
while x < y do
	x += x
end
seven,
# just comment
fifty = 5_0
ten = 10
?-
?\n
? foo : bar

def add(x, y)
	x + y
end
|
||

result = add(five, ten)
!-/ *%5;
+= -= *= /= %=
5 < 10 > 5
return
if 5 < 10 then
	true
else
	false
end

begin
rescue
end

10 == 10
10 != 9
10 <= 9
10 >= 9
10 <=> 9
10 << 9
""
"foobar"
'foobar'
"foo bar"
'foo bar'
:sym
:"sym"
:'sym'
.
&foo
&
&&
:dotAfter.

def nil?
end

def run!
end
[1, 2]
nil
self
Ten = 10
module Abc
end
class Abc
end
add { |x| x }
add do |x|
end
yield
while
A::B
=>
->
__FILE__
@
$foo,
$foo;
$Foo
$dotAfter.
$@
$a
..
...
1..2
1...2
/\//
/a/i
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.NEWLINE, "\n"},
		{token.WHILE, "while"},
		{token.IDENT, "x"},
		{token.LT, "<"},
		{token.IDENT, "y"},
		{token.DO, "do"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "x"},
		{token.ADDASSIGN, "+="},
		{token.IDENT, "x"},
		{token.NEWLINE, "\n"},
		{token.END, "end"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "seven"},
		{token.COMMA, ","},
		{token.NEWLINE, "\n"},
		{token.HASH, "#"},
		{token.STRING, " just comment"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "fifty"},
		{token.ASSIGN, "="},
		{token.INT, "5_0"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.NEWLINE, "\n"},
		{token.STRING, "-"},
		{token.NEWLINE, "\n"},
		{token.STRING, "\\n"},
		{token.NEWLINE, "\n"},
		{token.QMARK, "?"},
		{token.IDENT, "foo"},
		{token.COLON, ":"},
		{token.IDENT, "bar"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.DEF, "def"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.NEWLINE, "\n"},
		{token.END, "end"},
		{token.NEWLINE, "\n"},
		{token.PIPE, "|"},
		{token.NEWLINE, "\n"},
		{token.LOGICALOR, "||"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.NEWLINE, "\n"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.MODULO, "%"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.ADDASSIGN, "+="},
		{token.SUBASSIGN, "-="},
		{token.MULASSIGN, "*="},
		{token.DIVASSIGN, "/="},
		{token.MODASSIGN, "%="},
		{token.NEWLINE, "\n"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.NEWLINE, "\n"},
		{token.RETURN, "return"},
		{token.NEWLINE, "\n"},
		{token.IF, "if"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.THEN, "then"},
		{token.NEWLINE, "\n"},
		{token.TRUE, "true"},
		{token.NEWLINE, "\n"},
		{token.ELSE, "else"},
		{token.NEWLINE, "\n"},
		{token.FALSE, "false"},
		{token.NEWLINE, "\n"},
		{token.END, "end"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.BEGIN, "begin"},
		{token.NEWLINE, "\n"},
		{token.RESCUE, "rescue"},
		{token.NEWLINE, "\n"},
		{token.END, "end"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.NEWLINE, "\n"},
		{token.INT, "10"},
		{token.NOTEQ, "!="},
		{token.INT, "9"},
		{token.NEWLINE, "\n"},
		{token.INT, "10"},
		{token.LTE, "<="},
		{token.INT, "9"},
		{token.NEWLINE, "\n"},
		{token.INT, "10"},
		{token.GTE, ">="},
		{token.INT, "9"},
		{token.NEWLINE, "\n"},
		{token.INT, "10"},
		{token.SPACESHIP, "<=>"},
		{token.INT, "9"},
		{token.NEWLINE, "\n"},
		{token.INT, "10"},
		{token.LSHIFT, "<<"},
		{token.INT, "9"},
		{token.NEWLINE, "\n"},
		{token.STRING, ""},
		{token.NEWLINE, "\n"},
		{token.STRING, "foobar"},
		{token.NEWLINE, "\n"},
		{token.STRING, "foobar"},
		{token.NEWLINE, "\n"},
		{token.STRING, "foo bar"},
		{token.NEWLINE, "\n"},
		{token.STRING, "foo bar"},
		{token.NEWLINE, "\n"},
		{token.SYMBEG, ":"},
		{token.IDENT, "sym"},
		{token.NEWLINE, "\n"},
		{token.SYMBEG, ":"},
		{token.STRING, "sym"},
		{token.NEWLINE, "\n"},
		{token.SYMBEG, ":"},
		{token.STRING, "sym"},
		{token.NEWLINE, "\n"},
		{token.DOT, "."},
		{token.NEWLINE, "\n"},
		{token.CAPTURE, "&"},
		{token.IDENT, "foo"},
		{token.NEWLINE, "\n"},
		{token.AND, "&"},
		{token.NEWLINE, "\n"},
		{token.LOGICALAND, "&&"},
		{token.NEWLINE, "\n"},
		{token.SYMBEG, ":"},
		{token.IDENT, "dotAfter"},
		{token.DOT, "."},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.DEF, "def"},
		{token.IDENT, "nil?"},
		{token.NEWLINE, "\n"},
		{token.END, "end"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.DEF, "def"},
		{token.IDENT, "run!"},
		{token.NEWLINE, "\n"},
		{token.END, "end"},
		{token.NEWLINE, "\n"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.NEWLINE, "\n"},
		{token.NIL, "nil"},
		{token.NEWLINE, "\n"},
		{token.SELF, "self"},
		{token.NEWLINE, "\n"},
		{token.CONST, "Ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.NEWLINE, "\n"},
		{token.MODULE, "module"},
		{token.CONST, "Abc"},
		{token.NEWLINE, "\n"},
		{token.END, "end"},
		{token.NEWLINE, "\n"},
		{token.CLASS, "class"},
		{token.CONST, "Abc"},
		{token.NEWLINE, "\n"},
		{token.END, "end"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "add"},
		{token.LBRACE, "{"},
		{token.PIPE, "|"},
		{token.IDENT, "x"},
		{token.PIPE, "|"},
		{token.IDENT, "x"},
		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "add"},
		{token.DO, "do"},
		{token.PIPE, "|"},
		{token.IDENT, "x"},
		{token.PIPE, "|"},
		{token.NEWLINE, "\n"},
		{token.END, "end"},
		{token.NEWLINE, "\n"},
		{token.YIELD, "yield"},
		{token.NEWLINE, "\n"},
		{token.WHILE, "while"},
		{token.NEWLINE, "\n"},
		{token.CONST, "A"},
		{token.SCOPE, "::"},
		{token.CONST, "B"},
		{token.NEWLINE, "\n"},
		{token.HASHROCKET, "=>"},
		{token.NEWLINE, "\n"},
		{token.LAMBDAROCKET, "->"},
		{token.NEWLINE, "\n"},
		{token.KEYWORD__FILE__, "__FILE__"},
		{token.NEWLINE, "\n"},
		{token.AT, "@"},
		{token.NEWLINE, "\n"},
		{token.GLOBAL, "$foo"},
		{token.COMMA, ","},
		{token.NEWLINE, "\n"},
		{token.GLOBAL, "$foo"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.GLOBAL, "$Foo"},
		{token.NEWLINE, "\n"},
		{token.GLOBAL, "$dotAfter"},
		{token.DOT, "."},
		{token.NEWLINE, "\n"},
		{token.GLOBAL, "$@"},
		{token.NEWLINE, "\n"},
		{token.GLOBAL, "$a"},
		{token.NEWLINE, "\n"},
		{token.DDOT, ".."},
		{token.NEWLINE, "\n"},
		{token.DDDOT, "..."},
		{token.NEWLINE, "\n"},
		{token.INT, "1"},
		{token.DDOT, ".."},
		{token.INT, "2"},
		{token.NEWLINE, "\n"},
		{token.INT, "1"},
		{token.DDDOT, "..."},
		{token.INT, "2"},
		{token.NEWLINE, "\n"},
		{token.REGEX, "\\/"},
		{token.NEWLINE, "\n"},
		{token.REGEX, "a"},
		{token.REGEX_MODIFIER, "i"},
		{token.NEWLINE, "\n"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for pos, testCase := range tests {
		if !lexer.HasNext() {
			t.Logf("Unexpected EOF at %d\n", lexer.pos)
			t.FailNow()
		}
		token := lexer.NextToken()

		if token.Type != testCase.expectedType {
			t.Logf("Expected token with type %q at position %d, got type %q\n", testCase.expectedType, pos, token.Type)
			t.FailNow()
		}

		if token.Literal != testCase.expectedLiteral {
			t.Logf("Expected token with literal %q at position %d, got literal %q\n", testCase.expectedLiteral, pos, token.Literal)
			t.FailNow()
		}
	}
}

func TestLexSpecificCase(t *testing.T) {
	input := `"\"" => -> (a) { val_to_str a }`

	expected := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.STRING, "\\\""},
		{token.HASHROCKET, "=>"},
		{token.LAMBDAROCKET, "->"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "val_to_str"},
		{token.IDENT, "a"},
		{token.RBRACE, "}"},
	}

	lexer := New(input)
	for pos, testCase := range expected {
		if !lexer.HasNext() {
			t.Logf("Unexpected EOF at %d\n", lexer.pos)
			t.FailNow()
		}
		token := lexer.NextToken()

		if token.Type != testCase.expectedType {
			t.Logf("Expected token with type %q at position %d, got type %q\n", testCase.expectedType, pos, token.Type)
			t.FailNow()
		}

		if token.Literal != testCase.expectedLiteral {
			t.Logf("Expected token with literal %q at position %d, got literal %q\n", testCase.expectedLiteral, pos, token.Literal)
			t.FailNow()
		}
	}
}

// Tests that the lexer can handle the source of pyra.rb
// https://github.com/ConorOBrien-Foxx/Pyramid-Scheme/blob/master/pyra.rb
func TestLexPyraRb(t *testing.T) {
	input := `def indices(str, chr)
    (0 ... str.length).find_all { |i| str[i] == chr }
end

def unwrap(t)
    t.size == 1 ? t[0] : t
end

$TOP    = "^"
$BOTTOM = "-"
$L_SIDE = "/"
$R_SIDE = "\\"

def triangle_from(lines, ptr_inds = nil)
    raise "no triangle found" if !lines.first
    ptr_inds = ptr_inds || indices(lines.first, $TOP)
    row = ""
    ptr_inds.map { |pt|
        x1 = x2 = pt    # left and right sides
        y = 0
        data = []
        loop {
            x1 -= 1
            x2 += 1
            y += 1
            row = lines[y]
            raise "unexpected triangle end" if !row or x2 > row.size
            # are we done?
            if row[x1] != $L_SIDE
                # look for end
                if row[x2] == $R_SIDE # mismatched!
                    raise "left side too short"
                else
                    # both sides are exhausted--look for a bottom
                    # p (x1 + 1 .. x2 - 1).map { |e| row[e] }
                    # p [x1, x2, pt]
                    if (x1 + 1 .. x2 - 1).all? { |x| row[x] == $BOTTOM }
                        break
                    else
                        raise "malformed bottom"
                    end
                end
            elsif row[x2] != $R_SIDE
                # look for end
                if row[x1] == $L_SIDE # mismatched!
                    raise "right side too short"
                else
                    # both sides are exhausted--look for a bottom
                    if (x1 + 1 .. x2 - 1).all? { |x| row[x] == $BOTTOM }
                        break
                    else
                        raise "malformed bottom"
                    end
                end
            # elsif x1 == 0   # we must have found our bottom...
            end
            #todo: do for other side
            # we aren't done.
            data.push row[x1 + 1 .. x2 - 1]
        }
        op = data.join("").gsub(/\s+/, "")
        args = []
        if row[x1] == $TOP or row[x2] == $TOP
            next_inds = [x1, x2].find_all { |x| row[x] == $TOP }
            args.push triangle_from(lines[y..-1], next_inds)
        end
        unwrap [op, *args]
    }
end

$vars = {"eps" => ""}
$UNDEF = :UNDEF

def parse(str)
    # find ^s on first line
    lns = str.lines
    triangle_from(lns)
end

# converts a string to a pyramid value
def str_to_val(str)
    # todo: expand
    if $vars.has_key? str
        $vars[str]
    elsif str == "line" or str == "stdin" or str == "readline"
        $stdin.gets
    else
        str.to_f
    end
end

def val_to_str(val)
    sanatize(val).to_s
end

def falsey(val)
    [0, [], "", $UNDEF, "\x00", nil].include? val
end

def truthy(val)
    !falsey val
end

class TrueClass;  def to_i; 1; end; end
class FalseClass; def to_i; 0; end; end

$outted = false

$uneval_ops = {
    "set" => -> (left, right) {
        $vars[left] = eval_chain right
        $UNDEF
    },
    # condition: left
    # body: right
    "do" => -> (left, right) {
        loop {
            eval_chain right
            break unless truthy eval_chain left
        }
        $UNDEF
    },
    # condition: left
    # body: right
    "loop" => -> (left, right) {
        loop {
            break unless truthy eval_chain left
            eval_chain right
        }
        $UNDEF
    },
    # condition: left
    # body: right
    "?" => -> (left, right) {
        truthy(eval_chain left) ? eval_chain(right) : 0
    }
}

$ops = {
    "+" => -> (a, b) { a + b },
    "*" => -> (a, b) { a * b },
    "-" => -> (a, b) { a - b },
    "/" => -> (a, b) { 1.0 * a / b },
    "^" => -> (a, b) { a ** b },
    "=" => -> (a, b) { (a == b).to_i },
    "<=>" => -> (a, b) { a <=> b },
    "out" => -> (*a) { $outted = true; a.each { |e| print e }; },
    "chr" => -> (a) { a.to_i.chr },
    "arg" => -> (*a) { a.size == 1 ? ARGV[a[0]] : a[0][a[1]] },
    "#" => -> (a) { str_to_val a },
    "\"" => -> (a) { val_to_str a },
    "" => -> (*a) { unwrap a },
    "!" => -> (a) { falsey(a).to_i },
    "[" => -> (a, b) { a },
    "]" => -> (a, b) { b },
}

def eval_chain(chain)
    if chain.is_a? String
        return str_to_val chain
    else
        op, args = chain
        if $uneval_ops.has_key? op
            return $uneval_ops[op][*args]
        end
        raise "undefined operation ` + "`" + `#{op}` + "`" + `" unless $ops.has_key? op
        return sanatize $ops[op][*sanatize(args.map { |ch| eval_chain ch })]
    end
end

def sanatize(arg)
    if arg.is_a? Array
        arg.map { |e| sanatize e }
    elsif arg.is_a? Float
        arg == arg.to_i ? arg.to_i : arg
    else
        arg
    end
end

prog = File.read(ARGV[0]).gsub(/\r/, "")
parsed = parse(prog)
res = parsed.map { |ch| eval_chain ch }
res = res.reject { |e| e == $UNDEF } if res.is_a? Array
res = res.is_a?(Array) && res.length == 1 ? res.pop : res
to_print = sanatize(res)
unless $outted
    if ARGV[1] && ARGV[1][1] == "d"
        p to_print
    else
        puts to_print
    end
end

# p $vars
`
	lexer := New(input)

	for {
		tk := lexer.NextToken()
		if tk.Type == token.EOF {
			break
		}
		// fmt.Println(tk)
	}
}
