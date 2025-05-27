package ast

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/ast/infix"
)

// Node represents a node within the AST
//
// All node types implement the Node interface.
type Node interface {
	node() // marks this as a node
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Coder interface {
	Code() string
}

func getParentInfo(N int) (string, int) {
	parent, _, _, _ := runtime.Caller(1 + N)
	info := runtime.FuncForPC(parent)
	file, line := info.FileLine(parent)
	return file, line
}

func code_or_string(thing any) string {
	if thing, ok := thing.(Coder); ok {
		return thing.Code()
	}
	if thing, ok := thing.(fmt.Stringer); ok {
		file, line := getParentInfo(1)
		fmt.Printf("# AST %T does not implement Coder at %s:%d\n", thing, file, line)
		return thing.String()
	}
	panic(fmt.Sprintf("%T does not have a Code() or String() method", thing))
}

////////////////////////////////////////////////////////////////////////////////

// A Program node is the root node within the AST.
type Program struct {
	Statements []Statement
}

func (p *Program) node() {}

func (p *Program) String() string {
	stmts := make([]string, len(p.Statements))
	for i, s := range p.Statements {
		if s != nil {
			stmts[i] = s.String()
		}
	}
	return strings.Join(stmts, "\n")
}

var (
	_ Node = &Program{}
)

// A ReturnStatement represents a return node which yields another Expression.
type ReturnStatement struct {
	ReturnValue Expression
}

func (rs *ReturnStatement) node()          {}
func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) String() string {
	var out strings.Builder
	out.WriteString("return ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	return out.String()
}

var (
	_ Node      = &ReturnStatement{}
	_ Statement = &ReturnStatement{}
)

// An ExpressionStatement is a Statement wrapping an Expression
type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) node()          {}
func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (es *ExpressionStatement) Code() string {
	if es.Expression == nil {
		return ""
	}
	return code_or_string(es.Expression)
}

var (
	_ Node      = &ExpressionStatement{}
	_ Statement = &ExpressionStatement{}
	// _ Expression = &ExpressionStatement{}
	_ Coder = &ExpressionStatement{}
)

// BlockStatement represents a list of statements
type BlockStatement struct {
	Statements []Statement
}

func (bs *BlockStatement) node()          {}
func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) String() string {
	var out strings.Builder
	for _, s := range bs.Statements {
		if s != nil {
			out.WriteString(s.String())
		}
	}
	return out.String()
}

// Code returns the code representation of the block
func (bs *BlockStatement) Code() string {
	var out strings.Builder
	statement_strings := make([]string, 0)
	for _, s := range bs.Statements {
		if s == nil {
			continue
		}
		statement_string := code_or_string(s)
		statement_strings = append(statement_strings, statement_string)
	}

	if len(statement_strings) == 0 {
		return ""
	}

	out.WriteString(strings.Join(statement_strings, ";"))
	return out.String()
}

var (
	_ Node      = &BlockStatement{}
	_ Statement = &BlockStatement{}
	_ Coder     = &BlockStatement{}
)

// A BreakStatement represents a break statement
type BreakStatement struct {
	Condition Expression
	Unless    bool
}

func (bs *BreakStatement) node()          {}
func (bs *BreakStatement) statementNode() {}

func (bs *BreakStatement) String() string {
	return "<<<Break Statement>>>"
}

func (bs *BreakStatement) Code() string {
	var out strings.Builder
	out.WriteString("break ")
	if bs.Unless {
		out.WriteString("unless ")
	} else {
		out.WriteString("if ")
	}
	if bs.Condition != nil {
		out.WriteString(code_or_string(bs.Condition))
	}
	return out.String()
}

var (
	_ Node      = &BreakStatement{}
	_ Statement = &BreakStatement{}
	_ Coder     = &BreakStatement{}
)

// Assignment represents a generic assignment
type Assignment struct {
	Left  Expression
	Right Expression
}

func (a *Assignment) node()           {}
func (a *Assignment) expressionNode() {}

func (a *Assignment) String() string {
	var out strings.Builder
	out.WriteString(a.Left.String())
	out.WriteString(" = ")
	out.WriteString(a.Right.String())
	return out.String()
}

func (a *Assignment) Code() string {
	var out strings.Builder
	out.WriteString(code_or_string(a.Left))
	out.WriteString(" = ")
	out.WriteString(code_or_string(a.Right))
	return out.String()
}

var (
	_ Node       = &Assignment{}
	_ Expression = &Assignment{}
	_ Coder      = &Assignment{}
)

// MultiAssignment represents multiple variables on the left-hand side
type MultiAssignment struct {
	Variables []*Identifier
	Values    []Expression
}

func (m *MultiAssignment) node()           {}
func (m *MultiAssignment) expressionNode() {}

func (m *MultiAssignment) String() string {
	var out strings.Builder
	vars := make([]string, len(m.Variables))
	for i, v := range m.Variables {
		vars[i] = v.Value
	}
	out.WriteString(strings.Join(vars, ", "))
	out.WriteString(" = ")
	values := make([]string, len(m.Values))
	for i, v := range m.Values {
		values[i] = v.String()
	}
	out.WriteString(strings.Join(values, ", "))
	return out.String()
}

var (
	_ Node       = &MultiAssignment{}
	_ Expression = &MultiAssignment{}
)

type Identifier struct {
	Value string
}

func (i *Identifier) node()           {}
func (i *Identifier) expressionNode() {}

func (i *Identifier) String() string { return "<<<Identifier: " + i.Value + ">>>" }
func (i *Identifier) Code() string   { return i.Value }

var (
	_ Node       = &Identifier{}
	_ Expression = &Identifier{}
	_ Coder      = &Identifier{}
)

func (i *Identifier) IsConstant() bool {
	if len(i.Value) == 0 {
		return false
	}
	return i.Value[0] >= 'A' && i.Value[0] <= 'Z'
}

func (i *Identifier) IsGlobal() bool {
	if len(i.Value) == 0 {
		return false
	}
	return i.Value[0] == '$'
}

// IntegerLiteral represents an integer in the AST
type IntegerLiteral struct {
	Value int64
}

func (il *IntegerLiteral) node()           {}
func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) String() string  { return fmt.Sprintf("%d", il.Value) }
func (il *IntegerLiteral) Code() string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("%d", il.Value))
	return out.String()
}

var (
	_ Node       = &IntegerLiteral{}
	_ Expression = &IntegerLiteral{}
	_ Coder      = &IntegerLiteral{}
)

// FloatLiteral represents a float in the AST
type FloatLiteral struct {
	Value float64
}

func (fl *FloatLiteral) node()           {}
func (fl *FloatLiteral) expressionNode() {}

func (fl *FloatLiteral) String() string { return fmt.Sprintf("%f", fl.Value) }

var (
	_ Node       = &FloatLiteral{}
	_ Expression = &FloatLiteral{}
)

// StringLiteral represents a double quoted string in the AST
type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) node()           {}
func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) String() string { return sl.Value }
func (sl *StringLiteral) Code() string {
	var out strings.Builder
	out.WriteString("\"")
	out.WriteString(sl.Value)
	out.WriteString("\"")
	return out.String()
}

var (
	_ Node       = &StringLiteral{}
	_ Expression = &StringLiteral{}
	_ Coder      = &StringLiteral{}
)

// Comment represents a double quoted string in the AST
type Comment struct {
	Value string
}

func (c *Comment) node()          {}
func (c *Comment) statementNode() {}

func (c *Comment) String() string {
	if strings.HasPrefix(c.Value, "#") {
		return c.Value
	}
	return "# " + strings.TrimLeft(c.Value, " ")
}

var (
	_ Node      = &Comment{}
	_ Statement = &Comment{}
)

// SymbolLiteral represents a symbol within the AST
type SymbolLiteral struct {
	Value string
}

func (s *SymbolLiteral) node()           {}
func (s *SymbolLiteral) expressionNode() {}

func (s *SymbolLiteral) String() string { return ":" + s.Value }

var (
	_ Node       = &SymbolLiteral{}
	_ Expression = &SymbolLiteral{}
)

// ConditionalExpression represents an if expression within the AST
type ConditionalExpression struct {
	Unless      bool // true = unless, false = if
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ce *ConditionalExpression) node()           {}
func (ce *ConditionalExpression) expressionNode() {}

func (ce *ConditionalExpression) String() string {
	return "<<<ConditionalExpression>>>"
}

func (ce *ConditionalExpression) Code() string {
	var out strings.Builder
	if ce.Unless {
		out.WriteString("unless ")
	} else {
		out.WriteString("if ")
	}
	out.WriteString(code_or_string(ce.Condition))
	out.WriteString("; ")
	out.WriteString(code_or_string(ce.Consequence))
	if ce.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(code_or_string(ce.Alternative))
	}
	out.WriteString(" end")
	return out.String()
}

var (
	_ Node       = &ConditionalExpression{}
	_ Expression = &ConditionalExpression{}
	_ Coder      = &ConditionalExpression{}
)

// A LoopExpression represents a loop
type LoopExpression struct {
	// Condition Expression
	Block *BlockStatement
}

func (ce *LoopExpression) node()           {}
func (ce *LoopExpression) expressionNode() {}

func (ce *LoopExpression) String() string {
	return "<<<LoopExpression>>>"
}

func (ce *LoopExpression) Code() string {
	var out strings.Builder
	out.WriteString("loop {")
	out.WriteString(code_or_string(ce.Block))
	out.WriteString("}")
	return out.String()
}

var (
	_ Node       = &LoopExpression{}
	_ Expression = &LoopExpression{}
	_ Coder      = &LoopExpression{}
)

// ExpressionList represents a list of expressions within the AST divided by commas
type ExpressionList []Expression

func (el ExpressionList) node()           {}
func (el ExpressionList) expressionNode() {}

func (el ExpressionList) String() string {
	var out strings.Builder
	elements := []string{}
	for _, e := range el {
		elements = append(elements, e.String())
	}
	out.WriteString(strings.Join(elements, ", "))
	return out.String()
}

func (el ExpressionList) Code() string {
	var out strings.Builder
	elements := []string{}
	for _, e := range el {
		needs_parens := false
		// for now put everything except Identifiers in parentheses
		switch e.(type) {
		case *Identifier:
		case *StringLiteral:
		//
		default:
			needs_parens = true
		}
		var element string
		if needs_parens {
			element = "(" + code_or_string(e) + ")"
		} else {
			element = code_or_string(e)
		}
		elements = append(elements, element)
	}
	out.WriteString(strings.Join(elements, ", "))
	return out.String()
}

var (
	_ Node       = ExpressionList{}
	_ Expression = ExpressionList{}
	_ Coder      = ExpressionList{}
)

// ArrayLiteral represents an Array literal within the AST
type ArrayLiteral struct {
	Elements []Expression
}

func (al *ArrayLiteral) node()           {}
func (al *ArrayLiteral) expressionNode() {}

func (al *ArrayLiteral) String() string {
	var out strings.Builder
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

var (
	_ Node       = &ArrayLiteral{}
	_ Expression = &ArrayLiteral{}
)

// HashLiteral represents an Hash literal within the AST
type HashLiteral struct {
	Map map[Expression]Expression
}

func (hl *HashLiteral) node()           {}
func (hl *HashLiteral) expressionNode() {}

func (hl *HashLiteral) String() string {
	var out strings.Builder
	elements := []string{}
	for key, val := range hl.Map {
		elements = append(elements, fmt.Sprintf("%q => %q", key.String(), val.String()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")
	return out.String()
}

var (
	_ Node       = &HashLiteral{}
	_ Expression = &HashLiteral{}
)

// RangeLiteral represents a range literal within the AST
type RangeLiteral struct {
	Left      Expression
	Right     Expression
	Inclusive bool
}

func (rl *RangeLiteral) node()           {}
func (rl *RangeLiteral) expressionNode() {}

// String returns the string representation of the range
func (rl *RangeLiteral) String() string {
	var out strings.Builder
	out.WriteString("<<<RangeLiteral>>>")
	return out.String()
}

func (rl *RangeLiteral) Code() string {
	var out strings.Builder
	out.WriteString(code_or_string(rl.Left))
	out.WriteString(" ")
	if rl.Inclusive {
		out.WriteString("..")
	} else {
		out.WriteString("...")
	}
	out.WriteString(" ")
	out.WriteString(code_or_string(rl.Right))
	return out.String()
}

var (
	_ Node       = &RangeLiteral{}
	_ Expression = &RangeLiteral{}
	_ Coder      = &RangeLiteral{}
)

// A FunctionLiteral represents a function definition in the AST
type FunctionLiteral struct {
	Name       string
	Parameters []*FunctionParameter
	Body       *BlockStatement
}

func (fl *FunctionLiteral) node()           {}
func (fl *FunctionLiteral) expressionNode() {}

const BODY_STRING_LIMIT = 30

func (fl *FunctionLiteral) String() string {
	var out strings.Builder
	out.WriteString("<<<FunctionLiteral, name=\"")
	out.WriteString(fl.Name)
	args := []string{}
	for _, a := range fl.Parameters {
		args = append(args, a.String())
	}
	out.WriteString("\", parameters=[")
	if len(args) == 0 {
		out.WriteString("]")
	} else {
		out.WriteString("\"")
		out.WriteString(strings.Join(args, "\", \""))
		out.WriteString("\"]")
	}
	out.WriteString(", body=\"")
	body_string := fl.Body.String()
	if len(body_string) > BODY_STRING_LIMIT {
		body_string = body_string[:BODY_STRING_LIMIT] + "..."
	}
	out.WriteString(body_string)
	out.WriteString("\">>>")
	return out.String()
}

func (fl *FunctionLiteral) Code() string {
	var out strings.Builder
	if strings.HasPrefix(fl.Name, "__") {
		out.WriteString("-> ")
		out.WriteString("(")
		args := []string{}
		for _, a := range fl.Parameters {
			args = append(args, code_or_string(a))
		}
		out.WriteString(strings.Join(args, ", "))
		out.WriteString(") {")
		out.WriteString(code_or_string(fl.Body))
		out.WriteString("}")
	} else {
		out.WriteString("def ")
		out.WriteString(fl.Name)
		out.WriteString("(")
		args := []string{}
		for _, a := range fl.Parameters {
			args = append(args, code_or_string(a))
		}
		out.WriteString(strings.Join(args, ", "))
		out.WriteString(")")
		out.WriteString("\n")
		body_string := code_or_string(fl.Body)
		for _, line := range strings.Split(body_string, "\n") {
			out.WriteString("    ")
			out.WriteString(line)
			out.WriteString("\n")
		}
		out.WriteString("end")
	}
	return out.String()
}

var (
	_ Node       = &FunctionLiteral{}
	_ Expression = &FunctionLiteral{}
	_ Coder      = &FunctionLiteral{}
)

// A FunctionParameter represents a parameter in a function literal
type FunctionParameter struct {
	Name    string
	Default Expression
	Splat   bool
}

func (f *FunctionParameter) node()           {}
func (f *FunctionParameter) expressionNode() {}

func (f *FunctionParameter) String() string {
	return "<<<FunctionParameter, name=\"" + f.Name + "\">>>"
}
func (f *FunctionParameter) Code() string {
	var out strings.Builder
	if f.Splat {
		out.WriteString("*")
	}
	out.WriteString(f.Name)
	if f.Default != nil {
		out.WriteString(" = ")
		out.WriteString(f.Default.String())
	}
	return out.String()
}

var (
	_ Node       = &FunctionParameter{}
	_ Expression = &FunctionParameter{}
	_ Coder      = &FunctionParameter{}
)

// A Splat represents a splat operator in the AST
type Splat struct {
	Value Expression
}

func (s *Splat) node()           {}
func (s *Splat) expressionNode() {}

func (s *Splat) String() string {
	var out strings.Builder
	out.WriteString("*")
	out.WriteString(s.Value.String())
	return out.String()
}

var (
	_ Node       = &Splat{}
	_ Expression = &Splat{}
)

// An IndexExpression represents an array or hash access in the AST
type IndexExpression struct {
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) node()           {}
func (ie *IndexExpression) expressionNode() {}

func (ie *IndexExpression) String() string {
	var out strings.Builder
	out.WriteString("<<<IndexExpression>>>")
	return out.String()
}

func (ie *IndexExpression) Code() string {
	var out strings.Builder
	out.WriteString(code_or_string(ie.Left))
	out.WriteString("[")
	out.WriteString(code_or_string(ie.Index))
	out.WriteString("]")
	return out.String()
}

var (
	_ Node       = &IndexExpression{}
	_ Expression = &IndexExpression{}
	_ Coder      = &IndexExpression{}
)

// A ContextCallExpression represents a method call on a given Context
type ContextCallExpression struct {
	Context   Expression   // The left-hand side expression
	Function  string       // The function to call
	Arguments []Expression // Normal arguments
	Block     Expression   // Block argument, if any
}

func (ce *ContextCallExpression) node()           {}
func (ce *ContextCallExpression) expressionNode() {}

func (ce *ContextCallExpression) String() string {
	var out strings.Builder
	out.WriteString("ContextCallExpression<functions=\"")
	out.WriteString(ce.Function)
	out.WriteString("\">")
	return out.String()
}

func (ce *ContextCallExpression) Code() string {
	var out strings.Builder
	if ce.Context != nil {
		needs_parens := false
		if _, ok := ce.Context.(*Identifier); !ok {
			needs_parens = true
		}
		if needs_parens {
			out.WriteString("(")
		}
		out.WriteString(code_or_string(ce.Context))
		if needs_parens {
			out.WriteString(")")
		}
		out.WriteString(".")
	}
	out.WriteString(ce.Function)
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, code_or_string(a))
	}
	if len(args) > 0 {
		out.WriteString("(")
		out.WriteString(strings.Join(args, ", "))
		out.WriteString(")")
	}
	if ce.Block != nil {
		out.WriteString(" ")
		out.WriteString(code_or_string(ce.Block))
	}
	return out.String()
}

var (
	_ Node       = &ContextCallExpression{}
	_ Expression = &ContextCallExpression{}
	_ Coder      = &ContextCallExpression{}
)

// PrefixExpression represents a prefix operator
type PrefixExpression struct {
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) node()           {}
func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) String() string {
	var out strings.Builder
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// An InfixExpression represents an infix operator in the AST
type InfixExpression struct {
	Left     Expression
	Operator infix.Infix
	Right    Expression
}

func (oe *InfixExpression) node()           {}
func (oe *InfixExpression) expressionNode() {}

func (oe *InfixExpression) String() string {
	var out strings.Builder
	out.WriteString("<<<InfixExpression>>>")
	return out.String()
}

func (oe *InfixExpression) Code() string {
	var out strings.Builder
	out.WriteString(code_or_string(oe.Left))
	out.WriteString(" ")
	out.WriteString(code_or_string(oe.Operator))
	out.WriteString(" ")
	out.WriteString(code_or_string(oe.Right))
	return out.String()
}

var (
	_ Node       = &InfixExpression{}
	_ Expression = &InfixExpression{}
	_ Coder      = &InfixExpression{}
)
