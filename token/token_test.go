package token

import (
	"testing"
)

func findMissingSkipEmpty(present []string, arr []string) []string {
	present_map := make(map[string]bool)
	for _, repr := range present {
		present_map[repr] = true
	}

	// Find missing strings
	var missing []string
	for _, repr := range arr {
		if _, ok := present_map[repr]; !ok {
			if repr != "" { // skip empty strings
				missing = append(missing, repr)
			}
		}
	}

	return missing
}

func TestTypeSting(t *testing.T) {
	tests := []struct {
		tk   Type
		str  string
		repr string
	}{
		{tk: ILLEGAL, str: "ILLEGAL", repr: "ILLEGAL"},
		{tk: EOF, str: "EOF", repr: "EOF"},
		{tk: IDENT, str: "IDENT", repr: "IDENT"},
		{tk: CONST, str: "CONST", repr: "CONST"},
		{tk: GLOBAL, str: "GLOBAL", repr: "GLOBAL"},
		{tk: INT, str: "INT", repr: "INT"},
		{tk: FLOAT, str: "FLOAT", repr: "FLOAT"},
		{tk: STRING, str: "STRING", repr: "STRING"},
		{tk: REGEX, str: "REGEX", repr: "REGEX"},
		{tk: REGEX_MODIFIER, str: "REGEX_MODIFIER", repr: "REGEX_MODIFIER"},
		//
		{tk: ASSIGN, str: "ASSIGN", repr: "="},
		{tk: ADDASSIGN, str: "ADDASSIGN", repr: "+="},
		{tk: SUBASSIGN, str: "SUBASSIGN", repr: "-="},
		{tk: MULASSIGN, str: "MULASSIGN", repr: "*="},
		{tk: DIVASSIGN, str: "DIVASSIGN", repr: "/="},
		{tk: MODASSIGN, str: "MODASSIGN", repr: "%="},
		//
		{tk: PLUS, str: "PLUS", repr: "+"},
		{tk: MINUS, str: "MINUS", repr: "-"},
		{tk: BANG, str: "BANG", repr: "!"},
		{tk: ASTERISK, str: "ASTERISK", repr: "*"},
		{tk: SLASH, str: "SLASH", repr: "/"},
		{tk: MODULO, str: "MODULO", repr: "%"},
		{tk: AND, str: "AND", repr: "&"},
		{tk: LOGICALAND, str: "LOGICALAND", repr: "&&"},
		{tk: PIPE, str: "PIPE", repr: "|"},
		{tk: LOGICALOR, str: "LOGICALOR", repr: "||"},
		//
		{tk: LT, str: "LT", repr: "<"},
		{tk: LTE, str: "LTE", repr: "<="},
		{tk: GT, str: "GT", repr: ">"},
		{tk: GTE, str: "GTE", repr: ">="},
		{tk: EQ, str: "EQ", repr: "=="},
		{tk: NOTEQ, str: "NOTEQ", repr: "!="},
		{tk: SPACESHIP, str: "SPACESHIP", repr: "<=>"},
		{tk: LSHIFT, str: "LSHIFT", repr: "<<"},
		//
		{tk: HASHROCKET, str: "HASHROCKET", repr: "=>"},
		{tk: LAMBDAROCKET, str: "LAMBDAROCKET", repr: "->"},
		//
		{tk: NEWLINE, str: "NEWLINE", repr: "\\n"},
		{tk: COMMA, str: "COMMA", repr: ","},
		{tk: SEMICOLON, str: "SEMICOLON", repr: ";"},
		{tk: HASH, str: "HASH", repr: "#"},
		//
		{tk: CAPTURE, str: "CAPTURE", repr: "&"},
		{tk: DOT, str: "DOT", repr: "."},
		{tk: DDOT, str: "DDOT", repr: ".."},
		{tk: DDDOT, str: "DDDOT", repr: "..."},
		{tk: COLON, str: "COLON", repr: ":"},
		{tk: LPAREN, str: "LPAREN", repr: "("},
		{tk: RPAREN, str: "RPAREN", repr: ")"},
		{tk: LBRACE, str: "LBRACE", repr: "{"},
		{tk: RBRACE, str: "RBRACE", repr: "}"},
		{tk: LBRACKET, str: "LBRACKET", repr: "["},
		{tk: RBRACKET, str: "RBRACKET", repr: "]"},
		//
		{tk: SCOPE, str: "SCOPE", repr: "::"},
		{tk: AT, str: "AT", repr: "@"},
		//
		{tk: QMARK, str: "QMARK", repr: "?"},
		{tk: SQMARK, str: "SQMARK", repr: "_?"},
		{tk: SYMBEG, str: "SYMBEG", repr: ":"},
		//
		{tk: DEF, str: "DEF", repr: "def"},
		{tk: SELF, str: "SELF", repr: "self"},
		{tk: END, str: "END", repr: "end"},
		{tk: IF, str: "IF", repr: "if"},
		{tk: THEN, str: "THEN", repr: "then"},
		{tk: ELSE, str: "ELSE", repr: "else"},
		{tk: UNLESS, str: "UNLESS", repr: "unless"},
		{tk: TRUE, str: "TRUE", repr: "true"},
		{tk: FALSE, str: "FALSE", repr: "false"},
		{tk: RETURN, str: "RETURN", repr: "return"},
		{tk: NIL, str: "NIL", repr: "nil"},
		{tk: MODULE, str: "MODULE", repr: "module"},
		{tk: CLASS, str: "CLASS", repr: "class"},
		{tk: DO, str: "DO", repr: "do"},
		{tk: YIELD, str: "YIELD", repr: "yield"},
		{tk: BEGIN, str: "BEGIN", repr: "begin"},
		{tk: RESCUE, str: "RESCUE", repr: "rescue"},
		{tk: WHILE, str: "WHILE", repr: "while"},
		{tk: KEYWORD__FILE__, str: "KEYWORD__FILE__", repr: "__FILE__"},
	}

	seen := make(map[Type]bool)

	for _, test := range tests {
		if seen[test.tk] {
			t.Errorf("Duplicate token %s in tests!", test.tk)
		}
		seen[test.tk] = true

		if test.tk.String() != test.str {
			t.Errorf("Expected %s, got %s", test.str, test.tk.String())
		}
		if ToType(test.str) != test.tk {
			t.Errorf("Expected %s, got %s", test.str, ToType(test.str))
		}
		if test.repr != "" {
			if token_reprs[test.tk] != test.repr {
				t.Errorf("Expected %s, got %s", test.repr, token_reprs[test.tk])
			}
		} else {
			t.Errorf("Expected non-empty repr for %s", test.tk)
		}
	}

	test_strings := make([]string, 0)
	test_reprs := make([]string, 0)
	for _, test := range tests {
		test_strings = append(test_strings, test.str)
		test_reprs = append(test_reprs, test.repr)
	}

	missing := findMissingSkipEmpty(test_reprs, token_reprs[:])
	if len(missing) > 0 {
		t.Errorf("Missing tokens: %v", missing)
	}

	missing = findMissingSkipEmpty(test_strings, token_strings[:])
	if len(missing) > 0 {
		t.Errorf("Missing tokens: %v", missing)
	}

}
