package lexer

import (
	"os"
	"strings"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/token"
)

// loop {
// 	break unless x < 10
// }

// raise "foo" unless x < 10
type expected struct {
	typ token.Type
	lit string
}

func expect[T string | token.Token](tk T, literal string) expected {
	return expected{
		typ: token.ToType(tk),
		lit: literal,
	}
}

var NL = expect("NEWLINE", "\n")

func preprocessLines(lines string) string {
	slines := strings.Split(lines, "\n")
	new_slines := make([]string, 0)
	for _, line := range slines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		new_slines = append(new_slines, line)
	}
	return strings.Join(new_slines, "\n")
}

func allTokens(lexer *Lexer) []token.Token {
	tokens := []token.Token{}
	for {
		tk := lexer.NextToken()
		if tk.Type == token.EOF {
			break
		}
		tokens = append(tokens, tk)
	}
	return tokens
}

type test struct {
	desc  string
	lines string
	exp   []expected
}

func TestLex(t *testing.T) {
	tests := []test{
		{
			desc: "numbers",
			lines: `
				5
				5_0
				1.0
				123.456
				123.to_s
			`,
			exp: []expected{
				expect("INT", "5"),
				NL,
				expect("INT", "5_0"),
				NL,
				expect("FLOAT", "1.0"),
				NL,
				expect("FLOAT", "123.456"),
				NL,
				expect("INT", "123"),
				expect("DOT", "."),
				expect("IDENT", "to_s"),
			},
		},
		{
			desc: "assignments",
			lines: `
				five = 5
				Ten = 10
			`,
			exp: []expected{
				expect("IDENT", "five"),
				expect("ASSIGN", "="),
				expect("INT", "5"),
				NL,
				expect("CONST", "Ten"),
				expect("ASSIGN", "="),
				expect("INT", "10"),
			},
		},
		{
			desc: "strings",
			lines: `
				""
				"foobar"
				'foobar'
				"foo bar"
				'foo bar'
			`,
			exp: []expected{
				expect("STRING", ""),
				NL,
				expect("STRING", "foobar"),
				NL,
				expect("STRING", "foobar"),
				NL,
				expect("STRING", "foo bar"),
				NL,
				expect("STRING", "foo bar"),
			},
		},
		{
			desc: "symbols",
			lines: `
				:sym
				:"sym"
				:'sym'
				:dotAfter.
			`,
			exp: []expected{
				expect("SYMBEG", ":"),
				expect("IDENT", "sym"),
				NL,
				expect("SYMBEG", ":"),
				expect("STRING", "sym"),
				NL,
				expect("SYMBEG", ":"),
				expect("STRING", "sym"),
				NL,
				expect("SYMBEG", ":"),
				expect("IDENT", "dotAfter"),
				expect("DOT", "."),
			},
		},
		{
			desc: "ampersand",
			lines: `
				&foo
				&
				&&
			`,
			exp: []expected{
				expect("CAPTURE", "&"),
				expect("IDENT", "foo"),
				NL,
				expect("AND", "&"),
				NL,
				expect("LOGICALAND", "&&"),
			},
		},
		{
			desc: "operators",
			lines: `
				!-/ *%5;
				+= -= *= /= %=
				5 < 10 > 5
				10 == 10
				10 != 9
				10 <= 9
				10 >= 9
				10 <=> 9
				10 << 9
			`,
			exp: []expected{
				expect("BANG", "!"),
				expect("MINUS", "-"),
				expect("SLASH", "/"),
				expect("ASTERISK", "*"),
				expect("MODULO", "%"),
				expect("INT", "5"),
				expect("SEMICOLON", ";"),
				NL,
				expect("ADDASSIGN", "+="),
				expect("SUBASSIGN", "-="),
				expect("MULASSIGN", "*="),
				expect("DIVASSIGN", "/="),
				expect("MODASSIGN", "%="),
				NL,
				expect("INT", "5"),
				expect("LT", "<"),
				expect("INT", "10"),
				expect("GT", ">"),
				expect("INT", "5"),
				NL,
				expect("INT", "10"),
				expect("EQ", "=="),
				expect("INT", "10"),
				NL,
				expect("INT", "10"),
				expect("NOTEQ", "!="),
				expect("INT", "9"),
				NL,
				expect("INT", "10"),
				expect("LTE", "<="),
				expect("INT", "9"),
				NL,
				expect("INT", "10"),
				expect("GTE", ">="),
				expect("INT", "9"),
				NL,
				expect("INT", "10"),
				expect("SPACESHIP", "<=>"),
				expect("INT", "9"),
				NL,
				expect("INT", "10"),
				expect("LSHIFT", "<<"),
				expect("INT", "9"),
			},
		},
		{
			desc: "while",
			lines: `
				while x < y do
					x += x
				end
			`,
			exp: []expected{
				expect("WHILE", "while"),
				expect("IDENT", "x"),
				expect("LT", "<"),
				expect("IDENT", "y"),
				expect("DO", "do"),
				NL,
				expect("IDENT", "x"),
				expect("ADDASSIGN", "+="),
				expect("IDENT", "x"),
				NL,
				expect("END", "end"),
			},
		},
		{
			desc: "if",
			lines: `
				if x < y then
					true
				else
					false
				end
			`,
			exp: []expected{
				expect("IF", "if"),
				expect("IDENT", "x"),
				expect("LT", "<"),
				expect("IDENT", "y"),
				expect("THEN", "then"),
				NL,
				expect("TRUE", "true"),
				NL,
				expect("ELSE", "else"),
				NL,
				expect("FALSE", "false"),
				NL,
				expect("END", "end"),
			},
		},
		{
			desc: "qmark",
			lines: `
				foo.bar?
				foo.bar ? 1 : 2
				?-
				?\n
				? foo : bar
				`,
			exp: []expected{
				expect("IDENT", "foo"),
				expect("DOT", "."),
				expect("IDENT", "bar?"),
				NL,
				expect("IDENT", "foo"),
				expect("DOT", "."),
				expect("IDENT", "bar"),
				expect("SQMARK", " ?"),
				expect("INT", "1"),
				expect("COLON", ":"),
				expect("INT", "2"),
				NL,
				expect("STRING", "-"),
				NL,
				expect("STRING", "\\n"),
				NL,
				expect("QMARK", "?"),
				expect("IDENT", "foo"),
				expect("COLON", ":"),
				expect("IDENT", "bar"),
			},
		},
		{
			desc: "regex",
			lines: `
				/\//
				/a/i
			`,
			exp: []expected{
				expect("REGEX", "\\/"),
				NL,
				expect("REGEX", "a"),
				expect("REGEX_MODIFIER", "i"),
			},
		},
		{
			desc: "globals",
			lines: `
				$foo,
				$foo;
				$Foo
				$dotAfter.
				$@
			`,
			exp: []expected{
				expect("GLOBAL", "$foo"),
				expect("COMMA", ","),
				NL,
				expect("GLOBAL", "$foo"),
				expect("SEMICOLON", ";"),
				NL,
				expect("GLOBAL", "$Foo"),
				NL,
				expect("GLOBAL", "$dotAfter"),
				expect("DOT", "."),
				NL,
				expect("GLOBAL", "$@"),
			},
		},
		{
			desc: "defs",
			lines: `
				def add(x, y)
					return x + y
				end
				def nil?
				end
				def run!
				end
				module Abc
				end
				class Abc
				end
			`,
			exp: []expected{
				expect("DEF", "def"),
				expect("IDENT", "add"),
				expect("LPAREN", "("),
				expect("IDENT", "x"),
				expect("COMMA", ","),
				expect("IDENT", "y"),
				expect("RPAREN", ")"),
				NL,
				expect("RETURN", "return"),
				expect("IDENT", "x"),
				expect("PLUS", "+"),
				expect("IDENT", "y"),
				NL,
				expect("END", "end"),
				NL,
				expect("DEF", "def"),
				expect("IDENT", "nil?"),
				NL,
				expect("END", "end"),
				NL,
				expect("DEF", "def"),
				expect("IDENT", "run!"),
				NL,
				expect("END", "end"),
				NL,
				expect("MODULE", "module"),
				expect("CONST", "Abc"),
				NL,
				expect("END", "end"),
				NL,
				expect("CLASS", "class"),
				expect("CONST", "Abc"),
				NL,
				expect("END", "end"),
			},
		},
		{
			desc: "blocks",
			lines: `
				begin
					rescue
				end
				add { |x| x }
				add do |x|
				end
			`,
			exp: []expected{
				expect("BEGIN", "begin"),
				NL,
				expect("RESCUE", "rescue"),
				NL,
				expect("END", "end"),
				NL,
				expect("IDENT", "add"),
				expect("LBRACE", "{"),
				expect("PIPE", "|"),
				expect("IDENT", "x"),
				expect("PIPE", "|"),
				expect("IDENT", "x"),
				expect("RBRACE", "}"),
				NL,
				expect("IDENT", "add"),
				expect("DO", "do"),
				expect("PIPE", "|"),
				expect("IDENT", "x"),
				expect("PIPE", "|"),
				NL,
				expect("END", "end"),
			},
		},
		{
			desc: "comments",
			lines: `
				# just comment
				# just comment
			`,
			exp: []expected{
				expect("HASH", "#"),
				expect("STRING", " just comment"),
				NL,
				expect("HASH", "#"),
				expect("STRING", " just comment"),
			},
		},
		{
			desc: "pipe",
			lines: `
				|
				||
			`,
			exp: []expected{
				expect("PIPE", "|"),
				NL,
				expect("LOGICALOR", "||"),
			},
		},
		{
			desc: "misc",
			lines: `
				A::B
				=>
				->
				__FILE__
				@
				self
				nil
				yield
			`,
			exp: []expected{
				expect("CONST", "A"),
				expect("SCOPE", "::"),
				expect("CONST", "B"),
				NL,
				expect("HASHROCKET", "=>"),
				NL,
				expect("LAMBDAROCKET", "->"),
				NL,
				expect("KEYWORD__FILE__", "__FILE__"),
				NL,
				expect("AT", "@"),
				NL,
				expect("SELF", "self"),
				NL,
				expect("NIL", "nil"),
				NL,
				expect("YIELD", "yield"),
			},
		},
		{
			desc: "range",
			lines: `
				1..2
				1...2	
				..
				...
			`,
			exp: []expected{
				expect("INT", "1"),
				expect("DDOT", ".."),
				expect("INT", "2"),
				NL,
				expect("INT", "1"),
				expect("DDDOT", "..."),
				expect("INT", "2"),
				NL,
				expect("DDOT", ".."),
				NL,
				expect("DDDOT", "..."),
			},
		},
		{
			desc: "brackets",
			lines: `
				[1, 2]
			`,
			exp: []expected{
				expect("LBRACKET", "["),
				expect("INT", "1"),
				expect("COMMA", ","),
				expect("INT", "2"),
				expect("RBRACKET", "]"),
			},
		},
		{
			desc: "hash_of_lambdas",
			lines: `
				"\"" => -> (a) { val_to_str a }
			`,
			exp: []expected{
				expect("STRING", "\\\""),
				expect("HASHROCKET", "=>"),
				expect("LAMBDAROCKET", "->"),
				expect("LPAREN", "("),
				expect("IDENT", "a"),
				expect("RPAREN", ")"),
				expect("LBRACE", "{"),
				expect("IDENT", "val_to_str"),
				expect("IDENT", "a"),
				expect("RBRACE", "}"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			lexer := New(preprocessLines(test.lines))
			tokens := allTokens(lexer)
			if len(tokens) != len(test.exp) {
				t.Fatalf("Expected %d tokens, got %d", len(test.exp), len(tokens))
			}
			for i, exp := range test.exp {
				if tokens[i].Type != exp.typ {
					t.Fatalf("Expected token %d to be %q, got %q", i, exp.typ, tokens[i].Type)
				}
				if tokens[i].Literal != exp.lit {
					t.Fatalf("Expected token %d to be %q, got %q", i, exp.lit, tokens[i].Literal)
				}
			}
		})
	}
}

// Tests that the lexer can handle the source of pyra.rb
// https://github.com/ConorOBrien-Foxx/Pyramid-Scheme/blob/master/pyra.rb
func TestLexPyraRb(t *testing.T) {
	filename := "../pyra.rb"
	file, err := os.ReadFile(filename)
	if err != nil {
		t.Skip("Skipping test, file not found:", filename)
	}
	lexer := New(string(file))

	for {
		tk := lexer.NextToken()
		if tk.Type == token.EOF {
			break
		}
	}
}
