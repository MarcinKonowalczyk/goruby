package parser

import (
	"fmt"
	gotoken "go/token"
	"strconv"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/ast/infix"
	"github.com/MarcinKonowalczyk/goruby/lexer"
	"github.com/MarcinKonowalczyk/goruby/token"
	"github.com/pkg/errors"
)

// Possible precedence values
const (
	_ int = iota
	precLowest
	precBlockDo     // do
	precBlockBraces // { |x| }
	precIfUnless    // modifier-if, modifier-unless
	precAssignment  // x = 5
	precTernary     // ?, :
	precRange       // .., ...
	precLogicalOr   // ||
	precLogicalAnd  // &&
	precEquals      // ==, !=, <=>
	precLessGreater // >, <, >=, <=
	precOr          // |
	precAnd         // &
	precShift       // <<
	precSum         // + or -
	precProduct     // *, /, %
	precPower       // **
	precPrefix      // -X or !X
	precSplat       // x = [*y, 1]
	precCallArg     // func x
	precCall        // foo.myFunction(X)
	precIndex       // array[index]
	precCapture     // &block
	precSymbol      // :Symbol
	precHighest
)

var precedences = map[token.Type]int{
	token.IF:         precIfUnless,
	token.UNLESS:     precIfUnless,
	token.EQ:         precEquals,
	token.NOTEQ:      precEquals,
	token.SPACESHIP:  precEquals,
	token.LSHIFT:     precShift,
	token.QMARK:      precTernary,
	token.SQMARK:     precTernary,
	token.COLON:      precTernary,
	token.LT:         precLessGreater,
	token.GT:         precLessGreater,
	token.LTE:        precLessGreater,
	token.GTE:        precLessGreater,
	token.PLUS:       precSum,
	token.MINUS:      precSum,
	token.SLASH:      precProduct,
	token.ASTERISK:   precProduct,
	token.POW:        precPower,
	token.MODULO:     precProduct,
	token.ASSIGN:     precAssignment,
	token.ADDASSIGN:  precAssignment,
	token.SUBASSIGN:  precAssignment,
	token.MULASSIGN:  precAssignment,
	token.DIVASSIGN:  precAssignment,
	token.MODASSIGN:  precAssignment,
	token.LPAREN:     precCall,
	token.DOT:        precCall,
	token.DDOT:       precRange,
	token.DDDOT:      precRange,
	token.IDENT:      precCallArg,
	token.CONST:      precCallArg,
	token.GLOBAL:     precCallArg,
	token.INT:        precCallArg,
	token.STRING:     precCallArg,
	token.SLBRACKET:  precCallArg,
	token.LBRACKET:   precIndex,
	token.LBRACE:     precBlockBraces,
	token.DO:         precBlockDo,
	token.SYMBEG:     precSymbol,
	token.COMMA:      precAssignment,
	token.THEN:       precHighest,
	token.NEWLINE:    precHighest,
	token.PIPE:       precOr,
	token.AND:        precAnd,
	token.LOGICALOR:  precLogicalOr,
	token.LOGICALAND: precLogicalAnd,
	token.CAPTURE:    precCapture,
}

var tokensNotPossibleInCallArgs = []token.Type{
	token.ASSIGN,
	token.LT,
	token.LTE,
	token.GT,
	token.GTE,
	token.SPACESHIP,
	token.LSHIFT,
	token.EQ,
	token.NOTEQ,
	token.IF,
	token.UNLESS,
	token.COLON,
	token.RBRACKET,
	token.COMMA,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// A parser parses the token emitted by the provided lexer.Lexer and returns an
// AST describing the parsed program.
type parser struct {
	file   *gotoken.File
	l      *lexer.Lexer
	errors []error

	// Tracing/debugging
	mode   Mode // parsing mode
	trace  bool // == (mode & Trace != 0)
	indent int  // indentation used for tracing output

	pos       gotoken.Pos
	lastLine  string
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func (p *parser) init(fset *gotoken.FileSet, filename string, src []byte, mode Mode) {
	p.file = fset.AddFile(filename, -1, len(src))

	p.l = lexer.New(string(src))
	p.errors = []error{}

	p.mode = mode
	p.trace = mode&Trace != 0 // for convenience (p.trace is used frequently)

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.CONST, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.UNLESS, p.parseIfExpression)
	p.registerPrefix(token.WHILE, p.parseLoopExpression)
	p.registerPrefix(token.LOOP, p.parseLoopExpression)
	p.registerPrefix(token.DEF, p.parseFunctionLiteral)
	p.registerPrefix(token.SYMBEG, p.parseSymbolLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.SLBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.NIL, p.parseNilLiteral)
	p.registerPrefix(token.LBRACE, p.parseHash)
	p.registerPrefix(token.DO, p.parseBlock)
	p.registerPrefix(token.GLOBAL, p.parseGlobal)
	p.registerPrefix(token.KEYWORD__FILE__, p.parseKeyword__FILE__)
	p.registerPrefix(token.CAPTURE, p.parseBlockCapture)
	p.registerPrefix(token.LAMBDAROCKET, p.parseLambdaLiteral)
	p.registerPrefix(token.ASTERISK, p.parseSplat)

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.POW, p.parseInfixExpression)
	p.registerInfix(token.MODULO, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.PIPE, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOTEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.LOGICALOR, p.parseInfixExpression)
	p.registerInfix(token.LOGICALAND, p.parseInfixExpression)
	p.registerInfix(token.SPACESHIP, p.parseInfixExpression)
	p.registerInfix(token.LSHIFT, p.parseInfixExpression)
	p.registerInfix(token.DDOT, p.parseRangeLiteral)
	p.registerInfix(token.DDDOT, p.parseRangeLiteral)
	p.registerInfix(token.ASSIGN, p.parseAssignment)
	p.registerInfix(token.ADDASSIGN, p.parseAssignmentOperator)
	p.registerInfix(token.SUBASSIGN, p.parseAssignmentOperator)
	p.registerInfix(token.MULASSIGN, p.parseAssignmentOperator)
	p.registerInfix(token.DIVASSIGN, p.parseAssignmentOperator)
	p.registerInfix(token.MODASSIGN, p.parseAssignmentOperator)
	p.registerInfix(token.IF, p.parseModifierConditionalExpression)
	p.registerInfix(token.UNLESS, p.parseModifierConditionalExpression)
	p.registerInfix(token.QMARK, p.parseTernaryIfExpression)
	p.registerInfix(token.SQMARK, p.parseTernaryIfExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpressionWithParens)
	p.registerInfix(token.IDENT, p.parseCallArgument)
	p.registerInfix(token.CONST, p.parseCallArgument)
	p.registerInfix(token.GLOBAL, p.parseCallArgument)
	p.registerInfix(token.INT, p.parseCallArgument)
	p.registerInfix(token.STRING, p.parseCallArgument)
	p.registerInfix(token.SYMBEG, p.parseCallArgument)
	p.registerInfix(token.CAPTURE, p.parseCallArgument)
	p.registerInfix(token.LBRACE, p.parseCallBlock)
	p.registerInfix(token.DO, p.parseCallBlock)
	p.registerInfix(token.DOT, p.parseMethodCall)
	p.registerInfix(token.COMMA, p.parseExpressions)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.SLBRACKET, p.parseCallArgument)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
}

func (p *parser) printTrace(a ...interface{}) {
	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)
	pos := p.file.Position(p.pos)
	fmt.Printf("%5d:%3d: ", pos.Line, pos.Column)
	i := 2 * p.indent
	for i > n {
		fmt.Print(dots)
		i -= n
	}
	// i <= n
	fmt.Print(dots[0:i])
	fmt.Println(a...)
}

func trace(p *parser, msg string) *parser {
	p.printTrace(msg, "(")
	p.indent++
	return p
}

// Usage pattern: defer un(trace(p, "..."))
func un(p *parser) {
	p.indent--
	p.printTrace(")")
}

func (p *parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *parser) nextToken() {
	// The very first token (!p.pos.IsValid()) is not initialized
	// (it is token.ILLEGAL), so don't print it .
	if p.trace && p.pos.IsValid() {
		s := p.peekToken.Type.String()
		switch {
		case p.peekToken.IsLiteral():
			p.printTrace(s, p.peekToken.Literal)
		case p.peekToken.IsOperator(), p.peekToken.IsKeyword():
			p.printTrace("\"" + s + "\"")
		default:
			p.printTrace(s)
		}
	}
	p.curToken = p.peekToken
	p.pos = gotoken.Pos(p.curToken.Pos)
	p.lastLine += p.curToken.Literal
	if p.curToken.Type == token.NEWLINE {
		p.file.AddLine(int(p.pos))
		p.lastLine = ""
	}
	if p.l.HasNext() {
		p.peekToken = p.l.NextToken()
	} else {
		p.peekToken = token.NewToken(token.EOF, "", -1)
	}
}

func (p *parser) parseIdentifier() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseIdentifier"))
	}
	return &ast.Identifier{Constant: p.currentTokenIs(token.CONST), Value: p.curToken.Literal}
}

// Constant: p.currentTokenIs(token.CONST),

// Errors returns all errors which happened during the parsing of the input.
func (p *parser) Errors() []error {
	return p.errors
}

func (p *parser) Error(actual token.Type, description string, expected ...token.Type) {
	epos := p.file.Position(p.pos)
	err := &UnexpectedTokenError{
		Pos:            epos,
		ExpectedTokens: expected,
		ActualToken:    actual,
		Description:    description,
	}
	p.errors = append(p.errors, errors.WithStack(err))
}

// func (p *parser) peekErrorDesc(desc string, t ...token.Type) {
// 	epos := p.file.Position(p.pos)
// 	err := &UnexpectedTokenError{
// 		Pos:            epos,
// 		ExpectedTokens: t,
// 		ActualToken:    p.peekToken.Type,
// 		Description:    desc,
// 	}
// 	p.errors = append(p.errors, errors.WithStack(err))
// }

// func (p *parser) expectError(t ...token.Type) {
// 	epos := p.file.Position(p.pos)
// 	err := &UnexpectedTokenError{
// 		Pos:            epos,
// 		ExpectedTokens: t,
// 		ActualToken:    p.curToken.Type,
// 	}
// 	p.errors = append(p.errors, errors.WithStack(err))
// }

// func (p *parser) expectErrorDesc(desc string, t ...token.Type) {
// 	epos := p.file.Position(p.pos)
// 	err := &UnexpectedTokenError{
// 		Pos:            epos,
// 		ExpectedTokens: t,
// 		ActualToken:    p.curToken.Type,
// 		Description:    desc,
// 	}
// 	p.errors = append(p.errors, errors.WithStack(err))
// }

func (p *parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for type %s found", t)
	epos := p.file.Position(p.pos)
	if epos.Filename != "" || epos.IsValid() {
		msg = epos.String() + ": " + msg
	}
	p.errors = append(p.errors, errors.Errorf("%s", msg))

}

// ParseProgram returns the parsed program AST and all errors which occurred
// during the parse process. If the error is not nil the AST may be incomplete
// and callers should always check if they can handle the error with providing
// more input by checking with e.g. IsEOFError.
func (p *parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for !p.currentTokenIs(token.EOF) {
		if p.currentTokenIs(token.NEWLINE) {
			// Early exit
			p.nextToken()
			continue
		}
		// fmt.Println(p.curToken, p.errors)
		stmt := p.parseStatement()
		// fmt.Println(stmt, reflect.TypeOf(stmt), p.errors)
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	// fmt.Println("Last token:", p.curToken)
	if len(p.errors) != 0 {
		return program, NewErrors("Parsing errors", p.errors...)
	}
	// fmt.Println("Program:", program)
	return program, nil
}

func (p *parser) parseStatement() ast.Statement {
	if p.trace {
		defer un(trace(p, "parseStatement"))
	}
	switch p.curToken.Type {
	case token.ILLEGAL:
		msg := p.curToken.Literal
		epos := p.file.Position(p.pos)
		if epos.Filename != "" || epos.IsValid() {
			msg = epos.String() + ": " + msg
		}
		p.errors = append(p.errors, fmt.Errorf("%s", msg))
		return nil
	case token.EOF:
		p.Error(p.curToken.Type, "", token.NEWLINE)
		return nil
	case token.NEWLINE:
		return nil
	case token.RETURN:
		return p.parseReturnStatement()
	case token.COMMENT:
		return p.parseComment()
	case token.BREAK:
		return p.parseBreakStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *parser) parseReturnStatement() *ast.ReturnStatement {
	if p.trace {
		defer un(trace(p, "parseReturnStatement"))
	}
	stmt := &ast.ReturnStatement{}
	p.nextToken()

	if p.currentTokenOneOf(token.NEWLINE, token.SEMICOLON) {
		p.nextToken()
		return stmt
	}

	stmt.ReturnValue = p.parseExpression(precLowest)
	if list, ok := stmt.ReturnValue.(ast.ExpressionList); ok {
		stmt.ReturnValue = &ast.ArrayLiteral{Elements: list}
	}

	if p.peekTokenOneOf(token.NEWLINE, token.SEMICOLON) {
		p.nextToken()
		return stmt
	}
	if p.currentTokenOneOf(token.NEWLINE, token.SEMICOLON) {
		return stmt
	}

	if !p.peekTokenIs(token.COMMA) {
		p.Error(p.peekToken.Type, "", token.COMMA)
		return nil
	}

	arr := &ast.ArrayLiteral{Elements: []ast.Expression{stmt.ReturnValue}}
	for p.peekTokenIs(token.COMMA) {
		p.consume(token.COMMA)
		arr.Elements = append(arr.Elements, p.parseExpression(precLowest))
	}
	stmt.ReturnValue = arr

	if !p.acceptOneOf(token.NEWLINE, token.SEMICOLON) {
		return nil
	}
	return stmt
}

func (p *parser) parseExpressionStatement() *ast.ExpressionStatement {
	if p.trace {
		defer un(trace(p, "parseExpressionStatement"))
	}
	stmt := &ast.ExpressionStatement{}
	stmt.Expression = p.parseExpression(precLowest)
	if p.peekTokenOneOf(token.SEMICOLON, token.NEWLINE) {
		p.nextToken()
	}
	return stmt
}

func (p *parser) parseExpression(precedence int) ast.Expression {
	if p.trace {
		defer un(trace(p, "parseExpression"))
	}
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for precedence < p.peekPrecedence() {
		if leftExp == nil {
			return nil // fail early and stop parsing
		}
		if p.currentTokenOneOf(token.NEWLINE, token.SEMICOLON) {
			return leftExp
		}
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *parser) parseComment() ast.Statement {
	if p.trace {
		defer un(trace(p, "parseComment"))
	}
	comment := &ast.Comment{}
	comment.Value = p.curToken.Literal
	if !p.peekTokenOneOf(token.NEWLINE, token.EOF) {
		epos := p.file.Position(p.pos)
		msg := fmt.Errorf("%s: Expected newline or eof after comment", epos.String())
		p.errors = append(p.errors, msg)
		return nil
	}

	if p.mode&ParseComments == 0 {
		return nil
	}
	return comment
}

// LAMBDA          : "->" "(" CALL_ARGS ")" "{" COMPSTMT "}"
//                 | "->" "{" COMPSTMT "}";

func (p *parser) parseLambdaLiteral() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseLambdaLiteral"))
	}
	proc := &ast.ProcedureLiteral{}
	if p.peekTokenIs(token.LPAREN) {
		proc.Parameters = p.parseFunctionParameters(token.LPAREN, token.RPAREN)
	}
	if !p.accept(token.LBRACE) {
		return nil
	}
	proc.Body = p.parseBlockStatement(token.RBRACE)
	if proc.Body == nil {
		return nil
	}
	if !p.accept(token.RBRACE) {
		return nil
	}
	return proc
}

func (p *parser) parseSplat() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseSplat"))
	}
	splat := &ast.Splat{}
	// p.nextToken()
	// if p.curToken.Type != token.IDENT {
	// 	p.Error(p.curToken.Type, "", token.IDENT)
	// 	return nil
	// }
	p.nextToken()
	// TODO: what's supposed t be the prec here?
	// precLowest breaks for `x = [*y, 1]`
	expr := p.parseExpression(precSplat)
	// fmt.Println("Splat expression:", expr)
	// panic("TODO: splat expression")
	if expr == nil {
		return nil
	}
	splat.Value = expr
	return splat
}

func (p *parser) parseExpressions(left ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace(p, "parseExpressions"))
	}
	p.nextToken()
	elements := []ast.Expression{left}
	next := p.parseExpression(precAssignment)
	elements = append(elements, next)
	for p.peekTokenIs(token.COMMA) {
		p.consume(token.COMMA)
		next = p.parseExpression(precAssignment)
		elements = append(elements, next)
	}
	return ast.ExpressionList(elements)
}

func (p *parser) parseBlockCapture() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseBlockCapture"))
	}
	capture := &ast.BlockCapture{}
	if !p.accept(token.IDENT) {
		return nil
	}
	capture.Name = &ast.Identifier{Constant: p.currentTokenIs(token.CONST), Value: p.curToken.Literal}
	return capture
}

func (p *parser) parseAssignmentOperator(left ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace(p, "parseAssignmentOperator"))
	}
	assignIndex := strings.LastIndexByte(p.curToken.Literal, '=')
	if assignIndex < 0 {
		return nil
	}
	op := infix.InfixFromAssignmentOperator(p.curToken)
	if op == infix.ILLEGAL {
		epos := p.file.Position(p.pos)
		msg := fmt.Errorf("%s: invalid assignment operator %q", epos.String(), p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	newInf := &ast.InfixExpression{
		Left:     left,
		Operator: op,
	}
	assign := &ast.Assignment{
		Left: left,
	}
	p.nextToken()
	newInf.Right = p.parseExpression(precLowest)
	assign.Right = newInf
	return assign
}

func (p *parser) parseAssignment(left ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace(p, "parseAssignment"))
	}

	switch left.(type) {
	case *ast.Identifier:
	case *ast.Global:
	case *ast.IndexExpression:
	case ast.ExpressionList:
	case *ast.Keyword__FILE__:
		epos := p.file.Position(p.pos)
		msg := fmt.Errorf("%s: Can't assign to __FILE__", epos.String())
		p.errors = append(p.errors, msg)
		return nil
	default:
		p.Error(p.curToken.Type, "", token.EOF)
		return nil
	}

	assign := &ast.Assignment{
		Left: left,
	}
	p.nextToken()
	expr := p.parseExpression(precLowest)
	right, ok := expr.(*ast.ConditionalExpression)
	if !ok {
		assign.Right = expr
		return assign
	}
	expStmt, ok := right.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		p.errors = append(p.errors, fmt.Errorf("malformed AST in assignment"))
		return nil
	}
	assign.Right = expStmt.Expression
	cond := &ast.ConditionalExpression{
		Unless:    right.Unless,
		Condition: right.Condition,
		Consequence: &ast.BlockStatement{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{Expression: assign},
			},
		},
	}

	return cond
}

func (p *parser) parseNilLiteral() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseNilLiteral"))
	}
	return &ast.Nil{}
}

func (p *parser) parseGlobal() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseGlobal"))
	}
	return &ast.Global{Value: p.curToken.Literal}
}

func (p *parser) parseKeyword__FILE__() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseKeyword__FILE__"))
	}
	file := &ast.Keyword__FILE__{
		Filename: p.file.Name(),
	}
	return file
}

var integerLiteralReplacer = strings.NewReplacer("_", "")

func (p *parser) parseIntegerLiteral() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseIntegerLiteral"))
	}
	lit := &ast.IntegerLiteral{}
	value, err := strconv.ParseInt(integerLiteralReplacer.Replace(p.curToken.Literal), 0, 64)
	if err != nil {
		msg := fmt.Errorf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *parser) parseFloatLiteral() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseFloatLiteral"))
	}
	lit := &ast.FloatLiteral{}
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Errorf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit

}
func (p *parser) parseStringLiteral() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseStringLiteral"))
	}
	return &ast.StringLiteral{
		Value: p.curToken.Literal}
}

func (p *parser) parseSymbolLiteral() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseSymbolLiteral"))
	}
	symbol := &ast.SymbolLiteral{}
	if !p.acceptOneOf(token.IDENT, token.STRING, token.CONST) {
		return nil
	}
	val := p.parseExpression(precHighest)
	symbol.Value = val
	return symbol
}

func (p *parser) parseArrayLiteral() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseArrayLiteral"))
	}
	array := &ast.ArrayLiteral{}

	p.nextToken()
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *parser) parseBoolean() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseBoolean"))
	}
	return &ast.Boolean{Value: p.currentTokenIs(token.TRUE)}
}

func (p *parser) consumeNewlineOrComment() {
	for {
		if p.currentTokenIs(token.NEWLINE) {
			// consume the newline
			p.nextToken()
		} else if p.currentTokenIs(token.COMMENT) {
			// consume the comment
			p.nextToken()
			if p.currentTokenIs(token.STRING) {
				p.nextToken()
			}
		} else {
			break
		}
	}
}
func (p *parser) parseHash() ast.Expression {
	hash := &ast.HashLiteral{Map: make(map[ast.Expression]ast.Expression)}
	if p.trace {
		defer un(trace(p, "parseHash"))
	}
	p.nextToken()

	p.consumeNewlineOrComment()

	if p.currentTokenIs(token.RBRACE) {
		return hash
	}

	// parse the first key-value pair
	k, v, ok := p.parseKeyValue()
	if !ok {
		return nil
	}
	hash.Map[k] = v
	p.nextToken() // move past the end of the key-value pair
	p.consumeNewlineOrComment()

	for p.currentTokenIs(token.COMMA) {
		p.nextToken()
		p.consumeNewlineOrComment()
		if p.currentTokenIs(token.RBRACE) {
			break
		}
		k, v, ok := p.parseKeyValue()
		if !ok {
			return nil
		}
		hash.Map[k] = v
		p.nextToken() // move past the end of the key-value pair
		p.consumeNewlineOrComment()
	}
	if p.currentTokenIs(token.RBRACE) {
		// we've landed on a RBRACE
		// this is fine, no need to do anything
	} else if p.peekTokenIs(token.RBRACE) {
		// we've landed on the end of the last key-value pair (hopefully)
		// we need to accept the RBRACE
		if !p.accept(token.RBRACE) {
			return nil
		}
	}
	return hash
}

// func (p *parser) debugPrintState() {
// 	if p.trace {
// 		fmt.Println("DBG:", p.curToken, p.peekToken, p.errors)
// 	}
// }

func (p *parser) parseKeyValue() (ast.Expression, ast.Expression, bool) {
	if p.trace {
		defer un(trace(p, "parseKeyValue"))
	}
	key := p.parseExpression(precAssignment)
	if !p.consume(token.HASHROCKET) {
		return nil, nil, false
	}
	val := p.parseExpression(precAssignment)
	return key, val, true
}

func (p *parser) parseBlock() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseBlock"))
	}
	start_token := p.curToken
	block := &ast.BlockExpression{}
	if p.peekTokenIs(token.PIPE) {
		block.Parameters = p.parseFunctionParameters(token.PIPE, token.PIPE)
	}

	if p.peekTokenOneOf(token.NEWLINE, token.SEMICOLON) {
		p.acceptOneOf(token.NEWLINE, token.SEMICOLON)
	}

	endToken := token.RBRACE
	if start_token.Type == token.DO {
		endToken = token.END
	}

	block.Body = p.parseBlockStatement(endToken)
	p.nextToken()
	return block
}

func (p *parser) parsePrefixExpression() ast.Expression {
	if p.trace {
		defer un(trace(p, "parsePrefixExpression"))
	}
	expression := &ast.PrefixExpression{
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(precPrefix)
	return expression
}

func (p *parser) parseInfixExpression(left ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace(p, "parseInfixExpression"))
	}

	op := infix.InfixFromToken(p.curToken)
	if op == infix.ILLEGAL {
		p.errors = append(p.errors, fmt.Errorf("illegal infix operator %q", p.curToken.Literal))
		return nil
	}

	expression := &ast.InfixExpression{
		Operator: op,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *parser) parseRangeLiteral(left ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace(p, "parseRangeLiteral"))
	}
	expression := &ast.RangeLiteral{
		Left:      left,
		Inclusive: p.curToken.Type == token.DDOT,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *parser) parseIndexExpression(left ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace(p, "parseIndexExpression"))
	}
	exp := &ast.IndexExpression{Left: left}

	p.nextToken()
	content := p.parseExpression(precLowest)
	exp.Index = content

	// if elist, ok := exp.Index.(ast.ExpressionList); ok {
	// 	// fmt.Println("Got an ExpressionList for the IndexExpression:", elist)
	// 	exp.Index = elist
	// } else if splat, ok := exp.Index.(*ast.Splat); ok {
	// 	// Indexing with a splat, e.g. array[*a]
	// 	exp.Index = splat
	// }

	if !p.accept(token.RBRACKET) {
		return nil
	}
	// fmt.Println("IndexExpression:", exp)
	return exp
}

func (p *parser) parseGroupedExpression() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseGroupedExpression"))
	}
	p.nextToken()
	exp := p.parseExpression(precLowest)
	if !p.accept(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *parser) parseIfExpression() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseIfExpression"))
	}
	expression := &ast.ConditionalExpression{Unless: p.currentTokenIs(token.UNLESS)}
	p.nextToken()
	expression.Condition = p.parseExpression(precLowest)
	if p.peekTokenIs(token.THEN) {
		p.accept(token.THEN)
	}

	// there may be a comment here. gobble it up
	if p.peekTokenIs(token.COMMENT) {
		p.accept(token.COMMENT)
	}

	if !p.peekTokenOneOf(token.NEWLINE, token.SEMICOLON) {
		p.Error(p.peekToken.Type, "could not parse if expression", token.NEWLINE, token.SEMICOLON)
		return nil
	}
	p.acceptOneOf(token.NEWLINE, token.SEMICOLON)

	consequence := p.parseBlockStatement(token.ELSE, token.ELSIF)
	expression.Consequence = consequence
	parsed_elsif := false
	if p.peekTokenIs(token.ELSE) {
		p.accept(token.ELSE)
		p.accept(token.NEWLINE)
		expression.Alternative = p.parseBlockStatement() // parse until the END
	} else if p.peekTokenIs(token.ELSIF) {
		// start parsing a new if expression from here
		p.accept(token.ELSIF)
		alternative := p.parseIfExpression().(*ast.ConditionalExpression)
		alternative_block := &ast.BlockStatement{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Expression: alternative,
				},
			},
		}
		expression.Alternative = alternative_block
		// back up one token
		parsed_elsif = true
	}
	if parsed_elsif {
		// don't expect the END, since it has been gobbled up by the innermost
		// if expression. check that it's there, though.
		if !p.currentTokenIs(token.END) {
			p.Error(p.curToken.Type, "expected END after elsif", token.END)
		}
	} else {
		p.accept(token.END)
	}
	return expression
}

func (p *parser) parseTernaryIfExpression(condition ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace(p, "parseTernaryIfExpression"))
	}
	expression := &ast.ConditionalExpression{Unless: p.currentTokenIs(token.UNLESS)}
	p.nextToken()
	expression.Condition = condition
	expression.Consequence = &ast.BlockStatement{
		Statements: []ast.Statement{
			&ast.ExpressionStatement{
				Expression: p.parseExpression(precTernary),
			},
		},
	}
	p.consume(token.COLON)
	expression.Alternative = &ast.BlockStatement{
		Statements: []ast.Statement{
			&ast.ExpressionStatement{
				Expression: p.parseExpression(precLowest),
			},
		},
	}
	return expression
}

func (p *parser) parseModifierConditionalExpression(left ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace(p, "parseModifierConditionalExpression"))
	}
	expression := &ast.ConditionalExpression{Unless: p.currentTokenIs(token.UNLESS)}
	p.nextToken()
	expression.Condition = p.parseExpression(precLowest)

	expression.Consequence = &ast.BlockStatement{
		Statements: []ast.Statement{
			&ast.ExpressionStatement{Expression: left},
		},
	}
	return expression
}

func (p *parser) parseLoopExpression() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseLoopExpression"))
	}
	loop := &ast.LoopExpression{}
	if p.curToken.Type == token.WHILE {
		p.nextToken()
		loop.Condition = p.parseExpression(precBlockDo)
		if p.peekTokenIs(token.DO) {
			p.accept(token.DO)
		}
		loop.Block = p.parseBlockStatement(token.END)
		p.nextToken()
	} else if p.curToken.Type == token.LOOP {
		loop.Condition = &ast.Boolean{Value: true}
		if p.peekTokenIs(token.DO) {
			// body is a do block
			p.accept(token.DO)
			loop.Block = p.parseBlockStatement(token.END)
			p.nextToken()
		} else if p.peekTokenIs(token.LBRACE) {
			p.accept(token.LBRACE)
			loop.Block = p.parseBlockStatement(token.RBRACE)
			p.nextToken()
		}
	} else {
		p.Error(p.curToken.Type, "unexpected loop start", token.WHILE, token.LOOP)
		return nil
	}

	return loop
}

func (p *parser) parseBreakStatement() *ast.BreakStatement {
	if p.trace {
		defer un(trace(p, "parseBreakExpression"))
	}
	stmt := &ast.BreakStatement{}
	if p.peekTokenIs(token.IF) {
		p.accept(token.IF)
		p.nextToken()
		stmt.Condition = p.parseExpression(precLowest)
	} else if p.peekTokenIs(token.UNLESS) {
		p.accept(token.UNLESS)
		p.nextToken()
		stmt.Condition = p.parseExpression(precLowest)
		stmt.Unless = true
	}
	return stmt
}

func (p *parser) parseFunctionLiteral() ast.Expression {
	if p.trace {
		defer un(trace(p, "parseFunctionLiteral"))
	}
	lit := &ast.FunctionLiteral{}

	if !p.peekTokenOneOf(token.IDENT, token.CONST) && !p.peekToken.Type.IsOperator() {
		p.Error(p.peekToken.Type, "", token.IDENT, token.CONST)
		return nil
	}

	if p.peekTokenOneOf(token.IDENT, token.CONST) {
		p.acceptOneOf(token.IDENT, token.CONST)
		lit.Name = &ast.Identifier{Constant: p.currentTokenIs(token.CONST), Value: p.curToken.Literal}
	} else {
		p.nextToken()
		lit.Name = &ast.Identifier{Constant: p.currentTokenIs(token.CONST), Value: p.curToken.Literal}
	}

	lit.Parameters = p.parseFunctionParameters(token.LPAREN, token.RPAREN)

	if p.currentTokenOneOf(token.CAPTURE, token.AND) {
		capture := p.parseBlockCapture()
		if capture == nil {
			return nil
		}
		lit.CapturedBlock = capture.(*ast.BlockCapture)
		if p.peekTokenIs(token.RPAREN) {
			p.acceptOneOf(token.RPAREN)
		}
	}

	if !p.acceptOneOf(token.NEWLINE, token.SEMICOLON) {
		return nil
	}
	lit.Body = p.parseBlockStatement(token.END)
	if !p.accept(token.END) {
		return nil
	}
	inspect := func(n ast.Node) bool {
		if _, ok := n.(*ast.Splat); ok {
			// skip walking the splat since we don't want to evaluate it
			return false
		}
		x, ok := n.(*ast.Assignment)
		if !ok {
			return true
		}
		switch left := x.Left.(type) {
		case *ast.Identifier:
			if left.Constant {
				p.errors = append(p.errors, fmt.Errorf("dynamic constant assignment"))
			}
		case ast.ExpressionList:
			for _, expr := range left {
				if ident, ok := expr.(*ast.Identifier); ok {
					if ident.Constant {
						p.errors = append(p.errors, fmt.Errorf("dynamic constant assignment"))
					}
				}
			}
		}
		return true
	}
	ast.Inspect(lit.Body, inspect)
	return lit
}

func (p *parser) parseFunctionParameters(startToken, endToken token.Type) []*ast.FunctionParameter {
	if p.trace {
		defer un(trace(p, "parseFunctionParameters"))
	}
	hasDelimiters := false
	if p.peekTokenIs(startToken) {
		hasDelimiters = true
		p.accept(startToken)
	}

	identifiers := []*ast.FunctionParameter{}

	if !hasDelimiters && p.peekTokenIs(endToken) {
		p.Error(p.peekToken.Type, "", token.NEWLINE, token.SEMICOLON)
		return nil
	}

	if hasDelimiters && p.peekTokenIs(endToken) {
		p.accept(endToken)
		return identifiers
	}

	if !hasDelimiters && p.peekTokenOneOf(token.NEWLINE, token.SEMICOLON) {
		return identifiers
	}

	got_splat := false
	if p.peekTokenIs(token.ASTERISK) {
		got_splat = true
		p.accept(token.ASTERISK)
	}
	if p.peekTokenOneOf(token.CAPTURE, token.AND) {
		p.acceptOneOf(token.CAPTURE, token.AND)
		return identifiers
	}
	p.accept(token.IDENT)

	ident := &ast.FunctionParameter{Name: &ast.Identifier{Constant: p.currentTokenIs(token.CONST), Value: p.curToken.Literal}}
	if got_splat {
		ident.Splat = true
		got_splat = false
	}
	if p.peekTokenIs(token.ASSIGN) {
		p.consume(token.ASSIGN)
		ident.Default = p.parseExpression(precAssignment)
	}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.accept(token.COMMA)
		if p.peekTokenIs(token.ASTERISK) {
			got_splat = true
			p.accept(token.ASTERISK)
		}
		if p.peekTokenOneOf(token.CAPTURE, token.AND) {
			p.acceptOneOf(token.CAPTURE, token.AND)
			return identifiers
		}
		p.accept(token.IDENT)
		ident := &ast.FunctionParameter{Name: &ast.Identifier{Constant: p.currentTokenIs(token.CONST), Value: p.curToken.Literal}}
		if got_splat {
			ident.Splat = true
			got_splat = false
		}
		if p.peekTokenIs(token.ASSIGN) {
			p.consume(token.ASSIGN)
			ident.Default = p.parseExpression(precPrefix)
		}
		identifiers = append(identifiers, ident)
	}

	if !hasDelimiters && p.peekTokenIs(endToken) {
		p.Error(p.peekToken.Type, "no delimiters but end delimiter found", endToken)
		return nil
	}

	if hasDelimiters && p.peekTokenIs(endToken) {
		p.accept(endToken)
	}

	return identifiers
}

func (p *parser) parseBlockStatement(t ...token.Type) *ast.BlockStatement {
	if p.trace {
		defer un(trace(p, "parseBlockStatement"))
	}
	terminatorTokens := append(
		[]token.Type{
			token.END,
		},
		t...,
	)
	block := &ast.BlockStatement{}
	block.Statements = []ast.Statement{}

	for !p.peekTokenOneOf(terminatorTokens...) {
		if p.peekTokenIs(token.EOF) {
			p.Error(p.peekToken.Type, "", token.EOF)
			return block
		}
		p.nextToken()
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}

	return block
}

func (p *parser) parseMethodCall(context ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace(p, "parseMethodCall"))
	}
	contextCallExpression := &ast.ContextCallExpression{Context: context}

	p.nextToken()

	if !p.currentTokenOneOf(token.IDENT, token.CLASS) && !p.curToken.Type.IsOperator() {
		p.Error(p.curToken.Type, "", token.IDENT, token.CLASS)
		return nil
	}

	function := &ast.Identifier{Constant: p.currentTokenIs(token.CONST), Value: p.curToken.Literal}
	contextCallExpression.Function = function

	if p.peekTokenOneOf(token.SEMICOLON, token.NEWLINE, token.EOF, token.DOT, token.RPAREN, token.SQMARK) {
		contextCallExpression.Arguments = []ast.Expression{}
		return contextCallExpression
	}

	if p.peekTokenIs(token.LPAREN) {
		p.accept(token.LPAREN)
		p.nextToken()
		contextCallExpression.Arguments = p.parseExpressionList(token.RPAREN)
		if p.peekTokenOneOf(token.LBRACE, token.DO) {
			p.acceptOneOf(token.LBRACE, token.DO)
			contextCallExpression.Block = p.parseBlock().(*ast.BlockExpression)
		}
		return contextCallExpression
	}

	if p.peekTokenOneOf(append(tokensNotPossibleInCallArgs, token.RBRACE)...) {
		return contextCallExpression
	}

	p.nextToken()

	contextCallExpression.Arguments = p.parseCallArguments(
		token.LBRACE, token.DO,
	)
	if p.currentTokenOneOf(token.LBRACE, token.DO) {
		contextCallExpression.Block = p.parseBlock().(*ast.BlockExpression)
	}
	return contextCallExpression
}

func (p *parser) parseContextCallExpression(context ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace(p, "parseContextCallExpression"))
	}
	contextCallExpression := &ast.ContextCallExpression{Context: context}
	if p.currentTokenOneOf(token.DOT) {
		p.nextToken()
	}

	if !p.currentTokenOneOf(token.IDENT, token.CLASS) {
		p.Error(p.curToken.Type, "", token.IDENT, token.CLASS)
		return nil
	}

	function := p.parseIdentifier()
	ident := function.(*ast.Identifier)
	contextCallExpression.Function = ident

	if p.peekTokenOneOf(token.SEMICOLON, token.NEWLINE, token.DOT) {
		contextCallExpression.Arguments = []ast.Expression{}
		return contextCallExpression
	}

	if p.peekTokenIs(token.LPAREN) {
		p.accept(token.LPAREN)
		p.nextToken()
		contextCallExpression.Arguments = p.parseExpressionList(token.RPAREN)
		if p.peekTokenOneOf(token.LBRACE, token.DO) {
			p.acceptOneOf(token.LBRACE, token.DO)
			contextCallExpression.Block = p.parseBlock().(*ast.BlockExpression)
		}
		return contextCallExpression
	}

	if p.peekTokenOneOf(append(tokensNotPossibleInCallArgs, token.RBRACE)...) {
		return contextCallExpression
	}

	p.nextToken()
	contextCallExpression.Arguments = p.parseCallArguments(
		token.LBRACE, token.DO,
	)
	if p.currentTokenOneOf(token.LBRACE, token.DO) {
		contextCallExpression.Block = p.parseBlock().(*ast.BlockExpression)
	}
	return contextCallExpression
}

func (p *parser) parseCallArgument(function ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace(p, "parseCallArgument"))
	}
	ident, ok := function.(*ast.Identifier)
	if !ok {
		// method call on any other object
		return p.parseContextCallExpression(function)
	}
	exp := &ast.ContextCallExpression{Function: ident}
	if p.currentTokenOneOf(token.LBRACE, token.DO) {
		exp.Block = p.parseBlock().(*ast.BlockExpression)
		return exp
	}

	exp.Arguments = p.parseExpressionList(token.SEMICOLON, token.NEWLINE)
	if p.peekTokenOneOf(token.LBRACE, token.DO) {
		p.acceptOneOf(token.LBRACE, token.DO)
		exp.Block = p.parseBlock().(*ast.BlockExpression)
	}
	return exp
}

func (p *parser) parseCallBlock(function ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace(p, "parseCallBlock"))
	}

	exp := &ast.ContextCallExpression{}
	exp.Block = p.parseBlock().(*ast.BlockExpression)
	switch fn := function.(type) {
	case *ast.Identifier:
		exp.Function = fn
		return exp
	case *ast.InfixExpression:
		ident, ok := fn.Right.(*ast.Identifier)
		if !ok {
			break
		}
		exp.Function = ident
		fn.Right = exp
		return fn
	}
	msg := fmt.Errorf("could not parse call expression: expected identifier, got token '%T'", function)
	p.errors = append(p.errors, msg)
	return nil
}

// func (p *parser) parseCallExpression(function ast.Expression) ast.Expression {

func (p *parser) parseCallExpressionWithParens(function ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace(p, "parseCallExpressionWithParens"))
	}
	ident, ok := function.(*ast.Identifier)
	if !ok {
		msg := fmt.Errorf("could not parse call expression: expected identifier, got token '%T'", function)
		p.errors = append(p.errors, msg)
		return nil
	}
	exp := &ast.ContextCallExpression{Function: ident}
	p.nextToken()
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	if p.peekTokenOneOf(token.LBRACE, token.DO) {
		p.acceptOneOf(token.LBRACE, token.DO)
		exp.Block = p.parseBlock().(*ast.BlockExpression)
	}
	return exp
}

func (p *parser) parseCallArguments(end ...token.Type) []ast.Expression {
	if p.trace {
		defer un(trace(p, "parseCallArguments"))
	}
	list := []ast.Expression{}
	if p.currentTokenOneOf(end...) {
		return list
	}

	list = append(list, p.parseExpression(precAssignment))

	for p.peekTokenIs(token.COMMA) {
		p.consume(token.COMMA)
		list = append(list, p.parseExpression(precAssignment))
	}

	if p.peekTokenOneOf(end...) {
		p.acceptOneOf(end...)
	}

	return list
}

func (p *parser) parseExpressionList(end ...token.Type) []ast.Expression {
	if p.trace {
		defer un(trace(p, "parseExpressionList"))
	}
	list := []ast.Expression{}
	if p.currentTokenOneOf(end...) {
		return list
	}

	next := p.parseExpression(precIfUnless)
	if elist, ok := next.(ast.ExpressionList); ok {
		if p.peekTokenOneOf(end...) {
			p.acceptOneOf(end...)
		}
		return elist
	}
	list = append(list, next)

	if p.peekTokenOneOf(end...) {
		p.acceptOneOf(end...)
	}

	return list
}

func (p *parser) peekPrecedence() int {
	return precedenceForToken(p.peekToken.Type)
}

func (p *parser) curPrecedence() int {
	return precedenceForToken(p.curToken.Type)
}

func precedenceForToken(t token.Type) int {
	if prec, ok := precedences[t]; ok {
		return prec
	}
	return precLowest
}

func (p *parser) currentTokenOneOf(types ...token.Type) bool {
	for _, typ := range types {
		if p.curToken.Type == typ {
			return true
		}
	}
	return false
}

func (p *parser) currentTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *parser) peekTokenOneOf(types ...token.Type) bool {
	for _, typ := range types {
		if p.peekToken.Type == typ {
			return true
		}
	}
	return false
}

func (p *parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

// accept moves to the next Token
// if it's from the valid set.
func (p *parser) accept(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.Error(p.peekToken.Type, "", t)
	return false
}

// acceptOneOf moves to the next Token
// if it's from the valid set.
func (p *parser) acceptOneOf(t ...token.Type) bool {
	if p.peekTokenOneOf(t...) {
		p.nextToken()
		return true
	}

	p.Error(p.peekToken.Type, "", t...)
	return false
}

// consume consumes the next token
// if it's from the valid set.
func (p *parser) consume(t token.Type) bool {
	isRightToken := p.accept(t)
	if isRightToken {
		p.nextToken()
	}
	return isRightToken
}
