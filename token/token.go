package token

import (
	"bytes"
	"strconv"
	"unicode"
)

//go:generate stringer -type=Type

// Recognized token types
const (
	ILLEGAL Type = iota // An illegal/unknown character
	EOF                 // end of input

	// Identifier + literals
	literal_beg
	IDENT
	CONST
	GLOBAL
	INT
	FLOAT
	STRING
	literal_end

	// Operators
	operator_beg
	operator_assign_beg
	ASSIGN    // =
	ADDASSIGN // +=
	SUBASSIGN // -=
	MULASSIGN // *=
	DIVASSIGN // /=
	MODASSIGN // %=
	operator_assign_end

	PLUS       // +
	MINUS      // -
	BANG       // !
	ASTERISK   // *
	POW        // **
	SLASH      // /
	MODULO     // %
	AND        // &
	LOGICALAND // &&
	PIPE       // |
	LOGICALOR  // ||

	LT        // <
	LTE       // <=
	GT        // >
	GTE       // >=
	EQ        // ==
	NOTEQ     // !=
	SPACESHIP // <=>
	LSHIFT    // <<
	operator_end

	HASHROCKET   // =>
	LAMBDAROCKET // ->

	// Delimiters

	NEWLINE   // \n
	COMMA     // ,
	SEMICOLON // ;
	COMMENT   // # ...
	CLASS     // class # NOTE: treated as a comment. classes are not supported

	DOT       // .
	DDOT      // ..
	DDDOT     // ...
	COLON     // :
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	SLBRACKET // _[
	RBRACKET  // ]

	QMARK  // ?
	SQMARK // _?
	SYMBEG // :

	// Keywords
	keyword_beg
	DEF
	END
	IF
	THEN
	ELSE
	ELSIF
	UNLESS
	TRUE
	FALSE
	RETURN
	NIL
	WHILE
	LOOP
	BREAK
	KEYWORD__FILE__
	keyword_end
	types_end
)

var marker_types = [...]Type{
	literal_beg,
	literal_end,
	operator_beg,
	operator_end,
	operator_assign_beg,
	operator_assign_end,
	keyword_beg,
	keyword_end,
	types_end,
}

var type_count = 0

func isMarkerType(t Type) bool {
	for _, marker := range marker_types {
		if t == marker {
			return true
		}
	}
	return false
}

func init() {
	for i := 0; i < int(types_end); i++ {
		if isMarkerType(Type(i)) {
			continue
		}
		type_count++
	}
}

var type_strings = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	IDENT:  "IDENT",
	CONST:  "CONST",
	GLOBAL: "GLOBAL",
	INT:    "INT",
	FLOAT:  "FLOAT",
	STRING: "STRING",

	ASSIGN:    "ASSIGN",
	ADDASSIGN: "ADDASSIGN",
	SUBASSIGN: "SUBASSIGN",
	MULASSIGN: "MULASSIGN",
	DIVASSIGN: "DIVASSIGN",
	MODASSIGN: "MODASSIGN",

	PLUS:       "PLUS",
	MINUS:      "MINUS",
	BANG:       "BANG",
	ASTERISK:   "ASTERISK",
	POW:        "POW",
	SLASH:      "SLASH",
	MODULO:     "MODULO",
	AND:        "AND",
	LOGICALAND: "LOGICALAND",
	LOGICALOR:  "LOGICALOR",

	LT:        "LT",
	LTE:       "LTE",
	GT:        "GT",
	GTE:       "GTE",
	EQ:        "EQ",
	NOTEQ:     "NOTEQ",
	SPACESHIP: "SPACESHIP",
	LSHIFT:    "LSHIFT",

	NEWLINE:   "NEWLINE",
	COMMA:     "COMMA",
	SEMICOLON: "SEMICOLON",
	COMMENT:   "COMMENT",
	CLASS:     "CLASS",

	DOT:       "DOT",
	DDOT:      "DDOT",
	DDDOT:     "DDDOT",
	COLON:     "COLON",
	LPAREN:    "LPAREN",
	RPAREN:    "RPAREN",
	LBRACE:    "LBRACE",
	RBRACE:    "RBRACE",
	LBRACKET:  "LBRACKET",
	SLBRACKET: "SLBRACKET",
	RBRACKET:  "RBRACKET",
	PIPE:      "PIPE",

	HASHROCKET:   "HASHROCKET",
	LAMBDAROCKET: "LAMBDAROCKET",

	QMARK:  "QMARK",
	SQMARK: "SQMARK",
	SYMBEG: "SYMBEG",

	DEF:             "DEF",
	END:             "END",
	UNLESS:          "UNLESS",
	IF:              "IF",
	THEN:            "THEN",
	ELSE:            "ELSE",
	ELSIF:           "ELSIF",
	TRUE:            "TRUE",
	FALSE:           "FALSE",
	RETURN:          "RETURN",
	NIL:             "NIL",
	WHILE:           "WHILE",
	LOOP:            "LOOP",
	BREAK:           "BREAK",
	KEYWORD__FILE__: "KEYWORD__FILE__",
}

var type_reprs = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	IDENT:  "IDENT",
	CONST:  "CONST",
	GLOBAL: "GLOBAL",
	INT:    "INT",
	FLOAT:  "FLOAT",
	STRING: "STRING",

	ASSIGN:    "=",
	ADDASSIGN: "+=",
	SUBASSIGN: "-=",
	MULASSIGN: "*=",
	DIVASSIGN: "/=",
	MODASSIGN: "%=",

	PLUS:       "+",
	MINUS:      "-",
	BANG:       "!",
	ASTERISK:   "*",
	POW:        "**",
	SLASH:      "/",
	MODULO:     "%",
	AND:        "&",
	LOGICALAND: "&&",
	LOGICALOR:  "||",

	LT:        "<",
	LTE:       "<=",
	GT:        ">",
	GTE:       ">=",
	EQ:        "==",
	NOTEQ:     "!=",
	SPACESHIP: "<=>",
	LSHIFT:    "<<",

	NEWLINE:   "\\n",
	COMMA:     ",",
	SEMICOLON: ";",
	COMMENT:   "# ...",
	CLASS:     "class",

	DOT:       ".",
	DDOT:      "..",
	DDDOT:     "...",
	COLON:     ":",
	LPAREN:    "(",
	RPAREN:    ")",
	LBRACE:    "{",
	RBRACE:    "}",
	LBRACKET:  "[",
	SLBRACKET: "_[",
	RBRACKET:  "]",
	PIPE:      "|",

	HASHROCKET:   "=>",
	LAMBDAROCKET: "->",

	QMARK:  "?",
	SQMARK: "_?",
	SYMBEG: ":",

	DEF:             "def",
	END:             "end",
	UNLESS:          "unless",
	IF:              "if",
	THEN:            "then",
	ELSE:            "else",
	ELSIF:           "elsif",
	TRUE:            "true",
	FALSE:           "false",
	RETURN:          "return",
	NIL:             "nil",
	WHILE:           "while", // to remove?
	LOOP:            "loop",
	BREAK:           "break",
	KEYWORD__FILE__: "__FILE__",
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token ADD, the string is
// "+"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
func (tok Type) String() string {
	s := ""
	if 0 <= tok && tok < Type(len(type_strings)) {
		s = type_strings[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

func ToType(s any) Type {
	switch s := s.(type) {
	case string:
		return typeFromString(s)
	case Type:
		return s
	case int:
		if s < 0 || s >= len(type_strings) {
			return ILLEGAL
		}
		return Type(s)
	default:
		panic("ToType: unknown type")
	}
}

func typeFromString(s string) Type {
	for i := 0; i < len(type_strings); i++ {
		if type_strings[i] == s {
			return Type(i)
		}
	}
	return ILLEGAL
}

var keywords map[string]Type

func init() {
	keywords = make(map[string]Type)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[type_reprs[i]] = i
	}
}

// LookupIdent returns a keyword Type if ident is a keyword. If ident starts
// with an upper character it returns CONST. In any other case it returns IDENT
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	if ident == "class" {
		return COMMENT
	}
	if unicode.IsUpper(bytes.Runes([]byte(ident))[0]) {
		return CONST
	}
	return IDENT
}

// A Type represents a type of a known token
type Type int

// NewToken returns a new Token associated with the given Type typ, the Literal
// literal and the Position pos
func NewToken(typ Type, literal string, pos int) Token {
	return Token{typ, literal, pos}
}

// A Token represents a known token with its literal representation
type Token struct {
	Type    Type
	Literal string
	Pos     int
}

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
func (t Token) IsLiteral() bool {
	return t.Type.IsLiteral()
}

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
func (t Token) IsOperator() bool {
	return t.Type.IsOperator()
}

// IsAssignOperator returns true for tokens corresponding to assignment
// operators and delimiters; it returns false otherwise.
func (t Token) IsAssignOperator() bool {
	return t.Type.IsAssignOperator()
}

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
func (t Token) IsKeyword() bool {
	return t.Type.IsKeyword()
}

// Predicates

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
func (tok Type) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
func (tok Type) IsOperator() bool { return operator_beg < tok && tok < operator_end }

// IsAssignOperator returns true for tokens corresponding to assignment
// operators and delimiters; it returns false otherwise.
func (tok Type) IsAssignOperator() bool {
	return operator_assign_beg < tok && tok < operator_assign_end
}

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
func (tok Type) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }
